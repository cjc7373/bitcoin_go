// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.4
// source: internal/network/proto/protocol.proto

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

type NodeBroadcast struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nodes []*Node `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	TTL   uint32  `protobuf:"varint,2,opt,name=TTL,proto3" json:"TTL,omitempty"`
}

func (x *NodeBroadcast) Reset() {
	*x = NodeBroadcast{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_network_proto_protocol_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeBroadcast) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeBroadcast) ProtoMessage() {}

func (x *NodeBroadcast) ProtoReflect() protoreflect.Message {
	mi := &file_internal_network_proto_protocol_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeBroadcast.ProtoReflect.Descriptor instead.
func (*NodeBroadcast) Descriptor() ([]byte, []int) {
	return file_internal_network_proto_protocol_proto_rawDescGZIP(), []int{0}
}

func (x *NodeBroadcast) GetNodes() []*Node {
	if x != nil {
		return x.Nodes
	}
	return nil
}

func (x *NodeBroadcast) GetTTL() uint32 {
	if x != nil {
		return x.TTL
	}
	return 0
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_network_proto_protocol_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_internal_network_proto_protocol_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_internal_network_proto_protocol_proto_rawDescGZIP(), []int{1}
}

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Name    string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_network_proto_protocol_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_internal_network_proto_protocol_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_internal_network_proto_protocol_proto_rawDescGZIP(), []int{2}
}

func (x *Node) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Node) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_internal_network_proto_protocol_proto protoreflect.FileDescriptor

var file_internal_network_proto_protocol_proto_rawDesc = []byte{
	0x0a, 0x25, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3e, 0x0a, 0x0d, 0x4e, 0x6f, 0x64, 0x65, 0x42,
	0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05,
	0x6e, 0x6f, 0x64, 0x65, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x54, 0x54, 0x4c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x03, 0x54, 0x54, 0x4c, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x34, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x32, 0x56, 0x0a, 0x07, 0x42, 0x69, 0x74, 0x63, 0x6f, 0x69,
	0x6e, 0x12, 0x1f, 0x0a, 0x0c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4e, 0x6f, 0x64, 0x65,
	0x73, 0x12, 0x05, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x1a, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x12, 0x2a, 0x0a, 0x0e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x4e,
	0x6f, 0x64, 0x65, 0x73, 0x12, 0x0e, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x42, 0x72, 0x6f, 0x61, 0x64,
	0x63, 0x61, 0x73, 0x74, 0x1a, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x36,
	0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6a, 0x63,
	0x37, 0x33, 0x37, 0x33, 0x2f, 0x62, 0x69, 0x74, 0x63, 0x6f, 0x69, 0x6e, 0x5f, 0x67, 0x6f, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_network_proto_protocol_proto_rawDescOnce sync.Once
	file_internal_network_proto_protocol_proto_rawDescData = file_internal_network_proto_protocol_proto_rawDesc
)

func file_internal_network_proto_protocol_proto_rawDescGZIP() []byte {
	file_internal_network_proto_protocol_proto_rawDescOnce.Do(func() {
		file_internal_network_proto_protocol_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_network_proto_protocol_proto_rawDescData)
	})
	return file_internal_network_proto_protocol_proto_rawDescData
}

var file_internal_network_proto_protocol_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_internal_network_proto_protocol_proto_goTypes = []interface{}{
	(*NodeBroadcast)(nil), // 0: NodeBroadcast
	(*Empty)(nil),         // 1: Empty
	(*Node)(nil),          // 2: Node
}
var file_internal_network_proto_protocol_proto_depIdxs = []int32{
	2, // 0: NodeBroadcast.nodes:type_name -> Node
	2, // 1: Bitcoin.RequestNodes:input_type -> Node
	0, // 2: Bitcoin.BroadcastNodes:input_type -> NodeBroadcast
	1, // 3: Bitcoin.RequestNodes:output_type -> Empty
	1, // 4: Bitcoin.BroadcastNodes:output_type -> Empty
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_internal_network_proto_protocol_proto_init() }
func file_internal_network_proto_protocol_proto_init() {
	if File_internal_network_proto_protocol_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_network_proto_protocol_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeBroadcast); i {
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
		file_internal_network_proto_protocol_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_internal_network_proto_protocol_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Node); i {
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
			RawDescriptor: file_internal_network_proto_protocol_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_network_proto_protocol_proto_goTypes,
		DependencyIndexes: file_internal_network_proto_protocol_proto_depIdxs,
		MessageInfos:      file_internal_network_proto_protocol_proto_msgTypes,
	}.Build()
	File_internal_network_proto_protocol_proto = out.File
	file_internal_network_proto_protocol_proto_rawDesc = nil
	file_internal_network_proto_protocol_proto_goTypes = nil
	file_internal_network_proto_protocol_proto_depIdxs = nil
}
