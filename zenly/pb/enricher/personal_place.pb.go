// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: enricher/personal_place.proto

package enricher

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

type PersonalPlace_Type int32

const (
	PersonalPlace_PERSONAL_PLACE_HOME   PersonalPlace_Type = 0
	PersonalPlace_PERSONAL_PLACE_WORK   PersonalPlace_Type = 1
	PersonalPlace_PERSONAL_PLACE_SCHOOL PersonalPlace_Type = 2
)

// Enum value maps for PersonalPlace_Type.
var (
	PersonalPlace_Type_name = map[int32]string{
		0: "PERSONAL_PLACE_HOME",
		1: "PERSONAL_PLACE_WORK",
		2: "PERSONAL_PLACE_SCHOOL",
	}
	PersonalPlace_Type_value = map[string]int32{
		"PERSONAL_PLACE_HOME":   0,
		"PERSONAL_PLACE_WORK":   1,
		"PERSONAL_PLACE_SCHOOL": 2,
	}
)

func (x PersonalPlace_Type) Enum() *PersonalPlace_Type {
	p := new(PersonalPlace_Type)
	*p = x
	return p
}

func (x PersonalPlace_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PersonalPlace_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_enricher_personal_place_proto_enumTypes[0].Descriptor()
}

func (PersonalPlace_Type) Type() protoreflect.EnumType {
	return &file_enricher_personal_place_proto_enumTypes[0]
}

func (x PersonalPlace_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PersonalPlace_Type.Descriptor instead.
func (PersonalPlace_Type) EnumDescriptor() ([]byte, []int) {
	return file_enricher_personal_place_proto_rawDescGZIP(), []int{0, 0}
}

type PersonalPlace struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type PersonalPlace_Type `protobuf:"varint,1,opt,name=type,proto3,enum=enricher.PersonalPlace_Type" json:"type,omitempty"`
}

func (x *PersonalPlace) Reset() {
	*x = PersonalPlace{}
	if protoimpl.UnsafeEnabled {
		mi := &file_enricher_personal_place_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PersonalPlace) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PersonalPlace) ProtoMessage() {}

func (x *PersonalPlace) ProtoReflect() protoreflect.Message {
	mi := &file_enricher_personal_place_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PersonalPlace.ProtoReflect.Descriptor instead.
func (*PersonalPlace) Descriptor() ([]byte, []int) {
	return file_enricher_personal_place_proto_rawDescGZIP(), []int{0}
}

func (x *PersonalPlace) GetType() PersonalPlace_Type {
	if x != nil {
		return x.Type
	}
	return PersonalPlace_PERSONAL_PLACE_HOME
}

var File_enricher_personal_place_proto protoreflect.FileDescriptor

var file_enricher_personal_place_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x65, 0x6e, 0x72, 0x69, 0x63, 0x68, 0x65, 0x72, 0x2f, 0x70, 0x65, 0x72, 0x73, 0x6f,
	0x6e, 0x61, 0x6c, 0x5f, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x08, 0x65, 0x6e, 0x72, 0x69, 0x63, 0x68, 0x65, 0x72, 0x22, 0x96, 0x01, 0x0a, 0x0d, 0x50, 0x65,
	0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x12, 0x30, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x65, 0x6e, 0x72, 0x69,
	0x63, 0x68, 0x65, 0x72, 0x2e, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x61, 0x6c, 0x50, 0x6c, 0x61,
	0x63, 0x65, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x53, 0x0a,
	0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x13, 0x50, 0x45, 0x52, 0x53, 0x4f, 0x4e, 0x41,
	0x4c, 0x5f, 0x50, 0x4c, 0x41, 0x43, 0x45, 0x5f, 0x48, 0x4f, 0x4d, 0x45, 0x10, 0x00, 0x12, 0x17,
	0x0a, 0x13, 0x50, 0x45, 0x52, 0x53, 0x4f, 0x4e, 0x41, 0x4c, 0x5f, 0x50, 0x4c, 0x41, 0x43, 0x45,
	0x5f, 0x57, 0x4f, 0x52, 0x4b, 0x10, 0x01, 0x12, 0x19, 0x0a, 0x15, 0x50, 0x45, 0x52, 0x53, 0x4f,
	0x4e, 0x41, 0x4c, 0x5f, 0x50, 0x4c, 0x41, 0x43, 0x45, 0x5f, 0x53, 0x43, 0x48, 0x4f, 0x4f, 0x4c,
	0x10, 0x02, 0x42, 0x33, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x73, 0x68, 0x65, 0x6b, 0x68, 0x69, 0x72, 0x69, 0x6e, 0x2f, 0x7a, 0x65, 0x6e, 0x6c, 0x79,
	0x2d, 0x74, 0x61, 0x73, 0x6b, 0x2f, 0x7a, 0x65, 0x6e, 0x6c, 0x79, 0x2f, 0x70, 0x62, 0x2f, 0x65,
	0x6e, 0x72, 0x69, 0x63, 0x68, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_enricher_personal_place_proto_rawDescOnce sync.Once
	file_enricher_personal_place_proto_rawDescData = file_enricher_personal_place_proto_rawDesc
)

func file_enricher_personal_place_proto_rawDescGZIP() []byte {
	file_enricher_personal_place_proto_rawDescOnce.Do(func() {
		file_enricher_personal_place_proto_rawDescData = protoimpl.X.CompressGZIP(file_enricher_personal_place_proto_rawDescData)
	})
	return file_enricher_personal_place_proto_rawDescData
}

var file_enricher_personal_place_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_enricher_personal_place_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_enricher_personal_place_proto_goTypes = []interface{}{
	(PersonalPlace_Type)(0), // 0: enricher.PersonalPlace.Type
	(*PersonalPlace)(nil),   // 1: enricher.PersonalPlace
}
var file_enricher_personal_place_proto_depIdxs = []int32{
	0, // 0: enricher.PersonalPlace.type:type_name -> enricher.PersonalPlace.Type
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_enricher_personal_place_proto_init() }
func file_enricher_personal_place_proto_init() {
	if File_enricher_personal_place_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_enricher_personal_place_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PersonalPlace); i {
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
			RawDescriptor: file_enricher_personal_place_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_enricher_personal_place_proto_goTypes,
		DependencyIndexes: file_enricher_personal_place_proto_depIdxs,
		EnumInfos:         file_enricher_personal_place_proto_enumTypes,
		MessageInfos:      file_enricher_personal_place_proto_msgTypes,
	}.Build()
	File_enricher_personal_place_proto = out.File
	file_enricher_personal_place_proto_rawDesc = nil
	file_enricher_personal_place_proto_goTypes = nil
	file_enricher_personal_place_proto_depIdxs = nil
}
