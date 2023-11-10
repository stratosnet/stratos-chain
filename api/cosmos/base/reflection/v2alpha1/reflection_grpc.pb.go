// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package reflectionv2alpha1

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

// ReflectionServiceClient is the client API for ReflectionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReflectionServiceClient interface {
	// GetAuthnDescriptor returns information on how to authenticate transactions in the application
	// NOTE: this RPC is still experimental and might be subject to breaking changes or removal in
	// future releases of the cosmos-sdk.
	GetAuthnDescriptor(ctx context.Context, in *GetAuthnDescriptorRequest, opts ...grpc.CallOption) (*GetAuthnDescriptorResponse, error)
	// GetChainDescriptor returns the description of the chain
	GetChainDescriptor(ctx context.Context, in *GetChainDescriptorRequest, opts ...grpc.CallOption) (*GetChainDescriptorResponse, error)
	// GetCodecDescriptor returns the descriptor of the codec of the application
	GetCodecDescriptor(ctx context.Context, in *GetCodecDescriptorRequest, opts ...grpc.CallOption) (*GetCodecDescriptorResponse, error)
	// GetConfigurationDescriptor returns the descriptor for the sdk.Config of the application
	GetConfigurationDescriptor(ctx context.Context, in *GetConfigurationDescriptorRequest, opts ...grpc.CallOption) (*GetConfigurationDescriptorResponse, error)
	// GetQueryServicesDescriptor returns the available gRPC queryable services of the application
	GetQueryServicesDescriptor(ctx context.Context, in *GetQueryServicesDescriptorRequest, opts ...grpc.CallOption) (*GetQueryServicesDescriptorResponse, error)
	// GetTxDescriptor returns information on the used transaction object and available msgs that can be used
	GetTxDescriptor(ctx context.Context, in *GetTxDescriptorRequest, opts ...grpc.CallOption) (*GetTxDescriptorResponse, error)
}

type reflectionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReflectionServiceClient(cc grpc.ClientConnInterface) ReflectionServiceClient {
	return &reflectionServiceClient{cc}
}

func (c *reflectionServiceClient) GetAuthnDescriptor(ctx context.Context, in *GetAuthnDescriptorRequest, opts ...grpc.CallOption) (*GetAuthnDescriptorResponse, error) {
	out := new(GetAuthnDescriptorResponse)
	err := c.cc.Invoke(ctx, "/cosmos.base.reflection.v2alpha1.ReflectionService/GetAuthnDescriptor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reflectionServiceClient) GetChainDescriptor(ctx context.Context, in *GetChainDescriptorRequest, opts ...grpc.CallOption) (*GetChainDescriptorResponse, error) {
	out := new(GetChainDescriptorResponse)
	err := c.cc.Invoke(ctx, "/cosmos.base.reflection.v2alpha1.ReflectionService/GetChainDescriptor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reflectionServiceClient) GetCodecDescriptor(ctx context.Context, in *GetCodecDescriptorRequest, opts ...grpc.CallOption) (*GetCodecDescriptorResponse, error) {
	out := new(GetCodecDescriptorResponse)
	err := c.cc.Invoke(ctx, "/cosmos.base.reflection.v2alpha1.ReflectionService/GetCodecDescriptor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reflectionServiceClient) GetConfigurationDescriptor(ctx context.Context, in *GetConfigurationDescriptorRequest, opts ...grpc.CallOption) (*GetConfigurationDescriptorResponse, error) {
	out := new(GetConfigurationDescriptorResponse)
	err := c.cc.Invoke(ctx, "/cosmos.base.reflection.v2alpha1.ReflectionService/GetConfigurationDescriptor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reflectionServiceClient) GetQueryServicesDescriptor(ctx context.Context, in *GetQueryServicesDescriptorRequest, opts ...grpc.CallOption) (*GetQueryServicesDescriptorResponse, error) {
	out := new(GetQueryServicesDescriptorResponse)
	err := c.cc.Invoke(ctx, "/cosmos.base.reflection.v2alpha1.ReflectionService/GetQueryServicesDescriptor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reflectionServiceClient) GetTxDescriptor(ctx context.Context, in *GetTxDescriptorRequest, opts ...grpc.CallOption) (*GetTxDescriptorResponse, error) {
	out := new(GetTxDescriptorResponse)
	err := c.cc.Invoke(ctx, "/cosmos.base.reflection.v2alpha1.ReflectionService/GetTxDescriptor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReflectionServiceServer is the server API for ReflectionService service.
// All implementations must embed UnimplementedReflectionServiceServer
// for forward compatibility
type ReflectionServiceServer interface {
	// GetAuthnDescriptor returns information on how to authenticate transactions in the application
	// NOTE: this RPC is still experimental and might be subject to breaking changes or removal in
	// future releases of the cosmos-sdk.
	GetAuthnDescriptor(context.Context, *GetAuthnDescriptorRequest) (*GetAuthnDescriptorResponse, error)
	// GetChainDescriptor returns the description of the chain
	GetChainDescriptor(context.Context, *GetChainDescriptorRequest) (*GetChainDescriptorResponse, error)
	// GetCodecDescriptor returns the descriptor of the codec of the application
	GetCodecDescriptor(context.Context, *GetCodecDescriptorRequest) (*GetCodecDescriptorResponse, error)
	// GetConfigurationDescriptor returns the descriptor for the sdk.Config of the application
	GetConfigurationDescriptor(context.Context, *GetConfigurationDescriptorRequest) (*GetConfigurationDescriptorResponse, error)
	// GetQueryServicesDescriptor returns the available gRPC queryable services of the application
	GetQueryServicesDescriptor(context.Context, *GetQueryServicesDescriptorRequest) (*GetQueryServicesDescriptorResponse, error)
	// GetTxDescriptor returns information on the used transaction object and available msgs that can be used
	GetTxDescriptor(context.Context, *GetTxDescriptorRequest) (*GetTxDescriptorResponse, error)
	mustEmbedUnimplementedReflectionServiceServer()
}

// UnimplementedReflectionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedReflectionServiceServer struct {
}

func (UnimplementedReflectionServiceServer) GetAuthnDescriptor(context.Context, *GetAuthnDescriptorRequest) (*GetAuthnDescriptorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuthnDescriptor not implemented")
}
func (UnimplementedReflectionServiceServer) GetChainDescriptor(context.Context, *GetChainDescriptorRequest) (*GetChainDescriptorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetChainDescriptor not implemented")
}
func (UnimplementedReflectionServiceServer) GetCodecDescriptor(context.Context, *GetCodecDescriptorRequest) (*GetCodecDescriptorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCodecDescriptor not implemented")
}
func (UnimplementedReflectionServiceServer) GetConfigurationDescriptor(context.Context, *GetConfigurationDescriptorRequest) (*GetConfigurationDescriptorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfigurationDescriptor not implemented")
}
func (UnimplementedReflectionServiceServer) GetQueryServicesDescriptor(context.Context, *GetQueryServicesDescriptorRequest) (*GetQueryServicesDescriptorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetQueryServicesDescriptor not implemented")
}
func (UnimplementedReflectionServiceServer) GetTxDescriptor(context.Context, *GetTxDescriptorRequest) (*GetTxDescriptorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTxDescriptor not implemented")
}
func (UnimplementedReflectionServiceServer) mustEmbedUnimplementedReflectionServiceServer() {}

// UnsafeReflectionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReflectionServiceServer will
// result in compilation errors.
type UnsafeReflectionServiceServer interface {
	mustEmbedUnimplementedReflectionServiceServer()
}

func RegisterReflectionServiceServer(s grpc.ServiceRegistrar, srv ReflectionServiceServer) {
	s.RegisterService(&ReflectionService_ServiceDesc, srv)
}

func _ReflectionService_GetAuthnDescriptor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAuthnDescriptorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReflectionServiceServer).GetAuthnDescriptor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cosmos.base.reflection.v2alpha1.ReflectionService/GetAuthnDescriptor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReflectionServiceServer).GetAuthnDescriptor(ctx, req.(*GetAuthnDescriptorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReflectionService_GetChainDescriptor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetChainDescriptorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReflectionServiceServer).GetChainDescriptor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cosmos.base.reflection.v2alpha1.ReflectionService/GetChainDescriptor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReflectionServiceServer).GetChainDescriptor(ctx, req.(*GetChainDescriptorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReflectionService_GetCodecDescriptor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCodecDescriptorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReflectionServiceServer).GetCodecDescriptor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cosmos.base.reflection.v2alpha1.ReflectionService/GetCodecDescriptor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReflectionServiceServer).GetCodecDescriptor(ctx, req.(*GetCodecDescriptorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReflectionService_GetConfigurationDescriptor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetConfigurationDescriptorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReflectionServiceServer).GetConfigurationDescriptor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cosmos.base.reflection.v2alpha1.ReflectionService/GetConfigurationDescriptor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReflectionServiceServer).GetConfigurationDescriptor(ctx, req.(*GetConfigurationDescriptorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReflectionService_GetQueryServicesDescriptor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetQueryServicesDescriptorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReflectionServiceServer).GetQueryServicesDescriptor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cosmos.base.reflection.v2alpha1.ReflectionService/GetQueryServicesDescriptor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReflectionServiceServer).GetQueryServicesDescriptor(ctx, req.(*GetQueryServicesDescriptorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReflectionService_GetTxDescriptor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTxDescriptorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReflectionServiceServer).GetTxDescriptor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cosmos.base.reflection.v2alpha1.ReflectionService/GetTxDescriptor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReflectionServiceServer).GetTxDescriptor(ctx, req.(*GetTxDescriptorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ReflectionService_ServiceDesc is the grpc.ServiceDesc for ReflectionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReflectionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cosmos.base.reflection.v2alpha1.ReflectionService",
	HandlerType: (*ReflectionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAuthnDescriptor",
			Handler:    _ReflectionService_GetAuthnDescriptor_Handler,
		},
		{
			MethodName: "GetChainDescriptor",
			Handler:    _ReflectionService_GetChainDescriptor_Handler,
		},
		{
			MethodName: "GetCodecDescriptor",
			Handler:    _ReflectionService_GetCodecDescriptor_Handler,
		},
		{
			MethodName: "GetConfigurationDescriptor",
			Handler:    _ReflectionService_GetConfigurationDescriptor_Handler,
		},
		{
			MethodName: "GetQueryServicesDescriptor",
			Handler:    _ReflectionService_GetQueryServicesDescriptor_Handler,
		},
		{
			MethodName: "GetTxDescriptor",
			Handler:    _ReflectionService_GetTxDescriptor_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cosmos/base/reflection/v2alpha1/reflection.proto",
}
