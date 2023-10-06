package access_manager

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"github.com/libp2p/go-libp2p/core/peer"
	"google.golang.org/protobuf/proto"
)

type AllowNodeBridgeInput = models.StringMessage

func AllowNodeBridge(request proto.Message) (proto.Message, error) {
	peerMessage := request.(*AllowNodeBridgeInput)
	peerId, err := peer.Decode(peerMessage.Message)
	if err != nil {
		return nil, err
	}
	return nil, AllowNode(peerId)
}

type IsAllowedNodeBridgeInput = models.StringMessage
type IsAllowedNodeBridgeOutput = models.BoolMessage

func IsAllowedNodeBridge(request proto.Message) (proto.Message, error) {
	peerMessage := request.(*IsAllowedNodeBridgeInput)
	peerId, err := peer.Decode(peerMessage.Message)
	if err != nil {
		return nil, err
	}
	isAllowed, err := IsAllowedNode(peerId)
	if err != nil {
		return nil, err
	}
	return &IsAllowedNodeBridgeOutput{
		Message: isAllowed,
	}, nil
}

type GetAllowedNodesBridgeOutput = models.StringArrayMessage

func GetAllowedNodesBridge(proto.Message) (proto.Message, error) {
	allowedNodes, err := GetAllowedNodes()
	if err != nil {
		return nil, err
	}
	stringAllowedIds := make([]string, len(allowedNodes))
	for i, address := range allowedNodes {
		stringAllowedIds[i] = address.String()
	}
	return &GetAllowedNodesBridgeOutput{Message: stringAllowedIds}, nil
}

type BlockNodeBridgeInput = models.StringMessage

func BlockNodeBridge(request proto.Message) (proto.Message, error) {
	peerMessage := request.(*BlockNodeBridgeInput)
	peerId, err := peer.Decode(peerMessage.Message)
	if err != nil {
		return nil, err
	}
	return nil, BlockNode(peerId)
}

type IsBlockedNodeBridgeInput = models.StringMessage
type IsBlockedNodeBridgeOutput = models.BoolMessage

func IsBlockedNodeBridge(request proto.Message) (proto.Message, error) {
	peerMessage := request.(*IsBlockedNodeBridgeInput)
	peerId, err := peer.Decode(peerMessage.Message)
	if err != nil {
		return nil, err
	}
	isBlocked, err := IsBlockedNode(peerId)
	if err != nil {
		return nil, err
	}
	return &IsBlockedNodeBridgeOutput{
		Message: isBlocked,
	}, nil
}

type GetBlockedNodesBridgeOutput = models.StringArrayMessage

func GetBlockedNodesBridge(proto.Message) (proto.Message, error) {
	blockedNodes, err := GetBlockedNodes()
	if err != nil {
		return nil, err
	}
	stringBlockedIds := make([]string, len(blockedNodes))
	for i, address := range blockedNodes {
		stringBlockedIds[i] = address.String()
	}
	return &GetBlockedNodesBridgeOutput{Message: stringBlockedIds}, nil
}
