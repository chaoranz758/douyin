package service

import (
	"context"
	request1 "douyin/proto/follow/request"
	response1 "douyin/proto/follow/response"
	"douyin/proto/user/response"
	"douyin/web/client/grpc"
	"douyin/web/model/request"
	"go.uber.org/zap"
	"strconv"
)

const (
	errorFollowUser      = "follow user failed"
	errorGetFollowUser   = "get follower user information failed"
	errorGetFollowerUser = "get follower user information failed"
	errorGetFriendList   = "get friend list failed"
)

func FollowUser(followUserRequest request.FollowUserRequest, loginUserId int64) error {
	actionType, _ := strconv.Atoi(followUserRequest.ActionType)
	toUserId1, _ := strconv.ParseInt(followUserRequest.ToUserID, 10, 64)
	res, err := grpc.FollowClient.FollowUser(context.Background(), &request1.DouyinRelationActionRequest{
		ActionType:  int32(actionType),
		LoginUserId: loginUserId,
		ToUserId:    toUserId1,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return err
		}
		if res.Code == 2 {
			zap.L().Error(errorFollowUser, zap.Error(err))
			return err
		}
	}
	return nil
}

func GetFollowList(getFollowListRequest request.GetFollowListRequest, userId int64) ([]*response.User, error) {
	userId1, _ := strconv.ParseInt(getFollowListRequest.UserID, 10, 64)
	res, err := grpc.FollowClient.GetFollowList(context.Background(), &request1.DouyinRelationFollowListRequest{
		UserId:      userId1,
		LoginUserId: userId,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res.Code == 2 {
			zap.L().Error(errorGetFollowUser, zap.Error(err))
			return nil, err
		}
	}
	return res.UserList, nil
}

func GetFollowerList(getFollowerListRequest request.GetFollowerListRequest, userId int64) ([]*response.User, error) {
	userId1, _ := strconv.ParseInt(getFollowerListRequest.UserID, 10, 64)
	res, err := grpc.FollowClient.GetFollowerList(context.Background(), &request1.DouyinRelationFollowerListRequest{
		UserId:      userId1,
		LoginUserId: userId,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res.Code == 2 {
			zap.L().Error(errorGetFollowerUser, zap.Error(err))
			return nil, err
		}
	}
	return res.UserList, nil
}

func GetFriendList(getFriendListRequest request.GetFriendListRequest, userId int64) ([]*response1.FriendUser, error) {
	userId1, _ := strconv.ParseInt(getFriendListRequest.UserID, 10, 64)
	res, err := grpc.FollowClient.GetFriendList(context.Background(), &request1.DouyinRelationFriendListRequest{
		UserId:      userId1,
		LoginUserId: userId,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res.Code == 2 {
			zap.L().Error(errorGetFriendList, zap.Error(err))
			return nil, err
		}
	}
	return res.UserList, nil
}
