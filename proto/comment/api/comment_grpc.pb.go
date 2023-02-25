// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: comment/api/comment.proto

package api

import (
	context "context"
	request "douyin/proto/comment/request"
	response "douyin/proto/comment/response"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CommentClient is the client API for Comment service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommentClient interface {
	CommentVideo(ctx context.Context, in *request.DouyinCommentActionRequest, opts ...grpc.CallOption) (*response.DouyinCommentActionResponse, error)
	GetCommentVideoList(ctx context.Context, in *request.DouyinCommentListRequest, opts ...grpc.CallOption) (*response.DouyinCommentListResponse, error)
	PushVCommentBasicInfoInit(ctx context.Context, in *request.DouyinPushVCommentBasicInfoInitRequest, opts ...grpc.CallOption) (*response.DouyinPushVCommentBasicInfoInitResponse, error)
	GetCommentCount(ctx context.Context, in *request.DouyinCommentCountRequest, opts ...grpc.CallOption) (*response.DouyinCommentCountResponse, error)
}

type commentClient struct {
	cc grpc.ClientConnInterface
}

func NewCommentClient(cc grpc.ClientConnInterface) CommentClient {
	return &commentClient{cc}
}

func (c *commentClient) CommentVideo(ctx context.Context, in *request.DouyinCommentActionRequest, opts ...grpc.CallOption) (*response.DouyinCommentActionResponse, error) {
	out := new(response.DouyinCommentActionResponse)
	err := c.cc.Invoke(ctx, "/api.Comment/CommentVideo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentClient) GetCommentVideoList(ctx context.Context, in *request.DouyinCommentListRequest, opts ...grpc.CallOption) (*response.DouyinCommentListResponse, error) {
	out := new(response.DouyinCommentListResponse)
	err := c.cc.Invoke(ctx, "/api.Comment/GetCommentVideoList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentClient) PushVCommentBasicInfoInit(ctx context.Context, in *request.DouyinPushVCommentBasicInfoInitRequest, opts ...grpc.CallOption) (*response.DouyinPushVCommentBasicInfoInitResponse, error) {
	out := new(response.DouyinPushVCommentBasicInfoInitResponse)
	err := c.cc.Invoke(ctx, "/api.Comment/PushVCommentBasicInfoInit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *commentClient) GetCommentCount(ctx context.Context, in *request.DouyinCommentCountRequest, opts ...grpc.CallOption) (*response.DouyinCommentCountResponse, error) {
	out := new(response.DouyinCommentCountResponse)
	err := c.cc.Invoke(ctx, "/api.Comment/GetCommentCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommentServer is the server API for Comment service.
// All implementations must embed UnimplementedCommentServer
// for forward compatibility
type CommentServer interface {
	CommentVideo(context.Context, *request.DouyinCommentActionRequest) (*response.DouyinCommentActionResponse, error)
	GetCommentVideoList(context.Context, *request.DouyinCommentListRequest) (*response.DouyinCommentListResponse, error)
	PushVCommentBasicInfoInit(context.Context, *request.DouyinPushVCommentBasicInfoInitRequest) (*response.DouyinPushVCommentBasicInfoInitResponse, error)
	GetCommentCount(context.Context, *request.DouyinCommentCountRequest) (*response.DouyinCommentCountResponse, error)
	mustEmbedUnimplementedCommentServer()
}

// UnimplementedCommentServer must be embedded to have forward compatible implementations.
type UnimplementedCommentServer struct {
}

func (UnimplementedCommentServer) CommentVideo(context.Context, *request.DouyinCommentActionRequest) (*response.DouyinCommentActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommentVideo not implemented")
}
func (UnimplementedCommentServer) GetCommentVideoList(context.Context, *request.DouyinCommentListRequest) (*response.DouyinCommentListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCommentVideoList not implemented")
}
func (UnimplementedCommentServer) PushVCommentBasicInfoInit(context.Context, *request.DouyinPushVCommentBasicInfoInitRequest) (*response.DouyinPushVCommentBasicInfoInitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushVCommentBasicInfoInit not implemented")
}
func (UnimplementedCommentServer) GetCommentCount(context.Context, *request.DouyinCommentCountRequest) (*response.DouyinCommentCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCommentCount not implemented")
}
func (UnimplementedCommentServer) mustEmbedUnimplementedCommentServer() {}

// UnsafeCommentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommentServer will
// result in compilation errors.
type UnsafeCommentServer interface {
	mustEmbedUnimplementedCommentServer()
}

func RegisterCommentServer(s grpc.ServiceRegistrar, srv CommentServer) {
	s.RegisterService(&Comment_ServiceDesc, srv)
}

func _Comment_CommentVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinCommentActionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentServer).CommentVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Comment/CommentVideo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentServer).CommentVideo(ctx, req.(*request.DouyinCommentActionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Comment_GetCommentVideoList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinCommentListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentServer).GetCommentVideoList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Comment/GetCommentVideoList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentServer).GetCommentVideoList(ctx, req.(*request.DouyinCommentListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Comment_PushVCommentBasicInfoInit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinPushVCommentBasicInfoInitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentServer).PushVCommentBasicInfoInit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Comment/PushVCommentBasicInfoInit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentServer).PushVCommentBasicInfoInit(ctx, req.(*request.DouyinPushVCommentBasicInfoInitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Comment_GetCommentCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinCommentCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommentServer).GetCommentCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Comment/GetCommentCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommentServer).GetCommentCount(ctx, req.(*request.DouyinCommentCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Comment_ServiceDesc is the grpc.ServiceDesc for Comment service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Comment_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Comment",
	HandlerType: (*CommentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CommentVideo",
			Handler:    _Comment_CommentVideo_Handler,
		},
		{
			MethodName: "GetCommentVideoList",
			Handler:    _Comment_GetCommentVideoList_Handler,
		},
		{
			MethodName: "PushVCommentBasicInfoInit",
			Handler:    _Comment_PushVCommentBasicInfoInit_Handler,
		},
		{
			MethodName: "GetCommentCount",
			Handler:    _Comment_GetCommentCount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "comment/api/comment.proto",
}
