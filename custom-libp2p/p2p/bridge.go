package p2p

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"github.com/libp2p/go-libp2p/core/crypto"
	"google.golang.org/protobuf/proto"
)

type CreateNodeBridgeInput = models.CreateNodeRequest

func CreateNodeBridge(request proto.Message) (proto.Message, error) {
	createNodeRequest := request.(*CreateNodeBridgeInput)
	useInternet := createNodeRequest.UseInternet
	privateKey, err := crypto.UnmarshalPrivateKey(createNodeRequest.PrivateKey)
	if err != nil {
		return nil, err
	}
	err = CreateNode(privateKey, useInternet)
	return nil, err
}

func StopNodeBridge(proto.Message) (proto.Message, error) {
	err := StopNode()
	return nil, err
}

type GetNodeIdBridgeOutput = models.StringMessage

func GetNodeIdBridge(proto.Message) (proto.Message, error) {
	peerId, err := GetNodeId()
	if err != nil {
		return nil, err
	}
	return &GetNodeIdBridgeOutput{Message: peerId.String()}, err
}

type GetListenAddressBridgeOutput = models.StringArrayMessage

func GetListenAddressesBridge(proto.Message) (proto.Message, error) {
	addresses, err := GetListenAddresses()
	if err != nil {
		return nil, err
	}
	stringAddresses := make([]string, len(addresses))
	for i, address := range addresses {
		stringAddresses[i] = address.String()
	}
	return &GetListenAddressBridgeOutput{Message: stringAddresses}, nil
}
