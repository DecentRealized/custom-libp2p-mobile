// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package model

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type GetHelloMessageRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsGetHelloMessageRequest(buf []byte, offset flatbuffers.UOffsetT) *GetHelloMessageRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &GetHelloMessageRequest{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsGetHelloMessageRequest(buf []byte, offset flatbuffers.UOffsetT) *GetHelloMessageRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &GetHelloMessageRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *GetHelloMessageRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *GetHelloMessageRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *GetHelloMessageRequest) UserName() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func GetHelloMessageRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func GetHelloMessageRequestAddUserName(builder *flatbuffers.Builder, userName flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(userName), 0)
}
func GetHelloMessageRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
