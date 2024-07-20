// Package proto комментарий заглушка для обхода линтера
// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
//
//	protoc-gen-go v1.34.2
//	protoc        v5.27.1
//
// source: metrics/metrics.proto
package proto

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Metrics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type  string  `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Delta int64   `protobuf:"varint,3,opt,name=delta,proto3" json:"delta,omitempty"`
	Value float64 `protobuf:"fixed64,4,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Metrics) Reset() {
	*x = Metrics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metrics_metrics_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metrics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metrics) ProtoMessage() {}

func (x *Metrics) ProtoReflect() protoreflect.Message {
	mi := &file_metrics_metrics_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metrics.ProtoReflect.Descriptor instead.
func (*Metrics) Descriptor() ([]byte, []int) {
	return file_metrics_metrics_proto_rawDescGZIP(), []int{0}
}

func (x *Metrics) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Metrics) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Metrics) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *Metrics) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

type UpdateMetricsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metrics `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *UpdateMetricsRequest) Reset() {
	*x = UpdateMetricsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metrics_metrics_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricsRequest) ProtoMessage() {}

func (x *UpdateMetricsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_metrics_metrics_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricsRequest.ProtoReflect.Descriptor instead.
func (*UpdateMetricsRequest) Descriptor() ([]byte, []int) {
	return file_metrics_metrics_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateMetricsRequest) GetMetrics() []*Metrics {
	if x != nil {
		return x.Metrics
	}
	return nil
}

type UpdateMetricsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Error   string `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *UpdateMetricsResponse) Reset() {
	*x = UpdateMetricsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_metrics_metrics_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMetricsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMetricsResponse) ProtoMessage() {}

func (x *UpdateMetricsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_metrics_metrics_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMetricsResponse.ProtoReflect.Descriptor instead.
func (*UpdateMetricsResponse) Descriptor() ([]byte, []int) {
	return file_metrics_metrics_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateMetricsResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *UpdateMetricsResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_metrics_metrics_proto protoreflect.FileDescriptor

var file_metrics_metrics_proto_rawDesc = []byte{
	0x0a, 0x15, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x59, 0x0a, 0x07, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x22, 0x48, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x30, 0x0a, 0x07, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x73, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x22, 0x47, 0x0a, 0x15, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x14,
	0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x32, 0x6c, 0x0a, 0x0e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5a, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x23, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x6d,
	0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x1a, 0x5a, 0x18, 0x67, 0x6f, 0x2d, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x2d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_metrics_metrics_proto_rawDescOnce sync.Once
	file_metrics_metrics_proto_rawDescData = file_metrics_metrics_proto_rawDesc
)

func file_metrics_metrics_proto_rawDescGZIP() []byte {
	file_metrics_metrics_proto_rawDescOnce.Do(func() {
		file_metrics_metrics_proto_rawDescData = protoimpl.X.CompressGZIP(file_metrics_metrics_proto_rawDescData)
	})
	return file_metrics_metrics_proto_rawDescData
}

var file_metrics_metrics_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_metrics_metrics_proto_goTypes = []any{
	(*Metrics)(nil),               // 0: metrics_proto.Metrics
	(*UpdateMetricsRequest)(nil),  // 1: metrics_proto.UpdateMetricsRequest
	(*UpdateMetricsResponse)(nil), // 2: metrics_proto.UpdateMetricsResponse
}
var file_metrics_metrics_proto_depIdxs = []int32{
	0, // 0: metrics_proto.UpdateMetricsRequest.metrics:type_name -> metrics_proto.Metrics
	1, // 1: metrics_proto.MetricsService.UpdateMetrics:input_type -> metrics_proto.UpdateMetricsRequest
	2, // 2: metrics_proto.MetricsService.UpdateMetrics:output_type -> metrics_proto.UpdateMetricsResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_metrics_metrics_proto_init() }
func file_metrics_metrics_proto_init() {
	if File_metrics_metrics_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_metrics_metrics_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Metrics); i {
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
		file_metrics_metrics_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateMetricsRequest); i {
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
		file_metrics_metrics_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*UpdateMetricsResponse); i {
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
			RawDescriptor: file_metrics_metrics_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_metrics_metrics_proto_goTypes,
		DependencyIndexes: file_metrics_metrics_proto_depIdxs,
		MessageInfos:      file_metrics_metrics_proto_msgTypes,
	}.Build()
	File_metrics_metrics_proto = out.File
	file_metrics_metrics_proto_rawDesc = nil
	file_metrics_metrics_proto_goTypes = nil
	file_metrics_metrics_proto_depIdxs = nil
}
