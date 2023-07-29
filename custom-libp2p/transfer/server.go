package transfer

import (
	customLibP2P "github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/fileHandler"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/utils"
	gostream "github.com/libp2p/go-libp2p-gostream"
	"github.com/libp2p/go-libp2p/core/peer"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type server struct {
	listener         *net.Listener
	server           *http.Server
	servingMetafiles *sync.Map // key: file_SHA256, value: *FileMetadata
}

// initServer initializes the server
func initServer(node *customLibP2P.Node) {
	// TODO: Load serving from DB
	_instance := getInstance()
	_server := _instance.server
	_server.servingMetafiles = &sync.Map{}

	node.SetStreamHandler(holePunchSyncStreamProtocolID, handleHolePunchSyncStream)
	listener, err := gostream.Listen(node, protocolID)
	utils.CheckError(err)
	_server.listener = &listener

	if _server.server == nil {
		go func() {
			http.HandleFunc("/file", handleFileRequest)
			http.HandleFunc("/message", handleMessageRequest)
			_server.server = &http.Server{}
			err = _server.server.Serve(listener)
			log.Printf("Server Error: %v", err)
		}()
	} else {
		go func() {
			err = _server.server.Serve(listener)
			log.Printf("Server Error: %v", err)
		}()
	}
}

// closeServer closes the server
func closeServer() error {
	_instance := getInstance()
	err := _instance.server.server.Close()
	if err != nil {
		return err
	}
	listener := *(_instance.server.listener)
	err = listener.Close()
	if err != nil {
		return err
	}
	_instance.server.servingMetafiles = nil
	return nil
}

// ServeFile serves the file to peer, returns file SHA256Sum and error
func ServeFile(filePath string, peerId peer.ID) (string, error) {
	_instance := getInstance()
	if !_instance.running {
		return "", NotRunning
	}
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	sha256sum, err := fileHandler.GetSHA256Sum(file)
	if err != nil {
		return "", err
	}
	_, exists := _instance.server.servingMetafiles.Load(sha256sum)
	if !exists {
		fStat, err := os.Stat(filePath)
		if err != nil {
			return sha256sum, err
		}
		_instance.server.servingMetafiles.Store(
			sha256sum,
			&FileMetadata{
				_basePath:            filepath.Dir(filePath),
				_fileDescriptor:      file,
				fileName:             filepath.Base(filePath),
				fileSHA256:           sha256sum,
				fileSize:             uint64(fStat.Size()),
				_authorizedAccessors: make(map[peer.ID]bool),
			},
		)
	}
	value, _ := _instance.server.servingMetafiles.Load(sha256sum)
	metadata := value.(*FileMetadata)
	metadata.AllowPeerAccess(peerId)
	err = sendFileNotification(peerId, metadata)
	return sha256sum, err
}

// StopServingFile stops serving the file
func StopServingFile(fileSHA256 string) error {
	_instance := getInstance()
	if !_instance.running {
		return NotRunning
	}
	value, found := _instance.server.servingMetafiles.Load(fileSHA256)
	if !found {
		return FileNotServing
	}
	metadata := value.(*FileMetadata)
	err := metadata.CloseFile()
	if err != nil {
		return err
	}
	_instance.server.servingMetafiles.Delete(fileSHA256)
	return nil
}

// handleFileRequest handles file requests
func handleFileRequest(writer http.ResponseWriter, request *http.Request) {
	_instance := getInstance()
	_clientAddr, err := peer.Decode(request.RemoteAddr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	requestSHA256 := request.URL.Query().Get("sha256")
	value, exists := _instance.server.servingMetafiles.Load(requestSHA256)
	if !exists {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	metadata := value.(*FileMetadata)
	allowed, found := metadata._authorizedAccessors[_clientAddr]
	if !(found && allowed) {
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
func handleGetFile(writer http.ResponseWriter, request *http.Request, metadata *FileMetadata) {
	offset, err := strconv.Atoi(request.URL.Query().Get("offset"))
	if err != nil {
		log.Printf("Error in parsing offset: %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if offset < 0 {
		log.Printf("Invalid Offset")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if uint64(offset) >= metadata.fileSize {
		log.Printf("Invalid Offset")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	file, err := metadata.GetFile()
	if err != nil {
		log.Printf("Error getting local file %s.%s: %v", metadata._basePath, metadata.fileName, err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	fReader := io.NewSectionReader(file, int64(offset), int64(metadata.fileSize-uint64(offset)))
	n, err := io.CopyN(writer, fReader, int64(metadata.fileSize-uint64(offset)))
	if err != nil {
		log.Printf("Error in writing File: %v", err)
		return
	}
	log.Printf("Wrote File: %v bytes (%v)", n+int64(offset), metadata.fileSize)
}

// handleDeleteFile handles delete request for file
func handleDeleteFile(writer http.ResponseWriter, requestSHA256 string,
	metadata *FileMetadata, clientAddr peer.ID) {
	metadata.DenyPeerAccess(clientAddr)
	if len(metadata._authorizedAccessors) == 0 {
		err := StopServingFile(requestSHA256)
		if err != nil {
			log.Printf("Error stopping serving file %s: %v", requestSHA256, err)
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
	maxMsgLen := int64(maxMessageSize + headerSize + 1)
	bodyReader := io.LimitReader(request.Body, maxMsgLen)
	msgBuff := make([]byte, maxMsgLen)
	read, err := bodyReader.Read(msgBuff)
	if err != io.EOF && err != nil {
		log.Printf("Error reading message from peer %s: %v", peerId, err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("Read %d bytes from peer %s", read, peerId)
	message := &Message{}
	err = message.unmarshal(msgBuff[:read])
	if err != nil {
		log.Printf("Error unmarshaling message from peer %s: %v", peerId, err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if message.msgType == StringMessage {
		handleStringMessage(writer, message, peerId)
	} else if message.msgType == FileNotificationMessage {
		handleFileNotificationMessage(writer, message, peerId)
	}
}

// handleStringMessage handles string message
func handleStringMessage(writer http.ResponseWriter, message *Message, peerId peer.ID) {
	// TODO Write to DB
	log.Printf(">> %s: %s\n", peerId, string(message.msgRawData))
	writer.WriteHeader(http.StatusOK)
}

// handleFileNotificationMessage handles file notification message
func handleFileNotificationMessage(writer http.ResponseWriter, message *Message, peerId peer.ID) {
	// TODO Write to DB
	fileMetadata := &FileMetadata{}
	err := fileMetadata.unmarshal(message.msgRawData)
	if err != nil {
		log.Printf("Error unmarshaling file metadata from peer %s: %v", peerId, err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	fileMetadata = fileMetadata.Clone()
	fileMetadata._fileServer = peerId
	_instance := getInstance()
	err = _instance.client.addDownloadingMetafile(fileMetadata)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if autoDownloadEnabled {
		go func() {
			err := ResumeDownload(fileMetadata.fileSHA256, peerId)
			if err != nil {
				log.Printf("Error auto downloading of %s: %v", fileMetadata.fileName, err)
			}
		}()
	}
	log.Printf(">> File Notification %s: %s\n", peerId, fileMetadata.fileName)
	writer.WriteHeader(http.StatusOK)
}
