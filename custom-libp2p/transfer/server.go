package transfer

import (
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/config"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/file_handler"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
	gostream "github.com/libp2p/go-libp2p-gostream"
	"github.com/libp2p/go-libp2p/core/peer"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type server struct {
	listener         *net.Listener
	server           *http.Server
	servingMetafiles *sync.Map // key: file_SHA256, value: *models.FileMetadata
}

// initServer initializes the server
func initServer(node *models.Node) error {
	// TODO: Load serving from DB
	_server := _server
	_server.servingMetafiles = &sync.Map{}

	node.SetStreamHandler(holePunchSyncStreamProtocolID, handleHolePunchSyncStream)
	listener, err := gostream.Listen(node, protocolID)
	if err != nil {
		return err
	}
	_server.listener = &listener

	if _server.server == nil {
		go func() {
			http.HandleFunc("/file", handleFileDownloadRequest)
			http.HandleFunc("/message", handleMessageRequest)
			_server.server = &http.Server{}
			err = _server.server.Serve(listener)
			if err != nil {
				notifier.QueueWarning(&models.Warning{
					Error: err.Error(),
					Info:  fmt.Sprintf("Server %v", listener.Addr()),
				})
			}
		}()
	} else {
		go func() {
			_server.server = &http.Server{}
			err = _server.server.Serve(listener)
			if err != nil {
				notifier.QueueWarning(&models.Warning{
					Error: err.Error(),
					Info:  fmt.Sprintf("Server %v", listener.Addr()),
				})
			}
		}()
	}
	return nil
}

// closeServer closes the server
func closeServer() error {
	_serverLock.Lock()
	defer _serverLock.Unlock()
	if !ServerIsRunning() {
		return ErrServerNotRunning
	}
	err := _server.server.Close()
	if err != nil {
		return err
	}
	listener := *(_server.listener)
	err = listener.Close()
	if err != nil {
		return err
	}
	_server.servingMetafiles = nil
	return nil
}

// ServeFile serves the file to peer, returns file SHA256Sum and error
func ServeFile(filePath string, peerId peer.ID) (string, error) {
	if !ServerIsRunning() {
		return "", ErrServerNotRunning
	}
	file, err := file_handler.GetFile(filePath)
	if err != nil {
		return "", err
	}
	sha256sum, err := file_handler.GetSHA256Sum(file)
	if err != nil {
		return "", err
	}
	_, exists := _server.servingMetafiles.Load(sha256sum)
	if !exists {
		fSize, err := file_handler.GetFileSize(filePath)
		if err != nil {
			return sha256sum, err
		}
		_server.servingMetafiles.Store(
			sha256sum,
			&models.FileMetadata{
				FileName:   filepath.Base(filePath),
				FileSha256: sha256sum,
				FileSize:   fSize,
				SpecificData: &models.FileMetadata_ServerFileInfo{ServerFileInfo: &models.ServerFileInfo{
					BasePath:            filepath.Dir(filePath),
					AuthorizedAccessors: []string{},
				}},
			},
		)
	}
	value, _ := _server.servingMetafiles.Load(sha256sum)
	metadata := value.(*models.FileMetadata)
	metadata.GetServerFileInfo().AuthorizedAccessors = append(metadata.GetServerFileInfo().AuthorizedAccessors,
		peerId.String())
	err = sendFileMessage(peerId, metadata)
	return sha256sum, err
}

// StopServingFile stops serving the file
func StopServingFile(fileSHA256 string) error {
	if !ServerIsRunning() {
		return ErrServerNotRunning
	}
	value, found := _server.servingMetafiles.Load(fileSHA256)
	if !found {
		return ErrFileNotServing
	}
	metadata := value.(*models.FileMetadata)
	err := file_handler.CloseFile(getFilePath(metadata))
	if err != nil {
		return err
	}
	_server.servingMetafiles.Delete(fileSHA256)
	return nil
}

// handleFileDownloadRequest handles file download requests
func handleFileDownloadRequest(writer http.ResponseWriter, request *http.Request) {
	_clientAddr, err := peer.Decode(request.RemoteAddr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	requestSHA256 := request.URL.Query().Get("sha256")
	value, exists := _server.servingMetafiles.Load(requestSHA256)
	if !exists {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	metadata := value.(*models.FileMetadata)
	allowed := slices.Contains(metadata.GetServerFileInfo().AuthorizedAccessors, _clientAddr.String())
	if !allowed {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	if request.Method == http.MethodDelete {
		handleDeleteFile(writer, requestSHA256, metadata, _clientAddr)
		return
	} else if request.Method == http.MethodGet {
		handleGetFile(writer, request, metadata)
		return
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// handleGetFile handles get request for file
func handleGetFile(writer http.ResponseWriter, request *http.Request, metadata *models.FileMetadata) {
	var err error
	strOffset := request.URL.Query().Get("offset")
	offset := 0
	if strOffset != "" {
		offset, err = strconv.Atoi(strOffset)
		if err != nil {
			notifier.QueueWarning(&models.Warning{
				Error: err.Error(),
				Info:  "Error in parsing offset",
			})
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if offset < 0 {
		notifier.QueueWarning(&models.Warning{
			Info: fmt.Sprintf("Invalid Offset: %v", offset),
		})
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if uint64(offset) >= metadata.FileSize {
		notifier.QueueWarning(&models.Warning{
			Info: fmt.Sprintf("Invalid Offset: %v", offset),
		})
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	file, err := file_handler.GetFile(getFilePath(metadata))
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  fmt.Sprintf("Failed to serve file: %v to peer: %v", metadata.FileSha256, request.RemoteAddr),
		})
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	fReader := io.NewSectionReader(file, int64(offset), int64(metadata.FileSize-uint64(offset)))
	n, err := io.CopyN(writer, fReader, int64(metadata.FileSize-uint64(offset)))
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  fmt.Sprintf("Failed to serve file: %v to peer: %v", metadata.FileSha256, request.RemoteAddr),
		})
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	notifier.QueueInfo(fmt.Sprintf("Wrote File: %v bytes (%v)", n+int64(offset), metadata.FileSize))
}

// handleDeleteFile handles delete request for file
func handleDeleteFile(writer http.ResponseWriter, requestSHA256 string, metadata *models.FileMetadata,
	clientAddr peer.ID) {
	removePeer(metadata, clientAddr)
	if len(metadata.GetServerFileInfo().GetAuthorizedAccessors()) == 0 {
		err := StopServingFile(requestSHA256)
		if err != nil {
			notifier.QueueWarning(&models.Warning{
				Error: err.Error(),
				Info: fmt.Sprintf("Failed to stop serving file: %v to peer: %v", metadata.FileSha256,
					clientAddr),
			})
			writer.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			writer.WriteHeader(http.StatusOK)
			return
		}
	}
	writer.WriteHeader(http.StatusOK)
}

// handleMessageRequest handles message requests
func handleMessageRequest(writer http.ResponseWriter, request *http.Request) {
	// TODO: If Peer Is Not Authorized send 403
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	peerId, err := peer.Decode(request.RemoteAddr)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	maxMsgLen := int64(config.MaxMessageSize + 1)
	bodyReader := io.LimitReader(request.Body, maxMsgLen)
	msgBuff := make([]byte, maxMsgLen)
	read, err := bodyReader.Read(msgBuff)
	if err != io.EOF && err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  fmt.Sprintf("Failed to read message from peer: %v", peerId),
		})
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	notifier.QueueInfo(fmt.Sprintf("Read %d bytes from peer %s", read, peerId))
	messageData := &models.MessageData{}
	err = proto.Unmarshal(msgBuff[:read], messageData)
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  fmt.Sprintf("Failed to unmarshal message from peer: %v", peerId),
		})
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	message := &models.Message{
		Metadata: &models.MessageMetadata{
			From:      peerId.String(),
			To:        _node.ID().String(),
			Timestamp: uint64(time.Now().Unix()),
		},
		Data: messageData,
	}
	if message.GetData().GetStringMessage() != "" {
		handleStringMessage(writer, message)
	} else if message.GetData().GetFileMetadataMessage() != nil {
		handleFileMessage(writer, message, peerId)
	} else {
		notifier.QueueInfo(fmt.Sprintf("Recieved blank message from peer %s", peerId))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
}

// handleStringMessage handles string message
func handleStringMessage(writer http.ResponseWriter, message *models.Message) {
	notifier.QueueMessage(message)
	writer.WriteHeader(http.StatusOK)
}

// handleFileMessage handles file notification message
func handleFileMessage(writer http.ResponseWriter, message *models.Message, peerId peer.ID) {
	fileMetadata := message.GetData().GetFileMetadataMessage()
	fileMetadata.SpecificData = &models.FileMetadata_ClientFileInfo{
		ClientFileInfo: &models.ClientFileInfo{
			BasePath:   config.DefaultDownloadPath,
			FileServer: peerId.String(),
		},
	}
	err := _client.addDownloadingMetafile(fileMetadata)
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  fmt.Sprintf("Client %v failed to add file %v", peerId, fileMetadata.GetFileSha256()),
		})
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	notifier.QueueMessage(message)
	writer.WriteHeader(http.StatusOK)
}
