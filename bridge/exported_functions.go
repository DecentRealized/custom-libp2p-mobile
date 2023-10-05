package custom_libp2p_bridge

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/access_manager"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/crypto"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/database"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/example"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/file_handler"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/notifier"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/p2p"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/transfer"
	"google.golang.org/protobuf/proto"
)

type BridgeInfo struct {
	preCallWrapper  string                                     // Insert this before auto generated logic
	input           proto.Message                              // Used for inferring input type for decoding bytes
	function        func(proto.Message) (proto.Message, error) // Call This bridge func
	output          proto.Message                              // Used for inferring output type for generating dart utils
	postCallWrapper string                                     // Insert this after auto generated logic
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
		// Set Downloads and database directories
		preCallWrapper: "" +
			"    Directory? _defaultDownloadDir = await getDownloadsDirectory();\n" +
			"    if (_defaultDownloadDir == null) {\n" +
			"      _defaultDownloadDir = await getApplicationDocumentsDirectory();\n" +
			"    }\n" +
			"    var _tmp = models.StringMessage();\n" +
			"    _tmp.message = _defaultDownloadDir.path;\n" +
			"    await _setDownloadPath(_tmp);\n" +
			"    Directory _defaultDbBaseDir = await getApplicationSupportDirectory();\n" +
			"    var _defaultDbDirPath = p.join(_defaultDbBaseDir.path, '.db');\n" +
			"    Directory _defaultDbDir = await Directory(_defaultDbDirPath).create();\n" +
			"    _tmp.message = _defaultDbDir.path;\n" +
			"    await _setDatabaseDirectory(_tmp);",
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
	"connectToPeer": {
		input:    &p2p.ConnectToPeerBridgeInput{},
		function: p2p.ConnectToPeerBridge,
	},
	"checkConnectionStatus": {
		input:    &p2p.CheckConnectionStatusBridgeInput{},
		function: p2p.CheckConnectionStatusBridge,
		output:   &p2p.CheckConnectionStatusBridgeOutput{},
	},
	"allowNode": {
		input:    &access_manager.AllowNodeBridgeInput{},
		function: access_manager.AllowNodeBridge,
	},
	"isAllowedNode": {
		input:    &access_manager.IsAllowedNodeBridgeInput{},
		function: access_manager.IsAllowedNodeBridge,
		output:   &access_manager.IsAllowedNodeBridgeOutput{},
	},
	"getAllowedNodes": {
		function: access_manager.GetAllowedNodesBridge,
		output:   &access_manager.GetAllowedNodesBridgeOutput{},
	},
	"blockNode": {
		input:    &access_manager.BlockNodeBridgeInput{},
		function: access_manager.BlockNodeBridge,
	},
	"isBlockedNode": {
		input:    &access_manager.IsBlockedNodeBridgeInput{},
		function: access_manager.IsBlockedNodeBridge,
		output:   &access_manager.IsBlockedNodeBridgeOutput{},
	},
	"getBlockedNodes": {
		function: access_manager.GetBlockedNodesBridge,
		output:   &access_manager.GetBlockedNodesBridgeOutput{},
	},
	"_setDownloadPath": {
		input:    &file_handler.SetDownloadPathBridgeInput{},
		function: file_handler.SetDownloadPathBridge,
	},
	"getDownloadPath": {
		function: file_handler.GetDownloadPathBridge,
		output:   &file_handler.GetDownloadPathBridgeOutput{},
	},
	"_setDatabaseDirectory": {
		input:    &database.SetDatabaseDirectoryBridgeInput{},
		function: database.SetDatabaseDirectoryBridge,
	},
	"getDatabaseDirectory": {
		function: database.GetDatabaseDirectoryBridge,
		output:   &database.GetDatabaseDirectoryBridgeOutput{},
	},
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
	"getDownloadingFiles": {
		function: transfer.GetDownloadingFilesBridge,
		output:   &transfer.GetDownloadingFilesBridgeOutput{},
	},
	"getUploadingFiles": {
		function: transfer.GetUploadingFilesBridge,
		output:   &transfer.GetUploadingFilesBridgeOutput{},
	},
	"getDownloadStatus": {
		input:    &transfer.GetDownloadStatusBridgeInput{},
		function: transfer.GetDownloadStatusBridge,
		output:   &transfer.GetDownloadStatusBridgeOutput{},
	},
	"_flushNotifications": {
		function: notifier.FlushNotificationsBridge,
		output:   &notifier.FlushNotificationsBridgeOutput{},
	},
}
