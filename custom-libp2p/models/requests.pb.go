// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.12.4
// source: requests.proto

package models

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CreateNodeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UseInternet bool   `protobuf:"varint,1,opt,name=use_internet,json=useInternet,proto3" json:"use_internet,omitempty"`
	PrivateKey  []byte `protobuf:"bytes,2,opt,name=private_key,json=privateKey,proto3" json:"private_key,omitempty"`
}

func (x *CreateNodeRequest) Reset() {
	*x = CreateNodeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_requests_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateNodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateNodeRequest) ProtoMessage() {}

func (x *CreateNodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_requests_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateNodeRequest.ProtoReflect.Descriptor instead.
func (*CreateNodeRequest) Descriptor() ([]byte, []int) {
	return file_requests_proto_rawDescGZIP(), []int{0}
}

func (x *CreateNodeRequest) GetUseInternet() bool {
	if x != nil {
		return x.UseInternet
	}
	return false
}

func (x *CreateNodeRequest) GetPrivateKey() []byte {
	if x != nil {
		return x.PrivateKey
	}
	return nil
}

type ServeFileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FilePath string `protobuf:"bytes,1,opt,name=file_path,json=filePath,proto3" json:"file_path,omitempty"`
	PeerId   string `protobuf:"bytes,2,opt,name=peer_id,json=peerId,proto3" json:"peer_id,omitempty"`
}

func (x *ServeFileRequest) Reset() {
	*x = ServeFileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_requests_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServeFileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServeFileRequest) ProtoMessage() {}

func (x *ServeFileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_requests_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServeFileRequest.ProtoReflect.Descriptor instead.
func (*ServeFileRequest) Descriptor() ([]byte, []int) {
	return file_requests_proto_rawDescGZIP(), []int{1}
}

func (x *ServeFileRequest) GetFilePath() string {
	if x != nil {
		return x.FilePath
	}
	return ""
}

func (x *ServeFileRequest) GetPeerId() string {
	if x != nil {
		return x.PeerId
	}
	return ""
}

type SendMessageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PeerId  string `protobuf:"bytes,1,opt,name=peer_id,json=peerId,proto3" json:"peer_id,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *SendMessageRequest) Reset() {
	*x = SendMessageRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_requests_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessageRequest) ProtoMessage() {}

func (x *SendMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_requests_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMessageRequest.ProtoReflect.Descriptor instead.
func (*SendMessageRequest) Descriptor() ([]byte, []int) {
	return file_requests_proto_rawDescGZIP(), []int{2}
}

func (x *SendMessageRequest) GetPeerId() string {
	if x != nil {
		return x.PeerId
	}
	return ""
}

func (x *SendMessageRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type DownloadRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PeerId     string `protobuf:"bytes,1,opt,name=peer_id,json=peerId,proto3" json:"peer_id,omitempty"`
	FileSha256 string `protobuf:"bytes,2,opt,name=file_sha256,json=fileSha256,proto3" json:"file_sha256,omitempty"`
}

func (x *DownloadRequest) Reset() {
	*x = DownloadRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_requests_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadRequest) ProtoMessage() {}

func (x *DownloadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_requests_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadRequest.ProtoReflect.Descriptor instead.
func (*DownloadRequest) Descriptor() ([]byte, []int) {
	return file_requests_proto_rawDescGZIP(), []int{3}
}

func (x *DownloadRequest) GetPeerId() string {
	if x != nil {
		return x.PeerId
	}
	return ""
}

func (x *DownloadRequest) GetFileSha256() string {
	if x != nil {
		return x.FileSha256
	}
	return ""
}

var File_requests_proto protoreflect.FileDescriptor

var file_requests_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x06, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x22, 0x57, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a,
	0x0c, 0x75, 0x73, 0x65, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0b, 0x75, 0x73, 0x65, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74,
	0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65,
	0x79, 0x22, 0x48, 0x0a, 0x10, 0x53, 0x65, 0x72, 0x76, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x70, 0x61,
	0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x61,
	0x74, 0x68, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x65, 0x65, 0x72, 0x49, 0x64, 0x22, 0x47, 0x0a, 0x12, 0x53,
	0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x70, 0x65, 0x65, 0x72, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x4b, 0x0a, 0x0f, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x70, 0x65, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x65, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x1f, 0x0a, 0x0b, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x73, 0x68, 0x61, 0x32, 0x35, 0x36, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x68, 0x61, 0x32, 0x35,
	0x36, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_requests_proto_rawDescOnce sync.Once
	file_requests_proto_rawDescData = file_requests_proto_rawDesc
)

func file_requests_proto_rawDescGZIP() []byte {
	file_requests_proto_rawDescOnce.Do(func() {
		file_requests_proto_rawDescData = protoimpl.X.CompressGZIP(file_requests_proto_rawDescData)
	})
	return file_requests_proto_rawDescData
}

var file_requests_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_requests_proto_goTypes = []interface{}{
	(*CreateNodeRequest)(nil),  // 0: models.CreateNodeRequest
	(*ServeFileRequest)(nil),   // 1: models.ServeFileRequest
	(*SendMessageRequest)(nil), // 2: models.SendMessageRequest
	(*DownloadRequest)(nil),    // 3: models.DownloadRequest
}
var file_requests_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_requests_proto_init() }
func file_requests_proto_init() {
	if File_requests_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_requests_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateNodeRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_requests_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServeFileRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_requests_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendMessageRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_requests_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownloadRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_requests_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_requests_proto_goTypes,
		DependencyIndexes: file_requests_proto_depIdxs,
		MessageInfos:      file_requests_proto_msgTypes,
	}.Build()
	File_requests_proto = out.File
	file_requests_proto_rawDesc = nil
	file_requests_proto_goTypes = nil
	file_requests_proto_depIdxs = nil
}