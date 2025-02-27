// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v5.29.1
// source: libra-backend.proto

package pb

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

type QueryEmbedding struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Query         string                 `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
	Embedding     []float32              `protobuf:"fixed32,2,rep,packed,name=embedding,proto3" json:"embedding,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QueryEmbedding) Reset() {
	*x = QueryEmbedding{}
	mi := &file_libra_backend_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryEmbedding) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryEmbedding) ProtoMessage() {}

func (x *QueryEmbedding) ProtoReflect() protoreflect.Message {
	mi := &file_libra_backend_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryEmbedding.ProtoReflect.Descriptor instead.
func (*QueryEmbedding) Descriptor() ([]byte, []int) {
	return file_libra_backend_proto_rawDescGZIP(), []int{0}
}

func (x *QueryEmbedding) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *QueryEmbedding) GetEmbedding() []float32 {
	if x != nil {
		return x.Embedding
	}
	return nil
}

var File_libra_backend_proto protoreflect.FileDescriptor

var file_libra_backend_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6c, 0x69, 0x62, 0x72, 0x61, 0x2d, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22, 0x44, 0x0a, 0x0e, 0x51, 0x75, 0x65,
	0x72, 0x79, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x71,
	0x75, 0x65, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75, 0x65, 0x72,
	0x79, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x02, 0x52, 0x09, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x42,
	0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_libra_backend_proto_rawDescOnce sync.Once
	file_libra_backend_proto_rawDescData = file_libra_backend_proto_rawDesc
)

func file_libra_backend_proto_rawDescGZIP() []byte {
	file_libra_backend_proto_rawDescOnce.Do(func() {
		file_libra_backend_proto_rawDescData = protoimpl.X.CompressGZIP(file_libra_backend_proto_rawDescData)
	})
	return file_libra_backend_proto_rawDescData
}

var file_libra_backend_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_libra_backend_proto_goTypes = []any{
	(*QueryEmbedding)(nil), // 0: pb.QueryEmbedding
}
var file_libra_backend_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_libra_backend_proto_init() }
func file_libra_backend_proto_init() {
	if File_libra_backend_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_libra_backend_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_libra_backend_proto_goTypes,
		DependencyIndexes: file_libra_backend_proto_depIdxs,
		MessageInfos:      file_libra_backend_proto_msgTypes,
	}.Build()
	File_libra_backend_proto = out.File
	file_libra_backend_proto_rawDesc = nil
	file_libra_backend_proto_goTypes = nil
	file_libra_backend_proto_depIdxs = nil
}
