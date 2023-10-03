package transfer

import (
	"bytes"
	"context"
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/access_manager"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/database"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/file_handler"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
	"github.com/dgraph-io/badger/v4"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/peer"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type client struct {
	node      *models.Node
	client    *http.Client
	transport *http.Transport
	//downloadingMetafiles *sync.Map
	//isDownloading        *sync.Map // key = SHA256+fileServer, value = bool
}

const _downloadingMetafilesKeyBase = "transfer/client/downloadingMetafiles" // key = fileServer/SHA256, value = *FileMetadata
const _isDownloadingKeyBase = "transfer/client/isDownloading"               // key = fileServer/SHA256, value = bool

// initClient Initializes client
func initClient(node *models.Node) error {
	_client.transport = &http.Transport{}
	_client.transport.RegisterProtocol("libp2p", p2phttp.NewTransport(node, p2phttp.ProtocolOption(protocolID)))
	_client.client = &http.Client{Transport: _client.transport}
	// go Download all files isDownloading from database
	return nil
}

// closeClient Closes client
func closeClient() error {
	_clientLock.Lock()
	defer _clientLock.Unlock()
	if !ClientIsRunning() {
		return ErrClientNotRunning
	}
	_client.transport = nil
	_client.client = nil
	return nil
}

// SendMessage Sends message
func SendMessage(peerId peer.ID, message string) error {
	if !ClientIsRunning() {
		return ErrClientNotRunning
	}
	isAllowedNode, err := access_manager.IsAllowedNode(peerId)
	if err != nil {
		return err
	}
	if !isAllowedNode {
		return ErrNotAllowedNode
	}
	messageData := &models.MessageData{
		Data: &models.MessageData_StringMessage{StringMessage: message},
	}
	err = sendMessage(peerId, messageData)
	if err != nil {
		return err
	}
	return nil
}

// PauseDownload Pauses download
func PauseDownload(sha256Sum string, peerId peer.ID) error {
	if !ClientIsRunning() {
		return ErrClientNotRunning
	}
	metafileDbKey := []byte(filepath.Join(_downloadingMetafilesKeyBase, peerId.String(), sha256Sum))
	isDownloadingDbKey := []byte(filepath.Join(_isDownloadingKeyBase, peerId.String(), sha256Sum))
	_, err := database.Load(metafileDbKey)
	if err == badger.ErrKeyNotFound {
		return ErrFileMetadataNotAvailable
	} else if err != nil {
		return err
	}
	isDownloadingBytes, err := database.Load(isDownloadingDbKey)
	if err == badger.ErrKeyNotFound {
		return ErrFileNotDownloading
	} else if err != nil {
		return err
	}
	downloading := isDownloadingBytes[0] == 1
	if !downloading {
		return ErrFileNotDownloading
	}
	return database.Store(isDownloadingDbKey, []byte{0})
}

// ResumeDownload Resumes download
func ResumeDownload(sha256Sum string, peerId peer.ID) error {
	if !ClientIsRunning() {
		return ErrClientNotRunning
	}
	isAllowedNode, err := access_manager.IsAllowedNode(peerId)
	if err != nil {
		return err
	}
	if !isAllowedNode {
		return ErrNotAllowedNode
	}
	metafileDbKey := []byte(filepath.Join(_downloadingMetafilesKeyBase, peerId.String(), sha256Sum))
	isDownloadingDbKey := []byte(filepath.Join(_isDownloadingKeyBase, peerId.String(), sha256Sum))
	_, err = database.Load(metafileDbKey)
	if err == badger.ErrKeyNotFound {
		return ErrFileMetadataNotAvailable
	} else if err != nil {
		return err
	}
	isDownloadingBytes, err := database.Load(isDownloadingDbKey)
	if err == badger.ErrKeyNotFound {
		isDownloadingBytes = []byte{0}
	} else if err != nil {
		return err
	}
	downloading := isDownloadingBytes[0] == 1
	if downloading {
		return ErrAlreadyDownloadingFile
	}
	err = database.Store(isDownloadingDbKey, []byte{1})
	if err != nil {
		return err
	}
	go downloadFile(sha256Sum, peerId)
	return nil
}

// StopDownload Stops download
func StopDownload(sha256Sum string, peerId peer.ID) error {
	if !ClientIsRunning() {
		return ErrClientNotRunning
	}
	metafileDbKey := []byte(filepath.Join(_downloadingMetafilesKeyBase, peerId.String(), sha256Sum))
	isDownloadingDbKey := []byte(filepath.Join(_isDownloadingKeyBase, peerId.String(), sha256Sum))
	fileMetadataBytes, err := database.Load(metafileDbKey)
	if err == badger.ErrKeyNotFound {
		return ErrFileMetadataNotAvailable
	} else if err != nil {
		return err
	}
	fileMetadata := &models.FileMetadata{}
	err = proto.Unmarshal(fileMetadataBytes, fileMetadata)
	if err != nil {
		return err
	}
	err = database.Delete(metafileDbKey)
	if err != nil {
		return err
	}
	err = database.Delete(isDownloadingDbKey)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}
	return notifyServerStopDownloading(fileMetadata)
}

// GetDownloadStatus returns Download Status
func GetDownloadStatus(sha256Sum string, peerId peer.ID) (*models.DownloadStatus, error) {
	if !ClientIsRunning() {
		return nil, ErrClientNotRunning
	}
	metafileDbKey := []byte(filepath.Join(_downloadingMetafilesKeyBase, peerId.String(), sha256Sum))
	isDownloadingDbKey := []byte(filepath.Join(_isDownloadingKeyBase, peerId.String(), sha256Sum))
	metadataBytes, err := database.Load(metafileDbKey)
	if err == badger.ErrKeyNotFound {
		return nil, ErrFileMetadataNotAvailable
	} else if err != nil {
		return nil, err
	}
	metadata := &models.FileMetadata{}
	err = proto.Unmarshal(metadataBytes, metadata)
	if err != nil {
		return nil, err
	}
	isDownloadingBytes, err := database.Load(isDownloadingDbKey)
	if err == badger.ErrKeyNotFound {
		isDownloadingBytes = []byte{0}
	} else if err != nil {
		return nil, err
	}
	isDownloading := isDownloadingBytes[0] == 1
	fileSize, err := file_handler.GetFileSize(getPartDownloading(metadata))
	if err != nil {
		return nil, err
	}
	downloadStatus := &models.DownloadStatus{
		Downloading:     isDownloading,
		FileBasePath:    metadata.GetClientFileInfo().GetBasePath(),
		FileName:        metadata.GetFileName(),
		FileSha256:      metadata.GetFileSha256(),
		FullFileSize:    metadata.GetFileSize(),
		CurrentFileSize: fileSize,
	}
	return downloadStatus, nil
}

// sendFileMessage Sends file message
func sendFileMessage(peerId peer.ID, metadata *models.FileMetadata) error {
	clonedMetadata := proto.Clone(metadata).(*models.FileMetadata)
	clonedMetadata.SpecificData = nil // Mask Specific Data
	msg := &models.MessageData{
		Data: &models.MessageData_FileMetadataMessage{
			FileMetadataMessage: clonedMetadata,
		},
	}
	return sendMessage(peerId, msg)
}

// sendMessage Sends generic message
func sendMessage(peerId peer.ID, message *models.MessageData) error {
	msgBytes, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	msgReader := bytes.NewReader(msgBytes)
	err = _node.Connect(context.TODO(), peer.AddrInfo{ID: peerId})
	if err != nil {
		return err
	}
	url := getMessageUrl(peerId)
	post, err := _client.client.Post(url, "application/octet-stream", msgReader)
	if err != nil {
		return err
	}

	switch post.StatusCode {
	case http.StatusOK:
		notifier.QueueMessage(&models.Message{
			Metadata: &models.MessageMetadata{
				From:      _node.ID().String(),
				To:        peerId.String(),
				Timestamp: uint64(time.Now().Unix()),
			},
			Data: message,
		})
		return nil
	case http.StatusForbidden:
		return ErrForbidden
	default:
		return ErrSendingMessage
	}
}

// addDownloadingMetafile Adds downloading metafile
func (c *client) addDownloadingMetafile(metafile *models.FileMetadata) error {
	metafileDbKey := []byte(filepath.Join(_downloadingMetafilesKeyBase, metafile.GetClientFileInfo().FileServer,
		metafile.GetFileSha256()))
	isDownloadingDbKey := []byte(filepath.Join(_isDownloadingKeyBase, metafile.GetClientFileInfo().FileServer,
		metafile.GetFileSha256()))
	_, err := database.Load(metafileDbKey)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	} else if err == nil {
		return ErrMetafileAlreadyExists
	}
	metadataBytes, err := proto.Marshal(metafile)
	if err != nil {
		return err
	}
	err = database.Store(metafileDbKey, metadataBytes)
	if err != nil {
		return err
	}
	return database.Store(isDownloadingDbKey, []byte{0})
}

// downloadFile Downloads file
func downloadFile(sha256Sum string, peerId peer.ID) {
	metafileDbKey := []byte(filepath.Join(_downloadingMetafilesKeyBase, peerId.String(), sha256Sum))
	isDownloadingDbKey := []byte(filepath.Join(_isDownloadingKeyBase, peerId.String(), sha256Sum))
	metadataBytes, err := database.Load(metafileDbKey)
	if err == badger.ErrKeyNotFound {
		notifier.QueueWarning(&models.Warning{Error: ErrFileMetadataNotAvailable.Error()})
		return
	} else if err != nil {
		notifier.QueueWarning(&models.Warning{Error: err.Error()})
		return
	}
	metadata := &models.FileMetadata{}
	err = proto.Unmarshal(metadataBytes, metadata)
	if err != nil {
		notifier.QueueWarning(&models.Warning{Error: err.Error()})
		return
	}
	file, err := file_handler.GetFile(getPartDownloading(metadata)) // For Downloading
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  "Cant load/create file",
		})
		return
	}
	for {
		// Check if downloading
		isDownloadingBytes, err := database.Load(isDownloadingDbKey)
		if err == badger.ErrKeyNotFound {
			notifier.QueueWarning(&models.Warning{Error: ErrFileMetadataNotAvailable.Error()})
			return
		}
		downloading := isDownloadingBytes[0] == 1
		if !downloading {
			return
		}
		// Check if downloaded
		fStat, err := os.Stat(getPartDownloading(metadata))
		if err != nil {
			notifier.QueueWarning(&models.Warning{
				Error: err.Error(),
				Info:  fmt.Sprintf("File path: %s", getPartDownloading(metadata)),
			})
			err := database.Store(isDownloadingDbKey, []byte{0})
			if err != nil {
				notifier.QueueWarning(&models.Warning{
					Error: err.Error(),
					Info:  "Can not store is downloading",
				})
				return
			}
			continue
		}
		if uint64(fStat.Size()) == metadata.FileSize {
			err := database.Store(isDownloadingDbKey, []byte{0})
			if err != nil {
				notifier.QueueWarning(&models.Warning{
					Error: err.Error(),
					Info:  "Can not store is downloading",
				})
				return
			}
			afterDownloaded(metadata, file)
			break
		}
		// Decode file server
		peerId, err := peer.Decode(metadata.GetClientFileInfo().GetFileServer())
		if err != nil {
			notifier.QueueWarning(&models.Warning{
				Error: err.Error(),
				Info:  "Can not decode file server",
			})
			err = database.Store(isDownloadingDbKey, []byte{0})
			if err != nil {
				notifier.QueueWarning(&models.Warning{
					Error: err.Error(),
					Info:  "Can not store is downloading",
				})
				return
			}
			continue
		}
		// Connect
		err = _node.Connect(context.TODO(), peer.AddrInfo{ID: peerId})
		if err != nil {
			notifier.QueueWarning(&models.Warning{
				Error: err.Error(),
				Info:  fmt.Sprintf("While downloading File Path: %s", getFilePath(metadata)),
			})
			err = database.Store(isDownloadingDbKey, []byte{0})
			if err != nil {
				notifier.QueueWarning(&models.Warning{
					Error: err.Error(),
					Info:  "Can not store is downloading",
				})
				return
			}
			continue
		}
		// Hole Punch
		if !connectedWithoutRelay(_node, peerId) {
			newHolePunchSyncStream(_node, peerId)
		}
		// Download Request
		url := getFileServeUrl(metadata)
		res, err := _client.client.Get(url)
		if err != nil {
			notifier.QueueWarning(&models.Warning{
				Error: err.Error(),
				Info:  fmt.Sprintf("File path: %s", url),
			})
			err = database.Store(isDownloadingDbKey, []byte{0})
			if err != nil {
				notifier.QueueWarning(&models.Warning{
					Error: err.Error(),
					Info:  "Can not store is downloading",
				})
				return
			}
			continue
		}
		// Check server response
		switch res.StatusCode {
		case http.StatusOK:
			n, err := io.CopyN(file, res.Body, int64(metadata.FileSize-uint64(fStat.Size())))
			notifier.QueueInfo(fmt.Sprintf("Read File(%v): %v bytes (%v) [stopped with %v]",
				metadata.FileName, fStat.Size()+n,
				metadata.FileSize, err))
			break
		case http.StatusForbidden:
			notifier.QueueWarning(&models.Warning{
				Error: ErrForbidden.Error(),
				Info: fmt.Sprintf("Failed To Download %v, from %v",
					metadata.GetFileName(), metadata.GetClientFileInfo().GetFileServer()),
			})
			err = database.Store(isDownloadingDbKey, []byte{0})
			if err != nil {
				notifier.QueueWarning(&models.Warning{
					Error: err.Error(),
					Info:  "Can not store is downloading",
				})
				return
			}
			break
		default:
			notifier.QueueWarning(&models.Warning{
				Info: fmt.Sprintf("Error in Downloading File %v Server %v responded with: %v",
					metadata.GetFileName(), metadata.GetClientFileInfo().GetFileServer(), res.Status),
			})
			err = database.Store(isDownloadingDbKey, []byte{0})
			if err != nil {
				notifier.QueueWarning(&models.Warning{
					Error: err.Error(),
					Info:  "Can not store is downloading",
				})
				return
			}
			break
		}
	}
}

// afterDownloaded do after file is downloaded
func afterDownloaded(metadata *models.FileMetadata, file *os.File) {
	err := verifySHA256(metadata, file)
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  fmt.Sprintf("File path: %s", getPartDownloading(metadata)),
		})
		return
	}
	err = os.Rename(getPartDownloading(metadata), getFilePath(metadata))
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  fmt.Sprintf("File path: %s", getPartDownloading(metadata)),
		})
		return
	}
	err = file_handler.CloseFile(getPartDownloading(metadata))
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  fmt.Sprintf("File path: %s", getFilePath(metadata)),
		})
	}
	err = file_handler.CloseFile(getFilePath(metadata))
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  fmt.Sprintf("File path: %s", getFilePath(metadata)),
		})
	}
	err = notifyServerStopDownloading(metadata)
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  fmt.Sprintf("File path: %s", getFilePath(metadata)),
		})
	}
	metafileDbKey := []byte(filepath.Join(_downloadingMetafilesKeyBase, metadata.GetClientFileInfo().FileServer,
		metadata.GetFileSha256()))
	isDownloadingDbKey := []byte(filepath.Join(_isDownloadingKeyBase, metadata.GetClientFileInfo().FileServer,
		metadata.GetFileSha256()))
	err = database.Delete(metafileDbKey)
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  "Can not delete metafile",
		})
		return
	}
	err = database.Delete(isDownloadingDbKey)
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  "Can not delete is downloading",
		})
		return
	}
	notifier.QueueInfo(fmt.Sprintf("File downloaded %v (%v)", metadata.FileName, metadata.FileSha256))
}

func notifyServerStopDownloading(metadata *models.FileMetadata) error {
	deleteUrl := getFileServeUrl(metadata)
	deleteFileRequest, err := http.NewRequest(http.MethodDelete, deleteUrl, nil)
	if err != nil {
		return err
	}
	do, err := _client.client.Do(deleteFileRequest)
	if err != nil {
		return err
	}
	if do.StatusCode != http.StatusOK {
		return err
	}
	return nil
}
