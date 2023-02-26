package job

import (
	"context"
	"douyin/proto/follow/request"
	request1 "douyin/proto/video/request"
	"douyin/service/user/dao/mysql"
	"douyin/service/user/dao/redis"
	"douyin/service/user/initialize/grpc_client"
	"douyin/service/user/model"
	"go.uber.org/zap"
	"strconv"
)

const (
	errorPushActiveSetAndString     = "push active set and string failed"
	errorGetActiveSet               = "get active set failed"
	errorExeVActiveSet              = "exe if request is V or active from redis set failed"
	errorPushActiveFollowInfoInit   = "push active follow info init failed"
	errorPushActiveFollowerInfoInit = "push active follower info init failed"
	errorPushActiveBasicInfoInit    = "push active basic info init failed"
	errorJudgeIsInActiveSet         = "judge is in active set failed"
	errorGetUserInfoListFailed      = "get user information list failed"
	errorConnectToGRPCServer        = "connect to grpc server failed"
	errorGetFollowFollowerIdList    = "get follow follower id list failed"
)

func PushActiveSetAndString() error {
	//redis合并集合得到合并列表，然后根据列表去对比活跃者集合中不存在的那些用户id
	result, err := redis.GetActiveSet()
	if err != nil {
		zap.L().Error(errorGetActiveSet, zap.Error(err))
		return err
	}
	if len(result) == 0 {
		zap.L().Info("执行了定时任务但没有要更新的活跃用户")
		return nil
	}
	var userId []int64
	if len(result) != 0 {
		for i := 0; i < len(result); i++ {
			userId1, _ := strconv.ParseInt(result[i], 10, 64)
			isV, _, err := redis.IsVActive(userId1)
			if err != nil {
				zap.L().Error(errorExeVActiveSet, zap.Error(err))
				return err
			}
			if isV == false {
				userId = append(userId, userId1)
			}
		}
		if len(userId) == 0 {
			zap.L().Info("执行了定时任务有活跃用户但该用户已经是大V了")
			return nil
		}
		var userId1 []int64
		for i := 0; i < len(userId); i++ {
			isInActive, err := redis.JudgeIsInActiveSet(userId[i])
			if err != nil {
				zap.L().Error(errorJudgeIsInActiveSet, zap.Error(err))
				return err
			}
			if isInActive == false {
				userId1 = append(userId1, userId[i])
			}
		}
		if len(userId1) == 0 {
			zap.L().Info("执行了定时任务有活跃用户不是大V但该用户曾经是活跃用户已经写入过了")
			return nil
		}
		//从mysql读出这些用户id读出用户基本信息
		users := make([]model.User, 0)
		if err = mysql.GetUserInfoList(userId, &users); err != nil {
			zap.L().Error(errorGetUserInfoListFailed, zap.Error(err))
			return err
		}
		//将用户基本信息存入redis 2个
		if err = redis.PushActiveSetAndString(users, userId); err != nil {
			zap.L().Error(errorPushActiveSetAndString, zap.Error(err))
			return err
		}
	}
	//从follow处获取用户关注者、粉丝id列表
	res, err := grpc_client.FollowClient.GetFollowFollowerIdList(context.Background(), &request.DouyinGetFollowFollowerIdListRequest{
		UserId: userId,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return err
		}
		if res.Code == 2 {
			zap.L().Error(errorGetFollowFollowerIdList, zap.Error(err))
			return err
		}
	}
	//从mysql读用户基本信息
	for i := 0; i < len(res.FollowList); i++ {
		var userIdList []int64
		for j := 0; j < len(res.FollowList[i].UserId); j++ {
			userIdList = append(userIdList, res.FollowList[i].UserId[j])
		}
		if len(userIdList) == 0 {
			zap.L().Error("当前用户没有关注者")
		} else {
			us := make([]model.User, 0, len(userIdList))
			if err := mysql.GetUserInfoList(userIdList, &us); err != nil {
				zap.L().Error(errorGetUserInfoListFailed, zap.Error(err))
				return err
			}
			//将这些信息存入redis
			if err := redis.PushActiveFollowInfoInit(us, userId[i]); err != nil {
				zap.L().Error(errorPushActiveFollowInfoInit, zap.Error(err))
				return err
			}
		}
	}
	for i := 0; i < len(res.FollowerList); i++ {
		var userIdList []int64
		for j := 0; j < len(res.FollowerList[i].UserId); j++ {
			userIdList = append(userIdList, res.FollowerList[i].UserId[j])
		}
		if len(userIdList) == 0 {
			zap.L().Error("当前用户没有粉丝")
		} else {
			us := make([]model.User, 0, len(userIdList))
			if err := mysql.GetUserInfoList(userIdList, &us); err != nil {
				zap.L().Error(errorGetUserInfoListFailed, zap.Error(err))
				return err
			}
			//将这些信息存入redis
			if err := redis.PushActiveFollowerInfoInit(us, userId[i]); err != nil {
				zap.L().Error(errorPushActiveFollowerInfoInit, zap.Error(err))
				return err
			}
		}
	}
	//存视频相关信息
	res1, err := grpc_client.VideoClient.PushActiveBasicInfoInit(context.Background(), &request1.DouyinPushActiveInfoInitRequest{
		UserId: userId,
	})
	if err != nil {
		if res1 == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return err
		}
		if res1.Code == 2 {
			zap.L().Error(errorPushActiveBasicInfoInit, zap.Error(err))
			return err
		}
	}
	zap.L().Info("执行定时任务成功")
	return nil
}
