package transfer

import (
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/file_handler"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"golang.org/x/exp/slices"
	"os"
	"path/filepath"
)

// connectedWithoutRelay check if peer is connected to node without relay
func connectedWithoutRelay(node models.Node, _peer peer.ID) bool {
	connectedDirect := false
	for _, c := range node.Network().ConnsToPeer(_peer) {
		connectedDirect = connectedDirect || !isRelay(c.RemoteMultiaddr())
	}
	return connectedDirect
}

// isRelay check if address is relayed
func isRelay(address multiaddr.Multiaddr) bool {
	_isRelay := false
	multiaddr.ForEach(address, func(c multiaddr.Component) bool {
		switch c.Protocol().Code {
		case multiaddr.P_CIRCUIT:
			_isRelay = true
			return false
		default:
			return true
		}
	})
	return _isRelay
}

// getPartDownloading filePath
func getPartDownloading(metadata *models.FileMetadata) string {
	if metadata.GetClientFileInfo() != nil {
		return filepath.Join(metadata.GetClientFileInfo().GetBasePath(), fmt.Sprintf(".%v_%v",
			metadata.GetClientFileInfo().GetFileServer(),
			metadata.GetFileSha256()))
	}
	return ""
}

// getFilePath returns filePath
func getFilePath(metadata *models.FileMetadata) string {
	if metadata.GetClientFileInfo() != nil {
		return filepath.Join(metadata.GetClientFileInfo().GetBasePath(), metadata.GetFileName())
	} else if metadata.GetServerFileInfo() != nil {
		return filepath.Join(metadata.GetServerFileInfo().GetBasePath(), metadata.GetFileName())
	}
	return ""
}

// getNthFilePath returns filePath
func getNthFilePath(metadata *models.FileMetadata, n uint) string {
	fileName := metadata.FileName[:len(metadata.FileName)-len(filepath.Ext(metadata.FileName))]
	fileName = fmt.Sprintf("%v(%v)%v", fileName, n, filepath.Ext(metadata.FileName))
	if metadata.GetClientFileInfo() != nil {
		return filepath.Join(metadata.GetClientFileInfo().GetBasePath(), fileName)
	}
	return ""
}

// getFileServeUrl path for downloading file
func getFileServeUrl(metadata *models.FileMetadata) string {
	size, err := file_handler.GetFileSize(getPartDownloading(metadata))
	if err != nil {
		notifier.QueueWarning(&models.Warning{
			Error: err.Error(),
			Info:  "Inferring Offset get filesize",
		})
		size = 0
	}
	return fmt.Sprintf("libp2p://%v/file?sha256=%v&offset=%v", metadata.GetClientFileInfo().FileServer,
		metadata.FileSha256, size)
}

// getFileDeleteUrl path for deleting file
func getFileDeleteUrl(metadata *models.FileMetadata) string {
	return fmt.Sprintf("libp2p://%v/file?sha256=%v", metadata.GetClientFileInfo().FileServer,
		metadata.FileSha256)
}

// getMessageUrl url for sending message
func getMessageUrl(peerId peer.ID) string {
	return fmt.Sprintf("libp2p://%s/message", peerId.String())
}

// removePeer denies peer from server metafile
func removePeer(metadata *models.FileMetadata, _peer peer.ID) {
	peerIndex := slices.Index(metadata.GetServerFileInfo().AuthorizedAccessors, _peer.String())
	if peerIndex != -1 {
		metadata.GetServerFileInfo().AuthorizedAccessors =
			append(metadata.GetServerFileInfo().AuthorizedAccessors[:peerIndex],
				metadata.GetServerFileInfo().AuthorizedAccessors[peerIndex+1:]...)
	}
}

// verifySHA256 verifies sha256
func verifySHA256(metadata *models.FileMetadata, file *os.File) error {
	sha256Sum, err := file_handler.GetSHA256Sum(file)
	if err != nil {
		return err
	}
	if sha256Sum != metadata.GetFileSha256() {
		return ErrSha256DoesNotMatch
	}
	return nil
}

// getNextAvailableFilePath returns next available file path for saving file
func getNextAvailableFilePath(metadata *models.FileMetadata) string {
	newFilePath := getFilePath(metadata)
	if _, err := os.Stat(newFilePath); err == nil {
		var offset uint = 0
		for err == nil {
			offset++
			_, err = os.Stat(getNthFilePath(metadata, offset))
		}
		newFilePath = getNthFilePath(metadata, offset)
	}
	return newFilePath
}
