package handler

import (
	"context"
	"douyin/proto/comment/api"
	"douyin/proto/comment/request"
	"douyin/proto/comment/response"
	"douyin/service/comment/service"
	"go.uber.org/zap"
)

const (
	errorCommentVideo              = "comment video failed"
	errorGetCommentVideoList       = "get comment video list failed"
	errorGetCommentCount           = "get comment count failed"
	errorPushVCommentBasicInfoInit = "push v comment basic info init failed"
)

type Comment struct {
	api.UnimplementedCommentServer
}

func (comment *Comment) CommentVideo(ctx context.Context, req *request.DouyinCommentActionRequest) (*response.DouyinCommentActionResponse, error) {
	data, err := service.CommentVideo(req)
	if err != nil {
		zap.L().Error(errorCommentVideo, zap.Error(err))
		return &response.DouyinCommentActionResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinCommentActionResponse{
		Code:    1,
		Comment: data,
	}, nil
}

func (comment *Comment) GetCommentVideoList(ctx context.Context, req *request.DouyinCommentListRequest) (*response.DouyinCommentListResponse, error) {
	data, err := service.GetCommentVideoList(req)
	if err != nil {
		zap.L().Error(errorGetCommentVideoList, zap.Error(err))
		return &response.DouyinCommentListResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinCommentListResponse{
		Code:        1,
		CommentList: data,
	}, nil
}

func (comment *Comment) GetCommentCount(ctx context.Context, req *request.DouyinCommentCountRequest) (*response.DouyinCommentCountResponse, error) {
	countList, err := service.GetCommentCount(req)
	if err != nil {
		zap.L().Error(errorGetCommentCount, zap.Error(err))
		return &response.DouyinCommentCountResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinCommentCountResponse{
		CommentCount: countList,
		Code:         1,
	}, err
}

func (comment *Comment) PushVCommentBasicInfoInit(ctx context.Context, req *request.DouyinPushVCommentBasicInfoInitRequest) (*response.DouyinPushVCommentBasicInfoInitResponse, error) {
	if err := service.PushVCommentBasicInfoInit(req); err != nil {
		zap.L().Error(errorPushVCommentBasicInfoInit, zap.Error(err))
		return &response.DouyinPushVCommentBasicInfoInitResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinPushVCommentBasicInfoInitResponse{
		Code: 1,
	}, nil
}
