// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: video/api/video.proto

package api

import (
	context "context"
	request "douyin/proto/video/request"
	response "douyin/proto/video/response"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// VideoClient is the client API for Video service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VideoClient interface {
	PublishVideo(ctx context.Context, in *request.DouyinPublishActionRequest, opts ...grpc.CallOption) (*response.DouyinPublishActionResponse, error)
	GetVideoList(ctx context.Context, in *request.DouyinFeedRequest, opts ...grpc.CallOption) (*response.DouyinFeedResponse, error)
	GetPublishVideo(ctx context.Context, in *request.DouyinPublishListRequest, opts ...grpc.CallOption) (*response.DouyinPublishListResponse, error)
	GetUserPublishVideoCount(ctx context.Context, in *request.DouyinGetUserPublishVideoCountRequest, opts ...grpc.CallOption) (*response.DouyinGetUserPublishVideoCountResponse, error)
	GetUserPublishVideoCountList(ctx context.Context, in *request.DouyinGetUserPublishVideoCountListRequest, opts ...grpc.CallOption) (*response.DouyinGetUserPublishVideoCountListResponse, error)
	PushVBasicInfoInit(ctx context.Context, in *request.DouyinPushVInfoInitRequest, opts ...grpc.CallOption) (*response.DouyinPushVInfoInitResponse, error)
	PushActiveBasicInfoInit(ctx context.Context, in *request.DouyinPushActiveInfoInitRequest, opts ...grpc.CallOption) (*response.DouyinPushActiveInfoInitResponse, error)
	JudgeVideoAuthor(ctx context.Context, in *request.DouyinJudgeVideoAuthorRequest, opts ...grpc.CallOption) (*response.DouyinJudgeVideoAuthorResponse, error)
	PushVActiveFavoriteVideo(ctx context.Context, in *request.DouyinPushVActiveFavoriteVideoRequest, opts ...grpc.CallOption) (*response.DouyinPushVActiveFavoriteVideoResponse, error)
	PushVActiveFavoriteVideoRevert(ctx context.Context, in *request.DouyinPushVActiveFavoriteVideoRequest, opts ...grpc.CallOption) (*response.DouyinPushVActiveFavoriteVideoResponse, error)
	DeleteVActiveFavoriteVideo(ctx context.Context, in *request.DouyinDeleteVActiveFavoriteVideoRequest, opts ...grpc.CallOption) (*response.DouyinDeleteVActiveFavoriteVideoResponse, error)
	GetVActiveFavoriteVideo(ctx context.Context, in *request.DouyinGetVActiveFavoriteVideoRequest, opts ...grpc.CallOption) (*response.DouyinGetVActiveFavoriteVideoResponse, error)
	GetVideoListInner(ctx context.Context, in *request.DouyinGetVideoListRequest, opts ...grpc.CallOption) (*response.DouyinGetVideoListResponse, error)
}

type videoClient struct {
	cc grpc.ClientConnInterface
}

func NewVideoClient(cc grpc.ClientConnInterface) VideoClient {
	return &videoClient{cc}
}

func (c *videoClient) PublishVideo(ctx context.Context, in *request.DouyinPublishActionRequest, opts ...grpc.CallOption) (*response.DouyinPublishActionResponse, error) {
	out := new(response.DouyinPublishActionResponse)
	err := c.cc.Invoke(ctx, "/api.Video/PublishVideo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) GetVideoList(ctx context.Context, in *request.DouyinFeedRequest, opts ...grpc.CallOption) (*response.DouyinFeedResponse, error) {
	out := new(response.DouyinFeedResponse)
	err := c.cc.Invoke(ctx, "/api.Video/GetVideoList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) GetPublishVideo(ctx context.Context, in *request.DouyinPublishListRequest, opts ...grpc.CallOption) (*response.DouyinPublishListResponse, error) {
	out := new(response.DouyinPublishListResponse)
	err := c.cc.Invoke(ctx, "/api.Video/GetPublishVideo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) GetUserPublishVideoCount(ctx context.Context, in *request.DouyinGetUserPublishVideoCountRequest, opts ...grpc.CallOption) (*response.DouyinGetUserPublishVideoCountResponse, error) {
	out := new(response.DouyinGetUserPublishVideoCountResponse)
	err := c.cc.Invoke(ctx, "/api.Video/GetUserPublishVideoCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) GetUserPublishVideoCountList(ctx context.Context, in *request.DouyinGetUserPublishVideoCountListRequest, opts ...grpc.CallOption) (*response.DouyinGetUserPublishVideoCountListResponse, error) {
	out := new(response.DouyinGetUserPublishVideoCountListResponse)
	err := c.cc.Invoke(ctx, "/api.Video/GetUserPublishVideoCountList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) PushVBasicInfoInit(ctx context.Context, in *request.DouyinPushVInfoInitRequest, opts ...grpc.CallOption) (*response.DouyinPushVInfoInitResponse, error) {
	out := new(response.DouyinPushVInfoInitResponse)
	err := c.cc.Invoke(ctx, "/api.Video/PushVBasicInfoInit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) PushActiveBasicInfoInit(ctx context.Context, in *request.DouyinPushActiveInfoInitRequest, opts ...grpc.CallOption) (*response.DouyinPushActiveInfoInitResponse, error) {
	out := new(response.DouyinPushActiveInfoInitResponse)
	err := c.cc.Invoke(ctx, "/api.Video/PushActiveBasicInfoInit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) JudgeVideoAuthor(ctx context.Context, in *request.DouyinJudgeVideoAuthorRequest, opts ...grpc.CallOption) (*response.DouyinJudgeVideoAuthorResponse, error) {
	out := new(response.DouyinJudgeVideoAuthorResponse)
	err := c.cc.Invoke(ctx, "/api.Video/JudgeVideoAuthor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) PushVActiveFavoriteVideo(ctx context.Context, in *request.DouyinPushVActiveFavoriteVideoRequest, opts ...grpc.CallOption) (*response.DouyinPushVActiveFavoriteVideoResponse, error) {
	out := new(response.DouyinPushVActiveFavoriteVideoResponse)
	err := c.cc.Invoke(ctx, "/api.Video/PushVActiveFavoriteVideo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) PushVActiveFavoriteVideoRevert(ctx context.Context, in *request.DouyinPushVActiveFavoriteVideoRequest, opts ...grpc.CallOption) (*response.DouyinPushVActiveFavoriteVideoResponse, error) {
	out := new(response.DouyinPushVActiveFavoriteVideoResponse)
	err := c.cc.Invoke(ctx, "/api.Video/PushVActiveFavoriteVideoRevert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) DeleteVActiveFavoriteVideo(ctx context.Context, in *request.DouyinDeleteVActiveFavoriteVideoRequest, opts ...grpc.CallOption) (*response.DouyinDeleteVActiveFavoriteVideoResponse, error) {
	out := new(response.DouyinDeleteVActiveFavoriteVideoResponse)
	err := c.cc.Invoke(ctx, "/api.Video/DeleteVActiveFavoriteVideo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) GetVActiveFavoriteVideo(ctx context.Context, in *request.DouyinGetVActiveFavoriteVideoRequest, opts ...grpc.CallOption) (*response.DouyinGetVActiveFavoriteVideoResponse, error) {
	out := new(response.DouyinGetVActiveFavoriteVideoResponse)
	err := c.cc.Invoke(ctx, "/api.Video/GetVActiveFavoriteVideo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *videoClient) GetVideoListInner(ctx context.Context, in *request.DouyinGetVideoListRequest, opts ...grpc.CallOption) (*response.DouyinGetVideoListResponse, error) {
	out := new(response.DouyinGetVideoListResponse)
	err := c.cc.Invoke(ctx, "/api.Video/GetVideoListInner", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VideoServer is the server API for Video service.
// All implementations must embed UnimplementedVideoServer
// for forward compatibility
type VideoServer interface {
	PublishVideo(context.Context, *request.DouyinPublishActionRequest) (*response.DouyinPublishActionResponse, error)
	GetVideoList(context.Context, *request.DouyinFeedRequest) (*response.DouyinFeedResponse, error)
	GetPublishVideo(context.Context, *request.DouyinPublishListRequest) (*response.DouyinPublishListResponse, error)
	GetUserPublishVideoCount(context.Context, *request.DouyinGetUserPublishVideoCountRequest) (*response.DouyinGetUserPublishVideoCountResponse, error)
	GetUserPublishVideoCountList(context.Context, *request.DouyinGetUserPublishVideoCountListRequest) (*response.DouyinGetUserPublishVideoCountListResponse, error)
	PushVBasicInfoInit(context.Context, *request.DouyinPushVInfoInitRequest) (*response.DouyinPushVInfoInitResponse, error)
	PushActiveBasicInfoInit(context.Context, *request.DouyinPushActiveInfoInitRequest) (*response.DouyinPushActiveInfoInitResponse, error)
	JudgeVideoAuthor(context.Context, *request.DouyinJudgeVideoAuthorRequest) (*response.DouyinJudgeVideoAuthorResponse, error)
	PushVActiveFavoriteVideo(context.Context, *request.DouyinPushVActiveFavoriteVideoRequest) (*response.DouyinPushVActiveFavoriteVideoResponse, error)
	PushVActiveFavoriteVideoRevert(context.Context, *request.DouyinPushVActiveFavoriteVideoRequest) (*response.DouyinPushVActiveFavoriteVideoResponse, error)
	DeleteVActiveFavoriteVideo(context.Context, *request.DouyinDeleteVActiveFavoriteVideoRequest) (*response.DouyinDeleteVActiveFavoriteVideoResponse, error)
	GetVActiveFavoriteVideo(context.Context, *request.DouyinGetVActiveFavoriteVideoRequest) (*response.DouyinGetVActiveFavoriteVideoResponse, error)
	GetVideoListInner(context.Context, *request.DouyinGetVideoListRequest) (*response.DouyinGetVideoListResponse, error)
	mustEmbedUnimplementedVideoServer()
}

// UnimplementedVideoServer must be embedded to have forward compatible implementations.
type UnimplementedVideoServer struct {
}

func (UnimplementedVideoServer) PublishVideo(context.Context, *request.DouyinPublishActionRequest) (*response.DouyinPublishActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishVideo not implemented")
}
func (UnimplementedVideoServer) GetVideoList(context.Context, *request.DouyinFeedRequest) (*response.DouyinFeedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVideoList not implemented")
}
func (UnimplementedVideoServer) GetPublishVideo(context.Context, *request.DouyinPublishListRequest) (*response.DouyinPublishListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPublishVideo not implemented")
}
func (UnimplementedVideoServer) GetUserPublishVideoCount(context.Context, *request.DouyinGetUserPublishVideoCountRequest) (*response.DouyinGetUserPublishVideoCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserPublishVideoCount not implemented")
}
func (UnimplementedVideoServer) GetUserPublishVideoCountList(context.Context, *request.DouyinGetUserPublishVideoCountListRequest) (*response.DouyinGetUserPublishVideoCountListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserPublishVideoCountList not implemented")
}
func (UnimplementedVideoServer) PushVBasicInfoInit(context.Context, *request.DouyinPushVInfoInitRequest) (*response.DouyinPushVInfoInitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushVBasicInfoInit not implemented")
}
func (UnimplementedVideoServer) PushActiveBasicInfoInit(context.Context, *request.DouyinPushActiveInfoInitRequest) (*response.DouyinPushActiveInfoInitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushActiveBasicInfoInit not implemented")
}
func (UnimplementedVideoServer) JudgeVideoAuthor(context.Context, *request.DouyinJudgeVideoAuthorRequest) (*response.DouyinJudgeVideoAuthorResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JudgeVideoAuthor not implemented")
}
func (UnimplementedVideoServer) PushVActiveFavoriteVideo(context.Context, *request.DouyinPushVActiveFavoriteVideoRequest) (*response.DouyinPushVActiveFavoriteVideoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushVActiveFavoriteVideo not implemented")
}
func (UnimplementedVideoServer) PushVActiveFavoriteVideoRevert(context.Context, *request.DouyinPushVActiveFavoriteVideoRequest) (*response.DouyinPushVActiveFavoriteVideoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushVActiveFavoriteVideoRevert not implemented")
}
func (UnimplementedVideoServer) DeleteVActiveFavoriteVideo(context.Context, *request.DouyinDeleteVActiveFavoriteVideoRequest) (*response.DouyinDeleteVActiveFavoriteVideoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteVActiveFavoriteVideo not implemented")
}
func (UnimplementedVideoServer) GetVActiveFavoriteVideo(context.Context, *request.DouyinGetVActiveFavoriteVideoRequest) (*response.DouyinGetVActiveFavoriteVideoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVActiveFavoriteVideo not implemented")
}
func (UnimplementedVideoServer) GetVideoListInner(context.Context, *request.DouyinGetVideoListRequest) (*response.DouyinGetVideoListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVideoListInner not implemented")
}
func (UnimplementedVideoServer) mustEmbedUnimplementedVideoServer() {}

// UnsafeVideoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VideoServer will
// result in compilation errors.
type UnsafeVideoServer interface {
	mustEmbedUnimplementedVideoServer()
}

func RegisterVideoServer(s grpc.ServiceRegistrar, srv VideoServer) {
	s.RegisterService(&Video_ServiceDesc, srv)
}

func _Video_PublishVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinPublishActionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).PublishVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Video/PublishVideo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).PublishVideo(ctx, req.(*request.DouyinPublishActionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_GetVideoList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinFeedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).GetVideoList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Video/GetVideoList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).GetVideoList(ctx, req.(*request.DouyinFeedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_GetPublishVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinPublishListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).GetPublishVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Video/GetPublishVideo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).GetPublishVideo(ctx, req.(*request.DouyinPublishListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_GetUserPublishVideoCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinGetUserPublishVideoCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).GetUserPublishVideoCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Video/GetUserPublishVideoCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).GetUserPublishVideoCount(ctx, req.(*request.DouyinGetUserPublishVideoCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_GetUserPublishVideoCountList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinGetUserPublishVideoCountListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).GetUserPublishVideoCountList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Video/GetUserPublishVideoCountList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).GetUserPublishVideoCountList(ctx, req.(*request.DouyinGetUserPublishVideoCountListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_PushVBasicInfoInit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinPushVInfoInitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).PushVBasicInfoInit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Video/PushVBasicInfoInit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).PushVBasicInfoInit(ctx, req.(*request.DouyinPushVInfoInitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_PushActiveBasicInfoInit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinPushActiveInfoInitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).PushActiveBasicInfoInit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Video/PushActiveBasicInfoInit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).PushActiveBasicInfoInit(ctx, req.(*request.DouyinPushActiveInfoInitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_JudgeVideoAuthor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinJudgeVideoAuthorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).JudgeVideoAuthor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Video/JudgeVideoAuthor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).JudgeVideoAuthor(ctx, req.(*request.DouyinJudgeVideoAuthorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_PushVActiveFavoriteVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinPushVActiveFavoriteVideoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).PushVActiveFavoriteVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Video/PushVActiveFavoriteVideo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).PushVActiveFavoriteVideo(ctx, req.(*request.DouyinPushVActiveFavoriteVideoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_PushVActiveFavoriteVideoRevert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinPushVActiveFavoriteVideoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).PushVActiveFavoriteVideoRevert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Video/PushVActiveFavoriteVideoRevert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).PushVActiveFavoriteVideoRevert(ctx, req.(*request.DouyinPushVActiveFavoriteVideoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_DeleteVActiveFavoriteVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinDeleteVActiveFavoriteVideoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).DeleteVActiveFavoriteVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Video/DeleteVActiveFavoriteVideo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).DeleteVActiveFavoriteVideo(ctx, req.(*request.DouyinDeleteVActiveFavoriteVideoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_GetVActiveFavoriteVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinGetVActiveFavoriteVideoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).GetVActiveFavoriteVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Video/GetVActiveFavoriteVideo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).GetVActiveFavoriteVideo(ctx, req.(*request.DouyinGetVActiveFavoriteVideoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Video_GetVideoListInner_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinGetVideoListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideoServer).GetVideoListInner(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Video/GetVideoListInner",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideoServer).GetVideoListInner(ctx, req.(*request.DouyinGetVideoListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Video_ServiceDesc is the grpc.ServiceDesc for Video service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Video_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Video",
	HandlerType: (*VideoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PublishVideo",
			Handler:    _Video_PublishVideo_Handler,
		},
		{
			MethodName: "GetVideoList",
			Handler:    _Video_GetVideoList_Handler,
		},
		{
			MethodName: "GetPublishVideo",
			Handler:    _Video_GetPublishVideo_Handler,
		},
		{
			MethodName: "GetUserPublishVideoCount",
			Handler:    _Video_GetUserPublishVideoCount_Handler,
		},
		{
			MethodName: "GetUserPublishVideoCountList",
			Handler:    _Video_GetUserPublishVideoCountList_Handler,
		},
		{
			MethodName: "PushVBasicInfoInit",
			Handler:    _Video_PushVBasicInfoInit_Handler,
		},
		{
			MethodName: "PushActiveBasicInfoInit",
			Handler:    _Video_PushActiveBasicInfoInit_Handler,
		},
		{
			MethodName: "JudgeVideoAuthor",
			Handler:    _Video_JudgeVideoAuthor_Handler,
		},
		{
			MethodName: "PushVActiveFavoriteVideo",
			Handler:    _Video_PushVActiveFavoriteVideo_Handler,
		},
		{
			MethodName: "PushVActiveFavoriteVideoRevert",
			Handler:    _Video_PushVActiveFavoriteVideoRevert_Handler,
		},
		{
			MethodName: "DeleteVActiveFavoriteVideo",
			Handler:    _Video_DeleteVActiveFavoriteVideo_Handler,
		},
		{
			MethodName: "GetVActiveFavoriteVideo",
			Handler:    _Video_GetVActiveFavoriteVideo_Handler,
		},
		{
			MethodName: "GetVideoListInner",
			Handler:    _Video_GetVideoListInner_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "video/api/video.proto",
}
