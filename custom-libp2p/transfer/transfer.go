package transfer

import (
	customLibP2P "github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p"
	"sync"
)

var _instanceLock = &sync.Mutex{}

type instance struct {
	node    *customLibP2P.Node
	server  *server
	client  *client
	running bool
}

var singletonInstance *instance

// getInstance get Instance of transfer
func getInstance() *instance {
	if singletonInstance == nil {
		_instanceLock.Lock()
		defer _instanceLock.Unlock()
		if singletonInstance == nil {
			singletonInstance = &instance{
				server: &server{},
				client: &client{},
			}
		}
	}
	return singletonInstance
}

// Init initializes transfer service
func Init(node *customLibP2P.Node) {
	_instance := getInstance()
	_instanceLock.Lock()
	defer _instanceLock.Unlock()
	_instance.node = node
	initServer(node)
	initClient(node)
	_instance.running = true
}

// Close closes transfer service
func Close() error {
	_instanceLock.Lock()
	defer _instanceLock.Unlock()
	if singletonInstance != nil {
		err := closeServer()
		if err != nil {
			return err
		}
		err = closeClient()
		if err != nil {
			return err
		}
	}
	singletonInstance.running = false
	return nil
}
