// automatically generated by the FlatBuffers compiler, do not modify

import * as flatbuffers from 'flatbuffers';

export class CreateNodeRequest {
  bb: flatbuffers.ByteBuffer|null = null;
  bb_pos = 0;
__init(i:number, bb:flatbuffers.ByteBuffer):CreateNodeRequest {
  this.bb_pos = i;
  this.bb = bb;
  return this;
}

static getRootAsCreateNodeRequest(bb:flatbuffers.ByteBuffer, obj?:CreateNodeRequest):CreateNodeRequest {
  return (obj || new CreateNodeRequest()).__init(bb.readInt32(bb.position()) + bb.position(), bb);
}

static getSizePrefixedRootAsCreateNodeRequest(bb:flatbuffers.ByteBuffer, obj?:CreateNodeRequest):CreateNodeRequest {
  bb.setPosition(bb.position() + flatbuffers.SIZE_PREFIX_LENGTH);
  return (obj || new CreateNodeRequest()).__init(bb.readInt32(bb.position()) + bb.position(), bb);
}

privateKey(index: number):number|null {
  const offset = this.bb!.__offset(this.bb_pos, 4);
  return offset ? this.bb!.readUint8(this.bb!.__vector(this.bb_pos + offset) + index) : 0;
}

privateKeyLength():number {
  const offset = this.bb!.__offset(this.bb_pos, 4);
  return offset ? this.bb!.__vector_len(this.bb_pos + offset) : 0;
}

privateKeyArray():Uint8Array|null {
  const offset = this.bb!.__offset(this.bb_pos, 4);
  return offset ? new Uint8Array(this.bb!.bytes().buffer, this.bb!.bytes().byteOffset + this.bb!.__vector(this.bb_pos + offset), this.bb!.__vector_len(this.bb_pos + offset)) : null;
}

useInternet():boolean {
  const offset = this.bb!.__offset(this.bb_pos, 6);
  return offset ? !!this.bb!.readInt8(this.bb_pos + offset) : false;
}

mutate_use_internet(value:boolean):boolean {
  const offset = this.bb!.__offset(this.bb_pos, 6);

  if (offset === 0) {
    return false;
  }

  this.bb!.writeInt8(this.bb_pos + offset, +value);
  return true;
}

static startCreateNodeRequest(builder:flatbuffers.Builder) {
  builder.startObject(2);
}

static addPrivateKey(builder:flatbuffers.Builder, privateKeyOffset:flatbuffers.Offset) {
  builder.addFieldOffset(0, privateKeyOffset, 0);
}

static createPrivateKeyVector(builder:flatbuffers.Builder, data:number[]|Uint8Array):flatbuffers.Offset {
  builder.startVector(1, data.length, 1);
  for (let i = data.length - 1; i >= 0; i--) {
    builder.addInt8(data[i]!);
  }
  return builder.endVector();
}

static startPrivateKeyVector(builder:flatbuffers.Builder, numElems:number) {
  builder.startVector(1, numElems, 1);
}

static addUseInternet(builder:flatbuffers.Builder, useInternet:boolean) {
  builder.addFieldInt8(1, +useInternet, +false);
}

static endCreateNodeRequest(builder:flatbuffers.Builder):flatbuffers.Offset {
  const offset = builder.endObject();
  return offset;
}

static createCreateNodeRequest(builder:flatbuffers.Builder, privateKeyOffset:flatbuffers.Offset, useInternet:boolean):flatbuffers.Offset {
  CreateNodeRequest.startCreateNodeRequest(builder);
  CreateNodeRequest.addPrivateKey(builder, privateKeyOffset);
  CreateNodeRequest.addUseInternet(builder, useInternet);
  return CreateNodeRequest.endCreateNodeRequest(builder);
}
}