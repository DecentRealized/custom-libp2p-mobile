package p2p

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"sync"
)

var _nodeLock = &sync.Mutex{}
var node models.Node
var mdnsService mdns.Service
