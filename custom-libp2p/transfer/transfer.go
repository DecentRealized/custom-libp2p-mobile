package transfer

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"sync"
)

var _node *models.Node
var _server = &server{}
var _serverLock = &sync.Mutex{}
var _client = &client{}
var _clientLock = &sync.Mutex{}

// Init initializes transfer service
func Init(node *models.Node) error {
	_clientLock.Lock()
	_serverLock.Lock()
	defer _clientLock.Unlock()
	defer _serverLock.Unlock()
	if ServerIsRunning() && ClientIsRunning() {
		_clientLock.Unlock()
		_serverLock.Unlock()
		err := Close()
		_clientLock.Lock()
		_serverLock.Lock()
		if err != nil {
			return err
		}
	}
	_node = node
	err := initServer(node)
	if err != nil {
		return err
	}
	err = initClient(node)
	if err != nil {
		return err
	}
	return nil
}

// Close closes transfer service (client and server)
func Close() error {
	err := closeServer()
	if err != nil {
		return err
	}
	err = closeClient()
	if err != nil {
		return err
	}
	return nil
}

func ServerIsRunning() bool {
	return !(_server == nil || _server.server == nil)
}

func ClientIsRunning() bool {
	return !(_client == nil || _client.client == nil)
}
