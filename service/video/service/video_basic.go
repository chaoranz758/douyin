package service

import (
	"context"
	request3 "douyin/proto/comment/request"
	request2 "douyin/proto/favorite/request"
	request1 "douyin/proto/user/request"
	"douyin/proto/video/request"
	response2 "douyin/proto/video/response"
	"douyin/service/video/dao/mysql"
	"douyin/service/video/dao/redis"
	"douyin/service/video/initialize/grpc_client"
	"douyin/service/video/model"
	"douyin/service/video/util"
	"encoding/json"
	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmgrpc/workflow"
	"github.com/lithammer/shortuuid/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"strconv"
	"time"
)

const (
	wfName1 = "workflow-publish video"
)

func PublishVideoDtm(req *request.DouyinPublishActionRequest) error {
	wfName := wfName1 + shortuuid.New()
	err := workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
		var request request1.DouyinUserVideoCountSetRequest
		if err := proto.Unmarshal(data, &request); err != nil {
			zap.L().Error(errorJsonMarshal, zap.Error(err))
			return err
		}
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			//检查一下在不在集合里面，不再集合里面就删除
			res1, err := grpc_client.UserClientDtm.AddUserVideoCountSetRevert(wf.Context, &request)
			if err != nil {
				if res1 == nil {
					zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
					return err
				}
				if res1.Code == 2 {
					zap.L().Error(errorAddUserVideoCountSet, zap.Error(err))
					return err
				}
			}
			//zap.L().Info("发布视频接口调用用户加入用户发布视频集合回滚")
			return nil
		})
		//在分支事务中自己先检查一下数量 如果数量等于4就调用用户服务
		count, err := redis.GetUserVideoCount(req.UserId)
		if err != nil {
			zap.L().Error(errorGetUserVideoCount, zap.Error(err))
			return err
		}
		if count == 4 {
			res1, err := grpc_client.UserClientDtm.AddUserVideoCountSet(wf.Context, &request)
			if err != nil {
				if res1 == nil {
					zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
					return err
				}
				if res1.Code == 2 {
					zap.L().Error(errorAddUserVideoCountSet, zap.Error(err))
					return err
				}
			}
		}
		//zap.L().Info("发布视频接口调用用户加入用户发布视频集合成功")
		_, err = wf.NewBranch().Do(func(bb *dtmcli.BranchBarrier) ([]byte, error) {
			v := new(model.Video)
			//生成视频ID
			videoID := util.GenID()
			res, err := grpc_client.UserClient.UserIsInfluencerActiver(context.Background(), &request1.DouyinUserIsInfluencerActiverRequest{
				UserId: req.UserId,
			})
			if err != nil {
				if res == nil {
					zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
				if res.Code == 2 {
					zap.L().Error(errorExeVActiveSet, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
			}
			v.VideoID = videoID
			v.UserID = req.UserId
			v.Title = req.Title
			v.VideoUrl = req.VideoUrl
			v.CoverUrl = req.CoverUrl
			v.CreatedAt = time.Now()
			//将视频信息写入mysql视频表
			if err := mysql.CreateVideoInfo(v); err != nil {
				zap.L().Error(errorCreateVideo, zap.Error(err))
				return nil, status.Error(codes.Aborted, err.Error())
			}
			//redis中用户视频数加一并返回当前数量
			_, err = redis.IncrUserVideoCount(v.UserID)
			if err != nil {
				zap.L().Error(errorIncrUserVideoCount, zap.Error(err))
				//将视频信息从mysql表删除
				if err := mysql.DeleteVideoInfo(v.VideoID); err != nil {
					zap.L().Error(errorDeleteVideo, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
				//zap.L().Info("视频信息从mysql表删除成功")
				return nil, status.Error(codes.Aborted, err.Error())
			}
			//若投稿用户是大V就将投稿信息写入大V发布的视频信息hash中并将其加入大V发布的userID:videoID的hash
			if res.IsInfluencer == true {
				if err := redis.PublishVUserVideo(v); err != nil {
					zap.L().Error(errorPublishVUserVideo, zap.Error(err))
					//将视频信息从mysql表删除
					if err := mysql.DeleteVideoInfo(v.VideoID); err != nil {
						zap.L().Error(errorDeleteVideo, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					//zap.L().Info("视频信息从mysql表删除成功")
					//用户视频数减一
					if err := redis.SubUserVideoCount(v.UserID); err != nil {
						zap.L().Error(errorSubUserVideoCount, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					//zap.L().Info("用户视频数减一成功")
					return nil, status.Error(codes.Aborted, err.Error())
				}
			}
			//若投稿用户是活跃用户就将投稿信息写入活跃用户发布的视频信息hash中并将其加入活跃用户发布的userID:videoID的hash
			if res.IsActiver == true {
				if err := redis.PublishActiveUserVideo(v); err != nil {
					zap.L().Error(errorPublishActiveUserVideo, zap.Error(err))
					//将视频信息从mysql表删除
					if err := mysql.DeleteVideoInfo(v.VideoID); err != nil {
						zap.L().Error(errorDeleteVideo, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					//zap.L().Info("视频信息从mysql表删除成功")
					//用户视频数减一
					if err := redis.SubUserVideoCount(v.UserID); err != nil {
						zap.L().Error(errorSubUserVideoCount, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					//zap.L().Info("用户视频数减一成功")
					return nil, status.Error(codes.Aborted, err.Error())
				}
			}
			//将用户信息写入300个视频list并修剪使其长度满足300
			if err := redis.PublishVideoInfoZSet(v, v.CreatedAt.UnixMilli()); err != nil {
				zap.L().Error(errorPushUserInfoList)
				if res.IsInfluencer == true {
					if err := redis.DeleteVUserVideo(v.UserID, v.VideoID); err != nil {
						zap.L().Error(errorDeleteVUserVideo, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					//zap.L().Info("删除存入大V的发布视频信息")
				}
				if res.IsActiver == true {
					if err := redis.DeleteActiveUserVideo(v.UserID, v.VideoID); err != nil {
						zap.L().Error(errorDeleteActiveUserVideo, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					//zap.L().Info("删除存入活跃用户的发布视频信息")
				}
				//将视频信息从mysql表删除
				if err := mysql.DeleteVideoInfo(v.VideoID); err != nil {
					zap.L().Error(errorDeleteVideo, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
				//zap.L().Info("视频信息从mysql表删除成功")
				//用户视频数减一
				if err := redis.SubUserVideoCount(v.UserID); err != nil {
					zap.L().Error(errorSubUserVideoCount, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
				//zap.L().Info("用户视频数减一成功")
				return nil, status.Error(codes.Aborted, err.Error())
			}
			//zap.L().Info("用户发布视频本地事务执行完成")
			return nil, nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		zap.L().Error(errorRegisterWorkflow, zap.Error(err))
		return err
	}
	request := request1.DouyinUserVideoCountSetRequest{
		UserId: req.UserId,
	}
	data, err := proto.Marshal(&request)
	if err != nil {
		zap.L().Error(errorJsonMarshal, zap.Error(err))
		return err
	}
	if err = workflow.Execute(wfName, shortuuid.New(), data); err != nil {
		zap.L().Error(errorExecuteWorkflow, zap.Error(err))
		return err
	}
	return nil
}

func GetPublishVideo(req *request.DouyinPublishListRequest) ([]*response2.Video, error) {
	vs := make([]model.Video, 0)
	res, err := grpc_client.UserClient.UserIsInfluencerActiver(context.Background(), &request1.DouyinUserIsInfluencerActiverRequest{
		UserId: req.UserId,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res.Code == 2 {
			zap.L().Error(errorExeVActiveSet, zap.Error(err))
			return nil, err
		}
	}
	if res.IsInfluencer == true {
		m, err := redis.GetVVideoInfo(req.UserId)
		if err != nil {
			zap.L().Error(errorGetVVideoInfo, zap.Error(err))
			return nil, err
		}
		if len(m) == 0 {
			return nil, nil
		}
		for i := 0; i < len(m); i++ {
			var v1 model.Video
			if err := json.Unmarshal([]byte(m[i]), &v1); err != nil {
				zap.L().Error(errorJsonMarshal, zap.Error(err))
				return nil, err
			}
			vs = append(vs, v1)
		}
	}
	if res.IsActiver == true {
		m, err := redis.GetActiveVideoInfo(req.UserId)
		if err != nil {
			zap.L().Error(errorGetActiveVideoInfo, zap.Error(err))
			return nil, err
		}
		if len(m) == 0 {
			return nil, nil
		}
		for i := 0; i < len(m); i++ {
			var v1 model.Video
			if err := json.Unmarshal([]byte(m[i]), &v1); err != nil {
				zap.L().Error(errorJsonMarshal, zap.Error(err))
				return nil, err
			}
			vs = append(vs, v1)
		}
	}
	if err := mysql.GetVideoInfo(req.UserId, &vs); err != nil {
		zap.L().Error(errorGetVideoInfo, zap.Error(err))
		return nil, err
	}
	var idList []int64
	for i := 0; i < len(vs); i++ {
		idList = append(idList, vs[i].VideoID)
	}
	//从点赞service读视频的点赞数并查看用户是否对视频点赞
	res1, err := grpc_client.FavoriteClient.GetFavoriteCount(context.Background(), &request2.DouyinFavoriteCountRequest{
		VideoId: idList,
		UserId:  req.LoginId,
	})
	if err != nil {
		if res1 == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res1.Code == 2 {
			zap.L().Error(errorGetFavoriteCount, zap.Error(err))
			return nil, err
		}
	}
	//从评论service读视频的评论数
	res2, err := grpc_client.CommentClient.GetCommentCount(context.Background(), &request3.DouyinCommentCountRequest{
		VideoId: idList,
	})
	if err != nil {
		if res2 == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res2.Code == 2 {
			zap.L().Error(errorGetCommentCount, zap.Error(err))
			return nil, err
		}
	}
	//从用户服务获取用户信息,只需获取一个，因为每个视频的作者都是同一个人
	res3, err := grpc_client.UserClient.GetUserInfo(context.Background(), &request1.DouyinUserRequest{
		UserId:      req.UserId,
		LoginUserId: req.LoginId,
	})
	if err != nil {
		if res3 == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res3.Code == 2 {
			zap.L().Error(errorGetUserInfo, zap.Error(err))
			return nil, err
		}
	}
	var videos []*response2.Video
	for i := 0; i < len(vs); i++ {
		var video = response2.Video{
			Id:            vs[i].VideoID,
			Author:        res3.User,
			PlayUrl:       vs[i].VideoUrl,
			CoverUrl:      vs[i].CoverUrl,
			FavoriteCount: res1.FavoriteCount[i],
			CommentCount:  res2.CommentCount[i],
			IsFavorite:    res1.IsFavorite[i],
			Title:         vs[i].Title,
		}
		videos = append(videos, &video)
	}
	return videos, nil
}

func GetVideoList(req *request.DouyinFeedRequest) ([]*response2.Video, int64, error) {
	var results []*response2.Video
	var userIdList []int64
	var videoIdList []int64
	var lastTime int64
	var lastTime1 int64
	t := strconv.FormatInt(req.LatestTime, 10)
	countFirst, err := redis.JudgeFeedIsFirst(t)
	if err != nil {
		zap.L().Error(errorJudgeFeedIsFirst, zap.Error(err))
		return nil, 0, err
	}
	if countFirst == 0 {
		data, err := redis.GetVideoInfoZSetInitial()
		if err != nil {
			zap.L().Error(errorGetVideoInfoZSet, zap.Error(err))
			return nil, 0, err
		}
		for i := 0; i < len(data); i++ {
			var video model.Video
			if err := json.Unmarshal([]byte(data[i]), &video); err != nil {
				zap.L().Error(errorJsonUnMarshal, zap.Error(err))
				return nil, 0, err
			}
			var result = response2.Video{
				Id:       video.VideoID,
				PlayUrl:  video.VideoUrl,
				CoverUrl: video.CoverUrl,
				Title:    video.Title,
			}
			results = append(results, &result)
			userIdList = append(userIdList, video.UserID)
			videoIdList = append(videoIdList, video.VideoID)
			if i == len(data)-1 {
				lastTime = video.CreatedAt.UnixMilli()
			}
		}
	}
	if countFirst != 0 {
		times := strconv.FormatInt(req.LatestTime, 10)
		countVideo, data, err := redis.GetVideoInfoZSet(times)
		if err != nil {
			zap.L().Error(errorGetVideoInfoZSet, zap.Error(err))
			return nil, 0, err
		}
		for i := 1; i < len(data); i++ {
			var video model.Video
			if err := json.Unmarshal([]byte(data[i]), &video); err != nil {
				zap.L().Error(errorJsonUnMarshal, zap.Error(err))
				return nil, 0, err
			}
			var result = response2.Video{
				Id:       video.VideoID,
				PlayUrl:  video.VideoUrl,
				CoverUrl: video.CoverUrl,
				Title:    video.Title,
			}
			results = append(results, &result)
			userIdList = append(userIdList, video.UserID)
			videoIdList = append(videoIdList, video.VideoID)
			if i == len(data)-1 {
				lastTime = video.CreatedAt.UnixMilli()
				lastTime1 = video.CreatedAt.UnixMilli()
			}
		}
		count := len(data)
		if count < 31 && countVideo > 300 {
			videos := make([]model.Video, 0, 30)
			if err := mysql.GetVideoListByCreateTime(&videos, lastTime1, 31-count); err != nil {
				zap.L().Error(errorGetVideoListByCreateTime, zap.Error(err))
				return nil, 0, err
			}
			for i := 0; i < len(videos); i++ {
				var result = response2.Video{
					Id:       videos[i].VideoID,
					PlayUrl:  videos[i].VideoUrl,
					CoverUrl: videos[i].CoverUrl,
					Title:    videos[i].Title,
				}
				results = append(results, &result)
				userIdList = append(userIdList, videos[i].UserID)
				videoIdList = append(videoIdList, videos[i].VideoID)
			}
			lastTime = videos[len(videos)-1].CreatedAt.UnixMilli()
		}
	}
	if len(videoIdList) == 0 {
		return nil, 0, nil
	}
	//从点赞service读视频的点赞数并查看用户是否对视频点赞
	res1, err := grpc_client.FavoriteClient.GetFavoriteCount(context.Background(), &request2.DouyinFavoriteCountRequest{
		VideoId: videoIdList,
		UserId:  req.LoginUserId,
	})
	if err != nil {
		if res1 == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, 0, err
		}
		if res1.Code == 2 {
			zap.L().Error(errorGetFavoriteCount, zap.Error(err))
			return nil, 0, err
		}
	}
	//从评论service读视频的评论数
	res2, err := grpc_client.CommentClient.GetCommentCount(context.Background(), &request3.DouyinCommentCountRequest{
		VideoId: videoIdList,
	})
	if err != nil {
		if res2 == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, 0, err
		}
		if res2.Code == 2 {
			zap.L().Error(errorGetCommentCount, zap.Error(err))
			return nil, 0, err
		}
	}
	res3, err := grpc_client.UserClient.GetUserInfoList(context.Background(), &request1.DouyinUserListRequest{
		UserId:      userIdList,
		LoginUserId: req.LoginUserId,
	})
	if err != nil {
		if res3 == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, 0, err
		}
		if res3.Code == 2 {
			zap.L().Error(errorGetUserInfo, zap.Error(err))
			return nil, 0, err
		}
	}
	for i := 0; i < len(userIdList); i++ {
		results[i].Author = res3.User[i]
		results[i].FavoriteCount = res1.FavoriteCount[i]
		results[i].IsFavorite = res1.IsFavorite[i]
		results[i].CommentCount = res2.CommentCount[i]
	}
	return results, lastTime, nil
}
