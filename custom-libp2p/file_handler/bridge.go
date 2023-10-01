package file_handler

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"google.golang.org/protobuf/proto"
)

type SetDownloadPathBridgeInput = models.StringMessage

func SetDownloadPathBridge(message proto.Message) (proto.Message, error) {
	pathMsg := message.(*SetDownloadPathBridgeInput)
	return nil, SetDownloadPath(pathMsg.Message)
}

type GetDownloadPathBridgeOutput = models.StringMessage

func GetDownloadPathBridge(proto.Message) (proto.Message, error) {
	return &GetDownloadPathBridgeOutput{Message: GetDownloadPath()}, nil
}
