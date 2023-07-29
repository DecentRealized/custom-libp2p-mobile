package transfer

import (
	"bytes"
	"context"
	customLibP2P "github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p"
	"github.com/libp2p/go-libp2p/core"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"log"
)

// connectedWithoutRelay check if peer is connected to node without relay
func connectedWithoutRelay(node *customLibP2P.Node, _peer peer.ID) bool {
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
func newHolePunchSyncStream(node *customLibP2P.Node, peerId peer.ID) {
	stream, err := node.NewStream(context.TODO(), peerId, holePunchSyncStreamProtocolID)
	if err != nil {
		log.Printf("Error creating hole punch sync stream: %v", err)
	}
	buff := bytes.Repeat([]byte{1}, holePunchPacketSize)
	for i := 0; i < holePunchRetries; i++ {
		_, err := stream.Write(buff)
		if err != nil {
			log.Printf("Error writing to hole punch sync stream: %v", err)
		}
		_, err = stream.Read(buff)
		if err != nil {
			log.Printf("Error reading from hole punch sync stream: %v", err)
		}
	}
}

// handleHolePunchSyncStream handles new hole punch stream
func handleHolePunchSyncStream(stream core.Stream) {
	buff := bytes.Repeat([]byte{0}, holePunchPacketSize)
	for i := 0; i < holePunchRetries; i++ {
		_, err := stream.Read(buff)
		if err != nil {
			log.Printf("Error reading from hole punch sync stream: %v", err)
		}
		_, err = stream.Write(buff)
		if err != nil {
			log.Printf("Error Write to hole punch sync stream: %v", err)
		}
		log.Printf("Received: %v", buff)
	}
}
