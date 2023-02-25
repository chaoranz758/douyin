package service

import (
	"context"
	request1 "douyin/proto/favorite/request"
	"douyin/proto/video/response"
	"douyin/web/client/grpc"
	"douyin/web/model/request"
	"go.uber.org/zap"
	"strconv"
)

const (
	errorFavoriteVideo = "favorite video failed"
)

func FavoriteVideo(favoriteVideoRequest request.FavoriteVideoRequest, loginUserID int64) error {
	n, _ := strconv.Atoi(favoriteVideoRequest.ActionType)
	videoId1, _ := strconv.ParseInt(favoriteVideoRequest.VideoID, 10, 64)
	res, err := grpc.FavoriteClient.FavoriteVideo(context.Background(), &request1.DouyinFavoriteActionRequest{
		VideoId:     videoId1,
		ActionType:  int32(n),
		LoginUserId: loginUserID,
	})
	if err != nil {
		//fmt.Printf("%v\n", res.Code)
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return err
		}
		if res.Code == 2 {
			zap.L().Error(errorFavoriteVideo, zap.Error(err))
			return err
		}
	}
	return nil
}

func GetFavoriteVideoList(getFavoriteVideoListRequest request.GetFavoriteVideoListRequest, userID int64) ([]*response.Video, error) {
	userId1, _ := strconv.ParseInt(getFavoriteVideoListRequest.UserID, 10, 64)
	res, err := grpc.FavoriteClient.GetFavoriteVideoList(context.Background(), &request1.DouyinFavoriteListRequest{
		UserId:      userId1,
		LoginUserId: userID,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res.Code == 2 {
			zap.L().Error(errorFavoriteVideo, zap.Error(err))
			return nil, err
		}
	}
	return res.VideoList, nil
}
