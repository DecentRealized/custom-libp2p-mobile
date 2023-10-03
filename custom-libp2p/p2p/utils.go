package p2p

import (
	"context"
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/config"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/utils"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"
	"github.com/libp2p/go-libp2p/p2p/protocol/holepunch"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"time"
)

func getOptions(privateKey crypto.PrivKey, useInternet bool) []libp2p.Option {
	options := []libp2p.Option{libp2p.Identity(privateKey), libp2p.ForceReachabilityPrivate()}
	if useInternet {
		options = append(options,
			libp2p.EnableHolePunching(holepunch.WithTracer(&HolePunchEventTracer{})),
			libp2p.EnableAutoRelayWithStaticRelays(utils.GetDefaultBootstrapPeerAddrInfos(),
				autorelay.WithBackoff(5*time.Second),
				autorelay.WithMinInterval(5*time.Second),
				autorelay.WithBootDelay(0)))
	} else {
		options = append(options, libp2p.AddrsFactory(noInternetAddressFactory))
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

// newDHTRouting initializes DHT
func newDHTRouting(node *host.Host) (*dht.IpfsDHT, error) {
	ipfsDHT, err := dht.New(context.TODO(), *node, dht.Mode(dht.ModeAuto))
	if err != nil {
		return nil, err
	}
	return ipfsDHT, nil
}

// newMDNSService Initialize the MDNS service
func newMDNSService(node *host.Host) (mdns.Service, error) {
	mn := mdnsNotifee{node}
	ser := mdns.NewMdnsService(*node, config.MdnsRendezvous, mn)
	err := ser.Start()
	return ser, err
}

type mdnsNotifee struct {
	node *host.Host
}

func (mn mdnsNotifee) HandlePeerFound(foundPeer peer.AddrInfo) {
	node := *mn.node
	node.Peerstore().AddAddrs(foundPeer.ID, foundPeer.Addrs, peerstore.PermanentAddrTTL)
	notifier.QueueInfo(fmt.Sprintf("Found peer via MDNS: %v", foundPeer.ID))
}

type HolePunchEventTracer struct {
}

func (t HolePunchEventTracer) Trace(event *holepunch.Event) {
	notifier.QueueInfo(fmt.Sprintf("Hole punch Evt [%v] From [%v]", event.Type, event.Remote))
}
