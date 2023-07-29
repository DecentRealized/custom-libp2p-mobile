package butils

import (
	"github.com/DecentRealized/custom-libp2p-mobile/bridge/model"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

// StringRequest parses payload and returns string
func StringRequest(payload []byte) (string, error) {
	request := model.GetRootAsStringRequest(payload, 0)
	return ToString(request.Input()), nil
}

// CreateNodeRequest parses payload and returns useInternet, and private key
func CreateNodeRequest(payload []byte) (bool, crypto.PrivKey, error) {
	request := model.GetRootAsCreateNodeRequest(payload, 0)
	privKey, err := crypto.UnmarshalPrivateKey(request.PrivateKeyBytes())
	return request.UseInternet(), privKey, err
}

// ServeFileRequest parses payload and returns path, and peerId to whom to serve the file to
func ServeFileRequest(payload []byte) (string, peer.ID, error) {
	request := model.GetRootAsServeFileRequest(payload, 0)
	peerId, err := peer.Decode(ToString(request.PeerId()))
	return ToString(request.FilePath()), peerId, err
}

// SendMessageRequest parses payload and returns message and peerId
func SendMessageRequest(payload []byte) (string, peer.ID, error) {
	request := model.GetRootAsSendMessageRequest(payload, 0)
	peerId, err := peer.Decode(ToString(request.PeerId()))
	return ToString(request.Message()), peerId, err
}

// DownloadRequest parses payload and returns fileSHA256 and peerId
func DownloadRequest(payload []byte) (string, peer.ID, error) {
	request := model.GetRootAsDownloadRequest(payload, 0)
	peerId, err := peer.Decode(ToString(request.PeerId()))
	return ToString(request.FileSha256()), peerId, err
}
