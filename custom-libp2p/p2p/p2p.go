package p2p

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p"
	"sync"
)

var _instanceLock = &sync.Mutex{}

type instance struct {
	node *customLibP2P.Node
}

var singletonInstance *instance

func getInstance() *instance {
	if singletonInstance == nil {
		_instanceLock.Lock()
		defer _instanceLock.Unlock()
		if singletonInstance == nil {
			singletonInstance = &instance{}
		}
	}
	return singletonInstance
}
