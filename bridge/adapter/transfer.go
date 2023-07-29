package adapter

import (
	utils "github.com/DecentRealized/custom-libp2p-mobile/bridge/utils"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/transfer"
)

func ServeFile(payload []byte) []byte {
	filePath, peerId, err := utils.ServeFileRequest(payload)
	if err != nil {
		return utils.StringResponse("", err)
	}

	fileSHA256, err := transfer.ServeFile(filePath, peerId)
	return utils.StringResponse(fileSHA256, err)
}

func StopServingFile(payload []byte) []byte {
	sha256Sum, err := utils.StringRequest(payload)
	if err != nil {
		return utils.StringResponse("", err)
	}

	err = transfer.StopServingFile(sha256Sum)
	return utils.BoolResponse(err == nil, err)
}

func SendMessage(payload []byte) []byte {
	message, peerId, err := utils.SendMessageRequest(payload)
	if err != nil {
		return utils.BoolResponse(false, err)
	}

	err = transfer.SendMessage(peerId, message)
	return utils.BoolResponse(err == nil, err)
}

func PauseDownload(payload []byte) []byte {
	sha256Sum, peerId, err := utils.DownloadRequest(payload)
	if err != nil {
		return utils.BoolResponse(false, err)
	}

	err = transfer.PauseDownload(sha256Sum, peerId)
	return utils.BoolResponse(err == nil, err)
}

func ResumeDownload(payload []byte) []byte {
	sha256Sum, peerId, err := utils.DownloadRequest(payload)
	if err != nil {
		return utils.BoolResponse(false, err)
	}

	err = transfer.ResumeDownload(sha256Sum, peerId)
	return utils.BoolResponse(err == nil, err)
}

func StopDownload(payload []byte) []byte {
	sha256Sum, peerId, err := utils.DownloadRequest(payload)
	if err != nil {
		return utils.BoolResponse(false, err)
	}

	err = transfer.StopDownload(sha256Sum, peerId)
	return utils.BoolResponse(err == nil, err)
}

func GetDownloadStatus(payload []byte) []byte {
	sha256Sum, peerId, err := utils.DownloadRequest(payload)
	if err != nil {
		return utils.DownloadStatusResponse(&transfer.DownloadStatus{}, err)
	}

	downloadStatus, err := transfer.GetDownloadStatus(sha256Sum, peerId)
	return utils.DownloadStatusResponse(&downloadStatus, err)
}
