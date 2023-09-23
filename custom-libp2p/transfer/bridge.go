package transfer

import (
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/models"
	"github.com/libp2p/go-libp2p/core/peer"
	"google.golang.org/protobuf/proto"
)

type ServeFileBridgeInput = models.ServeFileRequest
type ServeFileBridgeOutput = models.StringMessage

func ServeFileBridge(input proto.Message) (proto.Message, error) {
	serveFileRequest := input.(*ServeFileBridgeInput)
	peerId, err := peer.Decode(serveFileRequest.GetPeerId())
	if err != nil {
		return nil, err
	}
	fileSha256, err := ServeFile(serveFileRequest.GetFilePath(), peerId)
	if err != nil {
		return nil, err
	}
	return &ServeFileBridgeOutput{Message: fileSha256}, nil
}

type StopServingFileBridgeInput = models.StringMessage

func StopServingFileBridge(input proto.Message) (proto.Message, error) {
	stopServingFileRequest := input.(*StopServingFileBridgeInput)
	err := StopServingFile(stopServingFileRequest.GetMessage())
	return nil, err
}

type SendMessageBridgeInput = models.SendMessageRequest

func SendMessageBridge(input proto.Message) (proto.Message, error) {
	sendMessageRequest := input.(*SendMessageBridgeInput)
	peerId, err := peer.Decode(sendMessageRequest.GetPeerId())
	if err != nil {
		return nil, err
	}
	err = SendMessage(peerId, sendMessageRequest.GetMessage())
	return nil, err
}

type PauseDownloadBridgeInput = models.DownloadRequest

func PauseDownloadBridge(input proto.Message) (proto.Message, error) {
	pauseDownloadRequest := input.(*PauseDownloadBridgeInput)
	peerId, err := peer.Decode(pauseDownloadRequest.GetPeerId())
	if err != nil {
		return nil, err
	}
	err = PauseDownload(pauseDownloadRequest.GetFileSha256(), peerId)
	return nil, err
}

type ResumeDownloadBridgeInput = models.DownloadRequest

func ResumeDownloadBridge(input proto.Message) (proto.Message, error) {
	resumeDownloadRequest := input.(*ResumeDownloadBridgeInput)
	peerId, err := peer.Decode(resumeDownloadRequest.GetPeerId())
	if err != nil {
		return nil, err
	}
	err = ResumeDownload(resumeDownloadRequest.GetFileSha256(), peerId)
	return nil, err
}

type StopDownloadBridgeInput = models.DownloadRequest

func StopDownloadBridge(input proto.Message) (proto.Message, error) {
	stopDownloadRequest := input.(*StopDownloadBridgeInput)
	peerId, err := peer.Decode(stopDownloadRequest.GetPeerId())
	if err != nil {
		return nil, err
	}
	err = StopDownload(stopDownloadRequest.GetFileSha256(), peerId)
	return nil, err
}

type GetDownloadStatusBridgeInput = models.DownloadRequest
type GetDownloadStatusBridgeOutput = models.DownloadStatus

func GetDownloadStatusBridge(input proto.Message) (proto.Message, error) {
	getDownloadStatusRequest := input.(*GetDownloadStatusBridgeInput)
	peerId, err := peer.Decode(getDownloadStatusRequest.GetPeerId())
	if err != nil {
		return nil, err
	}
	status, err := GetDownloadStatus(getDownloadStatusRequest.GetFileSha256(), peerId)
	return status, err
}
