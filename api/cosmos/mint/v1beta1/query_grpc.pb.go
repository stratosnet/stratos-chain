// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package mintv1beta1

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

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QueryClient interface {
	// Params returns the total set of minting parameters.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// Inflation returns the current minting inflation value.
	Inflation(ctx context.Context, in *QueryInflationRequest, opts ...grpc.CallOption) (*QueryInflationResponse, error)
	// AnnualProvisions current minting annual provisions value.
	AnnualProvisions(ctx context.Context, in *QueryAnnualProvisionsRequest, opts ...grpc.CallOption) (*QueryAnnualProvisionsResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/cosmos.mint.v1beta1.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Inflation(ctx context.Context, in *QueryInflationRequest, opts ...grpc.CallOption) (*QueryInflationResponse, error) {
	out := new(QueryInflationResponse)
	err := c.cc.Invoke(ctx, "/cosmos.mint.v1beta1.Query/Inflation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) AnnualProvisions(ctx context.Context, in *QueryAnnualProvisionsRequest, opts ...grpc.CallOption) (*QueryAnnualProvisionsResponse, error) {
	out := new(QueryAnnualProvisionsResponse)
	err := c.cc.Invoke(ctx, "/cosmos.mint.v1beta1.Query/AnnualProvisions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
// All implementations must embed UnimplementedQueryServer
// for forward compatibility
type QueryServer interface {
	// Params returns the total set of minting parameters.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// Inflation returns the current minting inflation value.
	Inflation(context.Context, *QueryInflationRequest) (*QueryInflationResponse, error)
	// AnnualProvisions current minting annual provisions value.
	AnnualProvisions(context.Context, *QueryAnnualProvisionsRequest) (*QueryAnnualProvisionsResponse, error)
	mustEmbedUnimplementedQueryServer()
}

// UnimplementedQueryServer must be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (UnimplementedQueryServer) Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (UnimplementedQueryServer) Inflation(context.Context, *QueryInflationRequest) (*QueryInflationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Inflation not implemented")
}
func (UnimplementedQueryServer) AnnualProvisions(context.Context, *QueryAnnualProvisionsRequest) (*QueryAnnualProvisionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AnnualProvisions not implemented")
}
func (UnimplementedQueryServer) mustEmbedUnimplementedQueryServer() {}

// UnsafeQueryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueryServer will
// result in compilation errors.
type UnsafeQueryServer interface {
	mustEmbedUnimplementedQueryServer()
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&Query_ServiceDesc, srv)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cosmos.mint.v1beta1.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Inflation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryInflationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Inflation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cosmos.mint.v1beta1.Query/Inflation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Inflation(ctx, req.(*QueryInflationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_AnnualProvisions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAnnualProvisionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).AnnualProvisions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cosmos.mint.v1beta1.Query/AnnualProvisions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).AnnualProvisions(ctx, req.(*QueryAnnualProvisionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Query_ServiceDesc is the grpc.ServiceDesc for Query service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Query_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cosmos.mint.v1beta1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "Inflation",
			Handler:    _Query_Inflation_Handler,
		},
		{
			MethodName: "AnnualProvisions",
			Handler:    _Query_AnnualProvisions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cosmos/mint/v1beta1/query.proto",
}
