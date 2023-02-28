package service

import (
	"context"
	request3 "douyin/proto/favorite/request"
	request1 "douyin/proto/follow/request"
	response1 "douyin/proto/follow/response"
	"douyin/proto/user/request"
	"douyin/proto/user/response"
	request2 "douyin/proto/video/request"
	"douyin/service/user/dao/mysql"
	"douyin/service/user/dao/redis"
	"douyin/service/user/initialize/grpc_client"
	"douyin/service/user/model"
	"go.uber.org/zap"
)

func GetUserInfoList(req *request.DouyinUserListRequest) ([]*response.User, error) {
	if len(req.UserId) == 0 {
		return nil, nil
	}
	result := make([]*response.User, 0, len(req.UserId))
	//批量读取用户粉丝数和关注数
	res, err := grpc_client.FollowClient.GetFollowFollowerList(context.Background(), &request1.DouyinFollowFollowerListCountRequest{
		UserId: req.UserId,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res.Code == 2 {
			zap.L().Error(errorGetFollowFollowerList, zap.Error(err))
			return nil, err
		}
	}
	//读取用户作品数量
	resGetUserPublishVideoCount, err := grpc_client.VideoClient.GetUserPublishVideoCountList(context.Background(), &request2.DouyinGetUserPublishVideoCountListRequest{
		UserId: req.UserId,
	})
	if err != nil {
		if resGetUserPublishVideoCount == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if resGetUserPublishVideoCount.Code == 2 {
			zap.L().Error(errorGetUserPublishVideoCount, zap.Error(err))
			return nil, err
		}
	}
	//读取用户点赞视频数量和获赞总数
	resGetUserFavoritedCount, err := grpc_client.FavoriteClient.GetUserFavoritedCount(context.Background(), &request3.DouyinGetUserFavoritedCountRequest{
		UserId: req.UserId,
	})
	if err != nil {
		if resGetUserFavoritedCount == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if resGetUserFavoritedCount.Code == 2 {
			zap.L().Error(errorGetUserFavoritedCount, zap.Error(err))
			return nil, err
		}
	}
	//判断请求查看的用户列表都是否为大V或者活跃用户
	isVList, isActiveList, err := redis.IsVActiveList(req.UserId)
	if err != nil {
		zap.L().Error(errorExeVActiveSet, zap.Error(err))
		return nil, err
	}
	//将既不是大V又不是活跃用户、是大V、是活跃用户的用户列表列出来 放在对应下标下,日后好拼接
	noVActiveList := make([]int64, 0, len(req.UserId))
	VList := make([]int64, 0, len(req.UserId))
	ActiveList := make([]int64, 0, len(req.UserId))
	noVActiveList01 := make([]int8, len(req.UserId), len(req.UserId))
	VList01 := make([]int8, len(req.UserId), len(req.UserId))
	ActiveList01 := make([]int8, len(req.UserId), len(req.UserId))
	for i := 0; i < len(isVList); i++ {
		if isVList[i] == false && isActiveList[i] == false {
			noVActiveList01[i] = 1
			noVActiveList = append(noVActiveList, req.UserId[i])
		}
		if isVList[i] == true {
			VList01[i] = 1
			VList = append(VList, req.UserId[i])
		}
		if isActiveList[i] == true {
			ActiveList01[i] = 1
			ActiveList = append(ActiveList, req.UserId[i])
		}
	}
	resultNotVActive := make([]*response.User, 0, len(req.UserId))
	if len(noVActiveList) != 0 {
		//对既不是大V又不是活跃用户的列表继续分类
		noVActiveSamePeopleList := make([]int64, 0, len(req.UserId))
		noVActiveNotSamePeopleList := make([]int64, 0, len(req.UserId))
		noVActiveSamePeopleList01 := make([]int8, len(req.UserId), len(req.UserId))
		noVActiveNotSamePeopleList01 := make([]int8, len(req.UserId), len(req.UserId))
		for i := 0; i < len(noVActiveList); i++ {
			if req.LoginUserId == noVActiveList[i] {
				noVActiveSamePeopleList01[i] = 1
				noVActiveSamePeopleList = append(noVActiveSamePeopleList, noVActiveList[i])
			} else {
				noVActiveNotSamePeopleList01[i] = 1
				noVActiveNotSamePeopleList = append(noVActiveNotSamePeopleList, noVActiveList[i])
			}
		}
		usReplace := make([]model.User, 0)
		us := make([]model.User, 0)
		if len(noVActiveSamePeopleList) != 0 {
			//1.既不是大V又不是活跃用户且是同一个人
			//从mysql表读取用户信息
			if err := mysql.GetUserInfoList(noVActiveSamePeopleList, &usReplace); err != nil {
				zap.L().Error(errorGetUserInfoListFailed, zap.Error(err))
				return nil, err
			}
			for i := 0; i < len(noVActiveSamePeopleList); i++ {
				for j := 0; j < len(usReplace); j++ {
					if noVActiveSamePeopleList[i] == usReplace[j].UserID {
						us = append(us, usReplace[j])
						break
					}
				}
			}
		}
		us1 := make([]model.User, 0)
		usX := make([]model.User, 0)
		var res1 *response1.DouyinGetFollowListResponse
		if len(noVActiveNotSamePeopleList) != 0 {
			//2.既不是大V又不是活跃用户但不是同一个人
			res1, err = grpc_client.FollowClient.GetFollowInfoList(context.Background(), &request1.DouyinGetFollowListRequest{
				UserId:   req.LoginUserId,
				ToUserId: noVActiveNotSamePeopleList,
			})
			if err != nil {
				if res1 == nil {
					zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
					return nil, err
				}
				if res1.Code == 2 {
					zap.L().Error(errorGetFollowInfoList, zap.Error(err))
					return nil, err
				}
			}
			//从mysql表读取用户信息
			if err := mysql.GetUserInfoList(noVActiveNotSamePeopleList, &usX); err != nil {
				zap.L().Error(errorGetUserInfoListFailed, zap.Error(err))
				return nil, err
			}
			for i := 0; i < len(noVActiveNotSamePeopleList); i++ {
				for j := 0; j < len(usX); j++ {
					if noVActiveNotSamePeopleList[i] == usX[j].UserID {
						us1 = append(us1, usX[j])
						break
					}
				}
			}
		}
		//部分信息拼凑
		k0 := 0
		k1 := 0
		for i := 0; i < len(noVActiveList); i++ {
			var uu response.User
			if noVActiveSamePeopleList01[i] == 1 {
				uu.Id = us[k0].UserID
				uu.Name = us[k0].UserName
				uu.Avatar = us[k0].Avatar
				uu.BackgroundImage = us[k0].BackgroundImage
				uu.Signature = us[k0].Signature
				uu.IsFollow = true
				k0++
			}
			if noVActiveNotSamePeopleList01[i] == 1 {
				uu.Id = us1[k1].UserID
				uu.Name = us1[k1].UserName
				uu.Avatar = us1[k1].Avatar
				uu.BackgroundImage = us1[k1].BackgroundImage
				uu.Signature = us1[k1].Signature
				uu.IsFollow = res1.IsFollow[k1]
				k1++
			}
			resultNotVActive = append(resultNotVActive, &uu)
		}
	}
	resultV := make([]*response.User, 0, len(req.UserId))
	if len(VList) != 0 {
		//3.是大V且是同一个人
		VSamePeopleList := make([]int64, 0, len(req.UserId))
		VNotSamePeopleList := make([]int64, 0, len(req.UserId))
		VSamePeopleList01 := make([]int8, len(req.UserId), len(req.UserId))
		VNotSamePeopleList01 := make([]int8, len(req.UserId), len(req.UserId))
		for i := 0; i < len(VList); i++ {
			if VList[i] == req.LoginUserId {
				VSamePeopleList01[i] = 1
				VSamePeopleList = append(VSamePeopleList, VList[i])
			} else {
				VNotSamePeopleList01[i] = 1
				VNotSamePeopleList = append(VNotSamePeopleList, VList[i])
			}
		}
		var us2 []model.User
		if len(VSamePeopleList) != 0 {
			us2, err = redis.GetVUserInfoList(VSamePeopleList)
			if err != nil {
				zap.L().Error(errGetVUserInfoList, zap.Error(err))
				return nil, err
			}
		}
		var us3 []model.User
		var res2 []bool
		if len(VNotSamePeopleList) != 0 {
			//4.是大V且不是同一个人
			us3, err = redis.GetVUserInfoList(VNotSamePeopleList)
			if err != nil {
				zap.L().Error(errGetVUserInfoList, zap.Error(err))
				return nil, err
			}
			res2, err = redis.GetVFollowerInfoList(VNotSamePeopleList, req.LoginUserId)
			if err != nil {
				zap.L().Error(errorGetVFollowerInfoList, zap.Error(err))
				return nil, err
			}
		}
		//部分信息拼凑
		k2 := 0
		k3 := 0
		for i := 0; i < len(VList); i++ {
			var uu response.User
			if VSamePeopleList01[i] == 1 {
				uu.Id = us2[k2].UserID
				uu.Name = us2[k2].UserName
				uu.Avatar = us2[k2].Avatar
				uu.BackgroundImage = us2[k2].BackgroundImage
				uu.Signature = us2[k2].Signature
				uu.IsFollow = true
				k2++
			}
			if VNotSamePeopleList01[i] == 1 {
				uu.Id = us3[k3].UserID
				uu.Name = us3[k3].UserName
				uu.Avatar = us3[k3].Avatar
				uu.BackgroundImage = us3[k3].BackgroundImage
				uu.Signature = us3[k3].Signature
				uu.IsFollow = res2[k3]
				k3++
			}
			resultV = append(resultV, &uu)
		}
	}
	resultActive := make([]*response.User, 0, len(req.UserId))
	if len(ActiveList) != 0 {
		//5.是活跃用户是同一个人
		ActiveSamePeopleList := make([]int64, 0, len(req.UserId))
		ActiveNotSamePeopleList := make([]int64, 0, len(req.UserId))
		ActiveSamePeopleList01 := make([]int8, len(req.UserId), len(req.UserId))
		ActiveNotSamePeopleList01 := make([]int8, len(req.UserId), len(req.UserId))
		for i := 0; i < len(ActiveList); i++ {
			if ActiveList[i] == req.LoginUserId {
				ActiveSamePeopleList01[i] = 1
				ActiveSamePeopleList = append(ActiveSamePeopleList, ActiveList[i])
			} else {
				ActiveNotSamePeopleList01[i] = 1
				ActiveNotSamePeopleList = append(ActiveNotSamePeopleList, ActiveList[i])
			}
		}
		var us4 []model.User
		if len(ActiveSamePeopleList) != 0 {
			us4, err = redis.GetActiveUserInfoList(ActiveSamePeopleList)
			if err != nil {
				zap.L().Error(errGetActiveUserInfoList, zap.Error(err))
				return nil, err
			}
		}
		var us5 []model.User
		var res3 []bool
		//6.是活跃用户且不是同一个人
		if len(ActiveNotSamePeopleList) != 0 {
			us5, err = redis.GetActiveUserInfoList(ActiveNotSamePeopleList)
			if err != nil {
				zap.L().Error(errGetActiveUserInfoList, zap.Error(err))
				return nil, err
			}
			res3, err = redis.GetActiveFollowerInfoList(ActiveNotSamePeopleList, req.LoginUserId)
			if err != nil {
				zap.L().Error(errorGetActiveFollowerInfoList, zap.Error(err))
				return nil, err
			}
		}
		//部分信息拼凑
		k4 := 0
		k5 := 0
		for i := 0; i < len(ActiveList); i++ {
			var uu response.User
			if ActiveSamePeopleList01[i] == 1 {
				uu.Id = us4[k4].UserID
				uu.Name = us4[k4].UserName
				uu.Avatar = us4[k4].Avatar
				uu.BackgroundImage = us4[k4].BackgroundImage
				uu.Signature = us4[k4].Signature
				uu.IsFollow = true
				k4++
			}
			if ActiveNotSamePeopleList01[i] == 1 {
				uu.Id = us5[k5].UserID
				uu.Name = us5[k5].UserName
				uu.Avatar = us5[k5].Avatar
				uu.BackgroundImage = us5[k5].BackgroundImage
				uu.Signature = us5[k5].Signature
				uu.IsFollow = res3[k5]
				k5++
			}
			resultActive = append(resultActive, &uu)
		}
	}
	//信息汇总
	k6 := 0
	k7 := 0
	k8 := 0
	for i := 0; i < len(req.UserId); i++ {
		if noVActiveList01[i] == 1 {
			result = append(result, resultNotVActive[k6])
			k6++
		}
		if VList01[i] == 1 {
			result = append(result, resultV[k7])
			k7++
		}
		if ActiveList01[i] == 1 {
			result = append(result, resultActive[k8])
			k8++
		}
	}
	for i := 0; i < len(result); i++ {
		result[i].FollowCount = res.FollowCount[i]
		result[i].FollowerCount = res.FollowerCount[i]
		result[i].TotalFavorited = resGetUserFavoritedCount.FavoritedCount[i]
		result[i].FavoriteCount = resGetUserFavoritedCount.FavoriteCount[i]
		result[i].WorkCount = resGetUserPublishVideoCount.Count[i]
	}
	return result, nil
}

func AddUserFollowUserCountSet(req *request.DouyinUserFollowCountSetRequest) error {
	return redis.AddUserFollowUserCountSet(req.UserId)
}

func AddUserFollowUserCountSetRevert(req *request.DouyinUserFollowCountSetRequest) error {
	return redis.DeleteUserFollowUserCountSet(req.UserId)
}

func AddUserFollowerUserCountSet(req *request.DouyinUserFollowerCountSetRequest) error {
	return redis.AddUserFollowerUserCountSet(req.UserId)
}

func AddUserFollowerUserCountSetRevert(req *request.DouyinUserFollowerCountSetRequest) error {
	return redis.DeleteUserFollowerUserCountSet(req.UserId)
}

func AddUserVideoCountSet(req *request.DouyinUserVideoCountSetRequest) error {
	return redis.AddUserVideoCountSet(req.UserId)
}

func AddUserVideoCountSetRevert(req *request.DouyinUserVideoCountSetRequest) error {
	return redis.DeletedUserVideoCountSet(req.UserId)
}

func AddUserFavoriteVideoCountSet(req *request.DouyinUserVideoCountSetRequest) error {
	return redis.AddUserFavoriteVideoCountSet(req.UserId)
}

func AddUserFavoriteVideoCountSetRevert(req *request.DouyinUserVideoCountSetRequest) error {
	return redis.DeleteUserFavoriteVideoCountSet(req.UserId)
}
