// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v5.27.3
// source: internal/block/proto/block.proto

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

type TXInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txid      []byte `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`                             // ID of tx this input refers
	VoutIndex int32  `protobuf:"varint,2,opt,name=vout_index,json=voutIndex,proto3" json:"vout_index,omitempty"` // index of an output in the tx
	Signature []byte `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	PubKey    []byte `protobuf:"bytes,4,opt,name=pub_key,json=pubKey,proto3" json:"pub_key,omitempty"`
}

func (x *TXInput) Reset() {
	*x = TXInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_block_proto_block_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TXInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TXInput) ProtoMessage() {}

func (x *TXInput) ProtoReflect() protoreflect.Message {
	mi := &file_internal_block_proto_block_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TXInput.ProtoReflect.Descriptor instead.
func (*TXInput) Descriptor() ([]byte, []int) {
	return file_internal_block_proto_block_proto_rawDescGZIP(), []int{0}
}

func (x *TXInput) GetTxid() []byte {
	if x != nil {
		return x.Txid
	}
	return nil
}

func (x *TXInput) GetVoutIndex() int32 {
	if x != nil {
		return x.VoutIndex
	}
	return 0
}

func (x *TXInput) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *TXInput) GetPubKey() []byte {
	if x != nil {
		return x.PubKey
	}
	return nil
}

type TXOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// stores the number of satoshis, which is 0.00000001 BTC.
	// this is the smallest unit of currency in Bitcoin
	Value int64 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	// we are not implementing the whole srcipt thing here, so just pubkey
	// pubkey hash is just pubkey hash, not an address
	PubKeyHash []byte `protobuf:"bytes,2,opt,name=pub_key_hash,json=pubKeyHash,proto3" json:"pub_key_hash,omitempty"`
}

func (x *TXOutput) Reset() {
	*x = TXOutput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_block_proto_block_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TXOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TXOutput) ProtoMessage() {}

func (x *TXOutput) ProtoReflect() protoreflect.Message {
	mi := &file_internal_block_proto_block_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TXOutput.ProtoReflect.Descriptor instead.
func (*TXOutput) Descriptor() ([]byte, []int) {
	return file_internal_block_proto_block_proto_rawDescGZIP(), []int{1}
}

func (x *TXOutput) GetValue() int64 {
	if x != nil {
		return x.Value
	}
	return 0
}

func (x *TXOutput) GetPubKeyHash() []byte {
	if x != nil {
		return x.PubKeyHash
	}
	return nil
}

type Transaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   []byte      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"` // hash of this tx
	VIn  []*TXInput  `protobuf:"bytes,2,rep,name=v_in,json=vIn,proto3" json:"v_in,omitempty"`
	VOut []*TXOutput `protobuf:"bytes,3,rep,name=v_out,json=vOut,proto3" json:"v_out,omitempty"`
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_block_proto_block_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_internal_block_proto_block_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_internal_block_proto_block_proto_rawDescGZIP(), []int{2}
}

func (x *Transaction) GetId() []byte {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *Transaction) GetVIn() []*TXInput {
	if x != nil {
		return x.VIn
	}
	return nil
}

func (x *Transaction) GetVOut() []*TXOutput {
	if x != nil {
		return x.VOut
	}
	return nil
}

// in bolt db, key will be block hash
type Block struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp     int64          `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Transactions  []*Transaction `protobuf:"bytes,2,rep,name=transactions,proto3" json:"transactions,omitempty"`
	PrevBlockHash []byte         `protobuf:"bytes,3,opt,name=prev_block_hash,json=prevBlockHash,proto3" json:"prev_block_hash,omitempty"`
	Hash          []byte         `protobuf:"bytes,4,opt,name=hash,proto3" json:"hash,omitempty"`
	Nonce         int64          `protobuf:"varint,5,opt,name=nonce,proto3" json:"nonce,omitempty"`
}

func (x *Block) Reset() {
	*x = Block{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_block_proto_block_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Block) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Block) ProtoMessage() {}

func (x *Block) ProtoReflect() protoreflect.Message {
	mi := &file_internal_block_proto_block_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Block.ProtoReflect.Descriptor instead.
func (*Block) Descriptor() ([]byte, []int) {
	return file_internal_block_proto_block_proto_rawDescGZIP(), []int{3}
}

func (x *Block) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Block) GetTransactions() []*Transaction {
	if x != nil {
		return x.Transactions
	}
	return nil
}

func (x *Block) GetPrevBlockHash() []byte {
	if x != nil {
		return x.PrevBlockHash
	}
	return nil
}

func (x *Block) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Block) GetNonce() int64 {
	if x != nil {
		return x.Nonce
	}
	return 0
}

var File_internal_block_proto_block_proto protoreflect.FileDescriptor

var file_internal_block_proto_block_proto_rawDesc = []byte{
	0x0a, 0x20, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x62, 0x6c, 0x6f, 0x63, 0x6b,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x73, 0x0a, 0x07, 0x54, 0x58, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x78, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x74, 0x78, 0x69,
	0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x76, 0x6f, 0x75, 0x74, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x76, 0x6f, 0x75, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x17,
	0x0a, 0x07, 0x70, 0x75, 0x62, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x06, 0x70, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x22, 0x42, 0x0a, 0x08, 0x54, 0x58, 0x4f, 0x75, 0x74,
	0x70, 0x75, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x20, 0x0a, 0x0c, 0x70, 0x75, 0x62,
	0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x0a, 0x70, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x48, 0x61, 0x73, 0x68, 0x22, 0x5a, 0x0a, 0x0b, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x04, 0x76, 0x5f,
	0x69, 0x6e, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x54, 0x58, 0x49, 0x6e, 0x70,
	0x75, 0x74, 0x52, 0x03, 0x76, 0x49, 0x6e, 0x12, 0x1e, 0x0a, 0x05, 0x76, 0x5f, 0x6f, 0x75, 0x74,
	0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x54, 0x58, 0x4f, 0x75, 0x74, 0x70, 0x75,
	0x74, 0x52, 0x04, 0x76, 0x4f, 0x75, 0x74, 0x22, 0xa9, 0x01, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12,
	0x30, 0x0a, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x26, 0x0a, 0x0f, 0x70, 0x72, 0x65, 0x76, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x5f,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x70, 0x72, 0x65, 0x76,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61, 0x73, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73,
	0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x14, 0x0a,
	0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6e, 0x6f,
	0x6e, 0x63, 0x65, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x63, 0x6a, 0x63, 0x37, 0x33, 0x37, 0x33, 0x2f, 0x62, 0x69, 0x74, 0x63, 0x6f, 0x69,
	0x6e, 0x5f, 0x67, 0x6f, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x62, 0x6c,
	0x6f, 0x63, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_internal_block_proto_block_proto_rawDescOnce sync.Once
	file_internal_block_proto_block_proto_rawDescData = file_internal_block_proto_block_proto_rawDesc
)

func file_internal_block_proto_block_proto_rawDescGZIP() []byte {
	file_internal_block_proto_block_proto_rawDescOnce.Do(func() {
		file_internal_block_proto_block_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_block_proto_block_proto_rawDescData)
	})
	return file_internal_block_proto_block_proto_rawDescData
}

var file_internal_block_proto_block_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_internal_block_proto_block_proto_goTypes = []interface{}{
	(*TXInput)(nil),     // 0: TXInput
	(*TXOutput)(nil),    // 1: TXOutput
	(*Transaction)(nil), // 2: Transaction
	(*Block)(nil),       // 3: Block
}
var file_internal_block_proto_block_proto_depIdxs = []int32{
	0, // 0: Transaction.v_in:type_name -> TXInput
	1, // 1: Transaction.v_out:type_name -> TXOutput
	2, // 2: Block.transactions:type_name -> Transaction
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_internal_block_proto_block_proto_init() }
func file_internal_block_proto_block_proto_init() {
	if File_internal_block_proto_block_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_block_proto_block_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TXInput); i {
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
		file_internal_block_proto_block_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TXOutput); i {
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
		file_internal_block_proto_block_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transaction); i {
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
		file_internal_block_proto_block_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Block); i {
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
			RawDescriptor: file_internal_block_proto_block_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_block_proto_block_proto_goTypes,
		DependencyIndexes: file_internal_block_proto_block_proto_depIdxs,
		MessageInfos:      file_internal_block_proto_block_proto_msgTypes,
	}.Build()
	File_internal_block_proto_block_proto = out.File
	file_internal_block_proto_block_proto_rawDesc = nil
	file_internal_block_proto_block_proto_goTypes = nil
	file_internal_block_proto_block_proto_depIdxs = nil
}
