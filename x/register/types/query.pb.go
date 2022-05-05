// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: stratos/register/v1/query.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// QueryResourceNodeRequest is request type for the Query/ResourceNode RPC method
type QueryResourceNodeRequest struct {
	// network_addr defines the node address to query for.
	NetworkAddr string `protobuf:"bytes,1,opt,name=network_addr,json=networkAddr,proto3" json:"network_addr,omitempty"`
}

func (m *QueryResourceNodeRequest) Reset()         { *m = QueryResourceNodeRequest{} }
func (m *QueryResourceNodeRequest) String() string { return proto.CompactTextString(m) }
func (*QueryResourceNodeRequest) ProtoMessage()    {}
func (*QueryResourceNodeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_59a612d1da8c0670, []int{0}
}
func (m *QueryResourceNodeRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryResourceNodeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryResourceNodeRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryResourceNodeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryResourceNodeRequest.Merge(m, src)
}
func (m *QueryResourceNodeRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryResourceNodeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryResourceNodeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryResourceNodeRequest proto.InternalMessageInfo

func (m *QueryResourceNodeRequest) GetNetworkAddr() string {
	if m != nil {
		return m.NetworkAddr
	}
	return ""
}

// QueryResourceNodeResponse is response type for the Query/ResourceNode RPC method
type QueryResourceNodeResponse struct {
	// node defines the the resourceNode info.
	Node ResourceNode `protobuf:"bytes,1,opt,name=node,proto3" json:"node"`
}

func (m *QueryResourceNodeResponse) Reset()         { *m = QueryResourceNodeResponse{} }
func (m *QueryResourceNodeResponse) String() string { return proto.CompactTextString(m) }
func (*QueryResourceNodeResponse) ProtoMessage()    {}
func (*QueryResourceNodeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_59a612d1da8c0670, []int{1}
}
func (m *QueryResourceNodeResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryResourceNodeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryResourceNodeResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryResourceNodeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryResourceNodeResponse.Merge(m, src)
}
func (m *QueryResourceNodeResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryResourceNodeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryResourceNodeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryResourceNodeResponse proto.InternalMessageInfo

func (m *QueryResourceNodeResponse) GetNode() ResourceNode {
	if m != nil {
		return m.Node
	}
	return ResourceNode{}
}

// QueryIndexingNodeRequest is request type for the Query/IndexingNode RPC method
type QueryIndexingNodeRequest struct {
	// network_addr defines the node address to query for.
	NetworkAddr string `protobuf:"bytes,1,opt,name=network_addr,json=networkAddr,proto3" json:"network_addr,omitempty"`
}

func (m *QueryIndexingNodeRequest) Reset()         { *m = QueryIndexingNodeRequest{} }
func (m *QueryIndexingNodeRequest) String() string { return proto.CompactTextString(m) }
func (*QueryIndexingNodeRequest) ProtoMessage()    {}
func (*QueryIndexingNodeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_59a612d1da8c0670, []int{2}
}
func (m *QueryIndexingNodeRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryIndexingNodeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryIndexingNodeRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryIndexingNodeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryIndexingNodeRequest.Merge(m, src)
}
func (m *QueryIndexingNodeRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryIndexingNodeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryIndexingNodeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryIndexingNodeRequest proto.InternalMessageInfo

func (m *QueryIndexingNodeRequest) GetNetworkAddr() string {
	if m != nil {
		return m.NetworkAddr
	}
	return ""
}

// QueryIndexingNodeResponse is response type for the Query/IndexingNode RPC method
type QueryIndexingNodeResponse struct {
	// node defines the the indexing info.
	Node IndexingNode `protobuf:"bytes,1,opt,name=node,proto3" json:"node"`
}

func (m *QueryIndexingNodeResponse) Reset()         { *m = QueryIndexingNodeResponse{} }
func (m *QueryIndexingNodeResponse) String() string { return proto.CompactTextString(m) }
func (*QueryIndexingNodeResponse) ProtoMessage()    {}
func (*QueryIndexingNodeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_59a612d1da8c0670, []int{3}
}
func (m *QueryIndexingNodeResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryIndexingNodeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryIndexingNodeResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryIndexingNodeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryIndexingNodeResponse.Merge(m, src)
}
func (m *QueryIndexingNodeResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryIndexingNodeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryIndexingNodeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryIndexingNodeResponse proto.InternalMessageInfo

func (m *QueryIndexingNodeResponse) GetNode() IndexingNode {
	if m != nil {
		return m.Node
	}
	return IndexingNode{}
}

func init() {
	proto.RegisterType((*QueryResourceNodeRequest)(nil), "stratos.register.v1.QueryResourceNodeRequest")
	proto.RegisterType((*QueryResourceNodeResponse)(nil), "stratos.register.v1.QueryResourceNodeResponse")
	proto.RegisterType((*QueryIndexingNodeRequest)(nil), "stratos.register.v1.QueryIndexingNodeRequest")
	proto.RegisterType((*QueryIndexingNodeResponse)(nil), "stratos.register.v1.QueryIndexingNodeResponse")
}

func init() { proto.RegisterFile("stratos/register/v1/query.proto", fileDescriptor_59a612d1da8c0670) }

var fileDescriptor_59a612d1da8c0670 = []byte{
	// 374 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0xcf, 0x4e, 0xea, 0x40,
	0x14, 0x87, 0x5b, 0xc2, 0xbd, 0xc9, 0x2d, 0xac, 0x7a, 0xef, 0x82, 0x4b, 0x4c, 0x91, 0xae, 0xdc,
	0xd0, 0x09, 0xc8, 0x4a, 0xe3, 0x42, 0x76, 0x6e, 0x48, 0xec, 0xca, 0xb8, 0x31, 0x85, 0x9e, 0x94,
	0x89, 0x3a, 0xa7, 0xcc, 0x4c, 0x11, 0x62, 0xdc, 0xf8, 0x04, 0x26, 0x3e, 0x85, 0x6f, 0xc2, 0x92,
	0xc4, 0x8d, 0x2b, 0x35, 0xe0, 0x83, 0x18, 0x86, 0x82, 0x35, 0x19, 0xf1, 0xcf, 0x6e, 0xd2, 0xf9,
	0x7d, 0xe7, 0x7c, 0xe7, 0x74, 0xac, 0x8a, 0x90, 0x3c, 0x90, 0x28, 0x08, 0x87, 0x88, 0x0a, 0x09,
	0x9c, 0x0c, 0xea, 0xa4, 0x9f, 0x00, 0x1f, 0x79, 0x31, 0x47, 0x89, 0xf6, 0xdf, 0x34, 0xe0, 0x2d,
	0x03, 0xde, 0xa0, 0x5e, 0xfe, 0x17, 0x61, 0x84, 0xea, 0x9e, 0xcc, 0x4f, 0x8b, 0x68, 0x79, 0x23,
	0x42, 0x8c, 0xce, 0x80, 0x04, 0x31, 0x25, 0x01, 0x63, 0x28, 0x03, 0x49, 0x91, 0x89, 0xf4, 0xd6,
	0xd5, 0x75, 0x5a, 0x15, 0x55, 0x19, 0x77, 0xcf, 0x2a, 0x1d, 0xce, 0x7b, 0xfb, 0x20, 0x30, 0xe1,
	0x5d, 0x68, 0x63, 0x08, 0x3e, 0xf4, 0x13, 0x10, 0xd2, 0xae, 0x5a, 0x45, 0x06, 0xf2, 0x02, 0xf9,
	0xe9, 0x49, 0x10, 0x86, 0xbc, 0x64, 0x6e, 0x9a, 0x5b, 0x7f, 0xfc, 0x42, 0xfa, 0x6d, 0x3f, 0x0c,
	0xb9, 0x7b, 0x64, 0xfd, 0xd7, 0xe0, 0x22, 0x46, 0x26, 0xc0, 0xde, 0xb5, 0xf2, 0x0c, 0x43, 0x50,
	0x5c, 0xa1, 0x51, 0xf5, 0x34, 0x73, 0x79, 0x59, 0xb0, 0x95, 0x1f, 0x3f, 0x56, 0x0c, 0x5f, 0x41,
	0x2b, 0xb1, 0x03, 0x16, 0xc2, 0x90, 0xb2, 0xe8, 0x87, 0x62, 0xef, 0xf1, 0x6f, 0x88, 0x65, 0xc1,
	0xac, 0x58, 0xe3, 0x29, 0x67, 0xfd, 0x52, 0xa5, 0xed, 0x3b, 0xd3, 0x2a, 0x66, 0xfd, 0xed, 0x9a,
	0xb6, 0xd2, 0x47, 0xfb, 0x2d, 0x7b, 0x5f, 0x8d, 0x2f, 0xb4, 0xdd, 0x9d, 0xeb, 0xfb, 0x97, 0xdb,
	0x5c, 0xd3, 0x6e, 0x10, 0xfd, 0x8f, 0x5d, 0x20, 0xb5, 0xb9, 0xa5, 0x20, 0x97, 0xd9, 0x0d, 0x5d,
	0x29, 0xd7, 0xec, 0x48, 0xeb, 0x5c, 0x35, 0x2b, 0x5f, 0xe7, 0xaa, 0x5b, 0xf1, 0x27, 0xae, 0x34,
	0x45, 0xb4, 0xae, 0xad, 0xf6, 0x78, 0xea, 0x98, 0x93, 0xa9, 0x63, 0x3e, 0x4f, 0x1d, 0xf3, 0x66,
	0xe6, 0x18, 0x93, 0x99, 0x63, 0x3c, 0xcc, 0x1c, 0xe3, 0xb8, 0x19, 0x51, 0xd9, 0x4b, 0x3a, 0x5e,
	0x17, 0xcf, 0x97, 0x75, 0x19, 0xc8, 0xe5, 0xb1, 0xd6, 0xed, 0x05, 0x94, 0x91, 0xe1, 0x5b, 0x2b,
	0x39, 0x8a, 0x41, 0x74, 0x7e, 0xab, 0xa7, 0xbe, 0xfd, 0x1a, 0x00, 0x00, 0xff, 0xff, 0x86, 0x3c,
	0x3f, 0x3a, 0x7a, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// ResourceNode queries ResourceNode info for given ResourceNode address.
	ResourceNode(ctx context.Context, in *QueryResourceNodeRequest, opts ...grpc.CallOption) (*QueryResourceNodeResponse, error)
	// IndexingNode queries IndexingNode info for given IndexingNode address.
	IndexingNode(ctx context.Context, in *QueryIndexingNodeRequest, opts ...grpc.CallOption) (*QueryIndexingNodeResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) ResourceNode(ctx context.Context, in *QueryResourceNodeRequest, opts ...grpc.CallOption) (*QueryResourceNodeResponse, error) {
	out := new(QueryResourceNodeResponse)
	err := c.cc.Invoke(ctx, "/stratos.register.v1.Query/ResourceNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) IndexingNode(ctx context.Context, in *QueryIndexingNodeRequest, opts ...grpc.CallOption) (*QueryIndexingNodeResponse, error) {
	out := new(QueryIndexingNodeResponse)
	err := c.cc.Invoke(ctx, "/stratos.register.v1.Query/IndexingNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// ResourceNode queries ResourceNode info for given ResourceNode address.
	ResourceNode(context.Context, *QueryResourceNodeRequest) (*QueryResourceNodeResponse, error)
	// IndexingNode queries IndexingNode info for given IndexingNode address.
	IndexingNode(context.Context, *QueryIndexingNodeRequest) (*QueryIndexingNodeResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) ResourceNode(ctx context.Context, req *QueryResourceNodeRequest) (*QueryResourceNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResourceNode not implemented")
}
func (*UnimplementedQueryServer) IndexingNode(ctx context.Context, req *QueryIndexingNodeRequest) (*QueryIndexingNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IndexingNode not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_ResourceNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryResourceNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).ResourceNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stratos.register.v1.Query/ResourceNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).ResourceNode(ctx, req.(*QueryResourceNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_IndexingNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryIndexingNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).IndexingNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stratos.register.v1.Query/IndexingNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).IndexingNode(ctx, req.(*QueryIndexingNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "stratos.register.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ResourceNode",
			Handler:    _Query_ResourceNode_Handler,
		},
		{
			MethodName: "IndexingNode",
			Handler:    _Query_IndexingNode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stratos/register/v1/query.proto",
}

func (m *QueryResourceNodeRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryResourceNodeRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryResourceNodeRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.NetworkAddr) > 0 {
		i -= len(m.NetworkAddr)
		copy(dAtA[i:], m.NetworkAddr)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.NetworkAddr)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryResourceNodeResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryResourceNodeResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryResourceNodeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Node.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryIndexingNodeRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryIndexingNodeRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryIndexingNodeRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.NetworkAddr) > 0 {
		i -= len(m.NetworkAddr)
		copy(dAtA[i:], m.NetworkAddr)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.NetworkAddr)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryIndexingNodeResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryIndexingNodeResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryIndexingNodeResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Node.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryResourceNodeRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.NetworkAddr)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryResourceNodeResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Node.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryIndexingNodeRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.NetworkAddr)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryIndexingNodeResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Node.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryResourceNodeRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryResourceNodeRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryResourceNodeRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NetworkAddr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NetworkAddr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryResourceNodeResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryResourceNodeResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryResourceNodeResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Node", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Node.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryIndexingNodeRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryIndexingNodeRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryIndexingNodeRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NetworkAddr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.NetworkAddr = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryIndexingNodeResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryIndexingNodeResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryIndexingNodeResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Node", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Node.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
