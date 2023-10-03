package models

import (
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	routedhost "github.com/libp2p/go-libp2p/p2p/host/routed"
)

type Node struct {
	routedhost.RoutedHost
	Mdns mdns.Service
}
