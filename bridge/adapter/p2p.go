package adapter

import (
	utils "github.com/DecentRealized/custom-libp2p-mobile/bridge/utils"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/p2p"
)

func CreateNode(payload []byte) []byte {
	useInternet, privKey, err := utils.CreateNodeRequest(payload)
	if err != nil {
		return utils.BoolResponse(false, err)
	}

	err = p2p.CreateNode(privKey, useInternet)
	if err != nil {
		return utils.BoolResponse(false, err)
	}

	return utils.BoolResponse(true, err)
}

func StopNode() []byte {
	err := p2p.StopNode()
	if err != nil {
		return utils.BoolResponse(false, err)
	}

	return utils.BoolResponse(true, err)
}

func GetNodeId() []byte {
	nodeId, err := p2p.GetNodeId()
	if err != nil {
		return utils.StringResponse("", err)
	}

	return utils.StringResponse(nodeId.String(), err)
}

func GetListenAddresses() []byte {
	listenAddrs, err := p2p.GetListenAddresses()
	if err != nil {
		return utils.StringArrayResponse([]string{}, err)
	}

	var listenAddrStrings []string
	for _, addr := range listenAddrs {
		listenAddrStrings = append(listenAddrStrings, addr.String())
	}

	return utils.StringArrayResponse(listenAddrStrings, err)
}
