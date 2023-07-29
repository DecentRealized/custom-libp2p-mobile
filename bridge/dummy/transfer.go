package dummy

import (
	"errors"
	"github.com/DecentRealized/custom-libp2p-mobile/bridge/model"
	utils "github.com/DecentRealized/custom-libp2p-mobile/bridge/utils"
	flatbuffers "github.com/google/flatbuffers/go"
)

var filePath = "./transfer.go"
var peerId = getNodeIdFromPrivateKey(1)
var message = "Hello World!"

// ServeFilePayload Serve GOPATH/README.md to second node in private keys
func ServeFilePayload() []byte {
	requestPayload := flatbuffers.NewBuilder(0)
	filePathFlatbuff := requestPayload.CreateString(filePath)
	peerIdFlatbuff := requestPayload.CreateString(peerId.String())
	model.ServeFileRequestStart(requestPayload)
	model.ServeFileRequestAddFilePath(requestPayload, filePathFlatbuff)
	model.ServeFileRequestAddPeerId(requestPayload, peerIdFlatbuff)
	requestPayload.Finish(model.ServeFileRequestEnd(requestPayload))
	return requestPayload.FinishedBytes()
}

// ServeFileExpectedResponse returns serve file expected response
func ServeFileExpectedResponse() []byte {
	return utils.BoolResponse(true, nil)
}

// PeerNotFoundResponse response when client is down
func PeerNotFoundResponse() []byte {
	return utils.BoolResponse(false, errors.New("routing: not found"))
}

// TransferNotRunningResponse response when client is down
func TransferNotRunningResponse() []byte {
	return utils.BoolResponse(false, errors.New("transfer not running"))
}

// SendMessagePayload Send var message to second node
func SendMessagePayload() []byte {
	requestPayload := flatbuffers.NewBuilder(0)
	messageFlatBuff := requestPayload.CreateString(message)
	peerIdFlatbuff := requestPayload.CreateString(peerId.String())
	model.SendMessageRequestStart(requestPayload)
	model.SendMessageRequestAddMessage(requestPayload, messageFlatBuff)
	model.SendMessageRequestAddPeerId(requestPayload, peerIdFlatbuff)
	requestPayload.Finish(model.SendMessageRequestEnd(requestPayload))
	return requestPayload.FinishedBytes()
}

// PauseDownloadPayload pause download from second node
func PauseDownloadPayload() []byte {
	// TODO
	return nil
}
