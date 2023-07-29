package butils

import (
	"github.com/DecentRealized/custom-libp2p-mobile/bridge/model"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/transfer"
	flatbuffers "github.com/google/flatbuffers/go"
)

func BoolResponse(output bool, err error) []byte {
	response := flatbuffers.NewBuilder(0)
	if err != nil {
		outputOffset := response.CreateString(err.Error())
		model.BoolResponseStart(response)
		model.BoolResponseAddError(response, outputOffset)
		response.Finish(model.BoolResponseEnd(response))
		return response.FinishedBytes()
	}
	model.BoolResponseStart(response)
	model.BoolResponseAddOutput(response, output)
	response.Finish(model.BoolResponseEnd(response))
	return response.FinishedBytes()
}

func StringResponse(output string, err error) []byte {
	response := flatbuffers.NewBuilder(0)
	if err != nil {
		outputOffset := response.CreateString(err.Error())
		model.StringResponseStart(response)
		model.StringResponseAddError(response, outputOffset)
		response.Finish(model.StringResponseEnd(response))
		return response.FinishedBytes()
	}
	outputOffset := response.CreateString(output)
	model.StringResponseStart(response)
	model.StringResponseAddOutput(response, outputOffset)
	response.Finish(model.StringResponseEnd(response))
	return response.FinishedBytes()
}

func StringArrayResponse(output []string, err error) []byte {
	response := flatbuffers.NewBuilder(0)
	if err != nil {
		outputOffset := response.CreateString(err.Error())
		model.StringArrayResponseStart(response)
		model.StringResponseAddError(response, outputOffset)
		response.Finish(model.StringResponseEnd(response))
		return response.FinishedBytes()
	}

	outputOffsets := make([]flatbuffers.UOffsetT, len(output))
	for i := range output {
		outputOffset := response.CreateString(output[i])
		outputOffsets[i] = outputOffset
	}

	model.StringArrayResponseStartOutputVector(response, len(output))
	for _, outputOffset := range outputOffsets {
		response.PrependUOffsetT(outputOffset)
	}
	outputVec := response.EndVector(len(output))

	model.StringArrayResponseStart(response)
	model.StringArrayResponseAddOutput(response, outputVec)

	response.Finish(model.StringResponseEnd(response))
	return response.FinishedBytes()
}

func KeyPairResponse(privKey []byte, pubKey []byte, err error) []byte {
	response := flatbuffers.NewBuilder(0)
	if err != nil {
		outputOffset := response.CreateString(err.Error())
		model.KeyPairResponseStart(response)
		model.KeyPairResponseAddError(response, outputOffset)
		response.Finish(model.KeyPairResponseEnd(response))
		return response.FinishedBytes()
	}

	publicKeyOffset := response.CreateByteVector(pubKey)
	privateKeyOffset := response.CreateByteVector(privKey)

	model.KeyPairStart(response)
	model.KeyPairAddPublicKey(response, publicKeyOffset)
	model.KeyPairAddPrivateKey(response, privateKeyOffset)
	KeyPair := model.KeyPairEnd(response)

	model.KeyPairResponseStart(response)
	model.KeyPairResponseAddOutput(response, KeyPair)
	response.Finish(model.KeyPairResponseEnd(response))
	return response.FinishedBytes()
}

func DownloadStatusResponse(status *transfer.DownloadStatus, err error) []byte {
	response := flatbuffers.NewBuilder(0)
	if err != nil {
		outputOffset := response.CreateString(err.Error())
		model.DownloadStatusResponseStart(response)
		model.DownloadStatusResponseAddError(response, outputOffset)
		response.Finish(model.DownloadStatusResponseEnd(response))
		return response.FinishedBytes()
	}

	fileBasePathOffset := response.CreateString(status.FileBasePath)
	fileNameOffset := response.CreateString(status.FileName)
	fileSHA256Sum := response.CreateString(status.FileSHA256)

	model.DownloadStatusStart(response)
	model.DownloadStatusAddDownloading(response, status.Downloading)
	model.DownloadStatusAddFileBasePath(response, fileBasePathOffset)
	model.DownloadStatusAddFileName(response, fileNameOffset)
	model.DownloadStatusAddFileSha256(response, fileSHA256Sum)
	model.DownloadStatusAddCurrentFileSize(response, status.CurrentFileSize)
	model.DownloadStatusAddFullFileSize(response, status.FullFileSize)
	response.Finish(model.DownloadStatusResponseEnd(response))
	return response.FinishedBytes()
}
