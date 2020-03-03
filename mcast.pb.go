// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mcast.proto

package main

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type JoinReq struct {
	Hierarchy            int32    `protobuf:"varint,1,opt,name=hierarchy,proto3" json:"hierarchy,omitempty"`
	Time                 int64    `protobuf:"varint,2,opt,name=time,proto3" json:"time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JoinReq) Reset()         { *m = JoinReq{} }
func (m *JoinReq) String() string { return proto.CompactTextString(m) }
func (*JoinReq) ProtoMessage()    {}
func (*JoinReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_17a4197c673499c8, []int{0}
}

func (m *JoinReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinReq.Unmarshal(m, b)
}
func (m *JoinReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinReq.Marshal(b, m, deterministic)
}
func (m *JoinReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinReq.Merge(m, src)
}
func (m *JoinReq) XXX_Size() int {
	return xxx_messageInfo_JoinReq.Size(m)
}
func (m *JoinReq) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinReq.DiscardUnknown(m)
}

var xxx_messageInfo_JoinReq proto.InternalMessageInfo

func (m *JoinReq) GetHierarchy() int32 {
	if m != nil {
		return m.Hierarchy
	}
	return 0
}

func (m *JoinReq) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

type JoinResp struct {
	Hierarchy            int32          `protobuf:"varint,1,opt,name=hierarchy,proto3" json:"hierarchy,omitempty"`
	RttMs                int64          `protobuf:"varint,2,opt,name=rtt_ms,json=rttMs,proto3" json:"rtt_ms,omitempty"`
	Time                 int64          `protobuf:"varint,3,opt,name=time,proto3" json:"time,omitempty"`
	Self                 *FingerEntry   `protobuf:"bytes,4,opt,name=self,proto3" json:"self,omitempty"`
	FEList               []*FingerEntry `protobuf:"bytes,5,rep,name=FEList,json=fEList,proto3" json:"FEList,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *JoinResp) Reset()         { *m = JoinResp{} }
func (m *JoinResp) String() string { return proto.CompactTextString(m) }
func (*JoinResp) ProtoMessage()    {}
func (*JoinResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_17a4197c673499c8, []int{1}
}

func (m *JoinResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinResp.Unmarshal(m, b)
}
func (m *JoinResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinResp.Marshal(b, m, deterministic)
}
func (m *JoinResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinResp.Merge(m, src)
}
func (m *JoinResp) XXX_Size() int {
	return xxx_messageInfo_JoinResp.Size(m)
}
func (m *JoinResp) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinResp.DiscardUnknown(m)
}

var xxx_messageInfo_JoinResp proto.InternalMessageInfo

func (m *JoinResp) GetHierarchy() int32 {
	if m != nil {
		return m.Hierarchy
	}
	return 0
}

func (m *JoinResp) GetRttMs() int64 {
	if m != nil {
		return m.RttMs
	}
	return 0
}

func (m *JoinResp) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *JoinResp) GetSelf() *FingerEntry {
	if m != nil {
		return m.Self
	}
	return nil
}

func (m *JoinResp) GetFEList() []*FingerEntry {
	if m != nil {
		return m.FEList
	}
	return nil
}

type NodeId struct {
	Ids                  []uint64 `protobuf:"varint,1,rep,packed,name=ids,proto3" json:"ids,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeId) Reset()         { *m = NodeId{} }
func (m *NodeId) String() string { return proto.CompactTextString(m) }
func (*NodeId) ProtoMessage()    {}
func (*NodeId) Descriptor() ([]byte, []int) {
	return fileDescriptor_17a4197c673499c8, []int{2}
}

func (m *NodeId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeId.Unmarshal(m, b)
}
func (m *NodeId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeId.Marshal(b, m, deterministic)
}
func (m *NodeId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeId.Merge(m, src)
}
func (m *NodeId) XXX_Size() int {
	return xxx_messageInfo_NodeId.Size(m)
}
func (m *NodeId) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeId.DiscardUnknown(m)
}

var xxx_messageInfo_NodeId proto.InternalMessageInfo

func (m *NodeId) GetIds() []uint64 {
	if m != nil {
		return m.Ids
	}
	return nil
}

type FingerEntry struct {
	Id *NodeId `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	//fixed64 address = 2;	// for now just the IPv4 address //TODO remove
	Hostname             string   `protobuf:"bytes,3,opt,name=hostname,proto3" json:"hostname,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FingerEntry) Reset()         { *m = FingerEntry{} }
func (m *FingerEntry) String() string { return proto.CompactTextString(m) }
func (*FingerEntry) ProtoMessage()    {}
func (*FingerEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_17a4197c673499c8, []int{3}
}

func (m *FingerEntry) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FingerEntry.Unmarshal(m, b)
}
func (m *FingerEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FingerEntry.Marshal(b, m, deterministic)
}
func (m *FingerEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FingerEntry.Merge(m, src)
}
func (m *FingerEntry) XXX_Size() int {
	return xxx_messageInfo_FingerEntry.Size(m)
}
func (m *FingerEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_FingerEntry.DiscardUnknown(m)
}

var xxx_messageInfo_FingerEntry proto.InternalMessageInfo

func (m *FingerEntry) GetId() *NodeId {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *FingerEntry) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

type EvalInfo struct {
	Hops                 int32    `protobuf:"varint,1,opt,name=hops,proto3" json:"hops,omitempty"`
	Time                 int64    `protobuf:"varint,2,opt,name=time,proto3" json:"time,omitempty"`
	Node                 *NodeId  `protobuf:"bytes,3,opt,name=node,proto3" json:"node,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EvalInfo) Reset()         { *m = EvalInfo{} }
func (m *EvalInfo) String() string { return proto.CompactTextString(m) }
func (*EvalInfo) ProtoMessage()    {}
func (*EvalInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_17a4197c673499c8, []int{4}
}

func (m *EvalInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EvalInfo.Unmarshal(m, b)
}
func (m *EvalInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EvalInfo.Marshal(b, m, deterministic)
}
func (m *EvalInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EvalInfo.Merge(m, src)
}
func (m *EvalInfo) XXX_Size() int {
	return xxx_messageInfo_EvalInfo.Size(m)
}
func (m *EvalInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_EvalInfo.DiscardUnknown(m)
}

var xxx_messageInfo_EvalInfo proto.InternalMessageInfo

func (m *EvalInfo) GetHops() int32 {
	if m != nil {
		return m.Hops
	}
	return 0
}

func (m *EvalInfo) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *EvalInfo) GetNode() *NodeId {
	if m != nil {
		return m.Node
	}
	return nil
}

type FwdPacket struct {
	Limit                *NodeId     `protobuf:"bytes,1,opt,name=limit,proto3" json:"limit,omitempty"`
	EvalList             []*EvalInfo `protobuf:"bytes,2,rep,name=evalList,proto3" json:"evalList,omitempty"`
	Payload              []byte      `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
	Src                  *NodeId     `protobuf:"bytes,4,opt,name=src,proto3" json:"src,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *FwdPacket) Reset()         { *m = FwdPacket{} }
func (m *FwdPacket) String() string { return proto.CompactTextString(m) }
func (*FwdPacket) ProtoMessage()    {}
func (*FwdPacket) Descriptor() ([]byte, []int) {
	return fileDescriptor_17a4197c673499c8, []int{5}
}

func (m *FwdPacket) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FwdPacket.Unmarshal(m, b)
}
func (m *FwdPacket) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FwdPacket.Marshal(b, m, deterministic)
}
func (m *FwdPacket) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FwdPacket.Merge(m, src)
}
func (m *FwdPacket) XXX_Size() int {
	return xxx_messageInfo_FwdPacket.Size(m)
}
func (m *FwdPacket) XXX_DiscardUnknown() {
	xxx_messageInfo_FwdPacket.DiscardUnknown(m)
}

var xxx_messageInfo_FwdPacket proto.InternalMessageInfo

func (m *FwdPacket) GetLimit() *NodeId {
	if m != nil {
		return m.Limit
	}
	return nil
}

func (m *FwdPacket) GetEvalList() []*EvalInfo {
	if m != nil {
		return m.EvalList
	}
	return nil
}

func (m *FwdPacket) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *FwdPacket) GetSrc() *NodeId {
	if m != nil {
		return m.Src
	}
	return nil
}

type GetFERequest struct {
	Hierarchy            int32    `protobuf:"varint,1,opt,name=hierarchy,proto3" json:"hierarchy,omitempty"`
	Limit                *NodeId  `protobuf:"bytes,2,opt,name=limit,proto3" json:"limit,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetFERequest) Reset()         { *m = GetFERequest{} }
func (m *GetFERequest) String() string { return proto.CompactTextString(m) }
func (*GetFERequest) ProtoMessage()    {}
func (*GetFERequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_17a4197c673499c8, []int{6}
}

func (m *GetFERequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetFERequest.Unmarshal(m, b)
}
func (m *GetFERequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetFERequest.Marshal(b, m, deterministic)
}
func (m *GetFERequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetFERequest.Merge(m, src)
}
func (m *GetFERequest) XXX_Size() int {
	return xxx_messageInfo_GetFERequest.Size(m)
}
func (m *GetFERequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetFERequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetFERequest proto.InternalMessageInfo

func (m *GetFERequest) GetHierarchy() int32 {
	if m != nil {
		return m.Hierarchy
	}
	return 0
}

func (m *GetFERequest) GetLimit() *NodeId {
	if m != nil {
		return m.Limit
	}
	return nil
}

type GetFEResponse struct {
	Hierarchy            int32        `protobuf:"varint,1,opt,name=hierarchy,proto3" json:"hierarchy,omitempty"`
	Limit                *NodeId      `protobuf:"bytes,2,opt,name=limit,proto3" json:"limit,omitempty"`
	Src                  *NodeId      `protobuf:"bytes,3,opt,name=src,proto3" json:"src,omitempty"`
	NewFE                *FingerEntry `protobuf:"bytes,4,opt,name=newFE,proto3" json:"newFE,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *GetFEResponse) Reset()         { *m = GetFEResponse{} }
func (m *GetFEResponse) String() string { return proto.CompactTextString(m) }
func (*GetFEResponse) ProtoMessage()    {}
func (*GetFEResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_17a4197c673499c8, []int{7}
}

func (m *GetFEResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetFEResponse.Unmarshal(m, b)
}
func (m *GetFEResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetFEResponse.Marshal(b, m, deterministic)
}
func (m *GetFEResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetFEResponse.Merge(m, src)
}
func (m *GetFEResponse) XXX_Size() int {
	return xxx_messageInfo_GetFEResponse.Size(m)
}
func (m *GetFEResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetFEResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetFEResponse proto.InternalMessageInfo

func (m *GetFEResponse) GetHierarchy() int32 {
	if m != nil {
		return m.Hierarchy
	}
	return 0
}

func (m *GetFEResponse) GetLimit() *NodeId {
	if m != nil {
		return m.Limit
	}
	return nil
}

func (m *GetFEResponse) GetSrc() *NodeId {
	if m != nil {
		return m.Src
	}
	return nil
}

func (m *GetFEResponse) GetNewFE() *FingerEntry {
	if m != nil {
		return m.NewFE
	}
	return nil
}

type Successor struct {
	FE                   *FingerEntry `protobuf:"bytes,1,opt,name=FE,json=fE,proto3" json:"FE,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Successor) Reset()         { *m = Successor{} }
func (m *Successor) String() string { return proto.CompactTextString(m) }
func (*Successor) ProtoMessage()    {}
func (*Successor) Descriptor() ([]byte, []int) {
	return fileDescriptor_17a4197c673499c8, []int{8}
}

func (m *Successor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Successor.Unmarshal(m, b)
}
func (m *Successor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Successor.Marshal(b, m, deterministic)
}
func (m *Successor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Successor.Merge(m, src)
}
func (m *Successor) XXX_Size() int {
	return xxx_messageInfo_Successor.Size(m)
}
func (m *Successor) XXX_DiscardUnknown() {
	xxx_messageInfo_Successor.DiscardUnknown(m)
}

var xxx_messageInfo_Successor proto.InternalMessageInfo

func (m *Successor) GetFE() *FingerEntry {
	if m != nil {
		return m.FE
	}
	return nil
}

type Predecessor struct {
	FE                   *FingerEntry `protobuf:"bytes,1,opt,name=FE,json=fE,proto3" json:"FE,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Predecessor) Reset()         { *m = Predecessor{} }
func (m *Predecessor) String() string { return proto.CompactTextString(m) }
func (*Predecessor) ProtoMessage()    {}
func (*Predecessor) Descriptor() ([]byte, []int) {
	return fileDescriptor_17a4197c673499c8, []int{9}
}

func (m *Predecessor) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Predecessor.Unmarshal(m, b)
}
func (m *Predecessor) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Predecessor.Marshal(b, m, deterministic)
}
func (m *Predecessor) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Predecessor.Merge(m, src)
}
func (m *Predecessor) XXX_Size() int {
	return xxx_messageInfo_Predecessor.Size(m)
}
func (m *Predecessor) XXX_DiscardUnknown() {
	xxx_messageInfo_Predecessor.DiscardUnknown(m)
}

var xxx_messageInfo_Predecessor proto.InternalMessageInfo

func (m *Predecessor) GetFE() *FingerEntry {
	if m != nil {
		return m.FE
	}
	return nil
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_17a4197c673499c8, []int{10}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

func init() {
	proto.RegisterType((*JoinReq)(nil), "main.JoinReq")
	proto.RegisterType((*JoinResp)(nil), "main.JoinResp")
	proto.RegisterType((*NodeId)(nil), "main.NodeId")
	proto.RegisterType((*FingerEntry)(nil), "main.FingerEntry")
	proto.RegisterType((*EvalInfo)(nil), "main.EvalInfo")
	proto.RegisterType((*FwdPacket)(nil), "main.FwdPacket")
	proto.RegisterType((*GetFERequest)(nil), "main.GetFERequest")
	proto.RegisterType((*GetFEResponse)(nil), "main.GetFEResponse")
	proto.RegisterType((*Successor)(nil), "main.Successor")
	proto.RegisterType((*Predecessor)(nil), "main.Predecessor")
	proto.RegisterType((*Empty)(nil), "main.Empty")
}

func init() { proto.RegisterFile("mcast.proto", fileDescriptor_17a4197c673499c8) }

var fileDescriptor_17a4197c673499c8 = []byte{
	// 531 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0x4f, 0x6f, 0xd3, 0x4e,
	0x10, 0xad, 0xff, 0xe5, 0xcf, 0x38, 0xcd, 0xef, 0xc7, 0x20, 0x24, 0x2b, 0xaa, 0x90, 0x59, 0xa9,
	0x6a, 0xe0, 0x10, 0xa1, 0x70, 0xec, 0xd9, 0xae, 0x8a, 0x00, 0x45, 0x0e, 0x77, 0x64, 0xec, 0x0d,
	0x59, 0x11, 0x7b, 0xdd, 0xdd, 0x6d, 0xa3, 0x7c, 0x10, 0x2e, 0xdc, 0xf9, 0x92, 0x9c, 0x90, 0xd7,
	0x9b, 0xc4, 0x8d, 0x12, 0x15, 0x89, 0x9b, 0x77, 0x66, 0xde, 0x9b, 0x37, 0x6f, 0x76, 0x0d, 0x7e,
	0x91, 0xa5, 0x52, 0x4d, 0x2a, 0xc1, 0x15, 0x47, 0xb7, 0x48, 0x59, 0x49, 0xae, 0xa1, 0xfb, 0x9e,
	0xb3, 0x32, 0xa1, 0x77, 0x78, 0x01, 0xfd, 0x25, 0xa3, 0x22, 0x15, 0xd9, 0x72, 0x13, 0x58, 0xa1,
	0x35, 0xf6, 0x92, 0x7d, 0x00, 0x11, 0x5c, 0xc5, 0x0a, 0x1a, 0xd8, 0xa1, 0x35, 0x76, 0x12, 0xfd,
	0x4d, 0x7e, 0x59, 0xd0, 0x6b, 0xd0, 0xb2, 0x7a, 0x02, 0xfe, 0x02, 0x3a, 0x42, 0xa9, 0x2f, 0x85,
	0x34, 0x04, 0x9e, 0x50, 0xea, 0xa3, 0xdc, 0xb1, 0x3a, 0x7b, 0x56, 0xbc, 0x04, 0x57, 0xd2, 0xd5,
	0x22, 0x70, 0x43, 0x6b, 0xec, 0x4f, 0x9f, 0x4d, 0x6a, 0x9d, 0x93, 0x98, 0x95, 0xdf, 0xa8, 0x88,
	0x4a, 0x25, 0x36, 0x89, 0x4e, 0xe3, 0x6b, 0xe8, 0xc4, 0xd1, 0x07, 0x26, 0x55, 0xe0, 0x85, 0xce,
	0xf1, 0xc2, 0xce, 0x42, 0x17, 0x90, 0x11, 0x74, 0x3e, 0xf1, 0x9c, 0xde, 0xe6, 0xf8, 0x3f, 0x38,
	0x2c, 0x97, 0x81, 0x15, 0x3a, 0x63, 0x37, 0xa9, 0x3f, 0xc9, 0x0d, 0xf8, 0x2d, 0x08, 0x5e, 0x80,
	0xcd, 0x72, 0x2d, 0xdf, 0x9f, 0x0e, 0x1a, 0xc6, 0x06, 0x9a, 0xd8, 0x2c, 0xc7, 0x11, 0xf4, 0x96,
	0x5c, 0xaa, 0x32, 0x35, 0x92, 0xfb, 0xc9, 0xee, 0x4c, 0x3e, 0x43, 0x2f, 0x7a, 0x48, 0x57, 0xb7,
	0xe5, 0x82, 0xd7, 0x63, 0x2d, 0x79, 0x25, 0x8d, 0x0d, 0xfa, 0xfb, 0x98, 0x81, 0x18, 0x82, 0x5b,
	0xf2, 0xbc, 0xe1, 0x3a, 0xec, 0xa7, 0x33, 0xe4, 0x87, 0x05, 0xfd, 0x78, 0x9d, 0xcf, 0xd2, 0xec,
	0x3b, 0x55, 0x48, 0xc0, 0x5b, 0xb1, 0x82, 0xa9, 0xa3, 0x02, 0x9b, 0x14, 0xbe, 0x81, 0x1e, 0x7d,
	0x48, 0x57, 0xda, 0x19, 0x5b, 0x3b, 0x33, 0x6c, 0xca, 0xb6, 0xea, 0x92, 0x5d, 0x1e, 0x03, 0xe8,
	0x56, 0xe9, 0x66, 0xc5, 0xd3, 0x5c, 0x4b, 0x18, 0x24, 0xdb, 0x23, 0xbe, 0x04, 0x47, 0x8a, 0xcc,
	0xec, 0xe0, 0x71, 0x9f, 0x3a, 0x41, 0x66, 0x30, 0xb8, 0xa1, 0x2a, 0x8e, 0x12, 0x7a, 0x77, 0x4f,
	0xa5, 0x7a, 0x62, 0xfb, 0x3b, 0xdd, 0xf6, 0x49, 0xdd, 0xe4, 0xa7, 0x05, 0xe7, 0x86, 0x52, 0x56,
	0xbc, 0x94, 0xf4, 0xdf, 0x39, 0xb7, 0x53, 0x38, 0x27, 0xa6, 0xc0, 0x2b, 0xf0, 0x4a, 0xba, 0x8e,
	0xa3, 0xd3, 0x77, 0xad, 0xc9, 0x93, 0x09, 0xf4, 0xe7, 0xf7, 0x59, 0x46, 0xa5, 0xe4, 0x02, 0x5f,
	0x81, 0x1d, 0x47, 0x66, 0x05, 0x47, 0x20, 0xf6, 0x22, 0x22, 0x6f, 0xc1, 0x9f, 0x09, 0x9a, 0xd3,
	0xbf, 0x47, 0x74, 0xc1, 0x8b, 0x8a, 0x4a, 0x6d, 0xa6, 0xbf, 0x2d, 0xe8, 0xea, 0x77, 0x4a, 0x05,
	0x5e, 0x81, 0x5b, 0xbf, 0x2f, 0x3c, 0x6f, 0x30, 0xe6, 0xa5, 0x8e, 0x86, 0xed, 0xa3, 0xac, 0xc8,
	0x19, 0x5e, 0x82, 0x13, 0xaf, 0x73, 0xfc, 0xcf, 0x70, 0x6f, 0x2f, 0xcc, 0xc8, 0x37, 0xab, 0xaf,
	0x99, 0xc9, 0x19, 0x5e, 0xc3, 0xb0, 0xb6, 0xb8, 0x75, 0xdf, 0xb1, 0x29, 0x68, 0xef, 0x72, 0xf4,
	0xfc, 0x51, 0xac, 0x59, 0x06, 0x39, 0xc3, 0x09, 0x0c, 0xe6, 0x54, 0xed, 0x6d, 0x30, 0xcd, 0x76,
	0x81, 0xc3, 0x66, 0x53, 0x18, 0xce, 0xa9, 0x6a, 0xdb, 0x60, 0x46, 0x6f, 0x85, 0x0e, 0x30, 0x5f,
	0x3b, 0xfa, 0xdf, 0xf4, 0xee, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb5, 0x17, 0x00, 0xfe, 0xaa,
	0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// McasterClient is the client API for Mcaster service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type McasterClient interface {
	Join(ctx context.Context, in *JoinReq, opts ...grpc.CallOption) (*JoinResp, error)
	Fwd(ctx context.Context, in *FwdPacket, opts ...grpc.CallOption) (*Empty, error)
	GetFingerEntry(ctx context.Context, in *GetFERequest, opts ...grpc.CallOption) (*GetFEResponse, error)
	SetSuccessor(ctx context.Context, in *Successor, opts ...grpc.CallOption) (*Empty, error)
	SetPredecessor(ctx context.Context, in *Predecessor, opts ...grpc.CallOption) (*Empty, error)
}

type mcasterClient struct {
	cc grpc.ClientConnInterface
}

func NewMcasterClient(cc grpc.ClientConnInterface) McasterClient {
	return &mcasterClient{cc}
}

func (c *mcasterClient) Join(ctx context.Context, in *JoinReq, opts ...grpc.CallOption) (*JoinResp, error) {
	out := new(JoinResp)
	err := c.cc.Invoke(ctx, "/main.mcaster/Join", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mcasterClient) Fwd(ctx context.Context, in *FwdPacket, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/main.mcaster/Fwd", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mcasterClient) GetFingerEntry(ctx context.Context, in *GetFERequest, opts ...grpc.CallOption) (*GetFEResponse, error) {
	out := new(GetFEResponse)
	err := c.cc.Invoke(ctx, "/main.mcaster/GetFingerEntry", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mcasterClient) SetSuccessor(ctx context.Context, in *Successor, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/main.mcaster/SetSuccessor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mcasterClient) SetPredecessor(ctx context.Context, in *Predecessor, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/main.mcaster/SetPredecessor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// McasterServer is the server API for Mcaster service.
type McasterServer interface {
	Join(context.Context, *JoinReq) (*JoinResp, error)
	Fwd(context.Context, *FwdPacket) (*Empty, error)
	GetFingerEntry(context.Context, *GetFERequest) (*GetFEResponse, error)
	SetSuccessor(context.Context, *Successor) (*Empty, error)
	SetPredecessor(context.Context, *Predecessor) (*Empty, error)
}

// UnimplementedMcasterServer can be embedded to have forward compatible implementations.
type UnimplementedMcasterServer struct {
}

func (*UnimplementedMcasterServer) Join(ctx context.Context, req *JoinReq) (*JoinResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (*UnimplementedMcasterServer) Fwd(ctx context.Context, req *FwdPacket) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Fwd not implemented")
}
func (*UnimplementedMcasterServer) GetFingerEntry(ctx context.Context, req *GetFERequest) (*GetFEResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFingerEntry not implemented")
}
func (*UnimplementedMcasterServer) SetSuccessor(ctx context.Context, req *Successor) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetSuccessor not implemented")
}
func (*UnimplementedMcasterServer) SetPredecessor(ctx context.Context, req *Predecessor) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetPredecessor not implemented")
}

func RegisterMcasterServer(s *grpc.Server, srv McasterServer) {
	s.RegisterService(&_Mcaster_serviceDesc, srv)
}

func _Mcaster_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(McasterServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.mcaster/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(McasterServer).Join(ctx, req.(*JoinReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mcaster_Fwd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FwdPacket)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(McasterServer).Fwd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.mcaster/Fwd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(McasterServer).Fwd(ctx, req.(*FwdPacket))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mcaster_GetFingerEntry_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFERequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(McasterServer).GetFingerEntry(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.mcaster/GetFingerEntry",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(McasterServer).GetFingerEntry(ctx, req.(*GetFERequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mcaster_SetSuccessor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Successor)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(McasterServer).SetSuccessor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.mcaster/SetSuccessor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(McasterServer).SetSuccessor(ctx, req.(*Successor))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mcaster_SetPredecessor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Predecessor)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(McasterServer).SetPredecessor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.mcaster/SetPredecessor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(McasterServer).SetPredecessor(ctx, req.(*Predecessor))
	}
	return interceptor(ctx, in, info, handler)
}

var _Mcaster_serviceDesc = grpc.ServiceDesc{
	ServiceName: "main.mcaster",
	HandlerType: (*McasterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Join",
			Handler:    _Mcaster_Join_Handler,
		},
		{
			MethodName: "Fwd",
			Handler:    _Mcaster_Fwd_Handler,
		},
		{
			MethodName: "GetFingerEntry",
			Handler:    _Mcaster_GetFingerEntry_Handler,
		},
		{
			MethodName: "SetSuccessor",
			Handler:    _Mcaster_SetSuccessor_Handler,
		},
		{
			MethodName: "SetPredecessor",
			Handler:    _Mcaster_SetPredecessor_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mcast.proto",
}
