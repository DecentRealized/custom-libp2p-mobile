package dummy

import (
	"github.com/DecentRealized/custom-libp2p-mobile/bridge/model"
	flatbuffers "github.com/google/flatbuffers/go"
)

var userName = "Test Caller 2!@#$%^&*()_+';./,><?><<>\\|\""

func GetHelloWorldPayload() []byte {
	requestPayload := flatbuffers.NewBuilder(0)
	userNameFlatbuff := requestPayload.CreateString(userName)
	model.StringRequestStart(requestPayload)
	model.StringRequestAddInput(requestPayload, userNameFlatbuff)
	requestPayload.Finish(model.StringRequestEnd(requestPayload))
	return requestPayload.FinishedBytes()
}

func GetHelloWorldExpectedResponse() []byte {
	return []byte("Hello " + userName + " this is a dummy function!")
}
