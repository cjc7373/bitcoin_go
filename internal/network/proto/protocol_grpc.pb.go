// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: internal/network/proto/protocol.proto

package proto

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
	Discovery_RequestNodes_FullMethodName   = "/Discovery/RequestNodes"
	Discovery_BroadcastNodes_FullMethodName = "/Discovery/BroadcastNodes"
)

// DiscoveryClient is the client API for Discovery service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DiscoveryClient interface {
	RequestNodes(ctx context.Context, in *Node, opts ...grpc.CallOption) (*Empty, error)
	BroadcastNodes(ctx context.Context, in *NodeBroadcast, opts ...grpc.CallOption) (*Empty, error)
}

type discoveryClient struct {
	cc grpc.ClientConnInterface
}

func NewDiscoveryClient(cc grpc.ClientConnInterface) DiscoveryClient {
	return &discoveryClient{cc}
}

func (c *discoveryClient) RequestNodes(ctx context.Context, in *Node, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, Discovery_RequestNodes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *discoveryClient) BroadcastNodes(ctx context.Context, in *NodeBroadcast, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, Discovery_BroadcastNodes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DiscoveryServer is the server API for Discovery service.
// All implementations must embed UnimplementedDiscoveryServer
// for forward compatibility
type DiscoveryServer interface {
	RequestNodes(context.Context, *Node) (*Empty, error)
	BroadcastNodes(context.Context, *NodeBroadcast) (*Empty, error)
	mustEmbedUnimplementedDiscoveryServer()
}

// UnimplementedDiscoveryServer must be embedded to have forward compatible implementations.
type UnimplementedDiscoveryServer struct {
}

func (UnimplementedDiscoveryServer) RequestNodes(context.Context, *Node) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestNodes not implemented")
}
func (UnimplementedDiscoveryServer) BroadcastNodes(context.Context, *NodeBroadcast) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BroadcastNodes not implemented")
}
func (UnimplementedDiscoveryServer) mustEmbedUnimplementedDiscoveryServer() {}

// UnsafeDiscoveryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DiscoveryServer will
// result in compilation errors.
type UnsafeDiscoveryServer interface {
	mustEmbedUnimplementedDiscoveryServer()
}

func RegisterDiscoveryServer(s grpc.ServiceRegistrar, srv DiscoveryServer) {
	s.RegisterService(&Discovery_ServiceDesc, srv)
}

func _Discovery_RequestNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Node)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscoveryServer).RequestNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Discovery_RequestNodes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscoveryServer).RequestNodes(ctx, req.(*Node))
	}
	return interceptor(ctx, in, info, handler)
}

func _Discovery_BroadcastNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodeBroadcast)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DiscoveryServer).BroadcastNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Discovery_BroadcastNodes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DiscoveryServer).BroadcastNodes(ctx, req.(*NodeBroadcast))
	}
	return interceptor(ctx, in, info, handler)
}

// Discovery_ServiceDesc is the grpc.ServiceDesc for Discovery service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Discovery_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Discovery",
	HandlerType: (*DiscoveryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RequestNodes",
			Handler:    _Discovery_RequestNodes_Handler,
		},
		{
			MethodName: "BroadcastNodes",
			Handler:    _Discovery_BroadcastNodes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/network/proto/protocol.proto",
}
