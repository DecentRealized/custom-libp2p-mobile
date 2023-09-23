package p2p

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"sync"
)

var _nodeLock = &sync.Mutex{}
var node *models.Node
