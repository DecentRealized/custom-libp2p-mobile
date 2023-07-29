package adapter

import (
	"github.com/DecentRealized/custom-libp2p-mobile/bridge/dummy"
	"github.com/DecentRealized/custom-libp2p-mobile/bridge/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO: Test IRL, cant test as node on one device is singleton

// TestServeFile Tests before notification, as notification needs more than 2 nodes
func TestServeFile(t *testing.T) {
	TestCreateNodeOnline(t)
	output := ServeFile(dummy.ServeFilePayload())
	response := model.GetRootAsBoolResponse(output, 0)
	assert.Equal(t, dummy.PeerNotFoundResponse(), response.Table().Bytes)
}

// TestServeFileWhenNotRunning When server is not running
func TestServeFileWhenNotRunning(t *testing.T) {
	StopNode()
	output := ServeFile(dummy.ServeFilePayload())
	response := model.GetRootAsBoolResponse(output, 0)
	t.Log(response.Error())
	assert.Equal(t, dummy.TransferNotRunningResponse(), response.Table().Bytes)
}

// TODO
//// TestStopServingFile stop serving file test
//func TestStopServingFile(t *testing.T) {
//	TestCreateNodeOnline(t)
//	output := StopServingFile(dummy.StopServingFilePayload())
//	response := model.GetRootAsBoolResponse(output, 0)
//	assert.Equal(t, dummy.PeerNotFoundResponse(), response.Table().Bytes)
//}
//
//// TestStopServingFileWhenNotServing When file is not being served
//func TestStopServingFileWhenNotServing(t *testing.T) {
//	TestCreateNodeOnline(t)
//	output := StopServingFile(dummy.StopServingFilePayload())
//	response := model.GetRootAsBoolResponse(output, 0)
//	assert.Equal(t, dummy.PeerNotFoundResponse(), response.Table().Bytes)
//}
//
//// TestStopServingFileWhenNotRunning When server is not running
//func TestStopServingFileWhenNotRunning(t *testing.T) {
//	TestCreateNodeOnline(t)
//	output := StopServingFile(dummy.StopServingFilePayload())
//	response := model.GetRootAsBoolResponse(output, 0)
//	assert.Equal(t, dummy.PeerNotFoundResponse(), response.Table().Bytes)
//}

// TestSendMessage
func TestSendMessage(t *testing.T) {
	TestCreateNodeOnline(t)
	output := SendMessage(dummy.SendMessagePayload())
	response := model.GetRootAsBoolResponse(output, 0)
	assert.Equal(t, dummy.PeerNotFoundResponse(), response.Table().Bytes)
}

// TestSendMessageWhenNotRunning
func TestSendMessageWhenNotRunning(t *testing.T) {
	StopNode()
	output := SendMessage(dummy.SendMessagePayload())
	response := model.GetRootAsBoolResponse(output, 0)
	t.Log(response.Error())
	assert.Equal(t, dummy.TransferNotRunningResponse(), response.Table().Bytes)
}

// TODO
//// TestPauseDownload
//func TestPauseDownload(t *testing.T) {
//	TestCreateNodeOnline(t)
//	output := PauseDownload(dummy.PauseDownloadPayload())
//	_ = model.GetRootAsBoolResponse(output, 0)
//	//assert.Equal(t, dummy.PeerNotFoundResponse(), response.Table().Bytes)
//}
//
//func TestPauseDownloadWhenNotRunning(t *testing.T) {
//	StopNode()
//	output := PauseDownload(dummy.SendMessagePayload())
//	response := model.GetRootAsBoolResponse(output, 0)
//	assert.Equal(t, dummy.PeerNotFoundResponse(), response.Table().Bytes)
//}
//
//func TestPauseDownloadWhenNotDownloading(t *testing.T) {
//	TestCreateNodeOnline(t)
//	output := PauseDownload(dummy.PauseDownloadPayload())
//	_ = model.GetRootAsBoolResponse(output, 0)
//	//assert.Equal(t, dummy.PeerNotFoundResponse(), response.Table().Bytes)
//}
//
//func TestResumeDownload(t *testing.T) {
//
//}
//
//func TestResumeDownloadWhenNotRunning(t *testing.T) {
//
//}
//
//func TestResumeDownloadWhenAlreadyDownloading(t *testing.T) {
//
//}
//
//func TestStopDownload(t *testing.T) {
//
//}
//
//func TestStopDownloadWhenNotRunning(t *testing.T) {
//
//}
//
//func TestStopDownloadWhenNotDownloading(t *testing.T) {
//
//}
