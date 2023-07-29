package adapter

import (
	utils "github.com/DecentRealized/custom-libp2p-mobile/bridge/utils"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/example"
)

func GetHelloWorld(payload []byte) []byte {
	userName, err := utils.StringRequest(payload)
	if err != nil {
		return utils.StringResponse("", err)
	}
	output, err := example.GetHelloMessage(userName)
	return utils.StringResponse(output, err)
}
