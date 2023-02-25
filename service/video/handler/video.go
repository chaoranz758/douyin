package handler

import (
	"context"
	"douyin/proto/video/api"
	"douyin/proto/video/request"
	"douyin/proto/video/response"
	"douyin/service/video/service"
	"go.uber.org/zap"
)

const (
	errorPublishVideo                   = "publish video failed"
	errorGetPublishVideo                = "get publish video failed"
	errorJudgeVideoAuthor               = "judge video author failed"
	errorPushVActiveFavoriteVideo       = "push v or active favorite video failed"
	errorPushVActiveFavoriteVideoRevert = "push v or active favorite video revert failed"
	errorGetVActiveFavoriteVideo        = "get v or active favorite video failed"
	errorDeleteVActiveFavoriteVideo     = "delete v or active favorite video failed"
	errorGetVideoListInner              = "get video list inner failed"
	errorGetVideoList                   = "get video list failed"
	errorPushVActiveBasicInfoInit       = "push v or active basic info init failed"
	errorPushActiveBasicInfoInit        = "push active basic info init failed"
	errorGetUserPublishVideoCountList   = "get user publish video count list failed"
	errorGetUserPublishVideoCount       = "get user publish video count failed"
)

type Video struct {
	api.UnimplementedVideoServer
}

func (video *Video) PublishVideo(ctx context.Context, req *request.DouyinPublishActionRequest) (*response.DouyinPublishActionResponse, error) {
	if err := service.PublishVideoDtm(req); err != nil {
		zap.L().Error(errorPublishVideo, zap.Error(err))
		return &response.DouyinPublishActionResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinPublishActionResponse{
		Code: 1,
	}, nil
}

func (video *Video) GetPublishVideo(ctx context.Context, req *request.DouyinPublishListRequest) (*response.DouyinPublishListResponse, error) {
	data, err := service.GetPublishVideo(req)
	if err != nil {
		zap.L().Error(errorGetPublishVideo, zap.Error(err))
		return &response.DouyinPublishListResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinPublishListResponse{
		Code:      1,
		VideoList: data,
	}, nil
}

func (video *Video) GetVideoList(ctx context.Context, req *request.DouyinFeedRequest) (*response.DouyinFeedResponse, error) {
	data, nextTime, err := service.GetVideoList(req)
	if err != nil {
		zap.L().Error(errorGetVideoList, zap.Error(err))
		return &response.DouyinFeedResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinFeedResponse{
		VideoList: data,
		NextTime:  nextTime,
		Code:      1,
	}, nil
}

func (video *Video) JudgeVideoAuthor(ctx context.Context, req *request.DouyinJudgeVideoAuthorRequest) (*response.DouyinJudgeVideoAuthorResponse, error) {
	authorId, isV, isActive, err := service.JudgeVideoAuthor(req)
	if err != nil {
		zap.L().Error(errorJudgeVideoAuthor, zap.Error(err))
		return &response.DouyinJudgeVideoAuthorResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinJudgeVideoAuthorResponse{
		Code:     1,
		AuthorId: authorId,
		IsV:      isV,
		IsActive: isActive,
	}, err
}

func (video *Video) PushVActiveFavoriteVideo(ctx context.Context, req *request.DouyinPushVActiveFavoriteVideoRequest) (*response.DouyinPushVActiveFavoriteVideoResponse, error) {
	if err := service.PushVActiveFavoriteVideo(req); err != nil {
		zap.L().Error(errorPushVActiveFavoriteVideo, zap.Error(err))
		return &response.DouyinPushVActiveFavoriteVideoResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinPushVActiveFavoriteVideoResponse{
		Code: 1,
	}, nil
}

func (video *Video) GetVActiveFavoriteVideo(ctx context.Context, req *request.DouyinGetVActiveFavoriteVideoRequest) (*response.DouyinGetVActiveFavoriteVideoResponse, error) {
	data, err := service.GetVActiveFavoriteVideo(req)
	if err != nil {
		zap.L().Error(errorGetVActiveFavoriteVideo, zap.Error(err))
		return &response.DouyinGetVActiveFavoriteVideoResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinGetVActiveFavoriteVideoResponse{
		VideoList: data,
		Code:      1,
	}, nil
}

func (video *Video) DeleteVActiveFavoriteVideo(ctx context.Context, req *request.DouyinDeleteVActiveFavoriteVideoRequest) (*response.DouyinDeleteVActiveFavoriteVideoResponse, error) {
	if err := service.DeleteVActiveFavoriteVideo(req); err != nil {
		zap.L().Error(errorDeleteVActiveFavoriteVideo, zap.Error(err))
		return &response.DouyinDeleteVActiveFavoriteVideoResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinDeleteVActiveFavoriteVideoResponse{
		Code: 1,
	}, nil
}

func (video *Video) GetVideoListInner(ctx context.Context, req *request.DouyinGetVideoListRequest) (*response.DouyinGetVideoListResponse, error) {
	data, err := service.GetVideoListInner(req)
	if err != nil {
		zap.L().Error(errorGetVideoListInner, zap.Error(err))
		return &response.DouyinGetVideoListResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinGetVideoListResponse{
		Code:      1,
		VideoList: data,
	}, nil
}

func (video *Video) PushVBasicInfoInit(ctx context.Context, req *request.DouyinPushVInfoInitRequest) (*response.DouyinPushVInfoInitResponse, error) {
	data, err := service.PushVActiveBasicInfoInit(req)
	if err != nil {
		zap.L().Error(errorPushVActiveBasicInfoInit, zap.Error(err))
		return &response.DouyinPushVInfoInitResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinPushVInfoInitResponse{
		Code:        1,
		VideoIdList: data,
	}, nil
}

func (video *Video) PushActiveBasicInfoInit(ctx context.Context, req *request.DouyinPushActiveInfoInitRequest) (*response.DouyinPushActiveInfoInitResponse, error) {
	if err := service.PushActiveBasicInfoInit(req); err != nil {
		zap.L().Error(errorPushActiveBasicInfoInit, zap.Error(err))
		return &response.DouyinPushActiveInfoInitResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinPushActiveInfoInitResponse{
		Code: 1,
	}, nil
}

func (video *Video) GetUserPublishVideoCount(ctx context.Context, req *request.DouyinGetUserPublishVideoCountRequest) (*response.DouyinGetUserPublishVideoCountResponse, error) {
	count, err := service.GetUserPublishVideoCount(req)
	if err != nil {
		zap.L().Error(errorGetUserPublishVideoCount, zap.Error(err))
		return &response.DouyinGetUserPublishVideoCountResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinGetUserPublishVideoCountResponse{
		Code:  1,
		Count: count,
	}, nil

}

func (video *Video) GetUserPublishVideoCountList(ctx context.Context, req *request.DouyinGetUserPublishVideoCountListRequest) (*response.DouyinGetUserPublishVideoCountListResponse, error) {
	count, err := service.GetUserPublishVideoCountList(req)
	if err != nil {
		zap.L().Error(errorGetUserPublishVideoCountList, zap.Error(err))
		return &response.DouyinGetUserPublishVideoCountListResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinGetUserPublishVideoCountListResponse{
		Count: count,
		Code:  1,
	}, nil
}

func (video *Video) PushVActiveFavoriteVideoRevert(ctx context.Context, req *request.DouyinPushVActiveFavoriteVideoRequest) (*response.DouyinPushVActiveFavoriteVideoResponse, error) {
	if err := service.PushVActiveFavoriteVideoRevert(req); err != nil {
		zap.L().Error(errorPushVActiveFavoriteVideoRevert, zap.Error(err))
		return &response.DouyinPushVActiveFavoriteVideoResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinPushVActiveFavoriteVideoResponse{
		Code: 1,
	}, nil
}
