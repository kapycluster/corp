// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.25.3
// source: kubeconfig.proto

package kubeconfig

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

type GetKubeConfigRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetKubeConfigRequest) Reset() {
	*x = GetKubeConfigRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kubeconfig_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetKubeConfigRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetKubeConfigRequest) ProtoMessage() {}

func (x *GetKubeConfigRequest) ProtoReflect() protoreflect.Message {
	mi := &file_kubeconfig_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetKubeConfigRequest.ProtoReflect.Descriptor instead.
func (*GetKubeConfigRequest) Descriptor() ([]byte, []int) {
	return file_kubeconfig_proto_rawDescGZIP(), []int{0}
}

type GetKubeConfigResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	KubeConfig string `protobuf:"bytes,1,opt,name=kubeConfig,proto3" json:"kubeConfig,omitempty"`
	Metadata   string `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *GetKubeConfigResponse) Reset() {
	*x = GetKubeConfigResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_kubeconfig_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetKubeConfigResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetKubeConfigResponse) ProtoMessage() {}

func (x *GetKubeConfigResponse) ProtoReflect() protoreflect.Message {
	mi := &file_kubeconfig_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetKubeConfigResponse.ProtoReflect.Descriptor instead.
func (*GetKubeConfigResponse) Descriptor() ([]byte, []int) {
	return file_kubeconfig_proto_rawDescGZIP(), []int{1}
}

func (x *GetKubeConfigResponse) GetKubeConfig() string {
	if x != nil {
		return x.KubeConfig
	}
	return ""
}

func (x *GetKubeConfigResponse) GetMetadata() string {
	if x != nil {
		return x.Metadata
	}
	return ""
}

var File_kubeconfig_proto protoreflect.FileDescriptor

var file_kubeconfig_proto_rawDesc = []byte{
	0x0a, 0x10, 0x6b, 0x75, 0x62, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x6b, 0x75, 0x62, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22, 0x16,
	0x0a, 0x14, 0x47, 0x65, 0x74, 0x4b, 0x75, 0x62, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x53, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x4b, 0x75, 0x62,
	0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x1e, 0x0a, 0x0a, 0x6b, 0x75, 0x62, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x6b, 0x75, 0x62, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12,
	0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x32, 0x69, 0x0a, 0x11, 0x4b,
	0x75, 0x62, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x54, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4b, 0x75, 0x62, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x12, 0x20, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x47,
	0x65, 0x74, 0x4b, 0x75, 0x62, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x6b, 0x75, 0x62, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x47, 0x65, 0x74, 0x4b, 0x75, 0x62, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x3a, 0x5a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x61, 0x70, 0x79, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72,
	0x2f, 0x6b, 0x61, 0x70, 0x79, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6b, 0x75, 0x62, 0x65, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_kubeconfig_proto_rawDescOnce sync.Once
	file_kubeconfig_proto_rawDescData = file_kubeconfig_proto_rawDesc
)

func file_kubeconfig_proto_rawDescGZIP() []byte {
	file_kubeconfig_proto_rawDescOnce.Do(func() {
		file_kubeconfig_proto_rawDescData = protoimpl.X.CompressGZIP(file_kubeconfig_proto_rawDescData)
	})
	return file_kubeconfig_proto_rawDescData
}

var file_kubeconfig_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_kubeconfig_proto_goTypes = []interface{}{
	(*GetKubeConfigRequest)(nil),  // 0: kubeconfig.GetKubeConfigRequest
	(*GetKubeConfigResponse)(nil), // 1: kubeconfig.GetKubeConfigResponse
}
var file_kubeconfig_proto_depIdxs = []int32{
	0, // 0: kubeconfig.KubeConfigService.GetKubeConfig:input_type -> kubeconfig.GetKubeConfigRequest
	1, // 1: kubeconfig.KubeConfigService.GetKubeConfig:output_type -> kubeconfig.GetKubeConfigResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_kubeconfig_proto_init() }
func file_kubeconfig_proto_init() {
	if File_kubeconfig_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_kubeconfig_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetKubeConfigRequest); i {
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
		file_kubeconfig_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetKubeConfigResponse); i {
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
			RawDescriptor: file_kubeconfig_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_kubeconfig_proto_goTypes,
		DependencyIndexes: file_kubeconfig_proto_depIdxs,
		MessageInfos:      file_kubeconfig_proto_msgTypes,
	}.Build()
	File_kubeconfig_proto = out.File
	file_kubeconfig_proto_rawDesc = nil
	file_kubeconfig_proto_goTypes = nil
	file_kubeconfig_proto_depIdxs = nil
}
