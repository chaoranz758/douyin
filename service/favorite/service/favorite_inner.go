package service

import (
	"douyin/proto/favorite/request"
	"douyin/proto/favorite/response"
	"douyin/service/favorite/dao/mysql"
	"douyin/service/favorite/dao/redis"
	"douyin/service/favorite/model"
	"errors"
	"go.uber.org/zap"
)

func GetFavoriteCount(req *request.DouyinFavoriteCountRequest) ([]int64, []bool, error) {
	if len(req.VideoId) == 0 {
		return nil, nil, nil
	}
	//读视频对应的点赞数
	favoriteCountList, err := redis.GetVideoFavoriteCount(req.VideoId)
	if err != nil {
		zap.L().Error(MGetVideoFavoriteCount, zap.Error(err))
		return nil, nil, err
	}
	if req.UserId == 0 {
		var userFavoriteBool []bool
		for i := 0; i < len(favoriteCountList); i++ {
			userFavoriteBool = append(userFavoriteBool, false)
		}
		return favoriteCountList, userFavoriteBool, nil
	}
	//mysql点赞表中查看用户是否对视频点赞
	if req.UserId != 0 {
		userFavoriteBool, err := mysql.GetUserFavoriteBool(req.UserId, req.VideoId)
		if err != nil {
			zap.L().Error(errorGetUserFavoriteBool, zap.Error(err))
			return nil, nil, err
		}
		return favoriteCountList, userFavoriteBool, nil
	}
	return nil, nil, errors.New(errorRpcInput)
}

func GetUserFavoriteVideoIdList(req *request.DouyinFavoriteIdListRequest) ([]int64, error) {
	if req.UserId == 0 {
		return nil, nil
	}
	f := make([]model.Favorite, 0)
	if err := mysql.GetUserFavoriteID(req.UserId, &f); err != nil {
		zap.L().Error(errorGetUserFavoriteID, zap.Error(err))
		return nil, err
	}
	if len(f) == 0 {
		return nil, nil
	}
	var idList []int64
	for i := 0; i < len(f); i++ {
		idList = append(idList, f[i].VideoID)
	}
	return idList, nil
}

func GetUserListFavoriteVideoIdList(req *request.DouyinFavoriteListIdListRequest) ([]*response.FavoriteVideoIdList, error) {
	if len(req.UserId) == 0 {
		return nil, nil
	}
	var result []*response.FavoriteVideoIdList
	for i := 0; i < len(req.UserId); i++ {
		f := make([]model.Favorite, 0)
		if err := mysql.GetUserFavoriteID(req.UserId[i], &f); err != nil {
			zap.L().Error(errorGetUserFavoriteID, zap.Error(err))
			return nil, err
		}
		var idList []int64
		for j := 0; j < len(f); j++ {
			idList = append(idList, f[j].VideoID)
		}
		var r1 = response.FavoriteVideoIdList{
			VideoId: idList,
		}
		result = append(result, &r1)
	}
	return result, nil
}

func GetUserFavoritedCount(req *request.DouyinGetUserFavoritedCountRequest) ([]int64, []int64, error) {
	if len(req.UserId) == 0 {
		return nil, nil, nil
	}
	countFavorite, countFavorited, err := redis.GetUserFavoritedCount(req.UserId)
	if err != nil {
		zap.L().Error(errorGetUserFavoritedCount, zap.Error(err))
		return nil, nil, err
	}
	return countFavorite, countFavorited, nil
}
