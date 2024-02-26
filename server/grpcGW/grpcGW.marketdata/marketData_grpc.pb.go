// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.3
// source: marketData.proto

package grpcGW_marketdata

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
	MarketData_GetCandles_FullMethodName       = "/marketData.MarketData/GetCandles"
	MarketData_SubscribeCandles_FullMethodName = "/marketData.MarketData/SubscribeCandles"
)

// MarketDataClient is the client API for MarketData service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MarketDataClient interface {
	// Название нашего эндпоинта
	GetCandles(ctx context.Context, in *GetCandlesRequest, opts ...grpc.CallOption) (*GetCandlesResponse, error)
	SubscribeCandles(ctx context.Context, in *SubscribeCandlesRequest, opts ...grpc.CallOption) (MarketData_SubscribeCandlesClient, error)
}

type marketDataClient struct {
	cc grpc.ClientConnInterface
}

func NewMarketDataClient(cc grpc.ClientConnInterface) MarketDataClient {
	return &marketDataClient{cc}
}

func (c *marketDataClient) GetCandles(ctx context.Context, in *GetCandlesRequest, opts ...grpc.CallOption) (*GetCandlesResponse, error) {
	out := new(GetCandlesResponse)
	err := c.cc.Invoke(ctx, MarketData_GetCandles_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketDataClient) SubscribeCandles(ctx context.Context, in *SubscribeCandlesRequest, opts ...grpc.CallOption) (MarketData_SubscribeCandlesClient, error) {
	stream, err := c.cc.NewStream(ctx, &MarketData_ServiceDesc.Streams[0], MarketData_SubscribeCandles_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &marketDataSubscribeCandlesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MarketData_SubscribeCandlesClient interface {
	Recv() (*OHLC, error)
	grpc.ClientStream
}

type marketDataSubscribeCandlesClient struct {
	grpc.ClientStream
}

func (x *marketDataSubscribeCandlesClient) Recv() (*OHLC, error) {
	m := new(OHLC)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MarketDataServer is the server API for MarketData service.
// All implementations must embed UnimplementedMarketDataServer
// for forward compatibility
type MarketDataServer interface {
	// Название нашего эндпоинта
	GetCandles(context.Context, *GetCandlesRequest) (*GetCandlesResponse, error)
	SubscribeCandles(*SubscribeCandlesRequest, MarketData_SubscribeCandlesServer) error
	mustEmbedUnimplementedMarketDataServer()
}

// UnimplementedMarketDataServer must be embedded to have forward compatible implementations.
type UnimplementedMarketDataServer struct {
}

func (UnimplementedMarketDataServer) GetCandles(context.Context, *GetCandlesRequest) (*GetCandlesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCandles not implemented")
}
func (UnimplementedMarketDataServer) SubscribeCandles(*SubscribeCandlesRequest, MarketData_SubscribeCandlesServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeCandles not implemented")
}
func (UnimplementedMarketDataServer) mustEmbedUnimplementedMarketDataServer() {}

// UnsafeMarketDataServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MarketDataServer will
// result in compilation errors.
type UnsafeMarketDataServer interface {
	mustEmbedUnimplementedMarketDataServer()
}

func RegisterMarketDataServer(s grpc.ServiceRegistrar, srv MarketDataServer) {
	s.RegisterService(&MarketData_ServiceDesc, srv)
}

func _MarketData_GetCandles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCandlesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketDataServer).GetCandles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MarketData_GetCandles_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketDataServer).GetCandles(ctx, req.(*GetCandlesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MarketData_SubscribeCandles_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeCandlesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MarketDataServer).SubscribeCandles(m, &marketDataSubscribeCandlesServer{stream})
}

type MarketData_SubscribeCandlesServer interface {
	Send(*OHLC) error
	grpc.ServerStream
}

type marketDataSubscribeCandlesServer struct {
	grpc.ServerStream
}

func (x *marketDataSubscribeCandlesServer) Send(m *OHLC) error {
	return x.ServerStream.SendMsg(m)
}

// MarketData_ServiceDesc is the grpc.ServiceDesc for MarketData service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MarketData_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "marketData.MarketData",
	HandlerType: (*MarketDataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCandles",
			Handler:    _MarketData_GetCandles_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeCandles",
			Handler:       _MarketData_SubscribeCandles_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "marketData.proto",
}
