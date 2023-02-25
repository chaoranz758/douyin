package handler

import (
	"context"
	"douyin/proto/favorite/api"
	"douyin/proto/favorite/request"
	"douyin/proto/favorite/response"
	"douyin/service/favorite/service"
	"go.uber.org/zap"
)

const (
	errorFavoriteVideo                  = "favorite video failed"
	errorGetFavoriteVideoList           = "get favorite video list failed"
	errorGetFavoriteCount               = "get favorite count failed"
	errorGetUserFavoriteVideoIdList     = "get user favorite video id list failed"
	errorGetUserListFavoriteVideoIdList = "get user list favorite video id list"
	errorGetUserFavoritedCount          = "get user favorited count failed"
)

type Favorite struct {
	api.UnimplementedFavoriteServer
}

func (favorite *Favorite) FavoriteVideo(ctx context.Context, req *request.DouyinFavoriteActionRequest) (*response.DouyinFavoriteActionResponse, error) {
	if err := service.FavoriteVideoDtm(req); err != nil {
		zap.L().Error(errorFavoriteVideo, zap.Error(err))
		return &response.DouyinFavoriteActionResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinFavoriteActionResponse{
		Code: 1,
	}, nil
}

func (favorite *Favorite) GetFavoriteVideoList(ctx context.Context, req *request.DouyinFavoriteListRequest) (*response.DouyinFavoriteListResponse, error) {
	data, err := service.GetFavoriteVideoList(req)
	if err != nil {
		zap.L().Error(errorGetFavoriteVideoList, zap.Error(err))
		return &response.DouyinFavoriteListResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinFavoriteListResponse{
		VideoList: data,
		Code:      1,
	}, nil
}

func (favorite *Favorite) GetFavoriteCount(ctx context.Context, req *request.DouyinFavoriteCountRequest) (*response.DouyinFavoriteCountResponse, error) {
	favoriteCountList, isFavoriteList, err := service.GetFavoriteCount(req)
	if err != nil {
		zap.L().Error(errorGetFavoriteCount, zap.Error(err))
		return &response.DouyinFavoriteCountResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinFavoriteCountResponse{
		Code:          1,
		FavoriteCount: favoriteCountList,
		IsFavorite:    isFavoriteList,
	}, nil
}

func (favorite *Favorite) GetUserFavoriteVideoIdList(ctx context.Context, req *request.DouyinFavoriteIdListRequest) (*response.DouyinFavoriteIdListResponse, error) {
	data, err := service.GetUserFavoriteVideoIdList(req)
	if err != nil {
		zap.L().Error(errorGetUserFavoriteVideoIdList, zap.Error(err))
		return &response.DouyinFavoriteIdListResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinFavoriteIdListResponse{
		VideoId: data,
		Code:    1,
	}, nil
}

func (favorite *Favorite) GetUserListFavoriteVideoIdList(ctx context.Context, req *request.DouyinFavoriteListIdListRequest) (*response.DouyinFavoriteListIdListResponse, error) {
	data, err := service.GetUserListFavoriteVideoIdList(req)
	if err != nil {
		zap.L().Error(errorGetUserListFavoriteVideoIdList, zap.Error(err))
		return &response.DouyinFavoriteListIdListResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinFavoriteListIdListResponse{
		FavoriteVideoIdList: data,
		Code:                1,
	}, nil
}

func (favorite *Favorite) GetUserFavoritedCount(ctx context.Context, req *request.DouyinGetUserFavoritedCountRequest) (*response.DouyinGetUserFavoritedCountResponse, error) {
	countFavorite, countFavorited, err := service.GetUserFavoritedCount(req)
	if err != nil {
		zap.L().Error(errorGetUserFavoritedCount, zap.Error(err))
		return &response.DouyinGetUserFavoritedCountResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinGetUserFavoritedCountResponse{
		FavoriteCount:  countFavorite,
		FavoritedCount: countFavorited,
		Code:           1,
	}, nil
}
