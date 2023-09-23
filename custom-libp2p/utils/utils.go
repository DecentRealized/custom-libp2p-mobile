package utils

import (
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/config"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
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
	ds := make([]peer.AddrInfo, 0, len(config.BootstrapAddrs))

	for i := range config.BootstrapAddrs {
		info, err := peer.AddrInfoFromString(config.BootstrapAddrs[i])
		if err != nil {
			notifier.QueueWarning(&models.Warning{
				Error: err.Error(),
				Info:  fmt.Sprintf("Failed to parse bootstrap address: %v", config.BootstrapAddrs[i]),
			})
			continue
		}
		ds = append(ds, *info)
	}
	return ds
}
