// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RpcNodeClient is the client API for RpcNode service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RpcNodeClient interface {
	SendRequest(ctx context.Context, in *TxRequest, opts ...grpc.CallOption) (*TxResponse, error)
	Subscribe(ctx context.Context, in *TxRequest, opts ...grpc.CallOption) (RpcNode_SubscribeClient, error)
	UpdateDebugConfig(ctx context.Context, in *DebugConfigRequest, opts ...grpc.CallOption) (*DebugConfigResponse, error)
	RefreshLogLevelsConfig(ctx context.Context, in *LogLevelsRequest, opts ...grpc.CallOption) (*LogLevelsResponse, error)
}

type rpcNodeClient struct {
	cc grpc.ClientConnInterface
}

func NewRpcNodeClient(cc grpc.ClientConnInterface) RpcNodeClient {
	return &rpcNodeClient{cc}
}

func (c *rpcNodeClient) SendRequest(ctx context.Context, in *TxRequest, opts ...grpc.CallOption) (*TxResponse, error) {
	out := new(TxResponse)
	err := c.cc.Invoke(ctx, "/pb.RpcNode/SendRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rpcNodeClient) Subscribe(ctx context.Context, in *TxRequest, opts ...grpc.CallOption) (RpcNode_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_RpcNode_serviceDesc.Streams[0], "/pb.RpcNode/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &rpcNodeSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RpcNode_SubscribeClient interface {
	Recv() (*SubscribeResult, error)
	grpc.ClientStream
}

type rpcNodeSubscribeClient struct {
	grpc.ClientStream
}

func (x *rpcNodeSubscribeClient) Recv() (*SubscribeResult, error) {
	m := new(SubscribeResult)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *rpcNodeClient) UpdateDebugConfig(ctx context.Context, in *DebugConfigRequest, opts ...grpc.CallOption) (*DebugConfigResponse, error) {
	out := new(DebugConfigResponse)
	err := c.cc.Invoke(ctx, "/pb.RpcNode/UpdateDebugConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rpcNodeClient) RefreshLogLevelsConfig(ctx context.Context, in *LogLevelsRequest, opts ...grpc.CallOption) (*LogLevelsResponse, error) {
	out := new(LogLevelsResponse)
	err := c.cc.Invoke(ctx, "/pb.RpcNode/RefreshLogLevelsConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RpcNodeServer is the server API for RpcNode service.
// All implementations must embed UnimplementedRpcNodeServer
// for forward compatibility
type RpcNodeServer interface {
	SendRequest(context.Context, *TxRequest) (*TxResponse, error)
	Subscribe(*TxRequest, RpcNode_SubscribeServer) error
	UpdateDebugConfig(context.Context, *DebugConfigRequest) (*DebugConfigResponse, error)
	RefreshLogLevelsConfig(context.Context, *LogLevelsRequest) (*LogLevelsResponse, error)
	mustEmbedUnimplementedRpcNodeServer()
}

// UnimplementedRpcNodeServer must be embedded to have forward compatible implementations.
type UnimplementedRpcNodeServer struct {
}

func (*UnimplementedRpcNodeServer) SendRequest(context.Context, *TxRequest) (*TxResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendRequest not implemented")
}
func (*UnimplementedRpcNodeServer) Subscribe(*TxRequest, RpcNode_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (*UnimplementedRpcNodeServer) UpdateDebugConfig(context.Context, *DebugConfigRequest) (*DebugConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDebugConfig not implemented")
}
func (*UnimplementedRpcNodeServer) RefreshLogLevelsConfig(context.Context, *LogLevelsRequest) (*LogLevelsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshLogLevelsConfig not implemented")
}
func (*UnimplementedRpcNodeServer) mustEmbedUnimplementedRpcNodeServer() {}

func RegisterRpcNodeServer(s *grpc.Server, srv RpcNodeServer) {
	s.RegisterService(&_RpcNode_serviceDesc, srv)
}

func _RpcNode_SendRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TxRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcNodeServer).SendRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.RpcNode/SendRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcNodeServer).SendRequest(ctx, req.(*TxRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RpcNode_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TxRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RpcNodeServer).Subscribe(m, &rpcNodeSubscribeServer{stream})
}

type RpcNode_SubscribeServer interface {
	Send(*SubscribeResult) error
	grpc.ServerStream
}

type rpcNodeSubscribeServer struct {
	grpc.ServerStream
}

func (x *rpcNodeSubscribeServer) Send(m *SubscribeResult) error {
	return x.ServerStream.SendMsg(m)
}

func _RpcNode_UpdateDebugConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DebugConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcNodeServer).UpdateDebugConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.RpcNode/UpdateDebugConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcNodeServer).UpdateDebugConfig(ctx, req.(*DebugConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RpcNode_RefreshLogLevelsConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogLevelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcNodeServer).RefreshLogLevelsConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.RpcNode/RefreshLogLevelsConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcNodeServer).RefreshLogLevelsConfig(ctx, req.(*LogLevelsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RpcNode_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.RpcNode",
	HandlerType: (*RpcNodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendRequest",
			Handler:    _RpcNode_SendRequest_Handler,
		},
		{
			MethodName: "UpdateDebugConfig",
			Handler:    _RpcNode_UpdateDebugConfig_Handler,
		},
		{
			MethodName: "RefreshLogLevelsConfig",
			Handler:    _RpcNode_RefreshLogLevelsConfig_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _RpcNode_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "rpc_node.proto",
}
