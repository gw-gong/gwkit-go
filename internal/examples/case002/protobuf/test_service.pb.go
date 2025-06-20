// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.2
// source: test_service.proto

package protobuf

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TestRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RequestName   string                 `protobuf:"bytes,1,opt,name=request_name,json=requestName,proto3" json:"request_name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TestRequest) Reset() {
	*x = TestRequest{}
	mi := &file_test_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TestRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestRequest) ProtoMessage() {}

func (x *TestRequest) ProtoReflect() protoreflect.Message {
	mi := &file_test_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestRequest.ProtoReflect.Descriptor instead.
func (*TestRequest) Descriptor() ([]byte, []int) {
	return file_test_service_proto_rawDescGZIP(), []int{0}
}

func (x *TestRequest) GetRequestName() string {
	if x != nil {
		return x.RequestName
	}
	return ""
}

type TestResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ResponseMsg   string                 `protobuf:"bytes,1,opt,name=response_msg,json=responseMsg,proto3" json:"response_msg,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TestResponse) Reset() {
	*x = TestResponse{}
	mi := &file_test_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TestResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestResponse) ProtoMessage() {}

func (x *TestResponse) ProtoReflect() protoreflect.Message {
	mi := &file_test_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestResponse.ProtoReflect.Descriptor instead.
func (*TestResponse) Descriptor() ([]byte, []int) {
	return file_test_service_proto_rawDescGZIP(), []int{1}
}

func (x *TestResponse) GetResponseMsg() string {
	if x != nil {
		return x.ResponseMsg
	}
	return ""
}

var File_test_service_proto protoreflect.FileDescriptor

const file_test_service_proto_rawDesc = "" +
	"\n" +
	"\x12test_service.proto\"0\n" +
	"\vTestRequest\x12!\n" +
	"\frequest_name\x18\x01 \x01(\tR\vrequestName\"1\n" +
	"\fTestResponse\x12!\n" +
	"\fresponse_msg\x18\x01 \x01(\tR\vresponseMsg28\n" +
	"\vTestService\x12)\n" +
	"\bTestFunc\x12\f.TestRequest\x1a\r.TestResponse\"\x00B\fZ\n" +
	".;protobufb\x06proto3"

var (
	file_test_service_proto_rawDescOnce sync.Once
	file_test_service_proto_rawDescData []byte
)

func file_test_service_proto_rawDescGZIP() []byte {
	file_test_service_proto_rawDescOnce.Do(func() {
		file_test_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_test_service_proto_rawDesc), len(file_test_service_proto_rawDesc)))
	})
	return file_test_service_proto_rawDescData
}

var file_test_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_test_service_proto_goTypes = []any{
	(*TestRequest)(nil),  // 0: TestRequest
	(*TestResponse)(nil), // 1: TestResponse
}
var file_test_service_proto_depIdxs = []int32{
	0, // 0: TestService.TestFunc:input_type -> TestRequest
	1, // 1: TestService.TestFunc:output_type -> TestResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_test_service_proto_init() }
func file_test_service_proto_init() {
	if File_test_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_test_service_proto_rawDesc), len(file_test_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_test_service_proto_goTypes,
		DependencyIndexes: file_test_service_proto_depIdxs,
		MessageInfos:      file_test_service_proto_msgTypes,
	}.Build()
	File_test_service_proto = out.File
	file_test_service_proto_goTypes = nil
	file_test_service_proto_depIdxs = nil
}
