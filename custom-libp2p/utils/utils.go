package utils

import (
	"fmt"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
	"os"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetDefaultBootstrapPeerAddrInfos() []peer.AddrInfo {
	return dht.GetDefaultBootstrapPeerAddrInfos()
}
