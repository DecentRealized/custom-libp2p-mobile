package access_manager

import (
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
	"github.com/libp2p/go-libp2p/core/peer"
	"sync"
)

var accessMap = &sync.Map{}

func AllowNode(id peer.ID) {
	accessMap.Store(id.String(), true)
	notifier.QueueInfo(fmt.Sprintf("Node: %v alowed", id.String()))
}

func BlockNode(id peer.ID) {
	accessMap.Store(id.String(), false)
	notifier.QueueInfo(fmt.Sprintf("Node: %v blocked", id.String()))
}

func IsAllowedNode(id peer.ID) bool {
	val, found := accessMap.Load(id.String())
	if !found {
		return false
	}
	allowed := val.(bool)
	return allowed
}

func IsBlockedNode(id peer.ID) bool {
	val, found := accessMap.Load(id.String())
	if !found {
		return false
	}
	allowed := val.(bool)
	return !allowed
}

func GetAllowedNodes() []peer.ID {
	var allowedNodes []peer.ID
	accessMap.Range(func(k, v any) bool {
		allowed := v.(bool)
		if allowed {
			peerId, err := peer.Decode(k.(string))
			if err != nil {
				return true
			}
			allowedNodes = append(allowedNodes, peerId)
		}
		return true
	})
	return allowedNodes
}

func GetBlockedNodes() []peer.ID {
	var allowedNodes []peer.ID
	accessMap.Range(func(k, v any) bool {
		allowed := v.(bool)
		if !allowed {
			peerId, err := peer.Decode(k.(string))
			if err != nil {
				return true
			}
			allowedNodes = append(allowedNodes, peerId)
		}
		return true
	})
	return allowedNodes
}

func Reset() error {
	accessMap = &sync.Map{}
	return nil
}
