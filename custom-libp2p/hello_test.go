package custom_libp2p

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCustomLibP2P_GetHelloMessage(t *testing.T) {
	instance := NewCustomLibP2P()
	output, err := instance.GetHelloMessage(testerUserName)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Hello "+testerUserName+" this is a dummy function!", output)
	t.Log("Output: " + output)
}
