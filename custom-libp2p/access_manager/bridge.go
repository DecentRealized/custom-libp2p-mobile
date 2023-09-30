package access_manager

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"github.com/libp2p/go-libp2p/core/peer"
	"google.golang.org/protobuf/proto"
)

type AllowNodeBridgeInput = models.StringMessage

func AllowNodeBridge(request proto.Message) (proto.Message, error) {
	peerMessage := request.(*AllowNodeBridgeInput)
	peerId, err := peer.Decode(peerMessage.String())
	if err != nil {
		return nil, err
	}
	AllowNode(peerId)
	return nil, nil
}

type IsAllowedNodeBridgeInput = models.StringMessage
type IsAllowedNodeBridgeOutput = models.BoolMessage

func IsAllowedNodeBridge(request proto.Message) (proto.Message, error) {
	peerMessage := request.(*IsAllowedNodeBridgeInput)
	peerId, err := peer.Decode(peerMessage.String())
	if err != nil {
		return nil, err
	}
	isAllowed := IsAllowedNode(peerId)
	return &IsAllowedNodeBridgeOutput{
		Message: isAllowed,
	}, nil
}

type GetAllowedNodesBridgeOutput = models.StringArrayMessage

func GetAllowedNodesBridge(proto.Message) (proto.Message, error) {
	allowedNodes := GetAllowedNodes()
	stringAllowedIds := make([]string, len(allowedNodes))
	for i, address := range allowedNodes {
		stringAllowedIds[i] = address.String()
	}
	return &GetAllowedNodesBridgeOutput{Message: stringAllowedIds}, nil
}

type BlockNodeBridgeInput = models.StringMessage

func BlockNodeBridge(request proto.Message) (proto.Message, error) {
	peerMessage := request.(*BlockNodeBridgeInput)
	peerId, err := peer.Decode(peerMessage.String())
	if err != nil {
		return nil, err
	}
	BlockNode(peerId)
	return nil, nil
}

type IsBlockedNodeBridgeInput = models.StringMessage
type IsBlockedNodeBridgeOutput = models.BoolMessage

func IsBlockedNodeBridge(request proto.Message) (proto.Message, error) {
	peerMessage := request.(*IsBlockedNodeBridgeInput)
	peerId, err := peer.Decode(peerMessage.String())
	if err != nil {
		return nil, err
	}
	isBlocked := IsBlockedNode(peerId)
	return &IsBlockedNodeBridgeOutput{
		Message: isBlocked,
	}, nil
}

type GetBlockedNodesBridgeOutput = models.StringArrayMessage

func GetBlockedNodesBridge(proto.Message) (proto.Message, error) {
	blockedNodes := GetBlockedNodes()
	stringBlockedIds := make([]string, len(blockedNodes))
	for i, address := range blockedNodes {
		stringBlockedIds[i] = address.String()
	}
	return &GetAllowedNodesBridgeOutput{Message: stringBlockedIds}, nil
}
