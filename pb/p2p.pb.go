// Code generated by protoc-gen-go. DO NOT EDIT.
// source: p2p.proto

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type VersionMessage struct {
	Version              uint64   `protobuf:"varint,1,opt,name=Version,proto3" json:"Version,omitempty"`
	PeerID               []byte   `protobuf:"bytes,2,opt,name=PeerID,proto3" json:"PeerID,omitempty"`
	PeerInfo             []byte   `protobuf:"bytes,3,opt,name=PeerInfo,proto3" json:"PeerInfo,omitempty"`
	GenesisHash          []byte   `protobuf:"bytes,4,opt,name=GenesisHash,proto3" json:"GenesisHash,omitempty"`
	Height               uint64   `protobuf:"varint,5,opt,name=Height,proto3" json:"Height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VersionMessage) Reset()         { *m = VersionMessage{} }
func (m *VersionMessage) String() string { return proto.CompactTextString(m) }
func (*VersionMessage) ProtoMessage()    {}
func (*VersionMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_p2p_9d3f6773f7300cbf, []int{0}
}
func (m *VersionMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VersionMessage.Unmarshal(m, b)
}
func (m *VersionMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VersionMessage.Marshal(b, m, deterministic)
}
func (dst *VersionMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VersionMessage.Merge(dst, src)
}
func (m *VersionMessage) XXX_Size() int {
	return xxx_messageInfo_VersionMessage.Size(m)
}
func (m *VersionMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_VersionMessage.DiscardUnknown(m)
}

var xxx_messageInfo_VersionMessage proto.InternalMessageInfo

func (m *VersionMessage) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *VersionMessage) GetPeerID() []byte {
	if m != nil {
		return m.PeerID
	}
	return nil
}

func (m *VersionMessage) GetPeerInfo() []byte {
	if m != nil {
		return m.PeerInfo
	}
	return nil
}

func (m *VersionMessage) GetGenesisHash() []byte {
	if m != nil {
		return m.GenesisHash
	}
	return nil
}

func (m *VersionMessage) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

type VerackMessage struct {
	Version              uint64   `protobuf:"varint,1,opt,name=Version,proto3" json:"Version,omitempty"`
	PeerID               []byte   `protobuf:"bytes,2,opt,name=PeerID,proto3" json:"PeerID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VerackMessage) Reset()         { *m = VerackMessage{} }
func (m *VerackMessage) String() string { return proto.CompactTextString(m) }
func (*VerackMessage) ProtoMessage()    {}
func (*VerackMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_p2p_9d3f6773f7300cbf, []int{1}
}
func (m *VerackMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VerackMessage.Unmarshal(m, b)
}
func (m *VerackMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VerackMessage.Marshal(b, m, deterministic)
}
func (dst *VerackMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VerackMessage.Merge(dst, src)
}
func (m *VerackMessage) XXX_Size() int {
	return xxx_messageInfo_VerackMessage.Size(m)
}
func (m *VerackMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_VerackMessage.DiscardUnknown(m)
}

var xxx_messageInfo_VerackMessage proto.InternalMessageInfo

func (m *VerackMessage) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *VerackMessage) GetPeerID() []byte {
	if m != nil {
		return m.PeerID
	}
	return nil
}

type PingMessage struct {
	Nonce                uint64   `protobuf:"varint,1,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PingMessage) Reset()         { *m = PingMessage{} }
func (m *PingMessage) String() string { return proto.CompactTextString(m) }
func (*PingMessage) ProtoMessage()    {}
func (*PingMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_p2p_9d3f6773f7300cbf, []int{2}
}
func (m *PingMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PingMessage.Unmarshal(m, b)
}
func (m *PingMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PingMessage.Marshal(b, m, deterministic)
}
func (dst *PingMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PingMessage.Merge(dst, src)
}
func (m *PingMessage) XXX_Size() int {
	return xxx_messageInfo_PingMessage.Size(m)
}
func (m *PingMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_PingMessage.DiscardUnknown(m)
}

var xxx_messageInfo_PingMessage proto.InternalMessageInfo

func (m *PingMessage) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

type PongMessage struct {
	Nonce                uint64   `protobuf:"varint,1,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PongMessage) Reset()         { *m = PongMessage{} }
func (m *PongMessage) String() string { return proto.CompactTextString(m) }
func (*PongMessage) ProtoMessage()    {}
func (*PongMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_p2p_9d3f6773f7300cbf, []int{3}
}
func (m *PongMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PongMessage.Unmarshal(m, b)
}
func (m *PongMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PongMessage.Marshal(b, m, deterministic)
}
func (dst *PongMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PongMessage.Merge(dst, src)
}
func (m *PongMessage) XXX_Size() int {
	return xxx_messageInfo_PongMessage.Size(m)
}
func (m *PongMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_PongMessage.DiscardUnknown(m)
}

var xxx_messageInfo_PongMessage proto.InternalMessageInfo

func (m *PongMessage) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

type RejectMessage struct {
	Message              string   `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RejectMessage) Reset()         { *m = RejectMessage{} }
func (m *RejectMessage) String() string { return proto.CompactTextString(m) }
func (*RejectMessage) ProtoMessage()    {}
func (*RejectMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_p2p_9d3f6773f7300cbf, []int{4}
}
func (m *RejectMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RejectMessage.Unmarshal(m, b)
}
func (m *RejectMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RejectMessage.Marshal(b, m, deterministic)
}
func (dst *RejectMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RejectMessage.Merge(dst, src)
}
func (m *RejectMessage) XXX_Size() int {
	return xxx_messageInfo_RejectMessage.Size(m)
}
func (m *RejectMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_RejectMessage.DiscardUnknown(m)
}

var xxx_messageInfo_RejectMessage proto.InternalMessageInfo

func (m *RejectMessage) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type GetBlockMessage struct {
	LocatorHashes        [][]byte `protobuf:"bytes,1,rep,name=LocatorHashes,proto3" json:"LocatorHashes,omitempty"`
	HashStop             []byte   `protobuf:"bytes,2,opt,name=HashStop,proto3" json:"HashStop,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetBlockMessage) Reset()         { *m = GetBlockMessage{} }
func (m *GetBlockMessage) String() string { return proto.CompactTextString(m) }
func (*GetBlockMessage) ProtoMessage()    {}
func (*GetBlockMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_p2p_9d3f6773f7300cbf, []int{5}
}
func (m *GetBlockMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetBlockMessage.Unmarshal(m, b)
}
func (m *GetBlockMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetBlockMessage.Marshal(b, m, deterministic)
}
func (dst *GetBlockMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetBlockMessage.Merge(dst, src)
}
func (m *GetBlockMessage) XXX_Size() int {
	return xxx_messageInfo_GetBlockMessage.Size(m)
}
func (m *GetBlockMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_GetBlockMessage.DiscardUnknown(m)
}

var xxx_messageInfo_GetBlockMessage proto.InternalMessageInfo

func (m *GetBlockMessage) GetLocatorHashes() [][]byte {
	if m != nil {
		return m.LocatorHashes
	}
	return nil
}

func (m *GetBlockMessage) GetHashStop() []byte {
	if m != nil {
		return m.HashStop
	}
	return nil
}

// Response to GetBlockMessage
type BlockMessage struct {
	Blocks               []*Block `protobuf:"bytes,1,rep,name=Blocks,proto3" json:"Blocks,omitempty"`
	LatestBlockHash      []byte   `protobuf:"bytes,2,opt,name=LatestBlockHash,proto3" json:"LatestBlockHash,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockMessage) Reset()         { *m = BlockMessage{} }
func (m *BlockMessage) String() string { return proto.CompactTextString(m) }
func (*BlockMessage) ProtoMessage()    {}
func (*BlockMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_p2p_9d3f6773f7300cbf, []int{6}
}
func (m *BlockMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockMessage.Unmarshal(m, b)
}
func (m *BlockMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockMessage.Marshal(b, m, deterministic)
}
func (dst *BlockMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockMessage.Merge(dst, src)
}
func (m *BlockMessage) XXX_Size() int {
	return xxx_messageInfo_BlockMessage.Size(m)
}
func (m *BlockMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockMessage.DiscardUnknown(m)
}

var xxx_messageInfo_BlockMessage proto.InternalMessageInfo

func (m *BlockMessage) GetBlocks() []*Block {
	if m != nil {
		return m.Blocks
	}
	return nil
}

func (m *BlockMessage) GetLatestBlockHash() []byte {
	if m != nil {
		return m.LatestBlockHash
	}
	return nil
}

type GetAddrMessage struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAddrMessage) Reset()         { *m = GetAddrMessage{} }
func (m *GetAddrMessage) String() string { return proto.CompactTextString(m) }
func (*GetAddrMessage) ProtoMessage()    {}
func (*GetAddrMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_p2p_9d3f6773f7300cbf, []int{7}
}
func (m *GetAddrMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAddrMessage.Unmarshal(m, b)
}
func (m *GetAddrMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAddrMessage.Marshal(b, m, deterministic)
}
func (dst *GetAddrMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAddrMessage.Merge(dst, src)
}
func (m *GetAddrMessage) XXX_Size() int {
	return xxx_messageInfo_GetAddrMessage.Size(m)
}
func (m *GetAddrMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAddrMessage.DiscardUnknown(m)
}

var xxx_messageInfo_GetAddrMessage proto.InternalMessageInfo

type AddrMessage struct {
	Addrs                [][]byte `protobuf:"bytes,1,rep,name=Addrs,proto3" json:"Addrs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddrMessage) Reset()         { *m = AddrMessage{} }
func (m *AddrMessage) String() string { return proto.CompactTextString(m) }
func (*AddrMessage) ProtoMessage()    {}
func (*AddrMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_p2p_9d3f6773f7300cbf, []int{8}
}
func (m *AddrMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddrMessage.Unmarshal(m, b)
}
func (m *AddrMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddrMessage.Marshal(b, m, deterministic)
}
func (dst *AddrMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddrMessage.Merge(dst, src)
}
func (m *AddrMessage) XXX_Size() int {
	return xxx_messageInfo_AddrMessage.Size(m)
}
func (m *AddrMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_AddrMessage.DiscardUnknown(m)
}

var xxx_messageInfo_AddrMessage proto.InternalMessageInfo

func (m *AddrMessage) GetAddrs() [][]byte {
	if m != nil {
		return m.Addrs
	}
	return nil
}

func init() {
	proto.RegisterType((*VersionMessage)(nil), "pb.VersionMessage")
	proto.RegisterType((*VerackMessage)(nil), "pb.VerackMessage")
	proto.RegisterType((*PingMessage)(nil), "pb.PingMessage")
	proto.RegisterType((*PongMessage)(nil), "pb.PongMessage")
	proto.RegisterType((*RejectMessage)(nil), "pb.RejectMessage")
	proto.RegisterType((*GetBlockMessage)(nil), "pb.GetBlockMessage")
	proto.RegisterType((*BlockMessage)(nil), "pb.BlockMessage")
	proto.RegisterType((*GetAddrMessage)(nil), "pb.GetAddrMessage")
	proto.RegisterType((*AddrMessage)(nil), "pb.AddrMessage")
}

func init() { proto.RegisterFile("p2p.proto", fileDescriptor_p2p_9d3f6773f7300cbf) }

var fileDescriptor_p2p_9d3f6773f7300cbf = []byte{
	// 315 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0x4d, 0x4f, 0xfa, 0x40,
	0x10, 0xc6, 0x53, 0xde, 0xfe, 0x7f, 0xa6, 0x05, 0x4c, 0x63, 0xcc, 0x86, 0x53, 0x2d, 0x1e, 0xea,
	0x85, 0x03, 0x7e, 0x02, 0x8c, 0x09, 0x98, 0xa0, 0x21, 0x25, 0xe1, 0xe2, 0xa9, 0x94, 0x11, 0xaa,
	0xb2, 0xb3, 0xd9, 0xdd, 0x0f, 0xe3, 0xc7, 0x35, 0xfb, 0x52, 0x04, 0x2f, 0x26, 0xde, 0xe6, 0xf7,
	0xcc, 0xf4, 0xe9, 0x33, 0x93, 0x85, 0xae, 0x98, 0x88, 0xb1, 0x90, 0xa4, 0x29, 0x6e, 0x88, 0xcd,
	0x30, 0x2a, 0xe9, 0x70, 0x20, 0xee, 0x94, 0xf4, 0x33, 0x80, 0xfe, 0x1a, 0xa5, 0xaa, 0x88, 0x3f,
	0xa1, 0x52, 0xc5, 0x0e, 0x63, 0x06, 0xff, 0xbc, 0xc2, 0x82, 0x24, 0xc8, 0x5a, 0x79, 0x8d, 0xf1,
	0x15, 0x74, 0x96, 0x88, 0xf2, 0xf1, 0x81, 0x35, 0x92, 0x20, 0x8b, 0x72, 0x4f, 0xf1, 0x10, 0xfe,
	0xdb, 0x8a, 0xbf, 0x12, 0x6b, 0xda, 0xce, 0x91, 0xe3, 0x04, 0xc2, 0x19, 0x72, 0x54, 0x95, 0x9a,
	0x17, 0x6a, 0xcf, 0x5a, 0xb6, 0x7d, 0x2a, 0x19, 0xd7, 0x39, 0x56, 0xbb, 0xbd, 0x66, 0x6d, 0xfb,
	0x3b, 0x4f, 0xe9, 0x14, 0x7a, 0x6b, 0x94, 0x45, 0xf9, 0xfe, 0xe7, 0x60, 0xe9, 0x08, 0xc2, 0x65,
	0xc5, 0x77, 0xb5, 0xc1, 0x25, 0xb4, 0x9f, 0x89, 0x97, 0xe8, 0x3f, 0x77, 0x60, 0x87, 0xe8, 0xb7,
	0xa1, 0x5b, 0xe8, 0xe5, 0xf8, 0x86, 0xa5, 0x3e, 0x09, 0xe3, 0x4b, 0x3b, 0xd8, 0xcd, 0x6b, 0x4c,
	0x57, 0x30, 0x98, 0xa1, 0xbe, 0xff, 0xa0, 0xef, 0xe4, 0x37, 0xd0, 0x5b, 0x50, 0x59, 0x68, 0x92,
	0x66, 0x63, 0x54, 0x2c, 0x48, 0x9a, 0x59, 0x94, 0x9f, 0x8b, 0xe6, 0x8c, 0xa6, 0x5a, 0x69, 0x12,
	0x7e, 0x8f, 0x23, 0xa7, 0x2f, 0x10, 0x9d, 0x39, 0x5e, 0x43, 0xc7, 0xb2, 0xb3, 0x0a, 0x27, 0xdd,
	0xb1, 0xd8, 0x8c, 0xad, 0x92, 0xfb, 0x46, 0x9c, 0xc1, 0x60, 0x51, 0x68, 0x54, 0x2e, 0x8a, 0xbd,
	0xbe, 0x73, 0xfd, 0x29, 0xa7, 0x17, 0xd0, 0x9f, 0xa1, 0x9e, 0x6e, 0xb7, 0xb2, 0xde, 0x61, 0x04,
	0xe1, 0x09, 0x9a, 0x9b, 0x18, 0xac, 0x73, 0x3b, 0xd8, 0x74, 0xec, 0x13, 0xba, 0xfb, 0x0a, 0x00,
	0x00, 0xff, 0xff, 0x4a, 0x6e, 0xe2, 0x0d, 0x61, 0x02, 0x00, 0x00,
}
