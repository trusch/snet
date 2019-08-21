// Code generated by protoc-gen-go. DO NOT EDIT.
// source: core.proto

package snet

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Message struct {
	Feed                 string               `protobuf:"bytes,1,opt,name=feed,proto3" json:"feed,omitempty"`
	Id                   string               `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	Previous             string               `protobuf:"bytes,3,opt,name=previous,proto3" json:"previous,omitempty"`
	Sequence             int64                `protobuf:"varint,4,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Type                 string               `protobuf:"bytes,6,opt,name=type,proto3" json:"type,omitempty"`
	Content              []byte               `protobuf:"bytes,7,opt,name=content,proto3" json:"content,omitempty"`
	Signature            string               `protobuf:"bytes,8,opt,name=signature,proto3" json:"signature,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,9,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_core_e9c66dafee7564d3, []int{0}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (dst *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(dst, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetFeed() string {
	if m != nil {
		return m.Feed
	}
	return ""
}

func (m *Message) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Message) GetPrevious() string {
	if m != nil {
		return m.Previous
	}
	return ""
}

func (m *Message) GetSequence() int64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *Message) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Message) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func (m *Message) GetSignature() string {
	if m != nil {
		return m.Signature
	}
	return ""
}

func (m *Message) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

type CreateRequest struct {
	Feed                 string   `protobuf:"bytes,1,opt,name=feed,proto3" json:"feed,omitempty"`
	Type                 string   `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Content              []byte   `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateRequest) Reset()         { *m = CreateRequest{} }
func (m *CreateRequest) String() string { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()    {}
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_core_e9c66dafee7564d3, []int{1}
}
func (m *CreateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateRequest.Unmarshal(m, b)
}
func (m *CreateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateRequest.Marshal(b, m, deterministic)
}
func (dst *CreateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateRequest.Merge(dst, src)
}
func (m *CreateRequest) XXX_Size() int {
	return xxx_messageInfo_CreateRequest.Size(m)
}
func (m *CreateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateRequest proto.InternalMessageInfo

func (m *CreateRequest) GetFeed() string {
	if m != nil {
		return m.Feed
	}
	return ""
}

func (m *CreateRequest) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *CreateRequest) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

type GetRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRequest) Reset()         { *m = GetRequest{} }
func (m *GetRequest) String() string { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()    {}
func (*GetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_core_e9c66dafee7564d3, []int{2}
}
func (m *GetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRequest.Unmarshal(m, b)
}
func (m *GetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRequest.Marshal(b, m, deterministic)
}
func (dst *GetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRequest.Merge(dst, src)
}
func (m *GetRequest) XXX_Size() int {
	return xxx_messageInfo_GetRequest.Size(m)
}
func (m *GetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRequest proto.InternalMessageInfo

func (m *GetRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type ListRequest struct {
	Feed                 string               `protobuf:"bytes,1,opt,name=feed,proto3" json:"feed,omitempty"`
	Type                 string               `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	StartSequence        int64                `protobuf:"varint,3,opt,name=start_sequence,json=startSequence,proto3" json:"start_sequence,omitempty"`
	EndSequence          int64                `protobuf:"varint,4,opt,name=end_sequence,json=endSequence,proto3" json:"end_sequence,omitempty"`
	StartTimestamp       *timestamp.Timestamp `protobuf:"bytes,5,opt,name=start_timestamp,json=startTimestamp,proto3" json:"start_timestamp,omitempty"`
	EndTimestamp         *timestamp.Timestamp `protobuf:"bytes,6,opt,name=end_timestamp,json=endTimestamp,proto3" json:"end_timestamp,omitempty"`
	Stream               bool                 `protobuf:"varint,7,opt,name=stream,proto3" json:"stream,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ListRequest) Reset()         { *m = ListRequest{} }
func (m *ListRequest) String() string { return proto.CompactTextString(m) }
func (*ListRequest) ProtoMessage()    {}
func (*ListRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_core_e9c66dafee7564d3, []int{3}
}
func (m *ListRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListRequest.Unmarshal(m, b)
}
func (m *ListRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListRequest.Marshal(b, m, deterministic)
}
func (dst *ListRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListRequest.Merge(dst, src)
}
func (m *ListRequest) XXX_Size() int {
	return xxx_messageInfo_ListRequest.Size(m)
}
func (m *ListRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListRequest proto.InternalMessageInfo

func (m *ListRequest) GetFeed() string {
	if m != nil {
		return m.Feed
	}
	return ""
}

func (m *ListRequest) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *ListRequest) GetStartSequence() int64 {
	if m != nil {
		return m.StartSequence
	}
	return 0
}

func (m *ListRequest) GetEndSequence() int64 {
	if m != nil {
		return m.EndSequence
	}
	return 0
}

func (m *ListRequest) GetStartTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.StartTimestamp
	}
	return nil
}

func (m *ListRequest) GetEndTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.EndTimestamp
	}
	return nil
}

func (m *ListRequest) GetStream() bool {
	if m != nil {
		return m.Stream
	}
	return false
}

func init() {
	proto.RegisterType((*Message)(nil), "Message")
	proto.RegisterType((*CreateRequest)(nil), "CreateRequest")
	proto.RegisterType((*GetRequest)(nil), "GetRequest")
	proto.RegisterType((*ListRequest)(nil), "ListRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CoreClient is the client API for Core service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CoreClient interface {
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*Message, error)
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*Message, error)
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (Core_ListClient, error)
}

type coreClient struct {
	cc *grpc.ClientConn
}

func NewCoreClient(cc *grpc.ClientConn) CoreClient {
	return &coreClient{cc}
}

func (c *coreClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/Core/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coreClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/Core/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coreClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (Core_ListClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Core_serviceDesc.Streams[0], "/Core/List", opts...)
	if err != nil {
		return nil, err
	}
	x := &coreListClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Core_ListClient interface {
	Recv() (*Message, error)
	grpc.ClientStream
}

type coreListClient struct {
	grpc.ClientStream
}

func (x *coreListClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// CoreServer is the server API for Core service.
type CoreServer interface {
	Get(context.Context, *GetRequest) (*Message, error)
	Create(context.Context, *CreateRequest) (*Message, error)
	List(*ListRequest, Core_ListServer) error
}

func RegisterCoreServer(s *grpc.Server, srv CoreServer) {
	s.RegisterService(&_Core_serviceDesc, srv)
}

func _Core_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoreServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Core/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoreServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Core_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoreServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Core/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoreServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Core_List_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CoreServer).List(m, &coreListServer{stream})
}

type Core_ListServer interface {
	Send(*Message) error
	grpc.ServerStream
}

type coreListServer struct {
	grpc.ServerStream
}

func (x *coreListServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

var _Core_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Core",
	HandlerType: (*CoreServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Core_Get_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _Core_Create_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "List",
			Handler:       _Core_List_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "core.proto",
}

func init() { proto.RegisterFile("core.proto", fileDescriptor_core_e9c66dafee7564d3) }

var fileDescriptor_core_e9c66dafee7564d3 = []byte{
	// 397 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0xc1, 0x8e, 0x94, 0x40,
	0x10, 0x86, 0xd3, 0x80, 0x0c, 0x14, 0x33, 0x63, 0xd2, 0x07, 0xd3, 0x21, 0x93, 0x88, 0x24, 0x26,
	0x9c, 0x7a, 0xcd, 0x7a, 0xf2, 0x64, 0x74, 0x0e, 0x7b, 0xd1, 0x83, 0xe8, 0xc9, 0xcb, 0x84, 0x85,
	0x5a, 0x82, 0x71, 0x68, 0xec, 0x2e, 0x4c, 0x7c, 0x18, 0xdf, 0xd0, 0x87, 0x30, 0x34, 0xdb, 0xb0,
	0x6b, 0x36, 0xd9, 0x78, 0xeb, 0xbf, 0xba, 0xf8, 0xfb, 0xaf, 0xaf, 0x00, 0xa8, 0x95, 0x46, 0x39,
	0x68, 0x45, 0x2a, 0x7d, 0xde, 0x2a, 0xd5, 0x7e, 0xc7, 0x0b, 0xab, 0xae, 0xc7, 0x9b, 0x0b, 0xea,
	0xce, 0x68, 0xa8, 0x3a, 0x0f, 0x73, 0x43, 0xfe, 0x87, 0xc1, 0xe6, 0x23, 0x1a, 0x53, 0xb5, 0xc8,
	0x39, 0x04, 0x37, 0x88, 0x8d, 0x60, 0x19, 0x2b, 0xe2, 0xd2, 0x9e, 0xf9, 0x1e, 0xbc, 0xae, 0x11,
	0x9e, 0xad, 0x78, 0x5d, 0xc3, 0x53, 0x88, 0x06, 0x8d, 0x3f, 0x3b, 0x35, 0x1a, 0xe1, 0xdb, 0xea,
	0xa2, 0xa7, 0x3b, 0x83, 0x3f, 0x46, 0xec, 0x6b, 0x14, 0x41, 0xc6, 0x0a, 0xbf, 0x5c, 0xf4, 0xe4,
	0x4d, 0xbf, 0x06, 0x14, 0xe1, 0xec, 0x3d, 0x9d, 0xb9, 0x80, 0x4d, 0xad, 0x7a, 0xc2, 0x9e, 0xc4,
	0x26, 0x63, 0xc5, 0xb6, 0x74, 0x92, 0x1f, 0x20, 0x36, 0x5d, 0xdb, 0x57, 0x34, 0x6a, 0x14, 0x91,
	0xfd, 0x64, 0x2d, 0xf0, 0x37, 0x00, 0xb5, 0xc6, 0x8a, 0xb0, 0x39, 0x55, 0x24, 0xe2, 0x8c, 0x15,
	0xc9, 0x65, 0x2a, 0xe7, 0x49, 0xa5, 0x9b, 0x54, 0x7e, 0x71, 0x93, 0x96, 0xf1, 0x6d, 0xf7, 0x3b,
	0xca, 0x3f, 0xc1, 0xee, 0x68, 0x45, 0x39, 0x05, 0x33, 0xf4, 0xe0, 0xcc, 0x2e, 0xab, 0xf7, 0x70,
	0x56, 0xff, 0x5e, 0xd6, 0xfc, 0x00, 0x70, 0x85, 0xe4, 0xfc, 0x66, 0x5e, 0xcc, 0xf1, 0xca, 0x7f,
	0x7b, 0x90, 0x7c, 0xe8, 0x0c, 0xfd, 0xef, 0x7b, 0x2f, 0x61, 0x6f, 0xa8, 0xd2, 0x74, 0x5a, 0x88,
	0xfa, 0x96, 0xe8, 0xce, 0x56, 0x3f, 0x3b, 0xac, 0x2f, 0x60, 0x8b, 0x7d, 0x73, 0xfa, 0x07, 0x7b,
	0x82, 0x7d, 0xb3, 0xb4, 0x1c, 0xe1, 0xe9, 0xec, 0xb4, 0xac, 0x5e, 0x3c, 0x79, 0x14, 0xd9, 0xfc,
	0xf8, 0xa2, 0xf9, 0x5b, 0xd8, 0x4d, 0xef, 0xac, 0x16, 0xe1, 0xa3, 0x16, 0x53, 0xb0, 0xd5, 0xe0,
	0x19, 0x84, 0x86, 0x34, 0x56, 0x67, 0xbb, 0xea, 0xa8, 0xbc, 0x55, 0x97, 0xdf, 0x20, 0x38, 0x2a,
	0x8d, 0xfc, 0x00, 0xfe, 0x15, 0x12, 0x4f, 0xe4, 0xca, 0x32, 0x8d, 0xa4, 0xfb, 0x33, 0x73, 0x08,
	0xe7, 0xb5, 0xf1, 0xbd, 0xbc, 0xb7, 0xbf, 0x3b, 0x3d, 0x19, 0x04, 0x13, 0x68, 0xbe, 0x95, 0x77,
	0x78, 0xaf, 0xf7, 0xaf, 0xd8, 0xfb, 0xf0, 0x6b, 0x60, 0x7a, 0xa4, 0xeb, 0xd0, 0xa6, 0x7d, 0xfd,
	0x37, 0x00, 0x00, 0xff, 0xff, 0x0d, 0xca, 0x23, 0x03, 0x29, 0x03, 0x00, 0x00,
}