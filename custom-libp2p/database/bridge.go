package database

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"google.golang.org/protobuf/proto"
)

type SetDatabaseDirectoryBridgeInput = models.StringMessage

func SetDatabaseDirectoryBridge(request proto.Message) (proto.Message, error) {
	setDatabaseBasePathRequest := request.(*SetDatabaseDirectoryBridgeInput)
	return nil, SetDatabaseDirectory(setDatabaseBasePathRequest.Message)
}

type GetDatabaseDirectoryBridgeOutput = models.StringMessage

func GetDatabaseDirectoryBridge(proto.Message) (proto.Message, error) {
	return &GetDatabaseDirectoryBridgeOutput{Message: GetDatabaseDirectory()}, nil
}
