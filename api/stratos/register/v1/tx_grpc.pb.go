// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: stratos/register/v1/tx.proto

package registerv1

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

const (
	Msg_HandleMsgCreateResourceNode_FullMethodName        = "/stratos.register.v1.Msg/HandleMsgCreateResourceNode"
	Msg_HandleMsgRemoveResourceNode_FullMethodName        = "/stratos.register.v1.Msg/HandleMsgRemoveResourceNode"
	Msg_HandleMsgUpdateResourceNode_FullMethodName        = "/stratos.register.v1.Msg/HandleMsgUpdateResourceNode"
	Msg_HandleMsgUpdateResourceNodeDeposit_FullMethodName = "/stratos.register.v1.Msg/HandleMsgUpdateResourceNodeDeposit"
	Msg_HandleMsgUpdateEffectiveDeposit_FullMethodName    = "/stratos.register.v1.Msg/HandleMsgUpdateEffectiveDeposit"
	Msg_HandleMsgCreateMetaNode_FullMethodName            = "/stratos.register.v1.Msg/HandleMsgCreateMetaNode"
	Msg_HandleMsgRemoveMetaNode_FullMethodName            = "/stratos.register.v1.Msg/HandleMsgRemoveMetaNode"
	Msg_HandleMsgUpdateMetaNode_FullMethodName            = "/stratos.register.v1.Msg/HandleMsgUpdateMetaNode"
	Msg_HandleMsgUpdateMetaNodeDeposit_FullMethodName     = "/stratos.register.v1.Msg/HandleMsgUpdateMetaNodeDeposit"
	Msg_HandleMsgMetaNodeRegistrationVote_FullMethodName  = "/stratos.register.v1.Msg/HandleMsgMetaNodeRegistrationVote"
	Msg_HandleMsgKickMetaNodeVote_FullMethodName          = "/stratos.register.v1.Msg/HandleMsgKickMetaNodeVote"
	Msg_UpdateParams_FullMethodName                       = "/stratos.register.v1.Msg/UpdateParams"
)

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	// CreateResourceNode defines a method for creating a new resource node.
	HandleMsgCreateResourceNode(ctx context.Context, in *MsgCreateResourceNode, opts ...grpc.CallOption) (*MsgCreateResourceNodeResponse, error)
	HandleMsgRemoveResourceNode(ctx context.Context, in *MsgRemoveResourceNode, opts ...grpc.CallOption) (*MsgRemoveResourceNodeResponse, error)
	HandleMsgUpdateResourceNode(ctx context.Context, in *MsgUpdateResourceNode, opts ...grpc.CallOption) (*MsgUpdateResourceNodeResponse, error)
	HandleMsgUpdateResourceNodeDeposit(ctx context.Context, in *MsgUpdateResourceNodeDeposit, opts ...grpc.CallOption) (*MsgUpdateResourceNodeDepositResponse, error)
	HandleMsgUpdateEffectiveDeposit(ctx context.Context, in *MsgUpdateEffectiveDeposit, opts ...grpc.CallOption) (*MsgUpdateEffectiveDepositResponse, error)
	HandleMsgCreateMetaNode(ctx context.Context, in *MsgCreateMetaNode, opts ...grpc.CallOption) (*MsgCreateMetaNodeResponse, error)
	HandleMsgRemoveMetaNode(ctx context.Context, in *MsgRemoveMetaNode, opts ...grpc.CallOption) (*MsgRemoveMetaNodeResponse, error)
	HandleMsgUpdateMetaNode(ctx context.Context, in *MsgUpdateMetaNode, opts ...grpc.CallOption) (*MsgUpdateMetaNodeResponse, error)
	HandleMsgUpdateMetaNodeDeposit(ctx context.Context, in *MsgUpdateMetaNodeDeposit, opts ...grpc.CallOption) (*MsgUpdateMetaNodeDepositResponse, error)
	HandleMsgMetaNodeRegistrationVote(ctx context.Context, in *MsgMetaNodeRegistrationVote, opts ...grpc.CallOption) (*MsgMetaNodeRegistrationVoteResponse, error)
	HandleMsgKickMetaNodeVote(ctx context.Context, in *MsgKickMetaNodeVote, opts ...grpc.CallOption) (*MsgKickMetaNodeVoteResponse, error)
	// UpdateParams defined a governance operation for updating the x/register module parameters.
	// The authority is hard-coded to the Cosmos SDK x/gov module account
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) HandleMsgCreateResourceNode(ctx context.Context, in *MsgCreateResourceNode, opts ...grpc.CallOption) (*MsgCreateResourceNodeResponse, error) {
	out := new(MsgCreateResourceNodeResponse)
	err := c.cc.Invoke(ctx, Msg_HandleMsgCreateResourceNode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) HandleMsgRemoveResourceNode(ctx context.Context, in *MsgRemoveResourceNode, opts ...grpc.CallOption) (*MsgRemoveResourceNodeResponse, error) {
	out := new(MsgRemoveResourceNodeResponse)
	err := c.cc.Invoke(ctx, Msg_HandleMsgRemoveResourceNode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) HandleMsgUpdateResourceNode(ctx context.Context, in *MsgUpdateResourceNode, opts ...grpc.CallOption) (*MsgUpdateResourceNodeResponse, error) {
	out := new(MsgUpdateResourceNodeResponse)
	err := c.cc.Invoke(ctx, Msg_HandleMsgUpdateResourceNode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) HandleMsgUpdateResourceNodeDeposit(ctx context.Context, in *MsgUpdateResourceNodeDeposit, opts ...grpc.CallOption) (*MsgUpdateResourceNodeDepositResponse, error) {
	out := new(MsgUpdateResourceNodeDepositResponse)
	err := c.cc.Invoke(ctx, Msg_HandleMsgUpdateResourceNodeDeposit_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) HandleMsgUpdateEffectiveDeposit(ctx context.Context, in *MsgUpdateEffectiveDeposit, opts ...grpc.CallOption) (*MsgUpdateEffectiveDepositResponse, error) {
	out := new(MsgUpdateEffectiveDepositResponse)
	err := c.cc.Invoke(ctx, Msg_HandleMsgUpdateEffectiveDeposit_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) HandleMsgCreateMetaNode(ctx context.Context, in *MsgCreateMetaNode, opts ...grpc.CallOption) (*MsgCreateMetaNodeResponse, error) {
	out := new(MsgCreateMetaNodeResponse)
	err := c.cc.Invoke(ctx, Msg_HandleMsgCreateMetaNode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) HandleMsgRemoveMetaNode(ctx context.Context, in *MsgRemoveMetaNode, opts ...grpc.CallOption) (*MsgRemoveMetaNodeResponse, error) {
	out := new(MsgRemoveMetaNodeResponse)
	err := c.cc.Invoke(ctx, Msg_HandleMsgRemoveMetaNode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) HandleMsgUpdateMetaNode(ctx context.Context, in *MsgUpdateMetaNode, opts ...grpc.CallOption) (*MsgUpdateMetaNodeResponse, error) {
	out := new(MsgUpdateMetaNodeResponse)
	err := c.cc.Invoke(ctx, Msg_HandleMsgUpdateMetaNode_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) HandleMsgUpdateMetaNodeDeposit(ctx context.Context, in *MsgUpdateMetaNodeDeposit, opts ...grpc.CallOption) (*MsgUpdateMetaNodeDepositResponse, error) {
	out := new(MsgUpdateMetaNodeDepositResponse)
	err := c.cc.Invoke(ctx, Msg_HandleMsgUpdateMetaNodeDeposit_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) HandleMsgMetaNodeRegistrationVote(ctx context.Context, in *MsgMetaNodeRegistrationVote, opts ...grpc.CallOption) (*MsgMetaNodeRegistrationVoteResponse, error) {
	out := new(MsgMetaNodeRegistrationVoteResponse)
	err := c.cc.Invoke(ctx, Msg_HandleMsgMetaNodeRegistrationVote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) HandleMsgKickMetaNodeVote(ctx context.Context, in *MsgKickMetaNodeVote, opts ...grpc.CallOption) (*MsgKickMetaNodeVoteResponse, error) {
	out := new(MsgKickMetaNodeVoteResponse)
	err := c.cc.Invoke(ctx, Msg_HandleMsgKickMetaNodeVote_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, Msg_UpdateParams_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	// CreateResourceNode defines a method for creating a new resource node.
	HandleMsgCreateResourceNode(context.Context, *MsgCreateResourceNode) (*MsgCreateResourceNodeResponse, error)
	HandleMsgRemoveResourceNode(context.Context, *MsgRemoveResourceNode) (*MsgRemoveResourceNodeResponse, error)
	HandleMsgUpdateResourceNode(context.Context, *MsgUpdateResourceNode) (*MsgUpdateResourceNodeResponse, error)
	HandleMsgUpdateResourceNodeDeposit(context.Context, *MsgUpdateResourceNodeDeposit) (*MsgUpdateResourceNodeDepositResponse, error)
	HandleMsgUpdateEffectiveDeposit(context.Context, *MsgUpdateEffectiveDeposit) (*MsgUpdateEffectiveDepositResponse, error)
	HandleMsgCreateMetaNode(context.Context, *MsgCreateMetaNode) (*MsgCreateMetaNodeResponse, error)
	HandleMsgRemoveMetaNode(context.Context, *MsgRemoveMetaNode) (*MsgRemoveMetaNodeResponse, error)
	HandleMsgUpdateMetaNode(context.Context, *MsgUpdateMetaNode) (*MsgUpdateMetaNodeResponse, error)
	HandleMsgUpdateMetaNodeDeposit(context.Context, *MsgUpdateMetaNodeDeposit) (*MsgUpdateMetaNodeDepositResponse, error)
	HandleMsgMetaNodeRegistrationVote(context.Context, *MsgMetaNodeRegistrationVote) (*MsgMetaNodeRegistrationVoteResponse, error)
	HandleMsgKickMetaNodeVote(context.Context, *MsgKickMetaNodeVote) (*MsgKickMetaNodeVoteResponse, error)
	// UpdateParams defined a governance operation for updating the x/register module parameters.
	// The authority is hard-coded to the Cosmos SDK x/gov module account
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) HandleMsgCreateResourceNode(context.Context, *MsgCreateResourceNode) (*MsgCreateResourceNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleMsgCreateResourceNode not implemented")
}
func (UnimplementedMsgServer) HandleMsgRemoveResourceNode(context.Context, *MsgRemoveResourceNode) (*MsgRemoveResourceNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleMsgRemoveResourceNode not implemented")
}
func (UnimplementedMsgServer) HandleMsgUpdateResourceNode(context.Context, *MsgUpdateResourceNode) (*MsgUpdateResourceNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleMsgUpdateResourceNode not implemented")
}
func (UnimplementedMsgServer) HandleMsgUpdateResourceNodeDeposit(context.Context, *MsgUpdateResourceNodeDeposit) (*MsgUpdateResourceNodeDepositResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleMsgUpdateResourceNodeDeposit not implemented")
}
func (UnimplementedMsgServer) HandleMsgUpdateEffectiveDeposit(context.Context, *MsgUpdateEffectiveDeposit) (*MsgUpdateEffectiveDepositResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleMsgUpdateEffectiveDeposit not implemented")
}
func (UnimplementedMsgServer) HandleMsgCreateMetaNode(context.Context, *MsgCreateMetaNode) (*MsgCreateMetaNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleMsgCreateMetaNode not implemented")
}
func (UnimplementedMsgServer) HandleMsgRemoveMetaNode(context.Context, *MsgRemoveMetaNode) (*MsgRemoveMetaNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleMsgRemoveMetaNode not implemented")
}
func (UnimplementedMsgServer) HandleMsgUpdateMetaNode(context.Context, *MsgUpdateMetaNode) (*MsgUpdateMetaNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleMsgUpdateMetaNode not implemented")
}
func (UnimplementedMsgServer) HandleMsgUpdateMetaNodeDeposit(context.Context, *MsgUpdateMetaNodeDeposit) (*MsgUpdateMetaNodeDepositResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleMsgUpdateMetaNodeDeposit not implemented")
}
func (UnimplementedMsgServer) HandleMsgMetaNodeRegistrationVote(context.Context, *MsgMetaNodeRegistrationVote) (*MsgMetaNodeRegistrationVoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleMsgMetaNodeRegistrationVote not implemented")
}
func (UnimplementedMsgServer) HandleMsgKickMetaNodeVote(context.Context, *MsgKickMetaNodeVote) (*MsgKickMetaNodeVoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleMsgKickMetaNodeVote not implemented")
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

func _Msg_HandleMsgCreateResourceNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateResourceNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).HandleMsgCreateResourceNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_HandleMsgCreateResourceNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).HandleMsgCreateResourceNode(ctx, req.(*MsgCreateResourceNode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_HandleMsgRemoveResourceNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRemoveResourceNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).HandleMsgRemoveResourceNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_HandleMsgRemoveResourceNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).HandleMsgRemoveResourceNode(ctx, req.(*MsgRemoveResourceNode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_HandleMsgUpdateResourceNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateResourceNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).HandleMsgUpdateResourceNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_HandleMsgUpdateResourceNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).HandleMsgUpdateResourceNode(ctx, req.(*MsgUpdateResourceNode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_HandleMsgUpdateResourceNodeDeposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateResourceNodeDeposit)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).HandleMsgUpdateResourceNodeDeposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_HandleMsgUpdateResourceNodeDeposit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).HandleMsgUpdateResourceNodeDeposit(ctx, req.(*MsgUpdateResourceNodeDeposit))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_HandleMsgUpdateEffectiveDeposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateEffectiveDeposit)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).HandleMsgUpdateEffectiveDeposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_HandleMsgUpdateEffectiveDeposit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).HandleMsgUpdateEffectiveDeposit(ctx, req.(*MsgUpdateEffectiveDeposit))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_HandleMsgCreateMetaNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateMetaNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).HandleMsgCreateMetaNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_HandleMsgCreateMetaNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).HandleMsgCreateMetaNode(ctx, req.(*MsgCreateMetaNode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_HandleMsgRemoveMetaNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRemoveMetaNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).HandleMsgRemoveMetaNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_HandleMsgRemoveMetaNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).HandleMsgRemoveMetaNode(ctx, req.(*MsgRemoveMetaNode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_HandleMsgUpdateMetaNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateMetaNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).HandleMsgUpdateMetaNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_HandleMsgUpdateMetaNode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).HandleMsgUpdateMetaNode(ctx, req.(*MsgUpdateMetaNode))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_HandleMsgUpdateMetaNodeDeposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateMetaNodeDeposit)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).HandleMsgUpdateMetaNodeDeposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_HandleMsgUpdateMetaNodeDeposit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).HandleMsgUpdateMetaNodeDeposit(ctx, req.(*MsgUpdateMetaNodeDeposit))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_HandleMsgMetaNodeRegistrationVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgMetaNodeRegistrationVote)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).HandleMsgMetaNodeRegistrationVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_HandleMsgMetaNodeRegistrationVote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).HandleMsgMetaNodeRegistrationVote(ctx, req.(*MsgMetaNodeRegistrationVote))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_HandleMsgKickMetaNodeVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgKickMetaNodeVote)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).HandleMsgKickMetaNodeVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Msg_HandleMsgKickMetaNodeVote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).HandleMsgKickMetaNodeVote(ctx, req.(*MsgKickMetaNodeVote))
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
		FullMethod: Msg_UpdateParams_FullMethodName,
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
	ServiceName: "stratos.register.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HandleMsgCreateResourceNode",
			Handler:    _Msg_HandleMsgCreateResourceNode_Handler,
		},
		{
			MethodName: "HandleMsgRemoveResourceNode",
			Handler:    _Msg_HandleMsgRemoveResourceNode_Handler,
		},
		{
			MethodName: "HandleMsgUpdateResourceNode",
			Handler:    _Msg_HandleMsgUpdateResourceNode_Handler,
		},
		{
			MethodName: "HandleMsgUpdateResourceNodeDeposit",
			Handler:    _Msg_HandleMsgUpdateResourceNodeDeposit_Handler,
		},
		{
			MethodName: "HandleMsgUpdateEffectiveDeposit",
			Handler:    _Msg_HandleMsgUpdateEffectiveDeposit_Handler,
		},
		{
			MethodName: "HandleMsgCreateMetaNode",
			Handler:    _Msg_HandleMsgCreateMetaNode_Handler,
		},
		{
			MethodName: "HandleMsgRemoveMetaNode",
			Handler:    _Msg_HandleMsgRemoveMetaNode_Handler,
		},
		{
			MethodName: "HandleMsgUpdateMetaNode",
			Handler:    _Msg_HandleMsgUpdateMetaNode_Handler,
		},
		{
			MethodName: "HandleMsgUpdateMetaNodeDeposit",
			Handler:    _Msg_HandleMsgUpdateMetaNodeDeposit_Handler,
		},
		{
			MethodName: "HandleMsgMetaNodeRegistrationVote",
			Handler:    _Msg_HandleMsgMetaNodeRegistrationVote_Handler,
		},
		{
			MethodName: "HandleMsgKickMetaNodeVote",
			Handler:    _Msg_HandleMsgKickMetaNodeVote_Handler,
		},
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stratos/register/v1/tx.proto",
}