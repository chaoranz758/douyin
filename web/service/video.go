package service

import (
	"context"
	request1 "douyin/proto/video/request"
	"douyin/proto/video/response"
	"douyin/web/client/grpc"
	"douyin/web/model/request"
	"go.uber.org/zap"
	"strconv"
)

const (
	errorPublishVideo    = "publish video failed"
	errorGetPublishVideo = "get publish video failed"
	errorGetVideoList    = "get video list failed"
)

func PublishVideo(publishVideoRequest request.PublishVideoRequest, userID int64, videoURL, coverURL string) error {
	res, err := grpc.VideoClient.PublishVideo(context.Background(), &request1.DouyinPublishActionRequest{
		VideoUrl: videoURL,
		CoverUrl: coverURL,
		Title:    publishVideoRequest.Title,
		UserId:   userID,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return err
		}
		if res.Code == 2 {
			zap.L().Error(errorPublishVideo, zap.Error(err))
			return err
		}
	}
	return nil
}

func GetPublishVideo(getPublishVideoRequest request.GetPublishVideoRequest, loginUserID int64) ([]*response.Video, error) {
	userId, _ := strconv.ParseInt(getPublishVideoRequest.UserID, 10, 64)
	res, err := grpc.VideoClient.GetPublishVideo(context.Background(), &request1.DouyinPublishListRequest{
		UserId:  userId,
		LoginId: loginUserID,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res.Code == 2 {
			zap.L().Error(errorGetPublishVideo, zap.Error(err))
			return nil, err
		}
	}
	return res.VideoList, nil
}

func GetVideoList(getVideoListRequest request.GetVideoListRequest, userId int64) ([]*response.Video, int64, error) {
	time1, _ := strconv.ParseInt(getVideoListRequest.LastTime, 10, 64)
	res, err := grpc.VideoClient.GetVideoList(context.Background(), &request1.DouyinFeedRequest{
		LatestTime:  time1,
		LoginUserId: userId,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, 0, err
		}
		if res.Code == 2 {
			zap.L().Error(errorGetVideoList, zap.Error(err))
			return nil, 0, err
		}
	}
	return res.VideoList, res.NextTime, nil
}
