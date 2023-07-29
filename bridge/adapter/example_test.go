package adapter

import (
	"github.com/DecentRealized/custom-libp2p-mobile/bridge/dummy"
	"github.com/DecentRealized/custom-libp2p-mobile/bridge/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetHelloWorld(t *testing.T) {
	output := GetHelloWorld(dummy.GetHelloWorldPayload())
	outputString := model.GetRootAsStringResponse(output, 0)
	assert.Equal(t, dummy.GetHelloWorldExpectedResponse(), outputString.Output())
}
