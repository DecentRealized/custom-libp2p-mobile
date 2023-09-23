package example

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"google.golang.org/protobuf/proto"
)

type GetHelloMessageBridgeInput = models.StringMessage
type GetHelloMessageBridgeOutput = models.StringMessage

func GetHelloMessageBridge(request proto.Message) (proto.Message, error) {
	username := request.(*GetHelloMessageBridgeInput)
	response := &GetHelloMessageBridgeOutput{
		Message: "Hello " + username.GetMessage() + " this is a dummy function!",
	}
	return response, nil
}
