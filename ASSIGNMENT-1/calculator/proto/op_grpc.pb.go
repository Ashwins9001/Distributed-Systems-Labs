// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.14.0
// source: proto/op.proto

package op_proto

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

// OperationClient is the client API for Operation service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OperationClient interface {
	// Sends a greeting
	DoOp(ctx context.Context, in *OpRequest, opts ...grpc.CallOption) (*OpReply, error)
}

type operationClient struct {
	cc grpc.ClientConnInterface
}

func NewOperationClient(cc grpc.ClientConnInterface) OperationClient {
	return &operationClient{cc}
}

func (c *operationClient) DoOp(ctx context.Context, in *OpRequest, opts ...grpc.CallOption) (*OpReply, error) {
	out := new(OpReply)
	err := c.cc.Invoke(ctx, "/op.Operation/DoOp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OperationServer is the server API for Operation service.
// All implementations must embed UnimplementedOperationServer
// for forward compatibility
type OperationServer interface {
	// Sends a greeting
	DoOp(context.Context, *OpRequest) (*OpReply, error)
	mustEmbedUnimplementedOperationServer()
}

// UnimplementedOperationServer must be embedded to have forward compatible implementations.
type UnimplementedOperationServer struct {
}

func (UnimplementedOperationServer) DoOp(context.Context, *OpRequest) (*OpReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DoOp not implemented")
}
func (UnimplementedOperationServer) mustEmbedUnimplementedOperationServer() {}

// UnsafeOperationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OperationServer will
// result in compilation errors.
type UnsafeOperationServer interface {
	mustEmbedUnimplementedOperationServer()
}

func RegisterOperationServer(s grpc.ServiceRegistrar, srv OperationServer) {
	s.RegisterService(&Operation_ServiceDesc, srv)
}

func _Operation_DoOp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationServer).DoOp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/op.Operation/DoOp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationServer).DoOp(ctx, req.(*OpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Operation_ServiceDesc is the grpc.ServiceDesc for Operation service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Operation_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "op.Operation",
	HandlerType: (*OperationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DoOp",
			Handler:    _Operation_DoOp_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/op.proto",
}
