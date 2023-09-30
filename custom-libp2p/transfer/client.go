package transfer

import (
	"bytes"
	"context"
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/access_manager"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/file_handler"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/peer"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type client struct {
	node                 *models.Node
	client               *http.Client
	transport            *http.Transport
	downloadingMetafiles *sync.Map // key = SHA256+fileServer, value = *FileMetadata
	isDownloading        *sync.Map // key = SHA256+fileServer, value = bool
}

// initClient Initializes client
func initClient(node *models.Node) error {
	// TODO: Load downloading Metafiles From DB
	_client.downloadingMetafiles = &sync.Map{}
	_client.isDownloading = &sync.Map{}

	_client.transport = &http.Transport{}
	_client.transport.RegisterProtocol("libp2p", p2phttp.NewTransport(node, p2phttp.ProtocolOption(protocolID)))
	_client.client = &http.Client{Transport: _client.transport}
	return nil
}

// closeClient Closes client
func closeClient() error {
	_clientLock.Lock()
	defer _clientLock.Unlock()
	if !ClientIsRunning() {
		return ErrClientNotRunning
	}
	_client.isDownloading.Range(func(key, value any) bool {
		_client.isDownloading.Store(key, false)
		return true
	})
	_client.downloadingMetafiles = nil
	return nil
}

// SendMessage Sends message
func SendMessage(peerId peer.ID, message string) error {
	if !ClientIsRunning() {
		return ErrClientNotRunning
	}
	if !access_manager.IsAllowedNode(peerId) {
		return ErrNotAllowedNode
	}
	messageData := &models.MessageData{
		Data: &models.MessageData_StringMessage{StringMessage: message},
	}
	err := sendMessage(peerId, messageData)
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
	key := sha256Sum + peerId.String()
	_, found := _client.downloadingMetafiles.Load(key)
	if !found {
		return ErrFileMetadataNotAvailable
	}
	value, found := _client.isDownloading.Load(key)
	downloading := value.(bool)
	if !(found && downloading) {
		return ErrFileNotDownloading
	}
	_client.isDownloading.Store(key, false)
	return nil
}

// ResumeDownload Resumes download
func ResumeDownload(sha256Sum string, peerId peer.ID) error {
	if !ClientIsRunning() {
		return ErrClientNotRunning
	}
	if !access_manager.IsAllowedNode(peerId) {
		return ErrNotAllowedNode
	}
	key := sha256Sum + peerId.String()
	value, found := _client.isDownloading.Load(key)
	if !found {
		return ErrFileMetadataNotAvailable
	}
	isDownloading := value.(bool)
	if isDownloading {
		return ErrAlreadyDownloadingFile
	}
	_client.isDownloading.Store(key, true)
	go downloadFile(key)
	return nil
}

// StopDownload Stops download
func StopDownload(sha256Sum string, peerId peer.ID) error {
	if !ClientIsRunning() {
		return ErrClientNotRunning
	}
	key := sha256Sum + peerId.String()
	value, found := _client.downloadingMetafiles.Load(key)
	if !found {
		return ErrFileMetadataNotAvailable
	}
	fileMetadata := value.(*models.FileMetadata)
	_client.downloadingMetafiles.Delete(key)
	_client.isDownloading.Delete(key)
	err := notifyServerStopDownloading(fileMetadata)
	if err != nil {
		return err
	}
	return nil
}

// GetDownloadStatus returns Download Status
func GetDownloadStatus(sha256Sum string, peerId peer.ID) (*models.DownloadStatus, error) {
	if !ClientIsRunning() {
		return nil, ErrClientNotRunning
	}
	key := sha256Sum + peerId.String()
	value, found := _client.downloadingMetafiles.Load(key)
	if !found {
		return nil, ErrFileMetadataNotAvailable
	}
	metafile := value.(*models.FileMetadata)
	value, found = _client.isDownloading.Load(key)
	if !found {
		return nil, ErrFileNotDownloading
	}
	isDownloading := value.(bool)
	fileSize, err := file_handler.GetFileSize(getFilePath(metafile))
	if err != nil {
		return nil, err
	}
	downloadStatus := &models.DownloadStatus{
		Downloading:     isDownloading,
		FileBasePath:    metafile.GetClientFileInfo().GetBasePath(),
		FileName:        metafile.GetFileName(),
		FileSha256:      metafile.GetFileSha256(),
		FullFileSize:    metafile.GetFileSize(),
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
	key := metafile.FileSha256 + metafile.GetClientFileInfo().GetFileServer()
	_, found := c.downloadingMetafiles.Load(key)
	if found {
		return ErrMetafileAlreadyExists
	}
	c.downloadingMetafiles.Store(key, metafile)
	c.isDownloading.Store(key, false)
	return nil
}

// downloadFile Downloads file
func downloadFile(key string) {
	value, found := _client.downloadingMetafiles.Load(key)
	if !found {
		notifier.QueueWarning(&models.Warning{Error: ErrFileMetadataNotAvailable.Error()})
		return
	}
	// TODO: Better download logic, if file exists, change name
	metadata := value.(*models.FileMetadata)
	file, err := file_handler.GetFile(getFilePath(metadata))
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  "Cant load/create file",
		})
		return
	}
	for {
		value, found = _client.isDownloading.Load(key)
		if !found {
			notifier.QueueWarning(&models.Warning{Error: ErrFileMetadataNotAvailable.Error()})
			return
		}
		downloading := value.(bool)
		if !downloading {
			return
		}
		peerId, err := peer.Decode(metadata.GetClientFileInfo().GetFileServer())
		if err != nil {
			notifier.QueueWarning(&models.Warning{
				Error: err.Error(),
				Info:  "Can not decode file server",
			})
		}
		err = _node.Connect(context.TODO(), peer.AddrInfo{ID: peerId})
		if err != nil {
			notifier.QueueWarning(&models.Warning{
				Error: err.Error(),
				Info:  fmt.Sprintf("While downloading to File Path: %s", getFilePath(metadata)),
			})
			_client.isDownloading.Store(key, false)
			continue
		}
		if !connectedWithoutRelay(_node, peerId) {
			newHolePunchSyncStream(_node, peerId)
		}
		fStat, err := os.Stat(getFilePath(metadata))
		if err != nil {
			notifier.QueueWarning(&models.Warning{
				Error: err.Error(),
				Info:  fmt.Sprintf("File path: %s", getFilePath(metadata)),
			})
			_client.isDownloading.Store(key, false)
			continue
		}
		if uint64(fStat.Size()) == metadata.FileSize {
			// TODO Handle after file download (verify sha256, ...)
			notifier.QueueInfo(fmt.Sprintf("File downloaded %v (%v)", metadata.FileName, metadata.FileSha256))
			_client.isDownloading.Store(key, false)
			afterDownloaded(metadata, file)
			continue
		}
		url := getFileServeUrl(metadata)
		res, err := _client.client.Get(url)
		if err != nil {
			notifier.QueueWarning(&models.Warning{
				Error: err.Error(),
				Info:  fmt.Sprintf("File path: %s", url),
			})
			_client.isDownloading.Store(key, false)
			continue
		}
		if res.StatusCode == http.StatusOK {
			n, _ := io.CopyN(file, res.Body, int64(metadata.FileSize-uint64(fStat.Size())))
			notifier.QueueInfo(fmt.Sprintf("Read File(%v): %v bytes (%v)", metadata.FileSize, fStat.Size()+n,
				metadata.FileSize))
			if uint64(n+fStat.Size()) == metadata.FileSize {
				_client.isDownloading.Store(key, false)
				afterDownloaded(metadata, file)
			}
		} else {
			notifier.QueueWarning(&models.Warning{
				Info: fmt.Sprintf("Error in Downloading File Server responded with: %v", res.Status),
			})
			_client.isDownloading.Store(key, false)
			continue
		}
	}
}

// afterDownloaded do after file is downloaded
func afterDownloaded(metadata *models.FileMetadata, file *os.File) {
	// TODO Check SHA256
	notifier.QueueInfo(fmt.Sprintf("File downloaded %v", metadata.FileName))
	err := file.Close()
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
