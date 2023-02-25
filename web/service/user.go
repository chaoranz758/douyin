package service

import (
	"context"
	request1 "douyin/proto/user/request"
	"douyin/proto/user/response"
	"douyin/web/client/grpc"
	"douyin/web/model/request"
	"go.uber.org/zap"
	"strconv"
)

const (
	errorConnectToGRPCServer = "connect to grpc server failed"
	errorUserRegister        = "rpc exe user register logic failed"
	errorUserLogin           = "rpc exe user login logic failed"
	errorGetUserInfo         = "rpc exe get user information logic failed"
)

func UserRegister(userRegisterRequest request.UserRegisterRequest) (int64, string, error) {
	//调用用户service注册逻辑
	res, err := grpc.UserGRPCClient.UserRegister(context.Background(), &request1.DouyinUserRegisterRequest{
		Username: userRegisterRequest.UserName,
		Password: userRegisterRequest.Password,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return 0, "", err
		}
		if res.Code == 2 {
			zap.L().Error(errorUserRegister, zap.Error(err))
			return 0, "", err
		}
	}
	return res.UserId, res.Token, nil
}

func UserLogin(userLoginRequest request.UserLoginRequest) (int64, string, error) {
	//调用用户service注册逻辑
	res, err := grpc.UserGRPCClient.UserLogin(context.Background(), &request1.DouyinUserLoginRequest{
		Username: userLoginRequest.UserName,
		Password: userLoginRequest.Password,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return 0, "", err
		}
		if res.Code == 2 {
			zap.L().Error(errorUserLogin, zap.Error(err))
			return 0, "", err
		}
	}
	return res.UserId, res.Token, nil
}

func GetUserInfo(getUserInfoRequest request.GetUserInfoRequest, loginUserID int64) (*response.User, error) {
	userId1, _ := strconv.ParseInt(getUserInfoRequest.UserID, 10, 64)
	res, err := grpc.UserGRPCClient.GetUserInfo(context.Background(), &request1.DouyinUserRequest{
		UserId:      userId1,
		LoginUserId: loginUserID,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res.Code == 2 {
			zap.L().Error(errorGetUserInfo, zap.Error(err))
			return nil, err
		}
	}
	return res.User, nil
}
