package service

import (
	"douyin/proto/follow/request"
	response1 "douyin/proto/follow/response"
	"douyin/service/follow/dao/mysql"
	"douyin/service/follow/dao/redis"
	"douyin/service/follow/model"
	"go.uber.org/zap"
)

func GetFollowInfo(req *request.DouyinGetFollowRequest) (bool, error) {
	b, err := mysql.JudgeUserIsFollow(req)
	if err != nil {
		zap.L().Error(errorJudgeUserIsFollow, zap.Error(err))
		return false, err
	}
	return b, nil
}

func GetFollowInfoList(req *request.DouyinGetFollowListRequest) ([]bool, error) {
	if len(req.ToUserId) == 0 {
		return nil, nil
	}
	boolList, err := mysql.JudgeUserIsFollowList(req)
	if err != nil {
		zap.L().Error(errorJudgeUserIsFollowList, zap.Error(err))
		return nil, err
	}
	return boolList, nil
}

func GetFollowFollower(req *request.DouyinFollowFollowerCountRequest) (int64, int64, error) {
	if req.UserId == 0 {
		return 0, 0, nil
	}
	count1, count2, err := redis.GetUserFollowFollowerCount(req.UserId)
	if err != nil {
		zap.L().Error(errorGetUserFollowFollowerCount, zap.Error(err))
		return 0, 0, err
	}
	return count1, count2, nil
}

func GetFollowFollowerList(req *request.DouyinFollowFollowerListCountRequest) ([]int64, []int64, error) {
	if len(req.UserId) == 0 {
		return nil, nil, nil
	}
	count1List, count2List, err := redis.GetUserFollowFollowerCountList(req.UserId)
	if err != nil {
		zap.L().Error(errorGetUserFollowFollowerCountList, zap.Error(err))
		return nil, nil, err
	}
	return count1List, count2List, nil
}

func GetFollowFollowerIdList(req *request.DouyinGetFollowFollowerIdListRequest) ([]*response1.FollowList, []*response1.FollowerList, error) {
	if len(req.UserId) == 0 {
		return nil, nil, nil
	}
	var result1 []*response1.FollowList
	var result2 []*response1.FollowerList
	for i := 0; i < len(req.UserId); i++ {
		fs1 := make([]model.Follow, 0)
		if err := mysql.GetUserFollowList(&fs1, req.UserId[i]); err != nil {
			zap.L().Error(errorGetUserFollowList, zap.Error(err))
			return nil, nil, err
		}
		var idList1 []int64
		for i := 0; i < len(fs1); i++ {
			idList1 = append(idList1, fs1[i].FollowID)
		}
		var r1 = response1.FollowList{
			UserId: idList1,
		}
		fs2 := make([]model.Follow, 0)
		if err := mysql.GetUserFollowerList(&fs2, req.UserId[i]); err != nil {
			zap.L().Error(errorGetUserFollowList, zap.Error(err))
			return nil, nil, err
		}
		var idList2 []int64
		for i := 0; i < len(fs2); i++ {
			idList2 = append(idList2, fs2[i].FollowerID)
		}
		var r2 = response1.FollowerList{
			UserId: idList2,
		}
		result1 = append(result1, &r1)
		result2 = append(result2, &r2)
	}
	return result1, result2, nil
}
