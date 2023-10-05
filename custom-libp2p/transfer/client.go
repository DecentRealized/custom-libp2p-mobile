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
}

const _downloadingMetafilesKeyBase = "transfer/client/downloadingMetafiles" // key = fileServer/SHA256, value = *FileMetadata

// initClient Initializes client
func initClient(node *models.Node) error {
	_client.transport = &http.Transport{}
	_client.transport.RegisterProtocol("libp2p", p2phttp.NewTransport(node, p2phttp.ProtocolOption(protocolID)))
	_client.client = &http.Client{Transport: _client.transport}
	err := database.Range([]byte(_downloadingMetafilesKeyBase), func(key, value []byte) bool {
		metadata := &models.FileMetadata{}
		err := proto.Unmarshal(value, metadata)
		if err != nil {
			notifier.QueueWarning(&models.Warning{
				Error: err.Error(),
				Info:  fmt.Sprintf("Error unmarshaling %s", key),
			})
			// Cleanup
			_ = database.Delete(key)
		}
		if metadata.GetClientFileInfo().IsDownloading {
			peerId, err := peer.Decode(metadata.GetClientFileInfo().FileServer)
			if err != nil {
				notifier.QueueWarning(&models.Warning{
					Error: err.Error(),
					Info:  fmt.Sprintf("Error decoding %s", key),
				})
				// Cleanup
				_ = database.Delete(key)
			} else {
				go downloadFile(metadata.FileSha256, peerId)
			}
		}
		return true
	})
	return err
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
	metadataBytes, err := database.Load(metafileDbKey)
	if err == badger.ErrKeyNotFound {
		return ErrFileMetadataNotAvailable
	} else if err != nil {
		return err
	}
	metadata := &models.FileMetadata{}
	err = proto.Unmarshal(metadataBytes, metadata)
	if err != nil {
		return err
	}
	downloading := metadata.GetClientFileInfo().IsDownloading
	if !downloading {
		return ErrFileNotDownloading
	}
	metadata.GetClientFileInfo().IsDownloading = false
	metadataBytes, err = proto.Marshal(metadata)
	if err != nil {
		return err
	}
	return database.Store(metafileDbKey, metadataBytes)
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
	metadataBytes, err := database.Load(metafileDbKey)
	if err == badger.ErrKeyNotFound {
		return ErrFileMetadataNotAvailable
	} else if err != nil {
		return err
	}
	metadata := &models.FileMetadata{}
	err = proto.Unmarshal(metadataBytes, metadata)
	if err != nil {
		return err
	}
	downloading := metadata.GetClientFileInfo().IsDownloading
	if downloading {
		return ErrAlreadyDownloadingFile
	}
	metadata.GetClientFileInfo().IsDownloading = true
	metadataBytes, err = proto.Marshal(metadata)
	if err != nil {
		return err
	}
	err = database.Store(metafileDbKey, metadataBytes)
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
	return notifyServerStopDownloading(fileMetadata)
}

// GetDownloadingFiles returns downloading files
func GetDownloadingFiles() ([]*models.FileMetadata, error) {
	if !ClientIsRunning() {
		return nil, ErrClientNotRunning
	}
	var unmarshalErr error
	var downloadingFiles []*models.FileMetadata
	rangeErr := database.Range([]byte(_downloadingMetafilesKeyBase), func(key []byte, value []byte) bool {
		downloadingFile := &models.FileMetadata{}
		err := proto.Unmarshal(value, downloadingFile)
		if err != nil {
			unmarshalErr = err
			return false
		}
		downloadingFiles = append(downloadingFiles, downloadingFile)
		return true
	})
	if rangeErr != nil {
		return nil, rangeErr
	}
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}
	return downloadingFiles, nil
}

// GetDownloadStatus returns Download Status
func GetDownloadStatus(sha256Sum string, peerId peer.ID) (*models.DownloadStatus, error) {
	if !ClientIsRunning() {
		return nil, ErrClientNotRunning
	}
	metafileDbKey := []byte(filepath.Join(_downloadingMetafilesKeyBase, peerId.String(), sha256Sum))
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
	isDownloading := metadata.GetClientFileInfo().IsDownloading
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
	return database.Store(metafileDbKey, metadataBytes)
}

// downloadFile Downloads file
func downloadFile(sha256Sum string, peerId peer.ID) {
	metafileDbKey := []byte(filepath.Join(_downloadingMetafilesKeyBase, peerId.String(), sha256Sum))
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
	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  "Cant load/create file",
		})
		return
	}
	for {
		// Re-Fetch Metadata
		metadataBytes, err = database.Load(metafileDbKey)
		if err == badger.ErrKeyNotFound {
			notifier.QueueWarning(&models.Warning{Error: ErrFileMetadataNotAvailable.Error()})
			return
		} else if err != nil {
			notifier.QueueWarning(&models.Warning{Error: err.Error()})
			return
		}
		err = proto.Unmarshal(metadataBytes, metadata)
		if err != nil {
			notifier.QueueWarning(&models.Warning{Error: err.Error()})
			return
		}
		downloading := metadata.GetClientFileInfo().IsDownloading
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
			err := PauseDownload(sha256Sum, peerId)
			if err != nil {
				notifier.QueueWarning(&models.Warning{
					Error: err.Error(),
					Info:  "Can not stop download",
				})
				return
			}
		}
		if uint64(fStat.Size()) == metadata.FileSize {
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
			err = PauseDownload(sha256Sum, peerId)
			if err != nil {
				notifier.QueueWarning(&models.Warning{
					Error: err.Error(),
					Info:  "Can not stop download",
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
			err = PauseDownload(sha256Sum, peerId)
			if err != nil {
				notifier.QueueWarning(&models.Warning{
					Error: err.Error(),
					Info:  "Can not stop download",
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
			err = PauseDownload(sha256Sum, peerId)
			if err != nil {
				notifier.QueueWarning(&models.Warning{
					Error: err.Error(),
					Info:  "Can not stop download",
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
			err = PauseDownload(sha256Sum, peerId)
			if err != nil {
				notifier.QueueWarning(&models.Warning{
					Error: err.Error(),
					Info:  "Can not stop download",
				})
				return
			}
			break
		default:
			notifier.QueueWarning(&models.Warning{
				Info: fmt.Sprintf("Error in Downloading File %v Server %v responded with: %v",
					metadata.GetFileName(), metadata.GetClientFileInfo().GetFileServer(), res.Status),
			})
			err = PauseDownload(sha256Sum, peerId)
			if err != nil {
				notifier.QueueWarning(&models.Warning{
					Error: err.Error(),
					Info:  "Can not stop download",
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
	err = database.Delete(metafileDbKey)
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  "Can not delete metafile",
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
