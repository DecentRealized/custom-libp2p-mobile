package custom_libp2p_bridge

import (
	"testing"
)

// TODO: Tests
func TestCall(t *testing.T) {
	// *****************************************************************************************************************
	t.Log("1. getHelloWorld")
	TestHelloWorldCall(t)
	// *****************************************************************************************************************
	t.Log("2. createKeyPair")
	TestCreateKeyPairCall(t)
	// *****************************************************************************************************************
	t.Log("3. createNode")
	TestCreateNodeCall(t)
	// *****************************************************************************************************************
	t.Log("4. getNodeId")
	TestGetNodeIdCall(t)
	// *****************************************************************************************************************
	t.Log("5. getListenAddresses")
	TestGetListenAddressesCall(t)
	// *****************************************************************************************************************
	t.Log("6. allowNode")
	TestAllowNodeCall(t)
	// *****************************************************************************************************************
	t.Log("7. isAllowedNode")
	TestIsAllowedNodeCall(t)
	// *****************************************************************************************************************
	t.Log("8. getAllowedNodes")
	TestGetAllowedNodesCall(t)
	// *****************************************************************************************************************
	t.Log("9. denyNode")
	TestDenyNodeCall(t)
	// *****************************************************************************************************************
	t.Log("10. sendFile")
	TestSendFileCall(t)
	// *****************************************************************************************************************
	t.Log("11. sendMessage")
	TestSendMessageCall(t)
	// *****************************************************************************************************************
	t.Log("12. stopNode")
	TestStopNodeCall(t)
	// *****************************************************************************************************************
}

func TestHelloWorldCall(t *testing.T) {
	a, b := Call("getHelloMessage", []byte{})
	t.Log(a, b)
	//outputString := model.GetRootAsStringResponse(output, 0)
	//assert.Nil(t, err)
	//assert.Equal(t, dummy.GetHelloWorldExpectedResponse(), outputString.Output())
}

func TestCreateKeyPairCall(t *testing.T) {
	t.Log("TODO")
	// assert.Equal(t, 1, 2)
}

func TestCreateNodeCall(t *testing.T) {
	t.Log("TODO")
	// assert.Equal(t, 1, 2)
}

func TestGetNodeIdCall(t *testing.T) {
	t.Log("TODO")
	// assert.Equal(t, 1, 2)
}

func TestGetListenAddressesCall(t *testing.T) {
	t.Log("TODO")
	// assert.Equal(t, 1, 2)
}

func TestAllowNodeCall(t *testing.T) {
	t.Log("TODO")
	// assert.Equal(t, 1, 2)
}

func TestIsAllowedNodeCall(t *testing.T) {
	t.Log("TODO")
	// assert.Equal(t, 1, 2)
}

func TestGetAllowedNodesCall(t *testing.T) {
	t.Log("TODO")
	// assert.Equal(t, 1, 2)
}

func TestDenyNodeCall(t *testing.T) {
	t.Log("TODO")
	// assert.Equal(t, 1, 2)
}

func TestSendFileCall(t *testing.T) {
	t.Log("TODO")
	// assert.Equal(t, 1, 2)
}

func TestSendMessageCall(t *testing.T) {
	t.Log("TODO")
	// assert.Equal(t, 1, 2)
}

func TestStopNodeCall(t *testing.T) {
	t.Log("TODO")
	// assert.Equal(t, 1, 2)
}
