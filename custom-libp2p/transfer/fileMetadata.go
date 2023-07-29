package transfer

import (
	"bytes"
	"encoding/binary"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/utils"
	"github.com/libp2p/go-libp2p/core/peer"
	"os"
	"path/filepath"
)

// TODO: Define FileTypes for use in frontend
const (
	ImageFileType = 0
	VideoFileType = 1
	AudioFileType = 2
	TextFileType  = 2
)

type FileMetadata struct {
	_basePath            string
	_fileDescriptor      *os.File
	_authorizedAccessors map[peer.ID]bool
	_fileServer          peer.ID
	fileName             string
	fileSHA256           string
	fileSize             uint64
	//fileType             int
}

// marshal Convert to bytes
func (fm *FileMetadata) marshal() ([]byte, error) {
	var buff bytes.Buffer
	err := utils.MarshalAndWriteString(fm._basePath, &buff)
	if err != nil {
		return nil, err
	}
	err = utils.MarshalAndWriteACL(fm._authorizedAccessors, &buff)
	if err != nil {
		return nil, err
	}
	err = utils.MarshalAndWritePeerId(fm._fileServer, &buff)
	if err != nil {
		return nil, err
	}
	err = utils.MarshalAndWriteString(fm.fileName, &buff)
	if err != nil {
		return nil, err
	}
	err = utils.MarshalAndWriteString(fm.fileSHA256, &buff)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buff, binary.LittleEndian, fm.fileSize)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// unmarshal Convert from bytes
func (fm *FileMetadata) unmarshal(data []byte) error {
	var err error
	reader := bytes.NewReader(data)
	fm._basePath, err = utils.UnMarshalStringFromReader(reader)
	if err != nil {
		return err
	}
	fm._authorizedAccessors, err = utils.UnMarshalACLFromReader(reader)
	if err != nil {
		return err
	}
	fm._fileServer, err = utils.UnMarshalPeerIdFromReader(reader)
	if err != nil {
		return err
	}
	fm.fileName, err = utils.UnMarshalStringFromReader(reader)
	if err != nil {
		return err
	}
	fm.fileSHA256, err = utils.UnMarshalStringFromReader(reader)
	if err != nil {
		return err
	}
	err = binary.Read(reader, binary.LittleEndian, &fm.fileSize)
	return err
}

// AllowPeerAccess Allow peer to access file
func (fm *FileMetadata) AllowPeerAccess(peerId peer.ID) {
	fm._authorizedAccessors[peerId] = true
}

// DenyPeerAccess Deny peer from accessing file
func (fm *FileMetadata) DenyPeerAccess(peerId peer.ID) {
	delete(fm._authorizedAccessors, peerId)
}

// IsAllowedAccess Check if peer is allowed to access file
func (fm *FileMetadata) IsAllowedAccess(peerId peer.ID) bool {
	val, ok := fm._authorizedAccessors[peerId]
	return val && ok
}

// GetFile open file if not opened, return file
func (fm *FileMetadata) GetFile() (*os.File, error) {
	if fm._fileDescriptor != nil {
		return fm._fileDescriptor, nil
	}
	file, err := os.OpenFile(filepath.Join(fm._basePath, fm.fileName), os.O_RDWR|os.O_CREATE, 0644)
	fm._fileDescriptor = file
	return file, err
}

// CloseFile close file if open
func (fm *FileMetadata) CloseFile() error {
	if fm._fileDescriptor == nil {
		return nil
	}
	err := fm._fileDescriptor.Close()
	fm._fileDescriptor = nil
	return err
}

// Clone clones struct without fields starting with '_'
func (fm *FileMetadata) Clone() *FileMetadata {
	return &FileMetadata{
		//fileType:   fm.fileType,
		fileName:   fm.fileName,
		fileSHA256: fm.fileSHA256,
		fileSize:   fm.fileSize,
	}
}
