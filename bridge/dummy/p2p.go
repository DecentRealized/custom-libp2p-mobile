package dummy

import (
	"errors"
	"github.com/DecentRealized/custom-libp2p-mobile/bridge/model"
	utils "github.com/DecentRealized/custom-libp2p-mobile/bridge/utils"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

var privateKeys = [][]byte{
	{
		8, 3, 18, 121, 48, 119, 2, 1, 1, 4, 32, 70, 201, 191, 214, 174, 124, 146, 122, 184, 140, 151, 43, 145, 25, 219, 51,
		253, 66, 63, 230, 26, 153, 182, 248, 159, 193, 150, 2, 8, 92, 206, 5, 160, 10, 6, 8, 42, 134, 72, 206, 61, 3, 1, 7,
		161, 68, 3, 66, 0, 4, 19, 65, 242, 195, 5, 198, 184, 105, 23, 233, 56, 115, 156, 103, 64, 114, 181, 34, 177, 44,
		125, 70, 250, 13, 55, 7, 37, 88, 155, 247, 188, 183, 54, 124, 10, 124, 99, 143, 82, 131, 8, 78, 33, 132, 66, 193,
		234, 173, 113, 23, 189, 51, 56, 7, 244, 106, 13, 119, 23, 148, 184, 172, 165, 28,
	},
	{
		8, 3, 18, 121, 48, 119, 2, 1, 1, 4, 32, 70, 201, 191, 214, 174, 124, 146, 122, 184, 140, 151, 43, 145, 25, 219, 51,
		252, 66, 63, 230, 26, 153, 182, 248, 159, 193, 150, 2, 8, 92, 206, 5, 160, 10, 6, 8, 42, 134, 72, 206, 61, 3, 1, 7,
		161, 68, 3, 66, 0, 4, 19, 65, 242, 195, 5, 198, 184, 105, 23, 233, 56, 115, 156, 103, 64, 114, 181, 34, 177, 44,
		125, 70, 250, 13, 55, 7, 37, 88, 155, 247, 188, 183, 54, 124, 10, 124, 99, 143, 82, 131, 8, 78, 33, 132, 66, 193,
		234, 173, 113, 23, 189, 51, 56, 7, 244, 106, 13, 119, 23, 148, 184, 172, 165, 28,
	},
	{
		8, 3, 18, 121, 48, 119, 2, 1, 1, 4, 32, 70, 201, 191, 214, 174, 124, 146, 122, 184, 140, 151, 43, 145, 25, 219, 51,
		251, 66, 63, 230, 26, 153, 182, 248, 159, 193, 150, 2, 8, 92, 206, 5, 160, 10, 6, 8, 42, 134, 72, 206, 61, 3, 1, 7,
		161, 68, 3, 66, 0, 4, 19, 65, 242, 195, 5, 198, 184, 105, 23, 233, 56, 115, 156, 103, 64, 114, 181, 34, 177, 44,
		125, 70, 250, 13, 55, 7, 37, 88, 155, 247, 188, 183, 54, 124, 10, 124, 99, 143, 82, 131, 8, 78, 33, 132, 66, 193,
		234, 173, 113, 23, 189, 51, 56, 7, 244, 106, 13, 119, 23, 148, 184, 172, 165, 28,
	},
	{
		8, 3, 18, 121, 48, 119, 2, 1, 1, 4, 32, 70, 201, 191, 214, 174, 124, 146, 122, 184, 140, 151, 43, 145, 25, 219, 51,
		250, 66, 63, 230, 26, 153, 182, 248, 159, 193, 150, 2, 8, 92, 206, 5, 160, 10, 6, 8, 42, 134, 72, 206, 61, 3, 1, 7,
		161, 68, 3, 66, 0, 4, 19, 65, 242, 195, 5, 198, 184, 105, 23, 233, 56, 115, 156, 103, 64, 114, 181, 34, 177, 44,
		125, 70, 250, 13, 55, 7, 37, 88, 155, 247, 188, 183, 54, 124, 10, 124, 99, 143, 82, 131, 8, 78, 33, 132, 66, 193,
		234, 173, 113, 23, 189, 51, 56, 7, 244, 106, 13, 119, 23, 148, 184, 172, 165, 28,
	},
	{
		8, 3, 18, 121, 48, 119, 2, 1, 1, 4, 32, 70, 201, 191, 214, 174, 124, 146, 122, 184, 140, 151, 43, 145, 25, 219, 51,
		249, 66, 63, 230, 26, 153, 182, 248, 159, 193, 150, 2, 8, 92, 206, 5, 160, 10, 6, 8, 42, 134, 72, 206, 61, 3, 1, 7,
		161, 68, 3, 66, 0, 4, 19, 65, 242, 195, 5, 198, 184, 105, 23, 233, 56, 115, 156, 103, 64, 114, 181, 34, 177, 44,
		125, 70, 250, 13, 55, 7, 37, 88, 155, 247, 188, 183, 54, 124, 10, 124, 99, 143, 82, 131, 8, 78, 33, 132, 66, 193,
		234, 173, 113, 23, 189, 51, 56, 7, 244, 106, 13, 119, 23, 148, 184, 172, 165, 28,
	},
}

// CreateNodePayload creates a dummy payload for node [keyNumber 0-4]
func CreateNodePayload(useInternet bool, keyNumber int) []byte {
	requestPayload := flatbuffers.NewBuilder(0)
	privateKeyFlatbuff := requestPayload.CreateByteVector(privateKeys[keyNumber])
	model.CreateNodeRequestStart(requestPayload)
	model.CreateNodeRequestAddPrivateKey(requestPayload, privateKeyFlatbuff)
	model.CreateNodeRequestAddUseInternet(requestPayload, useInternet)
	requestPayload.Finish(model.CreateNodeRequestEnd(requestPayload))
	return requestPayload.FinishedBytes()
}

// CreateNodeExpectedResponse creates a dummy response for node [keyNumber 0-4]
func CreateNodeExpectedResponse() []byte {
	return utils.BoolResponse(true, nil)
}

// StopNodeExpectedResponse creates a dummy response for node [keyNumber 0-4]
func StopNodeExpectedResponse() []byte {
	return utils.BoolResponse(true, nil)
}

// GetNodeIdExpectedResponse creates a dummy response for node [keyNumber 0-4]
func GetNodeIdExpectedResponse(keyNumber int) []byte {
	return utils.StringResponse(getNodeIdFromPrivateKey(keyNumber).String(), nil)
}

// NotRunningExpectedResponse creates a dummy response for node [keyNumber 0-4]
func NotRunningExpectedResponse() []byte {
	return utils.StringResponse("", errors.New("node does not exist"))
}

func getNodeIdFromPrivateKey(privateKeyIndex int) peer.ID {
	key, err := crypto.UnmarshalPrivateKey(privateKeys[privateKeyIndex])
	if err != nil {
		panic(err)
	}
	peerId, err := peer.IDFromPrivateKey(key)
	if err != nil {
		panic(err)
	}
	return peerId
}
