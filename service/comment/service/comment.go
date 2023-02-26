package service

import (
	"context"
	"douyin/proto/comment/request"
	"douyin/proto/comment/response"
	request1 "douyin/proto/user/request"
	request2 "douyin/proto/video/request"
	"douyin/service/comment/dao/mysql"
	"douyin/service/comment/dao/redis"
	"douyin/service/comment/initialize/grpc_client"
	"douyin/service/comment/model"
	"douyin/service/message/util"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	time2 "time"
)

const (
	errorConnectToGRPCServer       = "connect to grpc server failed"
	errorCommentVideo              = "comment video failed"
	errorPushCommentInfo           = "push comment information failed"
	errorAddVideoCommentCount      = "push video comment count failed"
	errorCreateComment             = "create comment failed"
	errorGetUserInfo               = "get user information failed"
	errorDeleteCommentInfo         = "delete comment information failed"
	errorSubVideoCommentCount      = "sub video comment count failed"
	errorGetCommentInfo            = "get comment information failed"
	errorJsonMarshal               = "json marshal failed"
	errorGetCommentCount           = "get comment count failed"
	errorPushVCommentBasicInfoInit = "push v comment basic info init failed"
)

func CommentVideo(req *request.DouyinCommentActionRequest) (*response.Comment, error) {
	if req.ActionType == 1 {
		//判断视频作者是否为大V并将视频作者ID返回
		res, err := grpc_client.VideoClient.JudgeVideoAuthor(context.Background(), &request2.DouyinJudgeVideoAuthorRequest{
			VideoId: req.VideoId,
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
		if res.IsV == true {
			time := time2.Now()
			commitID := util.GenID()
			var comment = model.Commit{
				CommitID:  commitID,
				VideoID:   req.VideoId,
				UserID:    req.LoginUserId,
				Message:   req.CommentText,
				CreatedAt: time2.Now(),
			}
			//评论写入mysql评论表
			if err := mysql.PushCommentInfo(comment); err != nil {
				zap.L().Error(errorCreateComment, zap.Error(err))
				return nil, err
			}
			//视频评论数加一
			if err := redis.AddVideoCommentCount(comment.VideoID); err != nil {
				zap.L().Error(errorAddVideoCommentCount, zap.Error(err))
				//将评论信息从mysql删除
				if err := mysql.DeleteCommentInfo(commitID); err != nil {
					zap.L().Error(errorDeleteCommentInfo, zap.Error(err))
					return nil, err
				}
				return nil, err
			}
			//将评论信息存到redis
			if err := redis.PushCommentInfo(comment, res.AuthorId); err != nil {
				zap.L().Error(errorPushCommentInfo, zap.Error(err))
				//将评论信息从mysql删除
				if err := mysql.DeleteCommentInfo(commitID); err != nil {
					zap.L().Error(errorDeleteCommentInfo, zap.Error(err))
					return nil, err
				}
				//视频评论数减一
				if err := redis.SubVideoCommentCount(comment.VideoID); err != nil {
					zap.L().Error(errorSubVideoCommentCount, zap.Error(err))
					return nil, err
				}
				return nil, err
			}
			//获取评论的用户信息
			res1, err := grpc_client.UserClient.GetUserInfo(context.Background(), &request1.DouyinUserRequest{
				UserId:      req.LoginUserId,
				LoginUserId: req.LoginUserId,
			})
			if err != nil {
				if res1 == nil {
					zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
					return nil, err
				}
				if res1.Code == 2 {
					zap.L().Error(errorGetUserInfo, zap.Error(err))
					return nil, err
				}
			}
			//返回信息
			var comments = response.Comment{
				Id:         commitID,
				User:       res1.User,
				Content:    req.CommentText,
				CreateDate: fmt.Sprintf("%v", time.Format("2006-01-02 15:04:05")),
			}
			return &comments, nil
		}
		time := time2.Now()
		commitID := util.GenID()
		var comment = model.Commit{
			CommitID:  commitID,
			VideoID:   req.VideoId,
			UserID:    req.LoginUserId,
			Message:   req.CommentText,
			CreatedAt: time2.Now(),
		}
		//评论写入mysql评论表
		if err := mysql.PushCommentInfo(comment); err != nil {
			zap.L().Error(errorCreateComment, zap.Error(err))
			return nil, err
		}
		//视频评论数加一
		if err := redis.AddVideoCommentCount(comment.VideoID); err != nil {
			zap.L().Error(errorAddVideoCommentCount, zap.Error(err))
			//将评论信息从mysql删除
			if err := mysql.DeleteCommentInfo(commitID); err != nil {
				zap.L().Error(errorDeleteCommentInfo, zap.Error(err))
				return nil, err
			}
			return nil, err
		}
		//获取评论的用户信息
		res1, err := grpc_client.UserClient.GetUserInfo(context.Background(), &request1.DouyinUserRequest{
			UserId:      req.LoginUserId,
			LoginUserId: req.LoginUserId,
		})
		if err != nil {
			if res1 == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return nil, err
			}
			if res1.Code == 2 {
				zap.L().Error(errorGetUserInfo, zap.Error(err))
				return nil, err
			}
		}
		//返回信息
		var comments = response.Comment{
			Id:         commitID,
			User:       res1.User,
			Content:    req.CommentText,
			CreateDate: fmt.Sprintf("%v", time.Format("2006-01-02 15:04:05")),
		}
		return &comments, nil
	}
	//判断视频作者是否为大V并将视频作者ID返回
	res, err := grpc_client.VideoClient.JudgeVideoAuthor(context.Background(), &request2.DouyinJudgeVideoAuthorRequest{
		VideoId: req.VideoId,
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
	if res.IsV == true {
		//评论从mysql评论表删除
		if err := mysql.DeleteCommentInfo(req.CommentId); err != nil {
			zap.L().Error(errorDeleteCommentInfo, zap.Error(err))
			return nil, err
		}
		//视频评论数减一
		if err := redis.SubVideoCommentCount(req.VideoId); err != nil {
			zap.L().Error(errorSubVideoCommentCount, zap.Error(err))
			//还原原来的信息 deleted_at
			if err := mysql.RevertComment(req.CommentId); err != nil {
				return nil, err
			}
			return nil, err
		}
		//将评论信息从redis删除
		if err := redis.DeleteCommentInfo(req.VideoId, res.AuthorId, req.CommentId); err != nil {
			zap.L().Error(errorDeleteCommentInfo, zap.Error(err))
			//还原原来的信息 deleted_at
			if err := mysql.RevertComment(req.CommentId); err != nil {
				return nil, err
			}
			//视频评论数加一
			if err := redis.AddVideoCommentCount(req.VideoId); err != nil {
				zap.L().Error(errorAddVideoCommentCount, zap.Error(err))
				return nil, err
			}
			return nil, err
		}
		//返回信息
		return nil, nil
	}
	//评论从mysql评论表删除
	if err := mysql.DeleteCommentInfo(req.CommentId); err != nil {
		zap.L().Error(errorDeleteCommentInfo, zap.Error(err))
		return nil, err
	}
	//视频评论数减一
	if err := redis.SubVideoCommentCount(req.VideoId); err != nil {
		//还原原来的信息 deleted_at
		if err := mysql.RevertComment(req.CommentId); err != nil {
			return nil, err
		}
		zap.L().Error(errorSubVideoCommentCount, zap.Error(err))
		return nil, err
	}
	//返回信息
	return nil, nil
}

func GetCommentVideoList(req *request.DouyinCommentListRequest) ([]*response.Comment, error) {
	//判断视频作者是否为大V并将视频作者ID返回
	res, err := grpc_client.VideoClient.JudgeVideoAuthor(context.Background(), &request2.DouyinJudgeVideoAuthorRequest{
		VideoId: req.VideoId,
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
	if res.IsV == true {
		//从redis中读大V视频评论
		comments, err := redis.GetCommentInfo(req.VideoId, res.AuthorId)
		if err != nil {
			zap.L().Error(errorGetCommentInfo, zap.Error(err))
			return nil, err
		}
		//如果没有评论,直接返回
		if len(comments) == 0 {
			return nil, nil
		}
		var cs []model.Commit
		var userId []int64
		for i := 0; i < len(comments); i++ {
			var c model.Commit
			if err := json.Unmarshal([]byte(comments[i]), &c); err != nil {
				zap.L().Error(errorJsonMarshal, zap.Error(err))
				return nil, err
			}
			userId = append(userId, c.UserID)
			cs = append(cs, c)
		}
		res1, err := grpc_client.UserClient.GetUserInfoList(context.Background(), &request1.DouyinUserListRequest{
			UserId:      userId,
			LoginUserId: req.LoginUserId,
		})
		var results []*response.Comment
		for i := 0; i < len(cs); i++ {
			var result = response.Comment{
				Id:         cs[i].CommitID,
				User:       res1.User[i],
				Content:    cs[i].Message,
				CreateDate: cs[i].CreatedAt.Format("2006-01-02 15:04:05"),
			}
			results = append(results, &result)
		}
		return results, nil
	}
	//从mysql中读大V视频评论
	cs := make([]model.Commit, 0)
	if err := mysql.GetCommentInfo(&cs, req.VideoId); err != nil {
		zap.L().Error(errorGetCommentInfo, zap.Error(err))
		return nil, err
	}
	if len(cs) == 0 {
		zap.L().Info("该视频没有评论")
		return nil, nil
	}
	var userId []int64
	for i := 0; i < len(cs); i++ {
		userId = append(userId, cs[i].UserID)
	}
	res1, err := grpc_client.UserClient.GetUserInfoList(context.Background(), &request1.DouyinUserListRequest{
		UserId:      userId,
		LoginUserId: req.LoginUserId,
	})
	var results []*response.Comment
	for i := 0; i < len(cs); i++ {
		var result = response.Comment{
			Id:         cs[i].CommitID,
			User:       res1.User[i],
			Content:    cs[i].Message,
			CreateDate: cs[i].CreatedAt.Format("2006-01-02 15:04:05"),
		}
		results = append(results, &result)
	}
	return results, nil
}

func GetCommentCount(req *request.DouyinCommentCountRequest) ([]int64, error) {
	if len(req.VideoId) == 0 {
		return nil, nil
	}
	countList, err := redis.GetCommentCount(req.VideoId)
	if err != nil {
		zap.L().Error(errorGetCommentCount, zap.Error(err))
		return nil, err
	}
	return countList, nil
}

func PushVCommentBasicInfoInit(req *request.DouyinPushVCommentBasicInfoInitRequest) error {
	for i := 0; i < len(req.VideoIdList); i++ {
		cs := make([]model.Commit, 0)
		if err := mysql.GetCommentInfo(&cs, req.VideoIdList[i]); err != nil {
			zap.L().Error(errorGetCommentInfo, zap.Error(err))
			return err
		}
		if len(cs) != 0 {
			if err := redis.PushVCommentBasicInfoInit(req.UserId, req.VideoIdList[i], cs); err != nil {
				zap.L().Error(errorPushVCommentBasicInfoInit, zap.Error(err))
			}
		}
	}
	return nil
}
