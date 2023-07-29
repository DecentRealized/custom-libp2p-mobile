package transfer

import (
	"bytes"
	"encoding/binary"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/utils"
	"os"
	"path/filepath"
)

type DownloadStatus struct {
	Downloading     bool
	FileBasePath    string
	FileName        string
	FileSHA256      string
	FullFileSize    uint64
	CurrentFileSize uint64
}

// marshal convert to bytes
func (ds *DownloadStatus) marshal() ([]byte, error) {
	var buff bytes.Buffer
	err := binary.Write(&buff, binary.LittleEndian, ds.Downloading)
	if err != nil {
		return nil, err
	}
	err = utils.MarshalAndWriteString(ds.FileBasePath, &buff)
	if err != nil {
		return nil, err
	}
	err = utils.MarshalAndWriteString(ds.FileName, &buff)
	if err != nil {
		return nil, err
	}
	err = utils.MarshalAndWriteString(ds.FileSHA256, &buff)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buff, binary.LittleEndian, ds.FullFileSize)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buff, binary.LittleEndian, ds.CurrentFileSize)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// unmarshal convert from bytes
func (ds *DownloadStatus) unmarshal(data []byte) error {
	var err error
	reader := bytes.NewReader(data)
	err = binary.Read(reader, binary.LittleEndian, &ds.Downloading)
	if err != nil {
		return err
	}
	ds.FileBasePath, err = utils.UnMarshalStringFromReader(reader)
	if err != nil {
		return err
	}
	ds.FileName, err = utils.UnMarshalStringFromReader(reader)
	if err != nil {
		return err
	}
	ds.FileSHA256, err = utils.UnMarshalStringFromReader(reader)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &ds.FullFileSize)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &ds.CurrentFileSize)
	return err
}

func downloadStatusFromDownloadingAndMetadata(downloading bool, metadata *FileMetadata) (DownloadStatus, error) {
	stat, err := os.Stat(filepath.Join(metadata._basePath, metadata.fileName))
	if err != nil {
		return DownloadStatus{}, err
	}
	return DownloadStatus{
		Downloading:     downloading,
		FileBasePath:    metadata._basePath,
		FileName:        metadata.fileName,
		FileSHA256:      metadata.fileSHA256,
		FullFileSize:    metadata.fileSize,
		CurrentFileSize: uint64(stat.Size()),
	}, nil
}
