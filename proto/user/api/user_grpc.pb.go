// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: user/api/user.proto

package api

import (
	context "context"
	request "douyin/proto/user/request"
	response "douyin/proto/user/response"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UserClient is the client API for User service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserClient interface {
	UserRegister(ctx context.Context, in *request.DouyinUserRegisterRequest, opts ...grpc.CallOption) (*response.DouyinUserRegisterResponse, error)
	UserLogin(ctx context.Context, in *request.DouyinUserLoginRequest, opts ...grpc.CallOption) (*response.DouyinUserLoginResponse, error)
	GetUserInfo(ctx context.Context, in *request.DouyinUserRequest, opts ...grpc.CallOption) (*response.DouyinUserResponse, error)
	GetUserInfoList(ctx context.Context, in *request.DouyinUserListRequest, opts ...grpc.CallOption) (*response.DouyinUserListResponse, error)
	PushVUserRelativeInfoInit(ctx context.Context, in *request.DouyinPushVSetRequest, opts ...grpc.CallOption) (*response.DouyinPushVSetResponse, error)
	UserIsInfluencerActiver(ctx context.Context, in *request.DouyinUserIsInfluencerActiverRequest, opts ...grpc.CallOption) (*response.DouyinUserIsInfluencerActiverResponse, error)
	AddUserVideoCountSet(ctx context.Context, in *request.DouyinUserVideoCountSetRequest, opts ...grpc.CallOption) (*response.DouyinUserVideoCountSetResponse, error)
	AddUserVideoCountSetRevert(ctx context.Context, in *request.DouyinUserVideoCountSetRequest, opts ...grpc.CallOption) (*response.DouyinUserVideoCountSetResponse, error)
	AddUserFavoriteVideoCountSet(ctx context.Context, in *request.DouyinUserVideoCountSetRequest, opts ...grpc.CallOption) (*response.DouyinUserVideoCountSetResponse, error)
	AddUserFavoriteVideoCountSetRevert(ctx context.Context, in *request.DouyinUserVideoCountSetRequest, opts ...grpc.CallOption) (*response.DouyinUserVideoCountSetResponse, error)
	AddUserFollowUserCountSet(ctx context.Context, in *request.DouyinUserFollowCountSetRequest, opts ...grpc.CallOption) (*response.DouyinUserFollowCountSetResponse, error)
	AddUserFollowerUserCountSet(ctx context.Context, in *request.DouyinUserFollowerCountSetRequest, opts ...grpc.CallOption) (*response.DouyinUserFollowerCountSetResponse, error)
	AddUserFollowUserCountSetRevert(ctx context.Context, in *request.DouyinUserFollowCountSetRequest, opts ...grpc.CallOption) (*response.DouyinUserFollowCountSetResponse, error)
	AddUserFollowerUserCountSetRevert(ctx context.Context, in *request.DouyinUserFollowerCountSetRequest, opts ...grpc.CallOption) (*response.DouyinUserFollowerCountSetResponse, error)
	PushVActiveFollowerUserinfo(ctx context.Context, in *request.DouyinVActiveFollowFollowerUserinfoRequest, opts ...grpc.CallOption) (*response.DouyinVActiveFollowFollowerUserinfoResponse, error)
	PushVActiveFollowerUserinfoRevert(ctx context.Context, in *request.DouyinVActiveFollowFollowerUserinfoRequest, opts ...grpc.CallOption) (*response.DouyinVActiveFollowFollowerUserinfoResponse, error)
	DeleteVActiveFollowerUserinfo(ctx context.Context, in *request.DouyinDeleteVActiveFollowFollowerUserinfoRequest, opts ...grpc.CallOption) (*response.DouyinDeleteVActiveFollowFollowerUserinfoResponse, error)
	GetVActiveFollowerUserinfo(ctx context.Context, in *request.DouyinGetVActiveFollowFollowerUserinfoRequest, opts ...grpc.CallOption) (*response.DouyinGetVActiveFollowFollowerUserinfoResponse, error)
}

type userClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClient(cc grpc.ClientConnInterface) UserClient {
	return &userClient{cc}
}

func (c *userClient) UserRegister(ctx context.Context, in *request.DouyinUserRegisterRequest, opts ...grpc.CallOption) (*response.DouyinUserRegisterResponse, error) {
	out := new(response.DouyinUserRegisterResponse)
	err := c.cc.Invoke(ctx, "/api.User/UserRegister", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) UserLogin(ctx context.Context, in *request.DouyinUserLoginRequest, opts ...grpc.CallOption) (*response.DouyinUserLoginResponse, error) {
	out := new(response.DouyinUserLoginResponse)
	err := c.cc.Invoke(ctx, "/api.User/UserLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetUserInfo(ctx context.Context, in *request.DouyinUserRequest, opts ...grpc.CallOption) (*response.DouyinUserResponse, error) {
	out := new(response.DouyinUserResponse)
	err := c.cc.Invoke(ctx, "/api.User/GetUserInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetUserInfoList(ctx context.Context, in *request.DouyinUserListRequest, opts ...grpc.CallOption) (*response.DouyinUserListResponse, error) {
	out := new(response.DouyinUserListResponse)
	err := c.cc.Invoke(ctx, "/api.User/GetUserInfoList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) PushVUserRelativeInfoInit(ctx context.Context, in *request.DouyinPushVSetRequest, opts ...grpc.CallOption) (*response.DouyinPushVSetResponse, error) {
	out := new(response.DouyinPushVSetResponse)
	err := c.cc.Invoke(ctx, "/api.User/PushVUserRelativeInfoInit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) UserIsInfluencerActiver(ctx context.Context, in *request.DouyinUserIsInfluencerActiverRequest, opts ...grpc.CallOption) (*response.DouyinUserIsInfluencerActiverResponse, error) {
	out := new(response.DouyinUserIsInfluencerActiverResponse)
	err := c.cc.Invoke(ctx, "/api.User/UserIsInfluencerActiver", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) AddUserVideoCountSet(ctx context.Context, in *request.DouyinUserVideoCountSetRequest, opts ...grpc.CallOption) (*response.DouyinUserVideoCountSetResponse, error) {
	out := new(response.DouyinUserVideoCountSetResponse)
	err := c.cc.Invoke(ctx, "/api.User/AddUserVideoCountSet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) AddUserVideoCountSetRevert(ctx context.Context, in *request.DouyinUserVideoCountSetRequest, opts ...grpc.CallOption) (*response.DouyinUserVideoCountSetResponse, error) {
	out := new(response.DouyinUserVideoCountSetResponse)
	err := c.cc.Invoke(ctx, "/api.User/AddUserVideoCountSetRevert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) AddUserFavoriteVideoCountSet(ctx context.Context, in *request.DouyinUserVideoCountSetRequest, opts ...grpc.CallOption) (*response.DouyinUserVideoCountSetResponse, error) {
	out := new(response.DouyinUserVideoCountSetResponse)
	err := c.cc.Invoke(ctx, "/api.User/AddUserFavoriteVideoCountSet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) AddUserFavoriteVideoCountSetRevert(ctx context.Context, in *request.DouyinUserVideoCountSetRequest, opts ...grpc.CallOption) (*response.DouyinUserVideoCountSetResponse, error) {
	out := new(response.DouyinUserVideoCountSetResponse)
	err := c.cc.Invoke(ctx, "/api.User/AddUserFavoriteVideoCountSetRevert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) AddUserFollowUserCountSet(ctx context.Context, in *request.DouyinUserFollowCountSetRequest, opts ...grpc.CallOption) (*response.DouyinUserFollowCountSetResponse, error) {
	out := new(response.DouyinUserFollowCountSetResponse)
	err := c.cc.Invoke(ctx, "/api.User/AddUserFollowUserCountSet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) AddUserFollowerUserCountSet(ctx context.Context, in *request.DouyinUserFollowerCountSetRequest, opts ...grpc.CallOption) (*response.DouyinUserFollowerCountSetResponse, error) {
	out := new(response.DouyinUserFollowerCountSetResponse)
	err := c.cc.Invoke(ctx, "/api.User/AddUserFollowerUserCountSet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) AddUserFollowUserCountSetRevert(ctx context.Context, in *request.DouyinUserFollowCountSetRequest, opts ...grpc.CallOption) (*response.DouyinUserFollowCountSetResponse, error) {
	out := new(response.DouyinUserFollowCountSetResponse)
	err := c.cc.Invoke(ctx, "/api.User/AddUserFollowUserCountSetRevert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) AddUserFollowerUserCountSetRevert(ctx context.Context, in *request.DouyinUserFollowerCountSetRequest, opts ...grpc.CallOption) (*response.DouyinUserFollowerCountSetResponse, error) {
	out := new(response.DouyinUserFollowerCountSetResponse)
	err := c.cc.Invoke(ctx, "/api.User/AddUserFollowerUserCountSetRevert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) PushVActiveFollowerUserinfo(ctx context.Context, in *request.DouyinVActiveFollowFollowerUserinfoRequest, opts ...grpc.CallOption) (*response.DouyinVActiveFollowFollowerUserinfoResponse, error) {
	out := new(response.DouyinVActiveFollowFollowerUserinfoResponse)
	err := c.cc.Invoke(ctx, "/api.User/PushVActiveFollowerUserinfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) PushVActiveFollowerUserinfoRevert(ctx context.Context, in *request.DouyinVActiveFollowFollowerUserinfoRequest, opts ...grpc.CallOption) (*response.DouyinVActiveFollowFollowerUserinfoResponse, error) {
	out := new(response.DouyinVActiveFollowFollowerUserinfoResponse)
	err := c.cc.Invoke(ctx, "/api.User/PushVActiveFollowerUserinfoRevert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) DeleteVActiveFollowerUserinfo(ctx context.Context, in *request.DouyinDeleteVActiveFollowFollowerUserinfoRequest, opts ...grpc.CallOption) (*response.DouyinDeleteVActiveFollowFollowerUserinfoResponse, error) {
	out := new(response.DouyinDeleteVActiveFollowFollowerUserinfoResponse)
	err := c.cc.Invoke(ctx, "/api.User/DeleteVActiveFollowerUserinfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetVActiveFollowerUserinfo(ctx context.Context, in *request.DouyinGetVActiveFollowFollowerUserinfoRequest, opts ...grpc.CallOption) (*response.DouyinGetVActiveFollowFollowerUserinfoResponse, error) {
	out := new(response.DouyinGetVActiveFollowFollowerUserinfoResponse)
	err := c.cc.Invoke(ctx, "/api.User/GetVActiveFollowerUserinfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServer is the server API for User service.
// All implementations must embed UnimplementedUserServer
// for forward compatibility
type UserServer interface {
	UserRegister(context.Context, *request.DouyinUserRegisterRequest) (*response.DouyinUserRegisterResponse, error)
	UserLogin(context.Context, *request.DouyinUserLoginRequest) (*response.DouyinUserLoginResponse, error)
	GetUserInfo(context.Context, *request.DouyinUserRequest) (*response.DouyinUserResponse, error)
	GetUserInfoList(context.Context, *request.DouyinUserListRequest) (*response.DouyinUserListResponse, error)
	PushVUserRelativeInfoInit(context.Context, *request.DouyinPushVSetRequest) (*response.DouyinPushVSetResponse, error)
	UserIsInfluencerActiver(context.Context, *request.DouyinUserIsInfluencerActiverRequest) (*response.DouyinUserIsInfluencerActiverResponse, error)
	AddUserVideoCountSet(context.Context, *request.DouyinUserVideoCountSetRequest) (*response.DouyinUserVideoCountSetResponse, error)
	AddUserVideoCountSetRevert(context.Context, *request.DouyinUserVideoCountSetRequest) (*response.DouyinUserVideoCountSetResponse, error)
	AddUserFavoriteVideoCountSet(context.Context, *request.DouyinUserVideoCountSetRequest) (*response.DouyinUserVideoCountSetResponse, error)
	AddUserFavoriteVideoCountSetRevert(context.Context, *request.DouyinUserVideoCountSetRequest) (*response.DouyinUserVideoCountSetResponse, error)
	AddUserFollowUserCountSet(context.Context, *request.DouyinUserFollowCountSetRequest) (*response.DouyinUserFollowCountSetResponse, error)
	AddUserFollowerUserCountSet(context.Context, *request.DouyinUserFollowerCountSetRequest) (*response.DouyinUserFollowerCountSetResponse, error)
	AddUserFollowUserCountSetRevert(context.Context, *request.DouyinUserFollowCountSetRequest) (*response.DouyinUserFollowCountSetResponse, error)
	AddUserFollowerUserCountSetRevert(context.Context, *request.DouyinUserFollowerCountSetRequest) (*response.DouyinUserFollowerCountSetResponse, error)
	PushVActiveFollowerUserinfo(context.Context, *request.DouyinVActiveFollowFollowerUserinfoRequest) (*response.DouyinVActiveFollowFollowerUserinfoResponse, error)
	PushVActiveFollowerUserinfoRevert(context.Context, *request.DouyinVActiveFollowFollowerUserinfoRequest) (*response.DouyinVActiveFollowFollowerUserinfoResponse, error)
	DeleteVActiveFollowerUserinfo(context.Context, *request.DouyinDeleteVActiveFollowFollowerUserinfoRequest) (*response.DouyinDeleteVActiveFollowFollowerUserinfoResponse, error)
	GetVActiveFollowerUserinfo(context.Context, *request.DouyinGetVActiveFollowFollowerUserinfoRequest) (*response.DouyinGetVActiveFollowFollowerUserinfoResponse, error)
	mustEmbedUnimplementedUserServer()
}

// UnimplementedUserServer must be embedded to have forward compatible implementations.
type UnimplementedUserServer struct {
}

func (UnimplementedUserServer) UserRegister(context.Context, *request.DouyinUserRegisterRequest) (*response.DouyinUserRegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserRegister not implemented")
}
func (UnimplementedUserServer) UserLogin(context.Context, *request.DouyinUserLoginRequest) (*response.DouyinUserLoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserLogin not implemented")
}
func (UnimplementedUserServer) GetUserInfo(context.Context, *request.DouyinUserRequest) (*response.DouyinUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfo not implemented")
}
func (UnimplementedUserServer) GetUserInfoList(context.Context, *request.DouyinUserListRequest) (*response.DouyinUserListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfoList not implemented")
}
func (UnimplementedUserServer) PushVUserRelativeInfoInit(context.Context, *request.DouyinPushVSetRequest) (*response.DouyinPushVSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushVUserRelativeInfoInit not implemented")
}
func (UnimplementedUserServer) UserIsInfluencerActiver(context.Context, *request.DouyinUserIsInfluencerActiverRequest) (*response.DouyinUserIsInfluencerActiverResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserIsInfluencerActiver not implemented")
}
func (UnimplementedUserServer) AddUserVideoCountSet(context.Context, *request.DouyinUserVideoCountSetRequest) (*response.DouyinUserVideoCountSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserVideoCountSet not implemented")
}
func (UnimplementedUserServer) AddUserVideoCountSetRevert(context.Context, *request.DouyinUserVideoCountSetRequest) (*response.DouyinUserVideoCountSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserVideoCountSetRevert not implemented")
}
func (UnimplementedUserServer) AddUserFavoriteVideoCountSet(context.Context, *request.DouyinUserVideoCountSetRequest) (*response.DouyinUserVideoCountSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserFavoriteVideoCountSet not implemented")
}
func (UnimplementedUserServer) AddUserFavoriteVideoCountSetRevert(context.Context, *request.DouyinUserVideoCountSetRequest) (*response.DouyinUserVideoCountSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserFavoriteVideoCountSetRevert not implemented")
}
func (UnimplementedUserServer) AddUserFollowUserCountSet(context.Context, *request.DouyinUserFollowCountSetRequest) (*response.DouyinUserFollowCountSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserFollowUserCountSet not implemented")
}
func (UnimplementedUserServer) AddUserFollowerUserCountSet(context.Context, *request.DouyinUserFollowerCountSetRequest) (*response.DouyinUserFollowerCountSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserFollowerUserCountSet not implemented")
}
func (UnimplementedUserServer) AddUserFollowUserCountSetRevert(context.Context, *request.DouyinUserFollowCountSetRequest) (*response.DouyinUserFollowCountSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserFollowUserCountSetRevert not implemented")
}
func (UnimplementedUserServer) AddUserFollowerUserCountSetRevert(context.Context, *request.DouyinUserFollowerCountSetRequest) (*response.DouyinUserFollowerCountSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserFollowerUserCountSetRevert not implemented")
}
func (UnimplementedUserServer) PushVActiveFollowerUserinfo(context.Context, *request.DouyinVActiveFollowFollowerUserinfoRequest) (*response.DouyinVActiveFollowFollowerUserinfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushVActiveFollowerUserinfo not implemented")
}
func (UnimplementedUserServer) PushVActiveFollowerUserinfoRevert(context.Context, *request.DouyinVActiveFollowFollowerUserinfoRequest) (*response.DouyinVActiveFollowFollowerUserinfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushVActiveFollowerUserinfoRevert not implemented")
}
func (UnimplementedUserServer) DeleteVActiveFollowerUserinfo(context.Context, *request.DouyinDeleteVActiveFollowFollowerUserinfoRequest) (*response.DouyinDeleteVActiveFollowFollowerUserinfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteVActiveFollowerUserinfo not implemented")
}
func (UnimplementedUserServer) GetVActiveFollowerUserinfo(context.Context, *request.DouyinGetVActiveFollowFollowerUserinfoRequest) (*response.DouyinGetVActiveFollowFollowerUserinfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVActiveFollowerUserinfo not implemented")
}
func (UnimplementedUserServer) mustEmbedUnimplementedUserServer() {}

// UnsafeUserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServer will
// result in compilation errors.
type UnsafeUserServer interface {
	mustEmbedUnimplementedUserServer()
}

func RegisterUserServer(s grpc.ServiceRegistrar, srv UserServer) {
	s.RegisterService(&User_ServiceDesc, srv)
}

func _User_UserRegister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinUserRegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).UserRegister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/UserRegister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).UserRegister(ctx, req.(*request.DouyinUserRegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_UserLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinUserLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).UserLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/UserLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).UserLogin(ctx, req.(*request.DouyinUserLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/GetUserInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetUserInfo(ctx, req.(*request.DouyinUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetUserInfoList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinUserListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetUserInfoList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/GetUserInfoList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetUserInfoList(ctx, req.(*request.DouyinUserListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_PushVUserRelativeInfoInit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinPushVSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).PushVUserRelativeInfoInit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/PushVUserRelativeInfoInit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).PushVUserRelativeInfoInit(ctx, req.(*request.DouyinPushVSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_UserIsInfluencerActiver_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinUserIsInfluencerActiverRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).UserIsInfluencerActiver(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/UserIsInfluencerActiver",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).UserIsInfluencerActiver(ctx, req.(*request.DouyinUserIsInfluencerActiverRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_AddUserVideoCountSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinUserVideoCountSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).AddUserVideoCountSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/AddUserVideoCountSet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).AddUserVideoCountSet(ctx, req.(*request.DouyinUserVideoCountSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_AddUserVideoCountSetRevert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinUserVideoCountSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).AddUserVideoCountSetRevert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/AddUserVideoCountSetRevert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).AddUserVideoCountSetRevert(ctx, req.(*request.DouyinUserVideoCountSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_AddUserFavoriteVideoCountSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinUserVideoCountSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).AddUserFavoriteVideoCountSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/AddUserFavoriteVideoCountSet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).AddUserFavoriteVideoCountSet(ctx, req.(*request.DouyinUserVideoCountSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_AddUserFavoriteVideoCountSetRevert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinUserVideoCountSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).AddUserFavoriteVideoCountSetRevert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/AddUserFavoriteVideoCountSetRevert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).AddUserFavoriteVideoCountSetRevert(ctx, req.(*request.DouyinUserVideoCountSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_AddUserFollowUserCountSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinUserFollowCountSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).AddUserFollowUserCountSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/AddUserFollowUserCountSet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).AddUserFollowUserCountSet(ctx, req.(*request.DouyinUserFollowCountSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_AddUserFollowerUserCountSet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinUserFollowerCountSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).AddUserFollowerUserCountSet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/AddUserFollowerUserCountSet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).AddUserFollowerUserCountSet(ctx, req.(*request.DouyinUserFollowerCountSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_AddUserFollowUserCountSetRevert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinUserFollowCountSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).AddUserFollowUserCountSetRevert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/AddUserFollowUserCountSetRevert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).AddUserFollowUserCountSetRevert(ctx, req.(*request.DouyinUserFollowCountSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_AddUserFollowerUserCountSetRevert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinUserFollowerCountSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).AddUserFollowerUserCountSetRevert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/AddUserFollowerUserCountSetRevert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).AddUserFollowerUserCountSetRevert(ctx, req.(*request.DouyinUserFollowerCountSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_PushVActiveFollowerUserinfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinVActiveFollowFollowerUserinfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).PushVActiveFollowerUserinfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/PushVActiveFollowerUserinfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).PushVActiveFollowerUserinfo(ctx, req.(*request.DouyinVActiveFollowFollowerUserinfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_PushVActiveFollowerUserinfoRevert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinVActiveFollowFollowerUserinfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).PushVActiveFollowerUserinfoRevert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/PushVActiveFollowerUserinfoRevert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).PushVActiveFollowerUserinfoRevert(ctx, req.(*request.DouyinVActiveFollowFollowerUserinfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_DeleteVActiveFollowerUserinfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinDeleteVActiveFollowFollowerUserinfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).DeleteVActiveFollowerUserinfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/DeleteVActiveFollowerUserinfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).DeleteVActiveFollowerUserinfo(ctx, req.(*request.DouyinDeleteVActiveFollowFollowerUserinfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetVActiveFollowerUserinfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(request.DouyinGetVActiveFollowFollowerUserinfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetVActiveFollowerUserinfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.User/GetVActiveFollowerUserinfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetVActiveFollowerUserinfo(ctx, req.(*request.DouyinGetVActiveFollowFollowerUserinfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// User_ServiceDesc is the grpc.ServiceDesc for User service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var User_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.User",
	HandlerType: (*UserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserRegister",
			Handler:    _User_UserRegister_Handler,
		},
		{
			MethodName: "UserLogin",
			Handler:    _User_UserLogin_Handler,
		},
		{
			MethodName: "GetUserInfo",
			Handler:    _User_GetUserInfo_Handler,
		},
		{
			MethodName: "GetUserInfoList",
			Handler:    _User_GetUserInfoList_Handler,
		},
		{
			MethodName: "PushVUserRelativeInfoInit",
			Handler:    _User_PushVUserRelativeInfoInit_Handler,
		},
		{
			MethodName: "UserIsInfluencerActiver",
			Handler:    _User_UserIsInfluencerActiver_Handler,
		},
		{
			MethodName: "AddUserVideoCountSet",
			Handler:    _User_AddUserVideoCountSet_Handler,
		},
		{
			MethodName: "AddUserVideoCountSetRevert",
			Handler:    _User_AddUserVideoCountSetRevert_Handler,
		},
		{
			MethodName: "AddUserFavoriteVideoCountSet",
			Handler:    _User_AddUserFavoriteVideoCountSet_Handler,
		},
		{
			MethodName: "AddUserFavoriteVideoCountSetRevert",
			Handler:    _User_AddUserFavoriteVideoCountSetRevert_Handler,
		},
		{
			MethodName: "AddUserFollowUserCountSet",
			Handler:    _User_AddUserFollowUserCountSet_Handler,
		},
		{
			MethodName: "AddUserFollowerUserCountSet",
			Handler:    _User_AddUserFollowerUserCountSet_Handler,
		},
		{
			MethodName: "AddUserFollowUserCountSetRevert",
			Handler:    _User_AddUserFollowUserCountSetRevert_Handler,
		},
		{
			MethodName: "AddUserFollowerUserCountSetRevert",
			Handler:    _User_AddUserFollowerUserCountSetRevert_Handler,
		},
		{
			MethodName: "PushVActiveFollowerUserinfo",
			Handler:    _User_PushVActiveFollowerUserinfo_Handler,
		},
		{
			MethodName: "PushVActiveFollowerUserinfoRevert",
			Handler:    _User_PushVActiveFollowerUserinfoRevert_Handler,
		},
		{
			MethodName: "DeleteVActiveFollowerUserinfo",
			Handler:    _User_DeleteVActiveFollowerUserinfo_Handler,
		},
		{
			MethodName: "GetVActiveFollowerUserinfo",
			Handler:    _User_GetVActiveFollowerUserinfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user/api/user.proto",
}
