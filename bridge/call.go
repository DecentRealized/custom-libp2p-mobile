package custom_libp2p_bridge

import (
	"google.golang.org/protobuf/proto"
	"log"
)

// Call ...
func Call(name string, payload []byte) ([]byte, error) {
	bridgeInfo, exists := bridgeMapping[name]
	if !exists {
		return nil, ErrMethodNotImplemented
	}
	bridgedPayload := proto.Clone(bridgeInfo.input)
	log.Println(bridgedPayload)
	log.Println(payload)
	if bridgeInfo.input != nil {
		err := proto.Unmarshal(payload, bridgedPayload)
		if err != nil {
			return nil, err
		}
	}
	response, err := bridgeInfo.function(bridgedPayload)
	if err != nil {
		responseBytes, _ := proto.Marshal(response)
		return responseBytes, err
	}
	return proto.Marshal(response)
}
