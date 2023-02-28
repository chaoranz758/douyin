package service

import (
	"context"
	request2 "douyin/proto/favorite/request"
	"douyin/proto/video/request"
	response2 "douyin/proto/video/response"
	"douyin/service/video/dao/mysql"
	"douyin/service/video/dao/redis"
	"douyin/service/video/initialize/grpc_client"
	"douyin/service/video/model"
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
	return nil, errors.New(errorRpcInput)
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
	return status.Error(codes.Aborted, errorRpcInput)
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
