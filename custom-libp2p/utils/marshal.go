package utils

import (
	"encoding/binary"
	"github.com/libp2p/go-libp2p/core/peer"
	"io"
)

// MarshalAndWriteBytes marshals bytes from reader
func MarshalAndWriteBytes(_bytes []byte, writer io.Writer) error {
	err := binary.Write(writer, binary.LittleEndian, uint32(len(_bytes)))
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.LittleEndian, _bytes)
	if err != nil {
		return err
	}
	return nil
}

// UnMarshalBytesFromReader un-marshals bytes from reader
func UnMarshalBytesFromReader(reader io.Reader) ([]byte, error) {
	var lenBytes uint32
	err := binary.Read(reader, binary.LittleEndian, &lenBytes)
	if err != nil {
		return nil, err
	}
	strB := make([]byte, lenBytes)
	err = binary.Read(reader, binary.LittleEndian, &strB)
	if err != nil {
		return nil, err
	}
	return strB, err
}

// MarshalAndWriteString marshals string and writes to writer
func MarshalAndWriteString(str string, writer io.Writer) error {
	strB := []byte(str)
	return MarshalAndWriteBytes(strB, writer)
}

// UnMarshalStringFromReader un-marshals string from reader
func UnMarshalStringFromReader(reader io.Reader) (string, error) {
	_bytes, err := UnMarshalBytesFromReader(reader)
	if err != nil {
		return "", err
	}
	return string(_bytes), nil
}

// MarshalAndWritePeerId marshals and writes peer Id
func MarshalAndWritePeerId(peerId peer.ID, writer io.Writer) error {
	marshalBinary, err := peerId.MarshalBinary()
	if err != nil {
		return err
	}
	return MarshalAndWriteBytes(marshalBinary, writer)
}

// UnMarshalPeerIdFromReader un-marshals peer id from reader
func UnMarshalPeerIdFromReader(reader io.Reader) (peer.ID, error) {
	_bytes, err := UnMarshalBytesFromReader(reader)
	if err != nil {
		return "", err
	}
	if len(_bytes) == 0 {
		return "", nil
	}
	var peerId peer.ID
	err = peerId.Unmarshal(_bytes)
	if err != nil {
		return "", err
	}
	return peerId, nil
}

// MarshalAndWriteACL marshals acl map
func MarshalAndWriteACL(acl map[peer.ID]bool, writer io.Writer) error {
	err := binary.Write(writer, binary.LittleEndian, uint32(len(acl)))
	if err != nil {
		return err
	}
	for k, v := range acl {
		err = MarshalAndWritePeerId(k, writer)
		err = binary.Write(writer, binary.LittleEndian, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// UnMarshalACLFromReader un-marshals acl from reader
func UnMarshalACLFromReader(reader io.Reader) (map[peer.ID]bool, error) {
	acl := make(map[peer.ID]bool)
	var aclLen uint32
	err := binary.Read(reader, binary.LittleEndian, &aclLen)
	if err != nil {
		return nil, err
	}
	for i := 0; uint32(i) < aclLen; i++ {
		var value bool
		peerId, err := UnMarshalPeerIdFromReader(reader)
		if err != nil {
			return nil, err
		}
		err = binary.Read(reader, binary.LittleEndian, &value)
		if err != nil {
			return nil, err
		}
		acl[peerId] = value
	}
	return acl, nil
}
