package transfer

import (
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/access_manager"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/config"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/database"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/file_handler"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
	"github.com/dgraph-io/badger/v4"
	gostream "github.com/libp2p/go-libp2p-gostream"
	"github.com/libp2p/go-libp2p/core/peer"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type server struct {
	listener *net.Listener
	server   *http.Server
}

const _servingMetafilesKeyBase = "transfer/server/servingMetafiles" // key: file_SHA256, value: *models.FileMetadata

// initServer initializes the server
func initServer(node *models.Node) error {
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
	return nil
}

// ServeFile serves the file to peer, returns file SHA256Sum and error
func ServeFile(filePath string, peerId peer.ID) (*models.FileMetadata, error) {
	if !ServerIsRunning() {
		return nil, ErrServerNotRunning
	}
	isAllowed, err := access_manager.IsAllowedNode(peerId)
	if err != nil {
		return nil, err
	}
	if !isAllowed {
		return nil, ErrNotAllowedNode
	}
	file, err := file_handler.GetFile(filePath)
	if err != nil {
		return nil, err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	sha256sum, err := file_handler.GetSHA256Sum(file)
	if err != nil {
		return nil, err
	}
	dbKey := []byte(filepath.Join(_servingMetafilesKeyBase, sha256sum))
	_, err = database.Load(dbKey)
	if err != nil && err != badger.ErrKeyNotFound {
		return nil, err
	}
	exists := err != badger.ErrKeyNotFound
	if !exists {
		fSize, err := file_handler.GetFileSize(filePath)
		if err != nil {
			return nil, err
		}
		createdMetadata := &models.FileMetadata{
			FileName:   filepath.Base(filePath),
			FileSha256: sha256sum,
			FileSize:   fSize,
			SpecificData: &models.FileMetadata_ServerFileInfo{ServerFileInfo: &models.ServerFileInfo{
				BasePath:            filepath.Dir(filePath),
				AuthorizedAccessors: []string{},
			}},
		}
		value, err := proto.Marshal(createdMetadata)
		if err != nil {
			return nil, err
		}
		err = database.Store(dbKey, value)
		if err != nil {
			return nil, err
		}
	}
	metadataBytes, err := database.Load(dbKey)
	if err != nil {
		return nil, err
	}
	metadata := &models.FileMetadata{}
	err = proto.Unmarshal(metadataBytes, metadata)
	if err != nil {
		return nil, err
	}
	metadata.GetServerFileInfo().AuthorizedAccessors = append(metadata.GetServerFileInfo().AuthorizedAccessors,
		peerId.String())
	metadataBytes, err = proto.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	err = database.Store(dbKey, metadataBytes)
	if err != nil {
		return nil, err
	}
	return metadata, sendFileMessage(peerId, metadata)
}

// StopServingFile stops serving the file
func StopServingFile(fileSHA256 string) error {
	if !ServerIsRunning() {
		return ErrServerNotRunning
	}
	dbKey := []byte(filepath.Join(_servingMetafilesKeyBase, fileSHA256))
	metadataBytes, err := database.Load(dbKey)
	if err == badger.ErrKeyNotFound {
		return ErrFileNotServing
	} else if err != nil {
		return err
	}
	metadata := &models.FileMetadata{}
	err = proto.Unmarshal(metadataBytes, metadata)
	if err != nil {
		return err
	}
	err = file_handler.CloseFile(getFilePath(metadata))
	if err != nil {
		return err
	}
	return database.Delete(dbKey)
}

// handleFileDownloadRequest handles file download requests
func handleFileDownloadRequest(writer http.ResponseWriter, request *http.Request) {
	_clientAddr, err := peer.Decode(request.RemoteAddr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	isAllowed, err := access_manager.IsAllowedNode(_clientAddr)
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  "Error in checking if node is allowed",
		})
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isAllowed {
		writer.WriteHeader(http.StatusForbidden)
		notifier.QueueInfo(fmt.Sprintf("Blocked node: %v tried to download file", _clientAddr.String()))
		return
	}
	requestSHA256 := request.URL.Query().Get("sha256")
	dbKey := []byte(filepath.Join(_servingMetafilesKeyBase, requestSHA256))
	metadataBytes, err := database.Load(dbKey)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			writer.WriteHeader(http.StatusNotFound)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
		}
		notifier.QueueWarning(&models.Warning{
			Info:  fmt.Sprintf("Client: %v tried to download file: %v", _clientAddr.String(), requestSHA256),
			Error: err.Error(),
		})
		return
	}
	metadata := &models.FileMetadata{}
	err = proto.Unmarshal(metadataBytes, metadata)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		notifier.QueueWarning(&models.Warning{
			Info:  fmt.Sprintf("Client: %v tried to download file: %v", _clientAddr.String(), requestSHA256),
			Error: err.Error(),
		})
		return
	}
	allowed := slices.Contains(metadata.GetServerFileInfo().AuthorizedAccessors, _clientAddr.String())
	if !allowed {
		writer.WriteHeader(http.StatusNotFound)
		notifier.QueueWarning(&models.Warning{
			Info:  fmt.Sprintf("Client: %v tried to download file: %v", _clientAddr.String(), requestSHA256),
			Error: err.Error(),
		})
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
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	peerId, err := peer.Decode(request.RemoteAddr)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	isAllowed, err := access_manager.IsAllowedNode(peerId)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  "Error in checking if node is allowed",
		})
		return
	}
	isBlocked, err := access_manager.IsBlockedNode(peerId)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  "Error in checking if node is blocked",
		})
		return
	}
	if !isAllowed && !isBlocked {
		// First timer (Honor this request, but block)
		err := access_manager.BlockNode(peerId)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			notifier.QueueWarning(&models.Warning{
				Error: err.Error(),
				Info:  "Error in blocking node",
			})
			return
		}
	} else if isBlocked {
		writer.WriteHeader(http.StatusForbidden)
		notifier.QueueInfo(fmt.Sprintf("Blocked node: %v tried to send message", peerId.String()))
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
			BasePath:   file_handler.GetDownloadPath(),
			FileServer: peerId.String(),
		},
	}
	fileMetadata.FileName = filepath.Base(getNextAvailableFilePath(fileMetadata)) // Find Next best name
	_, err := file_handler.GetFile(getFilePath(fileMetadata))                     // Reserve name now, use for renaming when downloaded
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  "Cant load/create file",
		})
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = _client.addDownloadingMetafile(fileMetadata)
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
