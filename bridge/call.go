package custom_libp2p_bridge

import (
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/bridge/model"
	customLibP2P "github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p"
	flatbuffers "github.com/google/flatbuffers/go"
)

// Call ...
func Call(name string, payload []byte) ([]byte, error) {

	instance := NewInstance()
	var output []byte
	switch name {
	case "getHelloWorld":
		output = instance.getHelloWorld(payload)
	default:
		return nil, fmt.Errorf("not implemented: %s", name)
	}

	return output, nil
}

type instance struct {
	instance *customLibP2P.CustomLibP2P
}

func NewInstance() *instance {
	return &instance{instance: customLibP2P.NewCustomLibP2P()}
}

func (m instance) getHelloWorld(payload []byte) []byte {
	response := flatbuffers.NewBuilder(0)
	request := model.GetRootAsGetHelloMessageRequest(payload, 0)
	output, err := m.instance.GetHelloMessage(m.toString(request.UserName()))
	return m._stringResponse(response, output, err)
}

func (m instance) _stringResponse(response *flatbuffers.Builder, output string, err error) []byte {
	if err != nil {
		outputOffset := response.CreateString(err.Error())
		model.StringResponseStart(response)
		model.StringResponseAddError(response, outputOffset)
		response.Finish(model.StringResponseEnd(response))
		return response.FinishedBytes()
	}
	outputOffset := response.CreateString(output)
	model.StringResponseStart(response)
	model.StringResponseAddOutput(response, outputOffset)
	response.Finish(model.StringResponseEnd(response))
	return response.FinishedBytes()
}

func (m instance) toString(input []byte) string {
	if input == nil {
		return ""
	}

	return string(input)
}
