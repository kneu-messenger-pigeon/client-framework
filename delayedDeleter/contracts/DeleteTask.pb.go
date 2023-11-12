// protoc --go_out=delayedDeleter delayedDeleter/DeleteTask.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: delayedDeleter/DeleteTask.proto

package contracts

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

type DeleteTask struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ScheduledAt int64 `protobuf:"varint,1,opt,name=scheduledAt,proto3" json:"scheduledAt,omitempty"`
	MessageId   int32 `protobuf:"varint,2,opt,name=MessageId,proto3" json:"MessageId,omitempty"`
	ChatId      int64 `protobuf:"varint,3,opt,name=ChatId,proto3" json:"ChatId,omitempty"`
}

func (x *DeleteTask) Reset() {
	*x = DeleteTask{}
	if protoimpl.UnsafeEnabled {
		mi := &file_delayedDeleter_DeleteTask_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteTask) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteTask) ProtoMessage() {}

func (x *DeleteTask) ProtoReflect() protoreflect.Message {
	mi := &file_delayedDeleter_DeleteTask_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteTask.ProtoReflect.Descriptor instead.
func (*DeleteTask) Descriptor() ([]byte, []int) {
	return file_delayedDeleter_DeleteTask_proto_rawDescGZIP(), []int{0}
}

func (x *DeleteTask) GetScheduledAt() int64 {
	if x != nil {
		return x.ScheduledAt
	}
	return 0
}

func (x *DeleteTask) GetMessageId() int32 {
	if x != nil {
		return x.MessageId
	}
	return 0
}

func (x *DeleteTask) GetChatId() int64 {
	if x != nil {
		return x.ChatId
	}
	return 0
}

var File_delayedDeleter_DeleteTask_proto protoreflect.FileDescriptor

var file_delayedDeleter_DeleteTask_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x65, 0x64, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x72,
	0x2f, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x66, 0x72, 0x61, 0x6d, 0x65, 0x77, 0x6f, 0x72, 0x6b, 0x22, 0x64, 0x0a, 0x0a,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x20, 0x0a, 0x0b, 0x73, 0x63,
	0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x64, 0x41, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0b, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1c, 0x0a, 0x09,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x09, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x43, 0x68,
	0x61, 0x74, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x43, 0x68, 0x61, 0x74,
	0x49, 0x64, 0x42, 0x13, 0x5a, 0x11, 0x2e, 0x2f, 0x3b, 0x64, 0x65, 0x6c, 0x61, 0x79, 0x65, 0x64,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_delayedDeleter_DeleteTask_proto_rawDescOnce sync.Once
	file_delayedDeleter_DeleteTask_proto_rawDescData = file_delayedDeleter_DeleteTask_proto_rawDesc
)

func file_delayedDeleter_DeleteTask_proto_rawDescGZIP() []byte {
	file_delayedDeleter_DeleteTask_proto_rawDescOnce.Do(func() {
		file_delayedDeleter_DeleteTask_proto_rawDescData = protoimpl.X.CompressGZIP(file_delayedDeleter_DeleteTask_proto_rawDescData)
	})
	return file_delayedDeleter_DeleteTask_proto_rawDescData
}

var file_delayedDeleter_DeleteTask_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_delayedDeleter_DeleteTask_proto_goTypes = []interface{}{
	(*DeleteTask)(nil), // 0: framework.DeleteTask
}
var file_delayedDeleter_DeleteTask_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_delayedDeleter_DeleteTask_proto_init() }
func file_delayedDeleter_DeleteTask_proto_init() {
	if File_delayedDeleter_DeleteTask_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_delayedDeleter_DeleteTask_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteTask); i {
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
			RawDescriptor: file_delayedDeleter_DeleteTask_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_delayedDeleter_DeleteTask_proto_goTypes,
		DependencyIndexes: file_delayedDeleter_DeleteTask_proto_depIdxs,
		MessageInfos:      file_delayedDeleter_DeleteTask_proto_msgTypes,
	}.Build()
	File_delayedDeleter_DeleteTask_proto = out.File
	file_delayedDeleter_DeleteTask_proto_rawDesc = nil
	file_delayedDeleter_DeleteTask_proto_goTypes = nil
	file_delayedDeleter_DeleteTask_proto_depIdxs = nil
}
