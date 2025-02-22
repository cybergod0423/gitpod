// Copyright (c) 2020 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: control.proto

package api

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

type ExposePortRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// local port
	Port uint32 `protobuf:"varint,1,opt,name=port,proto3" json:"port,omitempty"`
	// external port if missing the the same as port
	TargetPort uint32 `protobuf:"varint,2,opt,name=target_port,json=targetPort,proto3" json:"target_port,omitempty"`
}

func (x *ExposePortRequest) Reset() {
	*x = ExposePortRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExposePortRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExposePortRequest) ProtoMessage() {}

func (x *ExposePortRequest) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExposePortRequest.ProtoReflect.Descriptor instead.
func (*ExposePortRequest) Descriptor() ([]byte, []int) {
	return file_control_proto_rawDescGZIP(), []int{0}
}

func (x *ExposePortRequest) GetPort() uint32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *ExposePortRequest) GetTargetPort() uint32 {
	if x != nil {
		return x.TargetPort
	}
	return 0
}

type ExposePortResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ExposePortResponse) Reset() {
	*x = ExposePortResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_control_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExposePortResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExposePortResponse) ProtoMessage() {}

func (x *ExposePortResponse) ProtoReflect() protoreflect.Message {
	mi := &file_control_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExposePortResponse.ProtoReflect.Descriptor instead.
func (*ExposePortResponse) Descriptor() ([]byte, []int) {
	return file_control_proto_rawDescGZIP(), []int{1}
}

var File_control_proto protoreflect.FileDescriptor

var file_control_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0a, 0x73, 0x75, 0x70, 0x65, 0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x22, 0x48, 0x0a, 0x11, 0x45,
	0x78, 0x70, 0x6f, 0x73, 0x65, 0x50, 0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04,
	0x70, 0x6f, 0x72, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x70,
	0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x50, 0x6f, 0x72, 0x74, 0x22, 0x14, 0x0a, 0x12, 0x45, 0x78, 0x70, 0x6f, 0x73, 0x65, 0x50,
	0x6f, 0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x5f, 0x0a, 0x0e, 0x43,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4d, 0x0a,
	0x0a, 0x45, 0x78, 0x70, 0x6f, 0x73, 0x65, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x1d, 0x2e, 0x73, 0x75,
	0x70, 0x65, 0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x2e, 0x45, 0x78, 0x70, 0x6f, 0x73, 0x65, 0x50,
	0x6f, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x73, 0x75, 0x70,
	0x65, 0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x2e, 0x45, 0x78, 0x70, 0x6f, 0x73, 0x65, 0x50, 0x6f,
	0x72, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x2c, 0x5a, 0x2a,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x69, 0x74, 0x70, 0x6f,
	0x64, 0x2d, 0x69, 0x6f, 0x2f, 0x67, 0x69, 0x74, 0x70, 0x6f, 0x64, 0x2f, 0x73, 0x75, 0x70, 0x65,
	0x72, 0x76, 0x69, 0x73, 0x6f, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_control_proto_rawDescOnce sync.Once
	file_control_proto_rawDescData = file_control_proto_rawDesc
)

func file_control_proto_rawDescGZIP() []byte {
	file_control_proto_rawDescOnce.Do(func() {
		file_control_proto_rawDescData = protoimpl.X.CompressGZIP(file_control_proto_rawDescData)
	})
	return file_control_proto_rawDescData
}

var file_control_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_control_proto_goTypes = []interface{}{
	(*ExposePortRequest)(nil),  // 0: supervisor.ExposePortRequest
	(*ExposePortResponse)(nil), // 1: supervisor.ExposePortResponse
}
var file_control_proto_depIdxs = []int32{
	0, // 0: supervisor.ControlService.ExposePort:input_type -> supervisor.ExposePortRequest
	1, // 1: supervisor.ControlService.ExposePort:output_type -> supervisor.ExposePortResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_control_proto_init() }
func file_control_proto_init() {
	if File_control_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_control_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExposePortRequest); i {
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
		file_control_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExposePortResponse); i {
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
			RawDescriptor: file_control_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_control_proto_goTypes,
		DependencyIndexes: file_control_proto_depIdxs,
		MessageInfos:      file_control_proto_msgTypes,
	}.Build()
	File_control_proto = out.File
	file_control_proto_rawDesc = nil
	file_control_proto_goTypes = nil
	file_control_proto_depIdxs = nil
}
