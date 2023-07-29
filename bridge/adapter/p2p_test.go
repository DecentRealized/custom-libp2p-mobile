package adapter

import (
	"github.com/DecentRealized/custom-libp2p-mobile/bridge/dummy"
	"github.com/DecentRealized/custom-libp2p-mobile/bridge/model"
	utils "github.com/DecentRealized/custom-libp2p-mobile/bridge/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateNode(t *testing.T) {
	TestCreateNodeOnline(t)
	TestCreateNodeOffline(t)
}

func TestCreateNodeOnline(t *testing.T) {
	output := CreateNode(dummy.CreateNodePayload(true, 0))
	createNodeResponse := model.GetRootAsBoolResponse(output, 0)
	assert.Equal(t, dummy.CreateNodeExpectedResponse(), createNodeResponse.Table().Bytes)
}

func TestCreateNodeOffline(t *testing.T) {
	output := CreateNode(dummy.CreateNodePayload(false, 0))
	createNodeResponse := model.GetRootAsBoolResponse(output, 0)
	assert.Equal(t, dummy.CreateNodeExpectedResponse(), createNodeResponse.Table().Bytes)
}

func TestStopNode(t *testing.T) {
	TestCreateNodeOnline(t)
	output := StopNode()
	stopNodeResponse := model.GetRootAsBoolResponse(output, 0)
	assert.Equal(t, dummy.StopNodeExpectedResponse(), stopNodeResponse.Table().Bytes)

	TestCreateNodeOffline(t)
	output = StopNode()
	stopNodeResponse = model.GetRootAsBoolResponse(output, 0)
	assert.Equal(t, dummy.StopNodeExpectedResponse(), stopNodeResponse.Table().Bytes)
}

func TestStopNodeWhenNotRunning(t *testing.T) {
	StopNode()
	output := StopNode()
	stopNodeResponse := model.GetRootAsBoolResponse(output, 0)
	assert.Equal(t, dummy.NotRunningExpectedResponse(), stopNodeResponse.Table().Bytes)
}

func TestGetNodeId(t *testing.T) {
	TestCreateNodeOnline(t)
	output := GetNodeId()
	getNodeIdResponse := model.GetRootAsStringResponse(output, 0)
	assert.Equal(t, dummy.GetNodeIdExpectedResponse(0), getNodeIdResponse.Table().Bytes)

	TestCreateNodeOffline(t)
	output = GetNodeId()
	getNodeIdResponse = model.GetRootAsStringResponse(output, 0)
	assert.Equal(t, dummy.GetNodeIdExpectedResponse(0), getNodeIdResponse.Table().Bytes)
}

func TestGetNodeIdWhenNotRunning(t *testing.T) {
	StopNode()
	output := GetNodeId()
	getNodeIdResponse := model.GetRootAsStringResponse(output, 0)
	assert.Equal(t, dummy.NotRunningExpectedResponse(), getNodeIdResponse.Table().Bytes)
}

func TestGetListenAddresses(t *testing.T) {
	TestCreateNodeOnline(t)
	output := GetListenAddresses()
	getNodeIdResponse := model.GetRootAsStringArrayResponse(output, 0)
	assert.NotEqual(t, 0, getNodeIdResponse.OutputLength())
	assert.Equal(t, "", utils.ToString(getNodeIdResponse.Error()))
	for i := 0; i < getNodeIdResponse.OutputLength(); i++ {
		assert.NotEqual(t, 0, len(getNodeIdResponse.Output(i)))
	}

	TestCreateNodeOffline(t)
	output = GetListenAddresses()
	getNodeIdResponse = model.GetRootAsStringArrayResponse(output, 0)
	assert.NotEqual(t, 0, getNodeIdResponse.OutputLength())
	assert.Equal(t, "", utils.ToString(getNodeIdResponse.Error()))
	for i := 0; i < getNodeIdResponse.OutputLength(); i++ {
		assert.NotEqual(t, 0, len(getNodeIdResponse.Output(i)))
	}
}

func TestGetListenAddressesWhenNotRunning(t *testing.T) {
	StopNode()
	output := GetListenAddresses()
	getNodeIdResponse := model.GetRootAsStringArrayResponse(output, 0)
	assert.Equal(t, dummy.NotRunningExpectedResponse(), getNodeIdResponse.Table().Bytes)
}
