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
	"errors"
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
	errorConnectToGRPCServer            = "connect to grpc server failed"
	errorExeVActiveSet                  = "exe if request is V or active from redis set failed"
	errorPublishVUserVideo              = "publish v user video information to redis failed"
	errorPublishActiveUserVideo         = "publish active user video information to redis failed"
	errorPushUserInfoList               = "push user information to redis list failed"
	errorCreateVideo                    = "create video failed"
	errorDeleteVideo                    = "delete video failed"
	errorIncrUserVideoCount             = "incr user video count failed"
	errorSubUserVideoCount              = "sub user video count failed"
	errorGetUserVideoCount              = "get user video count failed"
	errorAddUserVideoCountSet           = "add user video count set failed"
	errorGetVVideoInfo                  = "get v video information from redis failed"
	errorGetActiveVideoInfo             = "get active video information from redis failed"
	errorJsonMarshal                    = "json marshal failed"
	errorGetVideoInfo                   = "get video information failed"
	errorGetFavoriteCount               = "get favorite count and judge isFavorite failed"
	errorGetCommentCount                = "get comment count failed"
	errorGetUserInfo                    = "get user information failed"
	errorJudgeVideoAuthor               = "judge video author failed"
	errorGetVideoByVideoId              = "get video information by video id failed"
	errorGetVVideoInfoByVideoId         = "get v video information by video id failed"
	errorGetActiveVideoInfoByVideoId    = "get active video information by video id failed"
	errorPushVFavoriteVideoInfo         = "push v favorite video information failed"
	errorPushActiveFavoriteVideoInfo    = "push active favorite video information failed"
	errorGetVFavoriteVideoInfo          = "get v favorite video information failed"
	errorGetActiveFavoriteVideoInfo     = "get active favorite video information failed"
	errorDeleteVFavoriteVideoInfo       = "delete v favorite video information failed"
	errorDeleteActiveFavoriteVideoInfo  = "delete active favorite video information failed"
	errorGetVideoListByVideoId          = "get video list by video id failed"
	errorGetVideoInfoZSet               = "get video information zSet failed"
	errorJsonUnMarshal                  = "json unmarshal failed"
	errorGetVideoListByCreateTime       = "get video list by create time failed"
	errorGetUserFavoriteVideoIdList     = "get user favorite video id list failed"
	errorPushVActiveBasicInfoInit       = "push v or active basic info init failed"
	errorGetUserListFavoriteVideoIdList = "get user list favorite video id list failed"
	errorJudgeFeedIsFirst               = "judge feed is first failed"
	errorDeleteActiveVideoInfo          = "delete active video information failed"
	errorGetUserPublishVideoCountList   = "get user publish video count list failed"
	errorGetUserPublishVideoCount       = "get user publish video count failed"
	errorDeleteVUserVideo               = "delete v user video failed"
	errorDeleteActiveUserVideo          = "delete active user video failed"
)

//以发布视频为例，测试一下分布式事务

func PublishVideoDtm(req *request.DouyinPublishActionRequest) error {
	wfName := "workflow-publish video" + shortuuid.New()
	err := workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
		var request request1.DouyinUserVideoCountSetRequest
		if err := proto.Unmarshal(data, &request); err != nil {
			zap.L().Error("proto unmarshal failed", zap.Error(err))
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
			zap.L().Info("发布视频接口调用用户加入用户发布视频集合回滚")
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
		zap.L().Info("发布视频接口调用用户加入用户发布视频集合成功")
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
				zap.L().Info("视频信息从mysql表删除成功")
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
					zap.L().Info("视频信息从mysql表删除成功")
					//用户视频数减一
					if err := redis.SubUserVideoCount(v.UserID); err != nil {
						zap.L().Error(errorSubUserVideoCount, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					zap.L().Info("用户视频数减一成功")
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
					zap.L().Info("视频信息从mysql表删除成功")
					//用户视频数减一
					if err := redis.SubUserVideoCount(v.UserID); err != nil {
						zap.L().Error(errorSubUserVideoCount, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					zap.L().Info("用户视频数减一成功")
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
					zap.L().Info("删除存入大V的发布视频信息")
				}
				if res.IsActiver == true {
					if err := redis.DeleteActiveUserVideo(v.UserID, v.VideoID); err != nil {
						zap.L().Error(errorDeleteActiveUserVideo, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					zap.L().Info("删除存入活跃用户的发布视频信息")
				}
				//将视频信息从mysql表删除
				if err := mysql.DeleteVideoInfo(v.VideoID); err != nil {
					zap.L().Error(errorDeleteVideo, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
				zap.L().Info("视频信息从mysql表删除成功")
				//用户视频数减一
				if err := redis.SubUserVideoCount(v.UserID); err != nil {
					zap.L().Error(errorSubUserVideoCount, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
				zap.L().Info("用户视频数减一成功")
				return nil, status.Error(codes.Aborted, err.Error())
			}
			zap.L().Info("用户发布视频本地事务执行完成")
			return nil, nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		zap.L().Error("workflow register failed", zap.Error(err))
		return err
	}
	request := request1.DouyinUserVideoCountSetRequest{
		UserId: req.UserId,
	}
	data, err := proto.Marshal(&request)
	if err != nil {
		zap.L().Error("proto unmarshal failed", zap.Error(err))
		return err
	}
	if err = workflow.Execute(wfName, shortuuid.New(), data); err != nil {
		zap.L().Error("result of workflow.Execute is", zap.Error(err))
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

func JudgeVideoAuthor(req *request.DouyinJudgeVideoAuthorRequest) (int64, bool, bool, error) {
	authorId, isV, isActive, err := redis.JudgeVideoAuthor(req.VideoId)
	authorId1, _ := strconv.ParseInt(authorId, 10, 64)
	if err != nil {
		zap.L().Error(errorJudgeVideoAuthor, zap.Error(err))
		return 0, false, false, err
	}
	if isV == false && isActive == false {
		var video model.Video
		if err := mysql.GetVideoByVideoId(&video, req.VideoId); err != nil {
			zap.L().Error(errorGetVideoByVideoId, zap.Error(err))
			return 0, false, false, err
		}
		return video.UserID, false, false, nil
	} else {
		return authorId1, isV, isActive, nil
	}
}

func PushVActiveFavoriteVideo(req *request.DouyinPushVActiveFavoriteVideoRequest) error {
	var data string
	if req.IsV == true {
		videoData, err := redis.GetVVideoInfoByVideoId(req.AuthorId, req.VideoId)
		if err != nil {
			zap.L().Error(errorGetVVideoInfoByVideoId, zap.Error(err))
			return status.Error(codes.Aborted, err.Error())
		}
		data = videoData
	}
	if req.IsActive == true {
		videoData, err := redis.GetActiveVideoInfoByVideoId(req.AuthorId, req.VideoId)
		if err != nil {
			zap.L().Error(errorGetActiveVideoInfoByVideoId, zap.Error(err))
			return status.Error(codes.Aborted, err.Error())
		}
		data = videoData
	}
	if req.IsV == false && req.IsActive == false {
		var video model.Video
		if err := mysql.GetVideoByVideoId(&video, req.VideoId); err != nil {
			zap.L().Error(errorGetVideoByVideoId, zap.Error(err))
			return status.Error(codes.Aborted, err.Error())
		}
		videoData, err := json.Marshal(video)
		if err != nil {
			zap.L().Error(errorJsonMarshal, zap.Error(err))
			return status.Error(codes.Aborted, err.Error())
		}
		data = string(videoData)
	}
	//登录用户为大V
	//首先要根据视频作者的身份获取视频基本信息，然后再存到redis中
	if req.LoginUserIsV == true {
		if err := redis.PushVFavoriteVideoInfo(data, req.LoginUserId, req.VideoId); err != nil {
			zap.L().Error(errorPushVFavoriteVideoInfo, zap.Error(err))
			return status.Error(codes.Aborted, err.Error())
		}
	}
	//登录用户为活跃用户
	if req.LoginUserIsActive == true {
		if err := redis.PushActiveFavoriteVideoInfo(data, req.LoginUserId, req.VideoId); err != nil {
			zap.L().Error(errorPushActiveFavoriteVideoInfo, zap.Error(err))
			return status.Error(codes.Aborted, err.Error())
		}
	}
	return nil
}

func PushVActiveFavoriteVideoRevert(req *request.DouyinPushVActiveFavoriteVideoRequest) error {
	if req.LoginUserIsV == true {
		if err := redis.DeleteVFavoriteVideoInfo(req.LoginUserId, req.VideoId); err != nil {
			zap.L().Error(errorDeleteVFavoriteVideoInfo, zap.Error(err))
			return status.Error(codes.Aborted, err.Error())
		}
	}
	//登录用户为活跃用户
	if req.LoginUserIsActive == true {
		if err := redis.DeleteActiveFavoriteVideoInfo(req.LoginUserId, req.VideoId); err != nil {
			zap.L().Error(errorDeleteActiveFavoriteVideoInfo, zap.Error(err))
			return status.Error(codes.Aborted, err.Error())
		}
	}
	return nil
}

func GetVActiveFavoriteVideo(req *request.DouyinGetVActiveFavoriteVideoRequest) ([]*response2.Simple_Video, error) {
	if req.IsV == true {
		data, err := redis.GetVFavoriteVideoInfo(req.UserId)
		if err != nil {
			zap.L().Error(errorGetVFavoriteVideoInfo, zap.Error(err))
			return nil, err
		}
		if len(data) == 0 {
			return nil, nil
		}
		var results []*response2.Simple_Video
		for i := 0; i < len(data); i++ {
			var result = response2.Simple_Video{
				Id:       data[i].VideoID,
				UserId:   data[i].UserID,
				PlayUrl:  data[i].VideoUrl,
				CoverUrl: data[i].CoverUrl,
				Title:    data[i].Title,
			}
			results = append(results, &result)
		}
		return results, nil
	}
	if req.IsActive == true {
		data, err := redis.GetActiveFavoriteVideoInfo(req.UserId)
		if err != nil {
			zap.L().Error(errorGetActiveFavoriteVideoInfo, zap.Error(err))
			return nil, err
		}
		if len(data) == 0 {
			return nil, nil
		}
		var results []*response2.Simple_Video
		for i := 0; i < len(data); i++ {
			var result = response2.Simple_Video{
				Id:       data[i].VideoID,
				UserId:   data[i].UserID,
				PlayUrl:  data[i].VideoUrl,
				CoverUrl: data[i].CoverUrl,
				Title:    data[i].Title,
			}
			results = append(results, &result)
		}
		return results, nil
	}
	return nil, errors.New("调用rpc输入参数错误")
}

func DeleteVActiveFavoriteVideo(req *request.DouyinDeleteVActiveFavoriteVideoRequest) error {
	if req.LoginUserIsV == true {
		if err := redis.DeleteVFavoriteVideoInfo(req.UserId, req.VideoId); err != nil {
			zap.L().Error(errorDeleteVFavoriteVideoInfo, zap.Error(err))
			return status.Error(codes.Aborted, err.Error())
		}
		return nil
	}
	if req.LoginUserIsActive == true {
		if err := redis.DeleteActiveFavoriteVideoInfo(req.UserId, req.VideoId); err != nil {
			zap.L().Error(errorDeleteActiveFavoriteVideoInfo, zap.Error(err))
			return status.Error(codes.Aborted, err.Error())
		}
		return nil
	}
	return status.Error(codes.Aborted, "调用rpc输入参数错误")
}

func GetVideoListInner(req *request.DouyinGetVideoListRequest) ([]*response2.Simple_Video, error) {
	videos := make([]model.Video, 0, len(req.VideoId))
	if err := mysql.GetVideoListByVideoId(&videos, req.VideoId); err != nil {
		zap.L().Error(errorGetVideoListByVideoId, zap.Error(err))
		return nil, err
	}
	var results []*response2.Simple_Video
	for i := 0; i < len(videos); i++ {
		var result = response2.Simple_Video{
			Id:       videos[i].VideoID,
			UserId:   videos[i].UserID,
			PlayUrl:  videos[i].VideoUrl,
			CoverUrl: videos[i].CoverUrl,
			Title:    videos[i].Title,
		}
		results = append(results, &result)
	}
	return results, nil
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

func PushVActiveBasicInfoInit(req *request.DouyinPushVInfoInitRequest) ([]int64, error) {
	//存大V
	//找发布的视频信息
	vs1 := make([]model.Video, 0)
	if err := mysql.GetVideoInfo(req.UserId, &vs1); err != nil {
		zap.L().Error(errorGetVideoInfo, zap.Error(err))
		return nil, err
	}
	var videoList1 []int64
	for i := 0; i < len(vs1); i++ {
		videoList1 = append(videoList1, vs1[i].VideoID)
	}
	if req.IsActive == true {
		if err := redis.DeleteActiveVideoInfo(req.UserId, videoList1); err != nil {
			zap.L().Error(errorDeleteActiveVideoInfo, zap.Error(err))
			return nil, err
		}
	}
	//找点赞的视频id
	res, err := grpc_client.FavoriteClient.GetUserFavoriteVideoIdList(context.Background(), &request2.DouyinFavoriteIdListRequest{
		UserId: req.UserId,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorGetUserFavoriteVideoIdList, zap.Error(err))
			return nil, err
		}
		if res.Code == 2 {
			zap.L().Error(errorGetUserInfo, zap.Error(err))
			return nil, err
		}
	}
	vs2 := make([]model.Video, 0)
	if len(res.VideoId) != 0 {
		if err = mysql.GetVideoListByVideoId(&vs2, res.VideoId); err != nil {
			zap.L().Error(errorGetVideoListByVideoId, zap.Error(err))
			return nil, err
		}
	}
	var videoList2 []int64
	for i := 0; i < len(vs2); i++ {
		videoList2 = append(videoList2, vs2[i].VideoID)
	}
	//存三个redis
	if err = redis.PushVBasicInfoInit(vs1, vs2, videoList1, videoList2, req.UserId); err != nil {
		zap.L().Error(errorPushVActiveBasicInfoInit, zap.Error(err))
		return nil, err
	}
	return videoList1, nil
}

func PushActiveBasicInfoInit(req *request.DouyinPushActiveInfoInitRequest) error {
	//找发布的视频信息
	var vss [][]model.Video
	for i := 0; i < len(req.UserId); i++ {
		vs1 := make([]model.Video, 0)
		if err := mysql.GetVideoInfo(req.UserId[i], &vs1); err != nil {
			zap.L().Error(errorGetVideoInfo, zap.Error(err))
			return err
		}
		vss = append(vss, vs1)
	}
	//找点赞的视频id
	res, err := grpc_client.FavoriteClient.GetUserListFavoriteVideoIdList(context.Background(), &request2.DouyinFavoriteListIdListRequest{
		UserId: req.UserId,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorGetUserFavoriteVideoIdList, zap.Error(err))
			return err
		}
		if res.Code == 2 {
			zap.L().Error(errorGetUserListFavoriteVideoIdList, zap.Error(err))
			return err
		}
	}
	for i := 0; i < len(req.UserId); i++ {
		var userPublishId []int64
		for j := 0; j < len(vss[i]); j++ {
			userPublishId = append(userPublishId, vss[i][j].VideoID)
		}
		var userFavoriteId []int64
		for k := 0; k < len(res.FavoriteVideoIdList[i].VideoId); k++ {
			userFavoriteId = append(userFavoriteId, res.FavoriteVideoIdList[i].VideoId[k])
		}
		vs2 := make([]model.Video, 0)
		if err = mysql.GetVideoListByVideoId(&vs2, userFavoriteId); err != nil {
			zap.L().Error(errorGetVideoListByVideoId, zap.Error(err))
			return err
		}
		//存三个redis
		if err = redis.PushActiveBasicInfoInit(vss[i], vs2, userPublishId, userFavoriteId, req.UserId[i]); err != nil {
			zap.L().Error(errorPushVActiveBasicInfoInit, zap.Error(err))
			return err
		}
	}
	return nil
}

func GetUserPublishVideoCount(req *request.DouyinGetUserPublishVideoCountRequest) (int64, error) {
	if req.UserId == 0 {
		return 0, nil
	}
	count, err := redis.GetUserPublishVideoCount(req.UserId)
	if err != nil {
		zap.L().Error(errorGetUserPublishVideoCount, zap.Error(err))
		return 0, err
	}
	return count, nil
}

func GetUserPublishVideoCountList(req *request.DouyinGetUserPublishVideoCountListRequest) ([]int64, error) {
	if len(req.UserId) == 0 {
		return nil, nil
	}
	counts, err := redis.GetUserPublishVideoCountList(req.UserId)
	if err != nil {
		zap.L().Error(errorGetUserPublishVideoCountList, zap.Error(err))
		return nil, err
	}
	return counts, nil
}
