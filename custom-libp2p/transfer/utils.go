package transfer

import (
	"bytes"
	"context"
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/file_handler"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
	"github.com/libp2p/go-libp2p/core"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"golang.org/x/exp/slices"
	"path/filepath"
)

// connectedWithoutRelay check if peer is connected to node without relay
func connectedWithoutRelay(node *models.Node, _peer peer.ID) bool {
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

// newHolePunchSyncStream creates new stream for syncing hole punching
func newHolePunchSyncStream(node *models.Node, peerId peer.ID) {
	stream, err := node.NewStream(context.TODO(), peerId, holePunchSyncStreamProtocolID)
	if err != nil {
		notifier.QueueWarning(&models.Warning{Error: err.Error(), Info: "Failed to create hole punch sync stream"})
	}
	buff := bytes.Repeat([]byte{1}, holePunchPacketSize)
	for i := 0; i < holePunchRetries; i++ {
		_, err := stream.Write(buff)
		if err != nil {
			notifier.QueueWarning(&models.Warning{
				Error: err.Error(),
				Info:  "Failed to write to hole punch sync stream",
			})
		}
		_, err = stream.Read(buff)
		if err != nil {
			notifier.QueueWarning(&models.Warning{
				Error: err.Error(),
				Info:  "Failed to read from hole punch sync stream",
			})
		}
	}
}

// handleHolePunchSyncStream handles new hole punch stream
func handleHolePunchSyncStream(stream core.Stream) {
	buff := bytes.Repeat([]byte{0}, holePunchPacketSize)
	for i := 0; i < holePunchRetries; i++ {
		_, err := stream.Read(buff)
		if err != nil {
			notifier.QueueWarning(&models.Warning{
				Error: err.Error(),
				Info:  "Failed to read from hole punch sync stream",
			})
		}
		_, err = stream.Write(buff)
		if err != nil {
			notifier.QueueWarning(&models.Warning{
				Error: err.Error(),
				Info:  "Failed to write to hole punch sync stream",
			})
		}
		notifier.QueueInfo(fmt.Sprintf("Received: %v", buff))
	}
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

// getFileServeUrl path for downloading file
func getFileServeUrl(metadata *models.FileMetadata) string {
	size, err := file_handler.GetFileSize(getFilePath(metadata))
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
