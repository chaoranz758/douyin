package service

import (
	"douyin/proto/comment/request"
	"douyin/service/comment/dao/mysql"
	"douyin/service/comment/dao/redis"
	"douyin/service/comment/model"
	"go.uber.org/zap"
)

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
