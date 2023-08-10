package p2p

import (
	"context"
	custom_libp2p "github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/transfer"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	routedhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/multiformats/go-multiaddr"
	"log"
)

// CreateNode Destroys existing node if present and creates new
func CreateNode(privateKey crypto.PrivKey, useInternet bool) error {
	getInstance()
	_instanceLock.Lock()
	defer _instanceLock.Unlock()
	if singletonInstance.node != nil {
		err := singletonInstance.node.Close()
		if err != nil {
			return err
		}
		singletonInstance.node = nil
	}
	_node, err := libp2p.New(getOptions(privateKey, useInternet)...) // Critical
	if err != nil {
		return err
	}
	dht, err := newDHTRouting(_node) // Critical
	if err != nil {
		return err
	}
	err = newMDNSService(_node) // Non-Critical TODO: Notification to frontend mechanism
	if err != nil {
		log.Printf("Error making mdns: %v", err)
		return err
	}
	node := &custom_libp2p.Node{RoutedHost: *routedhost.Wrap(_node, dht)}
	if err != nil { // Critical
		return err
	}
	if useInternet {
		connectToBootstrapNodes(node, context.TODO()) // Non-Critical
	}
	singletonInstance.node = node
	transfer.Init(node)
	return nil
}

// StopNode Destroys existing node
func StopNode() error {
	getInstance()
	_instanceLock.Lock()
	defer _instanceLock.Unlock()
	if singletonInstance.node == nil {
		return NodeDoesNotExist
	}
	err := singletonInstance.node.Close()
	if err != nil {
		return err
	}
	err = transfer.Close()
	if err != nil {
		return err
	}
	singletonInstance.node = nil
	return nil
}

// GetNodeId Returns node id
func GetNodeId() (peer.ID, error) {
	instance := getInstance()
	if instance.node != nil {
		return instance.node.ID(), nil
	}
	return "", NodeDoesNotExist
}

// GetListenAddresses Returns listen addresses
func GetListenAddresses() ([]multiaddr.Multiaddr, error) {
	instance := getInstance()
	if instance.node != nil {
		return instance.node.Addrs(), nil
	}
	return nil, NodeDoesNotExist
}
