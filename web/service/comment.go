package service

import (
	"context"
	request1 "douyin/proto/comment/request"
	"douyin/proto/comment/response"
	"douyin/web/client/grpc"
	"douyin/web/model/request"
	"go.uber.org/zap"
	"strconv"
)

const (
	errorCommentVideo        = "comment video failed"
	errorGetCommentVideoList = "get comment video list failed"
)

func CommentVideo(commentVideoRequest request.CommentVideoRequest, userId int64) (*response.Comment, error) {
	actionType, _ := strconv.Atoi(commentVideoRequest.ActionType)
	commentId, _ := strconv.ParseInt(commentVideoRequest.CommentID, 10, 64)
	videoId1, _ := strconv.ParseInt(commentVideoRequest.VideoID, 10, 64)
	res, err := grpc.CommentClient.CommentVideo(context.Background(), &request1.DouyinCommentActionRequest{
		VideoId:     videoId1,
		ActionType:  int32(actionType),
		CommentText: commentVideoRequest.CommentText,
		CommentId:   commentId,
		LoginUserId: userId,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res.Code == 2 {
			zap.L().Error(errorCommentVideo, zap.Error(err))
			return nil, err
		}
	}
	return res.Comment, nil
}

func GetCommentVideoList(getCommentVideoListRequest request.GetCommentVideoListRequest, userId int64) ([]*response.Comment, error) {
	videoId1, _ := strconv.ParseInt(getCommentVideoListRequest.VideoID, 10, 64)
	res, err := grpc.CommentClient.GetCommentVideoList(context.Background(), &request1.DouyinCommentListRequest{
		VideoId:     videoId1,
		LoginUserId: userId,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res.Code == 2 {
			zap.L().Error(errorGetCommentVideoList, zap.Error(err))
			return nil, err
		}
	}
	return res.CommentList, nil
}
