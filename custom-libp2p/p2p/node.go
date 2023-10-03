package p2p

import (
	"context"
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/database"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/file_handler"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/transfer"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	routedhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/multiformats/go-multiaddr"
)

// CreateNode Destroys existing node if present and creates new
func CreateNode(privateKey crypto.PrivKey, useInternet bool) error {
	// Create Node
	_nodeLock.Lock()
	defer _nodeLock.Unlock()
	if node != nil {
		_nodeLock.Unlock()
		err := StopNode()
		_nodeLock.Lock()
		if err != nil {
			return err
		}
	}
	_node, err := libp2p.New(getOptions(privateKey, useInternet)...) // Critical
	if err != nil {
		return err
	}
	dht, err := newDHTRouting(&_node) // Critical
	if err != nil {
		return err
	}
	mdns, err := newMDNSService(&_node)
	if err != nil {
		notifier.QueueWarning(&models.Warning{Error: err.Error(), Info: "Failed to create MDNS service"})
	}
	node = &models.Node{RoutedHost: *routedhost.Wrap(_node, dht), Mdns: mdns}
	// Initialize Other Modules
	err = transfer.Init(node)
	if err != nil {
		return err
	}
	err = database.Init(node)
	if err != nil {
		return err
	}
	notifier.QueueInfo(fmt.Sprintf("Node with ID %s created", node.ID()))
	return nil
}

// StopNode Destroys existing node
func StopNode() error {
	_nodeLock.Lock()
	defer _nodeLock.Unlock()
	if node == nil {
		return ErrNodeDoesNotExist
	}
	err := node.Mdns.Close()
	if err != nil {
		notifier.QueueWarning(&models.Warning{Error: err.Error(), Info: "Failed to close MDNS service"})
	}
	err = node.Close()
	if err != nil {
		return err
	}

	// Close/Reset All Modules
	err = notifier.Reset()
	if err != nil {
		return err
	}
	err = transfer.Close()
	if err != nil {
		return err
	}
	err = file_handler.Reset()
	if err != nil {
		return err
	}
	err = database.Close()
	if err != nil {
		return err
	}

	node = nil
	return nil
}

// GetNodeId Returns node id
func GetNodeId() (peer.ID, error) {
	if node != nil {
		return node.ID(), nil
	}
	return "", ErrNodeDoesNotExist
}

// GetListenAddresses Returns listen addresses
func GetListenAddresses() ([]multiaddr.Multiaddr, error) {
	if node != nil {
		return node.Addrs(), nil
	}
	return nil, ErrNodeDoesNotExist
}

// ConnectToPeer Connects to the given node
func ConnectToPeer(peerId peer.ID) error {
	if node != nil {
		return node.Connect(context.TODO(), peer.AddrInfo{ID: peerId})
	}
	return ErrNodeDoesNotExist
}

// CheckConnectionStatus Checks connection status
func CheckConnectionStatus(peerId peer.ID) ([]network.Conn, error) {
	if node != nil {
		return node.Network().ConnsToPeer(peerId), nil
	}
	return nil, ErrNodeDoesNotExist
}
