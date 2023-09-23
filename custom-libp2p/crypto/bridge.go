package crypto

import (
	"crypto/rand"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"github.com/libp2p/go-libp2p/core/crypto"
	"google.golang.org/protobuf/proto"
)

type CreateKeyPairBridgeOutput = models.KeyPair

func CreateKeyPairBridge(proto.Message) (proto.Message, error) {
	privateKey, publicKey, err := crypto.GenerateECDSAKeyPair(rand.Reader)
	if err != nil {
		return nil, err
	}
	privateKeyBytes, err := crypto.MarshalPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	publicKeyBytes, err := crypto.MarshalPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	response := &CreateKeyPairBridgeOutput{
		PrivateKey: privateKeyBytes,
		PublicKey:  publicKeyBytes,
	}
	return response, nil
}
