package access_manager

import (
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/database"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
	"github.com/dgraph-io/badger/v4"
	"github.com/libp2p/go-libp2p/core/peer"
	"path/filepath"
)

var _dbKeyBase = "access_manager/accessMap"

func AllowNode(id peer.ID) error {
	key := []byte(filepath.Join(_dbKeyBase, id.String()))
	value := []byte{1}
	err := database.Store(key, value)
	if err != nil {
		return err
	}
	notifier.QueueInfo(fmt.Sprintf("Node: %v alowed", id.String()))
	return nil
}

func BlockNode(id peer.ID) error {
	key := []byte(filepath.Join(_dbKeyBase, id.String()))
	value := []byte{0}
	err := database.Store(key, value)
	if err != nil {
		return err
	}
	notifier.QueueInfo(fmt.Sprintf("Node: %v blocked", id.String()))
	return nil
}

func IsAllowedNode(id peer.ID) (bool, error) {
	key := []byte(filepath.Join(_dbKeyBase, id.String()))
	value, err := database.Load(key)
	if err == badger.ErrKeyNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	found := value[0] == 1
	return found, nil
}

func IsBlockedNode(id peer.ID) (bool, error) {
	key := []byte(filepath.Join(_dbKeyBase, id.String()))
	value, err := database.Load(key)
	if err == badger.ErrKeyNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	found := value[0] == 0
	return found, nil
}

func GetAllowedNodes() ([]peer.ID, error) {
	var allowedNodes []peer.ID
	var decodeError error
	rangeError := database.Range([]byte(_dbKeyBase), func(key, value []byte) bool {
		allowed := value[0] == 1
		if allowed {
			peerIdStr := filepath.Base(string(key))
			peerId, err := peer.Decode(peerIdStr)
			if err != nil {
				decodeError = err
				return false
			}
			allowedNodes = append(allowedNodes, peerId)
		}
		return true
	})
	if decodeError != nil {
		return nil, decodeError
	}
	if rangeError != nil {
		return nil, rangeError
	}
	return allowedNodes, nil
}

func GetBlockedNodes() ([]peer.ID, error) {
	var blockedNodes []peer.ID
	var decodeError error
	rangeError := database.Range([]byte(_dbKeyBase), func(key, value []byte) bool {
		blocked := value[0] == 0
		if blocked {
			peerIdStr := filepath.Base(string(key))
			peerId, err := peer.Decode(peerIdStr)
			if err != nil {
				decodeError = err
				return false
			}
			blockedNodes = append(blockedNodes, peerId)
		}
		return true
	})
	if decodeError != nil {
		return nil, decodeError
	}
	if rangeError != nil {
		return nil, rangeError
	}
	return blockedNodes, nil
}
