package handler

import (
	"context"
	"douyin/proto/follow/api"
	"douyin/proto/follow/request"
	"douyin/proto/follow/response"
	"douyin/service/follow/service"
	"go.uber.org/zap"
)

type Follow struct {
	api.UnimplementedFollowServer
}

func (follow *Follow) FollowUser(ctx context.Context, req *request.DouyinRelationActionRequest) (*response.DouyinRelationActionResponse, error) {
	if err := service.FollowUserDtm(req); err != nil {
		zap.L().Error(errorFollowUser, zap.Error(err))
		return &response.DouyinRelationActionResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinRelationActionResponse{
		Code: 1,
	}, nil
}

func (follow *Follow) GetFollowList(ctx context.Context, req *request.DouyinRelationFollowListRequest) (*response.DouyinRelationFollowListResponse, error) {
	data, err := service.GetFollowList(req)
	if err != nil {
		zap.L().Error(errorGetFollowUser, zap.Error(err))
		return &response.DouyinRelationFollowListResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinRelationFollowListResponse{
		Code:     1,
		UserList: data,
	}, nil
}

func (follow *Follow) GetFollowerList(ctx context.Context, req *request.DouyinRelationFollowerListRequest) (*response.DouyinRelationFollowerListResponse, error) {
	data, err := service.GetFollowerList(req)
	if err != nil {
		zap.L().Error(errorGetFollowerUser, zap.Error(err))
		return &response.DouyinRelationFollowerListResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinRelationFollowerListResponse{
		Code:     1,
		UserList: data,
	}, nil
}

func (follow *Follow) GetFriendList(ctx context.Context, req *request.DouyinRelationFriendListRequest) (*response.DouyinRelationFriendListResponse, error) {
	data, err := service.GetFriendList(req)
	if err != nil {
		zap.L().Error(errorGetFriendList, zap.Error(err))
		return &response.DouyinRelationFriendListResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinRelationFriendListResponse{
		Code:     1,
		UserList: data,
	}, nil
}

func (follow *Follow) GetFollowFollower(ctx context.Context, req *request.DouyinFollowFollowerCountRequest) (*response.DouyinFollowFollowerCountResponse, error) {
	followCount, followerCount, err := service.GetFollowFollower(req)
	if err != nil {
		zap.L().Error(errorGetFollowFollower, zap.Error(err))
		return &response.DouyinFollowFollowerCountResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinFollowFollowerCountResponse{
		Code:          1,
		FollowCount:   followCount,
		FollowerCount: followerCount,
	}, nil
}

func (follow *Follow) GetFollowFollowerList(ctx context.Context, req *request.DouyinFollowFollowerListCountRequest) (*response.DouyinFollowFollowerListCountResponse, error) {
	followCountList, followerCountList, err := service.GetFollowFollowerList(req)
	if err != nil {
		zap.L().Error(errorGetFollowFollowerList, zap.Error(err))
		return &response.DouyinFollowFollowerListCountResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinFollowFollowerListCountResponse{
		FollowCount:   followCountList,
		FollowerCount: followerCountList,
		Code:          1,
	}, nil
}

func (follow *Follow) GetFollowInfo(ctx context.Context, req *request.DouyinGetFollowRequest) (*response.DouyinGetFollowResponse, error) {
	b, err := service.GetFollowInfo(req)
	if err != nil {
		zap.L().Error(errorGetFollowInfo, zap.Error(err))
		return &response.DouyinGetFollowResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinGetFollowResponse{
		Code:     1,
		IsFollow: b,
	}, nil
}

func (follow *Follow) GetFollowInfoList(ctx context.Context, req *request.DouyinGetFollowListRequest) (*response.DouyinGetFollowListResponse, error) {
	boolList, err := service.GetFollowInfoList(req)
	if err != nil {
		zap.L().Error(errorGetFollowInfoList)
		return &response.DouyinGetFollowListResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinGetFollowListResponse{
		Code:     1,
		IsFollow: boolList,
	}, nil
}

func (follow *Follow) GetFollowFollowerIdList(ctx context.Context, req *request.DouyinGetFollowFollowerIdListRequest) (*response.DouyinGetFollowFollowerIdListResponse, error) {
	followList, followerList, err := service.GetFollowFollowerIdList(req)
	if err != nil {
		zap.L().Error(errorGetFollowFollowerIdList, zap.Error(err))
		return &response.DouyinGetFollowFollowerIdListResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinGetFollowFollowerIdListResponse{
		Code:         1,
		FollowList:   followList,
		FollowerList: followerList,
	}, nil
}
