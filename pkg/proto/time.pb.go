// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.7
// source: pkg/proto/time.proto

package proto

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

// Weekday is a day in week
type Weekday int32

const (
	Weekday_SATURDAY  Weekday = 0
	Weekday_SUNDAY    Weekday = 1
	Weekday_MONDAY    Weekday = 2
	Weekday_TUESDAY   Weekday = 3
	Weekday_WEDNESDAY Weekday = 4
	Weekday_THURSDAY  Weekday = 5
	Weekday_FRIDAY    Weekday = 6
)

// Enum value maps for Weekday.
var (
	Weekday_name = map[int32]string{
		0: "SATURDAY",
		1: "SUNDAY",
		2: "MONDAY",
		3: "TUESDAY",
		4: "WEDNESDAY",
		5: "THURSDAY",
		6: "FRIDAY",
	}
	Weekday_value = map[string]int32{
		"SATURDAY":  0,
		"SUNDAY":    1,
		"MONDAY":    2,
		"TUESDAY":   3,
		"WEDNESDAY": 4,
		"THURSDAY":  5,
		"FRIDAY":    6,
	}
)

func (x Weekday) Enum() *Weekday {
	p := new(Weekday)
	*p = x
	return p
}

func (x Weekday) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Weekday) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_proto_time_proto_enumTypes[0].Descriptor()
}

func (Weekday) Type() protoreflect.EnumType {
	return &file_pkg_proto_time_proto_enumTypes[0]
}

func (x Weekday) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Weekday.Descriptor instead.
func (Weekday) EnumDescriptor() ([]byte, []int) {
	return file_pkg_proto_time_proto_rawDescGZIP(), []int{0}
}

// ClassTime contains a single time which class is held.
// An array of it holds the days and times which class is held.
type ClassTime struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// On what day?
	Day Weekday `protobuf:"varint,1,opt,name=day,proto3,enum=proto.Weekday" json:"day,omitempty"`
	// Starting minute from 00:00
	StartMinute uint32 `protobuf:"varint,2,opt,name=start_minute,json=startMinute,proto3" json:"start_minute,omitempty"`
	// Ending minute from 00:00
	EndMinute uint32 `protobuf:"varint,3,opt,name=end_minute,json=endMinute,proto3" json:"end_minute,omitempty"`
}

func (x *ClassTime) Reset() {
	*x = ClassTime{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_proto_time_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClassTime) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClassTime) ProtoMessage() {}

func (x *ClassTime) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_proto_time_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClassTime.ProtoReflect.Descriptor instead.
func (*ClassTime) Descriptor() ([]byte, []int) {
	return file_pkg_proto_time_proto_rawDescGZIP(), []int{0}
}

func (x *ClassTime) GetDay() Weekday {
	if x != nil {
		return x.Day
	}
	return Weekday_SATURDAY
}

func (x *ClassTime) GetStartMinute() uint32 {
	if x != nil {
		return x.StartMinute
	}
	return 0
}

func (x *ClassTime) GetEndMinute() uint32 {
	if x != nil {
		return x.EndMinute
	}
	return 0
}

var File_pkg_proto_time_proto protoreflect.FileDescriptor

var file_pkg_proto_time_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x69, 0x6d, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6f, 0x0a,
	0x09, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x03, 0x64, 0x61,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x57, 0x65, 0x65, 0x6b, 0x64, 0x61, 0x79, 0x52, 0x03, 0x64, 0x61, 0x79, 0x12, 0x21, 0x0a, 0x0c,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x6d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x0b, 0x73, 0x74, 0x61, 0x72, 0x74, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x12,
	0x1d, 0x0a, 0x0a, 0x65, 0x6e, 0x64, 0x5f, 0x6d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x09, 0x65, 0x6e, 0x64, 0x4d, 0x69, 0x6e, 0x75, 0x74, 0x65, 0x2a, 0x65,
	0x0a, 0x07, 0x57, 0x65, 0x65, 0x6b, 0x64, 0x61, 0x79, 0x12, 0x0c, 0x0a, 0x08, 0x53, 0x41, 0x54,
	0x55, 0x52, 0x44, 0x41, 0x59, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x53, 0x55, 0x4e, 0x44, 0x41,
	0x59, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x4d, 0x4f, 0x4e, 0x44, 0x41, 0x59, 0x10, 0x02, 0x12,
	0x0b, 0x0a, 0x07, 0x54, 0x55, 0x45, 0x53, 0x44, 0x41, 0x59, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09,
	0x57, 0x45, 0x44, 0x4e, 0x45, 0x53, 0x44, 0x41, 0x59, 0x10, 0x04, 0x12, 0x0c, 0x0a, 0x08, 0x54,
	0x48, 0x55, 0x52, 0x53, 0x44, 0x41, 0x59, 0x10, 0x05, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x52, 0x49,
	0x44, 0x41, 0x59, 0x10, 0x06, 0x42, 0x1c, 0x5a, 0x1a, 0x43, 0x6f, 0x75, 0x72, 0x73, 0x65, 0x45,
	0x6e, 0x72, 0x6f, 0x6c, 0x6c, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_proto_time_proto_rawDescOnce sync.Once
	file_pkg_proto_time_proto_rawDescData = file_pkg_proto_time_proto_rawDesc
)

func file_pkg_proto_time_proto_rawDescGZIP() []byte {
	file_pkg_proto_time_proto_rawDescOnce.Do(func() {
		file_pkg_proto_time_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_proto_time_proto_rawDescData)
	})
	return file_pkg_proto_time_proto_rawDescData
}

var file_pkg_proto_time_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pkg_proto_time_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_pkg_proto_time_proto_goTypes = []interface{}{
	(Weekday)(0),      // 0: proto.Weekday
	(*ClassTime)(nil), // 1: proto.ClassTime
}
var file_pkg_proto_time_proto_depIdxs = []int32{
	0, // 0: proto.ClassTime.day:type_name -> proto.Weekday
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pkg_proto_time_proto_init() }
func file_pkg_proto_time_proto_init() {
	if File_pkg_proto_time_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_proto_time_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClassTime); i {
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
			RawDescriptor: file_pkg_proto_time_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pkg_proto_time_proto_goTypes,
		DependencyIndexes: file_pkg_proto_time_proto_depIdxs,
		EnumInfos:         file_pkg_proto_time_proto_enumTypes,
		MessageInfos:      file_pkg_proto_time_proto_msgTypes,
	}.Build()
	File_pkg_proto_time_proto = out.File
	file_pkg_proto_time_proto_rawDesc = nil
	file_pkg_proto_time_proto_goTypes = nil
	file_pkg_proto_time_proto_depIdxs = nil
}