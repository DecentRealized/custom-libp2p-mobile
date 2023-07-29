package transfer

import "C"
import (
	"bytes"
	"encoding/binary"
	"github.com/DecentRealized/custom-libp2p-mobile/custom-libp2p/utils"
	"unsafe"
)

const (
	StringMessage           = uint8(1)
	FileNotificationMessage = uint8(2)
)

const headerSize = uint16(unsafe.Sizeof(uint8(0)) + unsafe.Sizeof(uint16(0)))

type Message struct {
	msgType    uint8
	msgRawData []byte
}

// marshal converts the message to cbor
func (m *Message) marshal() ([]byte, error) {
	var buff bytes.Buffer
	err := binary.Write(&buff, binary.LittleEndian, m.msgType)
	if err != nil {
		return nil, err
	}
	err = utils.MarshalAndWriteBytes(m.msgRawData, &buff)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// unmarshal converts the message from cbor
func (m *Message) unmarshal(data []byte) error {
	reader := bytes.NewReader(data)
	err := binary.Read(reader, binary.LittleEndian, &m.msgType)
	if err != nil {
		return err
	}
	m.msgRawData, err = utils.UnMarshalBytesFromReader(reader)
	return err
}

// NewStringMessage creates a new string message
func NewStringMessage(message string) (*Message, error) {
	if len(message) == 0 {
		return nil, MessageSizeZero
	}
	if len(message) > int(maxMessageSize) {
		return nil, MessageSizeTooLarge
	}
	return &Message{
		msgType:    StringMessage,
		msgRawData: []byte(message),
	}, nil
}

// NewFileNotificationMessage creates a new file notification message
func NewFileNotificationMessage(fileMetadata *FileMetadata) (*Message, error) {
	fmd := &FileMetadata{}
	marshal, err := fmd.marshal()
	if err != nil {
		return nil, err
	}
	fmdum := &FileMetadata{}
	err = fmdum.unmarshal(marshal)

	msgRawData, err := fileMetadata.Clone().marshal()
	if err != nil {
		return nil, err
	}
	if len(msgRawData) > int(maxMessageSize) {
		return nil, MessageSizeTooLarge
	}
	return &Message{
		msgType:    FileNotificationMessage,
		msgRawData: msgRawData,
	}, nil
}
