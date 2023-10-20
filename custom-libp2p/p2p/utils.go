package p2p

import (
	"context"
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/config"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"
	"github.com/libp2p/go-libp2p/p2p/protocol/holepunch"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"time"
)

func getOptions(privateKey crypto.PrivKey, useInternet bool) libp2p.Option {
	options := libp2p.ChainOptions(libp2p.Identity(privateKey))
	if useInternet {
		options = libp2p.ChainOptions(options,
			libp2p.Routing(onlineRouting),
			libp2p.EnableHolePunching(holepunch.WithTracer(&HolePunchEventTracer{})),
			libp2p.EnableAutoRelayWithStaticRelays(dht.GetDefaultBootstrapPeerAddrInfos(),
				autorelay.WithBackoff(5*time.Second),
				autorelay.WithMinInterval(5*time.Second),
				autorelay.WithBootDelay(0)))
	} else {
		options = libp2p.ChainOptions(options,
			libp2p.AddrsFactory(noInternetAddressFactory),
			libp2p.Routing(offlineRouting))
	}
	return options
}

// noInternetAddressFactory Filters public internet addresses
func noInternetAddressFactory(addresses []multiaddr.Multiaddr) []multiaddr.Multiaddr {
	var retVal []multiaddr.Multiaddr
	for _, addr := range addresses {
		if !manet.IsPublicAddr(addr) {
			retVal = append(retVal, addr)
		}
	}
	return retVal
}

// onlineRouting bootstrap routing from default peers
func onlineRouting(node host.Host) (routing.PeerRouting, error) {
	ipfsDHT, err := dht.New(
		context.TODO(),
		node,
		dht.Mode(dht.ModeAuto),
		dht.BootstrapPeersFunc(dht.GetDefaultBootstrapPeerAddrInfos))
	if err != nil {
		return nil, err
	}
	go func() {
		err := ipfsDHT.Bootstrap(context.TODO())
		if err != nil {
			notifier.QueueWarning(&models.Warning{Error: err.Error(), Info: "Failed to bootstrap DHT"})
		}
	}()
	return ipfsDHT, nil
}

// offlineRouting no bootstraps for routing
func offlineRouting(node host.Host) (routing.PeerRouting, error) {
	ipfsDHT, err := dht.New(
		context.TODO(),
		node,
		dht.Mode(dht.ModeClient))
	if err != nil {
		return nil, err
	}
	return ipfsDHT, nil
}

// startMdnsService Initialize the MDNS service
func startMdnsService(node host.Host) error {
	mn := mdnsNotifee{node}
	ser := mdns.NewMdnsService(node, config.MdnsRendezvous, mn)
	err := ser.Start()
	if err != nil {
		return err
	}
	mdnsService = ser
	return nil
}

// stopMdnsService Stop the MDNS service
func stopMdnsService() error {
	if mdnsService != nil {
		err := mdnsService.Close()
		if err != nil {
			return err
		}
		mdnsService = nil
	}
	return nil
}

type mdnsNotifee struct {
	node host.Host
}

func (mn mdnsNotifee) HandlePeerFound(foundPeer peer.AddrInfo) {
	node := mn.node
	node.Peerstore().AddAddrs(foundPeer.ID, foundPeer.Addrs, peerstore.PermanentAddrTTL)
	notifier.QueueInfo(fmt.Sprintf("Found peer via MDNS: %v", foundPeer.ID))
}

type HolePunchEventTracer struct {
}

func (t HolePunchEventTracer) Trace(event *holepunch.Event) {
	notifier.QueueInfo(fmt.Sprintf("Hole punch Event: %v\nEvt: %v", event, event.Evt))
}
