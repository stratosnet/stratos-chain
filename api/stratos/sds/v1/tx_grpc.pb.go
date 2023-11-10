// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package sdsv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	// HandleMsgFileUpload defines a method for file uploading
	HandleMsgFileUpload(ctx context.Context, in *MsgFileUpload, opts ...grpc.CallOption) (*MsgFileUploadResponse, error)
	// HandleMsgPrepay defines a method for prepay
	HandleMsgPrepay(ctx context.Context, in *MsgPrepay, opts ...grpc.CallOption) (*MsgPrepayResponse, error)
	// UpdateParams defined a governance operation for updating the x/sds module parameters.
	// The authority is hard-coded to the Cosmos SDK x/gov module account
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) HandleMsgFileUpload(ctx context.Context, in *MsgFileUpload, opts ...grpc.CallOption) (*MsgFileUploadResponse, error) {
	out := new(MsgFileUploadResponse)
	err := c.cc.Invoke(ctx, "/stratos.sds.v1.Msg/HandleMsgFileUpload", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) HandleMsgPrepay(ctx context.Context, in *MsgPrepay, opts ...grpc.CallOption) (*MsgPrepayResponse, error) {
	out := new(MsgPrepayResponse)
	err := c.cc.Invoke(ctx, "/stratos.sds.v1.Msg/HandleMsgPrepay", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, "/stratos.sds.v1.Msg/UpdateParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	// HandleMsgFileUpload defines a method for file uploading
	HandleMsgFileUpload(context.Context, *MsgFileUpload) (*MsgFileUploadResponse, error)
	// HandleMsgPrepay defines a method for prepay
	HandleMsgPrepay(context.Context, *MsgPrepay) (*MsgPrepayResponse, error)
	// UpdateParams defined a governance operation for updating the x/sds module parameters.
	// The authority is hard-coded to the Cosmos SDK x/gov module account
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) HandleMsgFileUpload(context.Context, *MsgFileUpload) (*MsgFileUploadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleMsgFileUpload not implemented")
}
func (UnimplementedMsgServer) HandleMsgPrepay(context.Context, *MsgPrepay) (*MsgPrepayResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleMsgPrepay not implemented")
}
func (UnimplementedMsgServer) UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_HandleMsgFileUpload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgFileUpload)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).HandleMsgFileUpload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stratos.sds.v1.Msg/HandleMsgFileUpload",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).HandleMsgFileUpload(ctx, req.(*MsgFileUpload))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_HandleMsgPrepay_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgPrepay)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).HandleMsgPrepay(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stratos.sds.v1.Msg/HandleMsgPrepay",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).HandleMsgPrepay(ctx, req.(*MsgPrepay))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/stratos.sds.v1.Msg/UpdateParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "stratos.sds.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HandleMsgFileUpload",
			Handler:    _Msg_HandleMsgFileUpload_Handler,
		},
		{
			MethodName: "HandleMsgPrepay",
			Handler:    _Msg_HandleMsgPrepay_Handler,
		},
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stratos/sds/v1/tx.proto",
}
