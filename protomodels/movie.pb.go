// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0-devel
// 	protoc        v3.11.4
// source: movie.proto

package protomodels

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Movie struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Actors []*Actor `protobuf:"bytes,2,rep,name=actors,proto3" json:"actors,omitempty"`
}

func (x *Movie) Reset() {
	*x = Movie{}
	if protoimpl.UnsafeEnabled {
		mi := &file_movie_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Movie) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Movie) ProtoMessage() {}

func (x *Movie) ProtoReflect() protoreflect.Message {
	mi := &file_movie_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Movie.ProtoReflect.Descriptor instead.
func (*Movie) Descriptor() ([]byte, []int) {
	return file_movie_proto_rawDescGZIP(), []int{0}
}

func (x *Movie) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Movie) GetActors() []*Actor {
	if x != nil {
		return x.Actors
	}
	return nil
}

var File_movie_proto protoreflect.FileDescriptor

var file_movie_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x61,
	0x72, 0x74, 0x65, 0x6d, 0x69, 0x73, 0x1a, 0x0b, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x43, 0x0a, 0x05, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x26, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0e, 0x2e, 0x61, 0x72, 0x74, 0x65, 0x6d, 0x69, 0x73, 0x2e, 0x41, 0x63, 0x74, 0x6f, 0x72,
	0x52, 0x06, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x42, 0x28, 0x5a, 0x26, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x72, 0x65, 0x65, 0x6e, 0x61, 0x63, 0x2f, 0x61,
	0x72, 0x74, 0x65, 0x6d, 0x69, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_movie_proto_rawDescOnce sync.Once
	file_movie_proto_rawDescData = file_movie_proto_rawDesc
)

func file_movie_proto_rawDescGZIP() []byte {
	file_movie_proto_rawDescOnce.Do(func() {
		file_movie_proto_rawDescData = protoimpl.X.CompressGZIP(file_movie_proto_rawDescData)
	})
	return file_movie_proto_rawDescData
}

var file_movie_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_movie_proto_goTypes = []interface{}{
	(*Movie)(nil), // 0: artemis.Movie
	(*Actor)(nil), // 1: artemis.Actor
}
var file_movie_proto_depIdxs = []int32{
	1, // 0: artemis.Movie.actors:type_name -> artemis.Actor
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_movie_proto_init() }
func file_movie_proto_init() {
	if File_movie_proto != nil {
		return
	}
	file_actor_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_movie_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Movie); i {
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
			RawDescriptor: file_movie_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_movie_proto_goTypes,
		DependencyIndexes: file_movie_proto_depIdxs,
		MessageInfos:      file_movie_proto_msgTypes,
	}.Build()
	File_movie_proto = out.File
	file_movie_proto_rawDesc = nil
	file_movie_proto_goTypes = nil
	file_movie_proto_depIdxs = nil
}
