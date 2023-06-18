package custom_libp2p_bridge

import (
	"github.com/DecentRealized/custom-libp2p-mobile/bridge/model"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCall(t *testing.T) {
	requestPayload := flatbuffers.NewBuilder(0)
	userName := "Test Caller 2"
	userNameFlatbuff := requestPayload.CreateString(userName)
	model.GetHelloMessageRequestStart(requestPayload)
	model.GetHelloMessageRequestAddUserName(requestPayload, userNameFlatbuff)
	requestPayload.Finish(model.GetHelloMessageRequestEnd(requestPayload))
	output, err := Call("getHelloWorld", requestPayload.FinishedBytes())
	outputString := model.GetRootAsStringResponse(output, 0)
	assert.Nil(t, err)
	assert.Equal(t, "Hello "+userName+" this is a dummy function!", string(outputString.Output()))
}
