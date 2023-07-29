package adapter

import (
	utils "github.com/DecentRealized/custom-libp2p-mobile/bridge/utils"
	"github.com/libp2p/go-libp2p/core/crypto"
)

func CreateKeyPair() []byte {
	privKey, pubKey, err := crypto.GenerateKeyPair(crypto.ECDSA, -1)
	if err != nil {
		return utils.KeyPairResponse(nil, nil, err)
	}

	privKeyBytes, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return utils.KeyPairResponse(nil, nil, err)
	}

	pubKeyBytes, err := crypto.MarshalPublicKey(pubKey)
	if err != nil {
		return utils.KeyPairResponse(privKeyBytes, nil, err)
	}

	return utils.KeyPairResponse(privKeyBytes, pubKeyBytes, nil)
}
