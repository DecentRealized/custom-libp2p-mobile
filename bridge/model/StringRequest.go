// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package model

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type StringRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsStringRequest(buf []byte, offset flatbuffers.UOffsetT) *StringRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &StringRequest{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsStringRequest(buf []byte, offset flatbuffers.UOffsetT) *StringRequest {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &StringRequest{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *StringRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *StringRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *StringRequest) Input() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func StringRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func StringRequestAddInput(builder *flatbuffers.Builder, input flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(input), 0)
}
func StringRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
