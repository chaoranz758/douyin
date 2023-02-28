package service

import (
	"douyin/proto/video/request"
	response2 "douyin/proto/video/response"
	"douyin/service/video/dao/mysql"
	"douyin/service/video/dao/redis"
	"douyin/service/video/model"
	"go.uber.org/zap"
	"strconv"
)

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
