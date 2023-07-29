// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package model

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type SendMessageRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsSendMessageRequest(buf []byte, offset flatbuffers.UOffsetT) *SendMessageRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &SendMessageRequest{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsSendMessageRequest(buf []byte, offset flatbuffers.UOffsetT) *SendMessageRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &SendMessageRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *SendMessageRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *SendMessageRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *SendMessageRequest) Message() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *SendMessageRequest) PeerId() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func SendMessageRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func SendMessageRequestAddMessage(builder *flatbuffers.Builder, message flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(message), 0)
}
func SendMessageRequestAddPeerId(builder *flatbuffers.Builder, peerId flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(peerId), 0)
}
func SendMessageRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
