// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package model

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type DownloadRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsDownloadRequest(buf []byte, offset flatbuffers.UOffsetT) *DownloadRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &DownloadRequest{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsDownloadRequest(buf []byte, offset flatbuffers.UOffsetT) *DownloadRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &DownloadRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *DownloadRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *DownloadRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *DownloadRequest) FileSha256() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *DownloadRequest) PeerId() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func DownloadRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func DownloadRequestAddFileSha256(builder *flatbuffers.Builder, fileSha256 flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(fileSha256), 0)
}
func DownloadRequestAddPeerId(builder *flatbuffers.Builder, peerId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(peerId), 0)
}
func DownloadRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
