package p2p

import (
	"context"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/libp2p/go-libp2p/p2p/protocol/holepunch"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
	"log"
	"sync"
)

func GetDefaultBootstrapPeerAddrInfos() []peer.AddrInfo {
	ds := make([]peer.AddrInfo, 0, len(bootstrapAddrs))

	for i := range bootstrapAddrs {
		info, err := peer.AddrInfoFromString(bootstrapAddrs[i])
		if err != nil {
			log.Printf("failed to convert bootstrapper address to peer addr info: %v", err)
			continue
		}
		ds = append(ds, *info)
	}
	return ds
}

func connectToBootstrapNodes(node host.Host, ctx context.Context) {
	var wg sync.WaitGroup
	for _, peerInfo := range GetDefaultBootstrapPeerAddrInfos() {
		wg.Add(1)
		peerInfo := peerInfo
		go func() {
			defer wg.Done()
			if err := node.Connect(ctx, peerInfo); err != nil {
				log.Printf("Failed to connect to bootstrap node (%v): %s", peerInfo, err)
			} else {
				log.Printf("Connection established with bootstrap node: %v", peerInfo)
			}
		}()
	}
	wg.Wait()
}

func getOptions(privateKey crypto.PrivKey, useInternet bool) []libp2p.Option {
	options := []libp2p.Option{libp2p.Identity(privateKey)}
	if useInternet {
		options = append(options,
			libp2p.EnableHolePunching(holepunch.WithTracer(&HolePunchEventTracer{})),
			libp2p.EnableAutoRelayWithStaticRelays(GetDefaultBootstrapPeerAddrInfos()))
	} else {
		options = append(options, libp2p.AddrsFactory(filterPublicAddresses))
	}
	return options
}

func filterPublicAddresses(addrs []multiaddr.Multiaddr) []multiaddr.Multiaddr {
	filteredAddrs := make([]multiaddr.Multiaddr, 0, len(addrs))
	for _, addr := range addrs {
		if !manet.IsPublicAddr(addr) {
			filteredAddrs = append(filteredAddrs, addr)
		}
	}
	return filteredAddrs
}

// newDHTRouting initializes DHT
func newDHTRouting(node host.Host) (*dht.IpfsDHT, error) {
	ipfsDHT, err := dht.New(context.TODO(), node, dht.Mode(dht.ModeAuto))
	if err != nil {
		return nil, err
	}
	return ipfsDHT, nil
}

// newMDNSService Initialize the MDNS service
func newMDNSService(node host.Host) error {
	mn := mdnsNotifee{peerStore: node.Peerstore()}
	ser := mdns.NewMdnsService(node, mdnsRendezvous, mn)
	err := ser.Start()
	return err
}

type mdnsNotifee struct {
	peerStore peerstore.Peerstore
}

func (mn mdnsNotifee) HandlePeerFound(foundPeer peer.AddrInfo) {
	mn.peerStore.AddAddrs(foundPeer.ID, foundPeer.Addrs, peerstore.PermanentAddrTTL)
	log.Println(mn.peerStore.PeerInfo(foundPeer.ID))
}

type HolePunchEventTracer struct {
}

func (t HolePunchEventTracer) Trace(event *holepunch.Event) {
	log.Printf("Hole punch Evt [%v] From [%v]: %v", event.Type, event.Remote, event.Evt)
}
