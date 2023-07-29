package adapter

import (
	"github.com/DecentRealized/custom-libp2p-mobile/bridge/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateKeyPair(t *testing.T) {
	output := CreateKeyPair()
	keyPairResponse := model.GetRootAsKeyPairResponse(output, 0)
	keyPair := model.KeyPair{}
	keyPairResponse.Output(&keyPair)
	assert.Equal(t, len(keyPairResponse.Error()), 0)
	assert.NotEqual(t, len(keyPair.PrivateKeyBytes()), 0)
	assert.NotEqual(t, len(keyPair.PublicKeyBytes()), 0)
}
