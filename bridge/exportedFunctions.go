package custom_libp2p_bridge

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/crypto"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/example"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/p2p"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/transfer"
	"google.golang.org/protobuf/proto"
)

type BridgeInfo struct {
	input    proto.Message                              // Used for inferring input type for decoding bytes
	function func(proto.Message) (proto.Message, error) // Call This bridge func
	output   proto.Message                              // Used for inferring output type for generating dart utils
}

var bridgeMapping = map[string]BridgeInfo{ // Maps flutter name to golang BridgeInfo
	"getHelloMessage": {
		input:    &example.GetHelloMessageBridgeInput{},
		function: example.GetHelloMessageBridge,
		output:   &example.GetHelloMessageBridgeOutput{},
	},
	"createKeyPair": {
		function: crypto.CreateKeyPairBridge,
		output:   &crypto.CreateKeyPairBridgeOutput{},
	},
	"createNode": {
		input:    &p2p.CreateNodeBridgeInput{},
		function: p2p.CreateNodeBridge,
	},
	"stopNode": {
		function: p2p.StopNodeBridge,
	},
	"getNodeId": {
		function: p2p.GetNodeIdBridge,
		output:   &p2p.GetNodeIdBridgeOutput{},
	},
	"getListenAddresses": {
		function: p2p.GetListenAddressesBridge,
		output:   &p2p.GetListenAddressBridgeOutput{},
	},
	//"allowNode":         {},
	//"isAllowedNode":     {},
	//"getAllowedNodes":   {},
	//"denyNode":          {},
	"serveFile": {
		input:    &transfer.ServeFileBridgeInput{},
		function: transfer.ServeFileBridge,
		output:   &transfer.ServeFileBridgeOutput{},
	},
	"stopServingFile": {
		input:    &transfer.StopServingFileBridgeInput{},
		function: transfer.StopServingFileBridge,
	},
	"sendMessage": {
		input:    &transfer.SendMessageBridgeInput{},
		function: transfer.SendMessageBridge,
	},
	"pauseDownload": {
		input:    &transfer.PauseDownloadBridgeInput{},
		function: transfer.PauseDownloadBridge,
	},
	"resumeDownload": {
		input:    &transfer.ResumeDownloadBridgeInput{},
		function: transfer.ResumeDownloadBridge,
	},
	"stopDownload": {
		input:    &transfer.StopDownloadBridgeInput{},
		function: transfer.StopDownloadBridge,
	},
	"getDownloadStatus": {
		input:    &transfer.GetDownloadStatusBridgeInput{},
		function: transfer.GetDownloadStatusBridge,
		output:   &transfer.GetDownloadStatusBridgeOutput{},
	},
	"flushNotifications": {
		function: notifier.FlushNotificationsBridge,
		output:   &notifier.FlushNotificationsBridgeOutput{},
	},
}
