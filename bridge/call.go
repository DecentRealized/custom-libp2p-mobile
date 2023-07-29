package custom_libp2p_bridge

import (
	"fmt"
	"github.com/DecentRealized/custom-libp2p-mobile/bridge/adapter"
)

// Call ...
func Call(name string, payload []byte) ([]byte, error) {
	var output []byte
	switch name {
	case "getHelloWorld":
		output = adapter.GetHelloWorld(payload)
	case "createKeyPair":
		output = adapter.CreateKeyPair()
	case "createNode":
		output = adapter.CreateNode(payload)
	case "getNodeId":
		output = adapter.GetNodeId()
	case "getListenAddresses":
		output = adapter.GetListenAddresses()
	case "allowNode":
		output = adapter.AllowNode(payload)
	case "isAllowedNode":
		output = adapter.IsAllowedNode(payload)
	case "getAllowedNodes":
		output = adapter.GetAllowedNodes()
	case "denyNode":
		output = adapter.DenyNode(payload)
	case "serveFile":
		output = adapter.ServeFile(payload)
	case "stopServingFile":
		output = adapter.StopServingFile(payload)
	case "sendMessage":
		output = adapter.SendMessage(payload)
	case "pauseDownload":
		output = adapter.PauseDownload(payload)
	case "resumeDownload":
		output = adapter.ResumeDownload(payload)
	case "stopDownload":
		output = adapter.StopDownload(payload)
	case "getDownloadStatus":
		output = adapter.GetDownloadStatus(payload)
	case "stopNode":
		output = adapter.StopNode()
	default:
		return nil, fmt.Errorf("not implemented: %s", name)
	}

	return output, nil
}
