package transfer

import (
	"bytes"
	"context"
	"fmt"
	customLibP2P "github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/peer"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type client struct {
	node                 *customLibP2P.Node
	client               *http.Client
	transport            *http.Transport
	downloadingMetafiles *sync.Map // key = SHA256+fileServer, value = *FileMetadata
	isDownloading        *sync.Map // key = SHA256+fileServer, value = bool
}

// initClient Initializes client
func initClient(node *customLibP2P.Node) {
	// TODO: Load downloading Metafiles From DB
	_instance := getInstance()
	_instance.client.downloadingMetafiles = &sync.Map{}
	_instance.client.isDownloading = &sync.Map{}

	_instance.client.transport = &http.Transport{}
	_instance.client.transport.RegisterProtocol("libp2p",
		p2phttp.NewTransport(node, p2phttp.ProtocolOption(protocolID)))
	_instance.client.client = &http.Client{Transport: _instance.client.transport}
}

// closeClient Closes client
func closeClient() error {
	_instance := getInstance()
	_instance.client.isDownloading.Range(func(key, value any) bool {
		_instance.client.isDownloading.Store(key, false)
		return true
	})
	_instance.client.downloadingMetafiles = nil
	return nil
}

// SendMessage Sends message
func SendMessage(peerId peer.ID, message string) error {
	_instance := getInstance()
	if !_instance.running {
		return NotRunning
	}
	stringMessage, err := NewStringMessage(message)
	if err != nil {
		return err
	}
	err = sendMessage(peerId, stringMessage)
	if err != nil {
		return err
	}
	return nil
}

// PauseDownload Pauses download
func PauseDownload(sha256Sum string, peerId peer.ID) error {
	_instance := getInstance()
	if !_instance.running {
		return NotRunning
	}
	key := sha256Sum + peerId.String()
	_, found := _instance.client.downloadingMetafiles.Load(key)
	if !found {
		return FileMetadataNotAvailable
	}
	value, found := _instance.client.isDownloading.Load(key)
	downloading := value.(bool)
	if !(found && downloading) {
		return FileNotDownloading
	}
	_instance.client.isDownloading.Store(key, false)
	return nil
}

// ResumeDownload Resumes download
func ResumeDownload(sha256Sum string, peerId peer.ID) error {
	_instance := getInstance()
	if !_instance.running {
		return NotRunning
	}
	key := sha256Sum + peerId.String()
	value, found := _instance.client.isDownloading.Load(key)
	if !found {
		return FileMetadataNotAvailable
	}
	isDownloading := value.(bool)
	if isDownloading {
		return AlreadyDownloadingFile
	}
	_instance.client.isDownloading.Store(key, true)
	go downloadFile(key)
	return nil
}

// StopDownload Stops download
func StopDownload(sha256Sum string, peerId peer.ID) error {
	_instance := getInstance()
	if !_instance.running {
		return NotRunning
	}
	key := sha256Sum + peerId.String()
	value, found := _instance.client.downloadingMetafiles.Load(key)
	if !found {
		return FileMetadataNotAvailable
	}
	fileMetadata := value.(*FileMetadata)
	_instance.client.downloadingMetafiles.Delete(key)
	_instance.client.isDownloading.Delete(key)
	err := notifyServerStopDownloading(fileMetadata)
	if err != nil {
		return err
	}
	return nil
}

// GetDownloadStatus returns Download Status
func GetDownloadStatus(sha256Sum string, peerId peer.ID) (DownloadStatus, error) {
	_instance := getInstance()
	if !_instance.running {
		return DownloadStatus{}, NotRunning
	}
	key := sha256Sum + peerId.String()
	value, found := _instance.client.downloadingMetafiles.Load(key)
	if !found {
		return DownloadStatus{}, FileMetadataNotAvailable
	}
	metafile := value.(*FileMetadata)
	value, found = _instance.client.isDownloading.Load(key)
	if !found {
		return DownloadStatus{}, FileNotDownloading
	}
	isDownloading := value.(bool)
	return downloadStatusFromDownloadingAndMetadata(isDownloading, metafile)
}

// sendFileNotification Sends file notification
func sendFileNotification(peerId peer.ID, metadata *FileMetadata) error {
	msg, err := NewFileNotificationMessage(metadata)
	if err != nil {
		return err
	}
	return sendMessage(peerId, msg)
}

func sendMessage(peerId peer.ID, message *Message) error {
	msgBytes, err := message.marshal()
	if err != nil {
		return err
	}
	msgReader := bytes.NewReader(msgBytes)
	_instance := getInstance()
	err = _instance.node.Connect(context.TODO(), peer.AddrInfo{ID: peerId})
	if err != nil {
		return err
	}
	url := fmt.Sprintf("libp2p://%s/message", peerId.String())
	post, err := _instance.client.client.Post(url, "application/octet-stream", msgReader)
	if err != nil {
		return err
	}
	if post.StatusCode != http.StatusOK {
		return ErrorSendingMessage
	}
	return nil
}

// addDownloadingMetafile Adds downloading metafile
func (c *client) addDownloadingMetafile(metafile *FileMetadata) error {
	key := metafile.fileSHA256 + metafile._fileServer.String()
	_, found := c.downloadingMetafiles.Load(key)
	if found {
		return MetafileAlreadyExists
	}
	c.downloadingMetafiles.Store(key, metafile)
	c.isDownloading.Store(key, false)
	return nil
}

// downloadFile Downloads file
func downloadFile(key string) {
	_instance := getInstance()
	value, found := _instance.client.downloadingMetafiles.Load(key)
	if !found {
		log.Printf("Error Downloading File: %v", FileMetadataNotAvailable)
		return
	}
	metadata := value.(*FileMetadata)
	file, err := metadata.GetFile()
	if err != nil {
		log.Printf("Error getting local file %s.%s: %v", metadata._basePath, metadata.fileName, err)
		_instance.client.isDownloading.Store(key, false)
		return
	}
	for {
		value, found = _instance.client.isDownloading.Load(key)
		if !found {
			log.Printf("Error Downloading File: %v", FileMetadataNotAvailable)
			return
		}
		downloading := value.(bool)
		if !downloading {
			return
		}
		err := _instance.node.Connect(context.TODO(), peer.AddrInfo{ID: metadata._fileServer})
		if err != nil {
			log.Printf("Error connecting to file server %s.%s: %v", metadata._basePath, metadata.fileName, err)
			_instance.client.isDownloading.Store(key, false)
			continue
		}
		if !connectedWithoutRelay(_instance.node, metadata._fileServer) {
			newHolePunchSyncStream(_instance.node, metadata._fileServer)
		}
		fStat, err := os.Stat(filepath.Join(metadata._basePath, metadata.fileName))
		if err != nil {
			log.Printf("Error getting metadata %s.%s: %v", metadata._basePath, metadata.fileName, err)
			_instance.client.isDownloading.Store(key, false)
			continue
		}
		if uint64(fStat.Size()) == metadata.fileSize {
			// TODO Handle after file download (verify sha256, ...)
			log.Printf("File downloaded %v", metadata.fileName)
			_instance.client.isDownloading.Store(key, false)
			afterDownloaded(metadata)
			continue
		}
		url := fmt.Sprintf("libp2p://%v/file?sha256=%v&offset=%v",
			metadata._fileServer.String(), metadata.fileSHA256, fStat.Size())
		res, err := _instance.client.client.Get(url)
		if err != nil {
			log.Printf("Error in Downloading GET: %v", err)
			_instance.client.isDownloading.Store(key, false)
			continue
		}
		if res.StatusCode == http.StatusOK {
			n, _ := io.CopyN(file, res.Body, int64(metadata.fileSize-uint64(fStat.Size())))
			log.Printf("Read File: %v bytes (%v)", fStat.Size()+n, metadata.fileSize)
			if uint64(n+fStat.Size()) == metadata.fileSize {
				_instance.client.isDownloading.Store(key, false)
				afterDownloaded(metadata)
			}
		} else {
			log.Printf("Error in Downloading File Server responded with: %v", res.Status)
			_instance.client.isDownloading.Store(key, false)
			continue
		}
	}
}

// afterDownloaded do after file is downloaded
func afterDownloaded(metadata *FileMetadata) {
	// TODO Check SHA256
	log.Printf("File downloaded %v", metadata.fileName)
	err := metadata.CloseFile()
	if err != nil {
		log.Printf("Error closing file %s.%s: %v", metadata._basePath, metadata.fileName, err)
	}
	err = notifyServerStopDownloading(metadata)
	if err != nil {
		log.Printf("Error notifying server to stop downloading %s.%s: %v", metadata._basePath, metadata.fileName, err)
	}
}

func notifyServerStopDownloading(metadata *FileMetadata) error {
	_instance := getInstance()
	deleteUrl := fmt.Sprintf("libp2p://%v/file?sha256=%v", metadata._fileServer, metadata.fileSHA256)
	deleteFileRequest, err := http.NewRequest(http.MethodDelete, deleteUrl, nil)
	if err != nil {
		return err
	}
	do, err := _instance.client.client.Do(deleteFileRequest)
	if err != nil {
		return err
	}
	if do.StatusCode != http.StatusOK {
		return err
	}
	return nil
}
