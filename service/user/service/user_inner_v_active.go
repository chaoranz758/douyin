package service

import (
	"context"
	request3 "douyin/proto/favorite/request"
	request1 "douyin/proto/follow/request"
	"douyin/proto/user/request"
	"douyin/proto/user/response"
	request2 "douyin/proto/video/request"
	"douyin/service/user/dao/mysql"
	"douyin/service/user/dao/redis"
	"douyin/service/user/initialize/grpc_client"
	"douyin/service/user/model"
	"errors"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UserIsInfluencerActiver(req *request.DouyinUserIsInfluencerActiverRequest) (bool, bool, error) {
	if req.UserId == 0 {
		return false, false, nil
	}
	isV, isActive, err := redis.IsVActive(req.UserId)
	if err != nil {
		zap.L().Error(errorExeVActiveSet, zap.Error(err))
		return false, false, err
	}
	return isV, isActive, nil
}

func PushVSet(req *request.DouyinPushVSetRequest) (bool, int32, error) {
	//首先判断该用户之前是不是活跃用户
	_, isActive, err := redis.IsVActive(req.UserId)
	if err != nil {
		zap.L().Error(errorExeVActiveSet, zap.Error(err))
		return false, 0, err
	}
	//若之前是活跃用户，则删除相关信息
	if isActive == true {
		if err = redis.DeleteActiveAllInfo(req.UserId); err != nil {
			zap.L().Error(errorDeleteActiveAllInfo, zap.Error(err))
			return false, 0, err
		}
	}
	//写入VSet
	b, err := redis.PushVSet(req.UserId)
	if err != nil {
		zap.L().Error(errorPushVSet, zap.Error(err))
		return false, 0, err
	}
	if b == true {
		return isActive, 1, nil
	}
	//写入大v用户信息redis
	if b == false {
		var user model.User
		if err := mysql.GetUserInfo(req.UserId, &user); err != nil {
			zap.L().Error(errorGetUserInfoFailed, zap.Error(err))
			return false, 0, err
		}
		if err := redis.PushVString(req.UserId, user); err != nil {
			zap.L().Error(errorPushVString, zap.Error(err))
			return false, 0, err
		}
	}
	//从mysql读用户基本信息
	usersFollow := make([]model.User, 0, len(req.FollowIdList))
	if len(req.FollowIdList) != 0 {
		if err := mysql.GetUserInfoList(req.FollowIdList, &usersFollow); err != nil {
			zap.L().Error(errorGetUserInfoListFailed, zap.Error(err))
			return false, 0, err
		}
	}
	usersFollower := make([]model.User, 0, len(req.FollowerIdList))
	if err := mysql.GetUserInfoList(req.FollowerIdList, &usersFollower); err != nil {
		zap.L().Error(errorGetUserInfoListFailed, zap.Error(err))
		return false, 0, err
	}
	//将这些信息存入大V关注者redis、大V粉丝redis
	if err := redis.PushVFollowFollowerInfoInit(usersFollow, usersFollower, req.UserId); err != nil {
		zap.L().Error(errorPushVFollowFollowerInfoInit, zap.Error(err))
		return false, 0, err
	}
	return isActive, 0, nil
}

func GetVActiveFollowerUserinfo(req *request.DouyinGetVActiveFollowFollowerUserinfoRequest) ([]*response.User, error) {
	if req.IsFollowFollower == 1 {
		var users []model.User
		if req.IsV == true {
			users1, err := redis.GetVFollowUserInfo(req.UserId)
			if err != nil {
				zap.L().Error(errorGetVFollowUserInfo, zap.Error(err))
				return nil, err
			}
			if len(users1) == 0 {
				return nil, nil
			}
			users = users1
		}
		if req.IsActive == true {
			users2, err := redis.GetActiveFollowUserInfo(req.UserId)
			if err != nil {
				zap.L().Error(errorGetVFollowUserInfo, zap.Error(err))
				return nil, err
			}
			if len(users2) == 0 {
				return nil, nil
			}
			users = users2
		}
		var userIdList []int64
		for i := 0; i < len(users); i++ {
			userIdList = append(userIdList, users[i].UserID)
		}
		//批量读取用户粉丝数和关注数
		res, err := grpc_client.FollowClient.GetFollowFollowerList(context.Background(), &request1.DouyinFollowFollowerListCountRequest{
			UserId: userIdList,
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
			UserId: userIdList,
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
			UserId: userIdList,
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
		isVList, isActiveList, err := redis.IsVActiveList(userIdList)
		if err != nil {
			zap.L().Error(errorExeVActiveSet, zap.Error(err))
			return nil, err
		}
		var noVActiveList []int64
		var VList []int64
		var ActiveList []int64
		for i := 0; i < len(userIdList); i++ {
			if isVList[i] == true {
				VList = append(VList, userIdList[i])
			}
			if isActiveList[i] == true {
				ActiveList = append(ActiveList, userIdList[i])
			}
			if isVList[i] == false && isActiveList[i] == false {
				noVActiveList = append(noVActiveList, userIdList[i])
			}
		}
		isVFollowList, err := redis.GetVFollowInfoList(VList, req.LoginUserId)
		if err != nil {
			zap.L().Error(errorGetVFollowerInfoList, zap.Error(err))
			return nil, err
		}
		isActiveFollowList, err := redis.GetActiveFollowInfoList(ActiveList, req.LoginUserId)
		if err != nil {
			zap.L().Error(errorGetActiveFollowerInfoList, zap.Error(err))
			return nil, err
		}
		res1, err := grpc_client.FollowClient.GetFollowInfoList(context.Background(), &request1.DouyinGetFollowListRequest{
			UserId:   req.LoginUserId,
			ToUserId: noVActiveList,
		})
		var results []*response.User
		k1 := 0
		k2 := 0
		k3 := 0
		for i := 0; i < len(userIdList); i++ {
			var result = response.User{
				Id:              users[i].UserID,
				Name:            users[i].UserName,
				Signature:       users[i].Signature,
				Avatar:          users[i].Avatar,
				BackgroundImage: users[i].BackgroundImage,
				FollowCount:     res.FollowCount[i],
				FollowerCount:   res.FollowerCount[i],
				TotalFavorited:  resGetUserFavoritedCount.FavoritedCount[i],
				FavoriteCount:   resGetUserFavoritedCount.FavoriteCount[i],
				WorkCount:       resGetUserPublishVideoCount.Count[i],
			}
			if isVList[i] == true {
				result.IsFollow = isVFollowList[k1]
				k1++
			}
			if isActiveList[i] == true {
				result.IsFollow = isActiveFollowList[k2]
				k2++
			}
			if isVList[i] == false && isActiveList[i] == false {
				result.IsFollow = res1.IsFollow[k3]
				k3++
			}
			results = append(results, &result)
		}
		return results, nil
	}
	if req.IsFollowFollower == 2 {
		var users []model.User
		if req.IsV == true {
			users1, err := redis.GetVFollowerUserInfo(req.UserId)
			if err != nil {
				zap.L().Error(errorGetVFollowerUserInfo, zap.Error(err))
				return nil, err
			}
			if len(users1) == 0 {
				return nil, nil
			}
			users = users1
		}
		if req.IsActive == true {
			users2, err := redis.GetActiveFollowerUserInfo(req.UserId)
			if err != nil {
				zap.L().Error(errorGetActiveFollowerUserInfo, zap.Error(err))
				return nil, err
			}
			if len(users2) == 0 {
				return nil, nil
			}
			users = users2
		}
		var userIdList []int64
		for i := 0; i < len(users); i++ {
			userIdList = append(userIdList, users[i].UserID)
		}
		//批量读取用户粉丝数和关注数
		res, err := grpc_client.FollowClient.GetFollowFollowerList(context.Background(), &request1.DouyinFollowFollowerListCountRequest{
			UserId: userIdList,
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
			UserId: userIdList,
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
			UserId: userIdList,
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
		isVList, isActiveList, err := redis.IsVActiveList(userIdList)
		if err != nil {
			zap.L().Error(errorExeVActiveSet, zap.Error(err))
			return nil, err
		}
		var noVActiveList []int64
		var VList []int64
		var ActiveList []int64
		for i := 0; i < len(userIdList); i++ {
			if isVList[i] == true {
				VList = append(VList, userIdList[i])
			}
			if isActiveList[i] == true {
				ActiveList = append(ActiveList, userIdList[i])
			}
			if isVList[i] == false && isActiveList[i] == false {
				noVActiveList = append(noVActiveList, userIdList[i])
			}
		}
		isVFollowList, err := redis.GetVFollowerInfoList(VList, req.LoginUserId)
		if err != nil {
			zap.L().Error(errorGetVFollowerInfoList, zap.Error(err))
			return nil, err
		}
		isActiveFollowList, err := redis.GetActiveFollowerInfoList(ActiveList, req.LoginUserId)
		if err != nil {
			zap.L().Error(errorGetActiveFollowerInfoList, zap.Error(err))
			return nil, err
		}
		res1, err := grpc_client.FollowClient.GetFollowInfoList(context.Background(), &request1.DouyinGetFollowListRequest{
			UserId:   req.LoginUserId,
			ToUserId: noVActiveList,
		})
		var results []*response.User
		k1 := 0
		k2 := 0
		k3 := 0
		for i := 0; i < len(userIdList); i++ {
			var result = response.User{
				Id:              users[i].UserID,
				Name:            users[i].UserName,
				Signature:       users[i].Signature,
				Avatar:          users[i].Avatar,
				BackgroundImage: users[i].BackgroundImage,
				FollowCount:     res.FollowCount[i],
				FollowerCount:   res.FollowerCount[i],
				TotalFavorited:  resGetUserFavoritedCount.FavoritedCount[i],
				FavoriteCount:   resGetUserFavoritedCount.FavoriteCount[i],
				WorkCount:       resGetUserPublishVideoCount.Count[i],
			}
			if isVList[i] == true {
				result.IsFollow = isVFollowList[k1]
				k1++
			}
			if isActiveList[i] == true {
				result.IsFollow = isActiveFollowList[k2]
				k2++
			}
			if isVList[i] == false && isActiveList[i] == false {
				result.IsFollow = res1.IsFollow[k3]
				k3++
			}
			results = append(results, &result)
		}
		return results, nil
	}
	return nil, errors.New(errorRpc)
}

func PushVActiveFollowerUserinfo(req *request.DouyinVActiveFollowFollowerUserinfoRequest) error {
	//存大V或活跃用户关注的粉丝信息 登录用户是大V
	if req.IsFollowFollower == 1 {
		//存大V
		if req.LoginIsV == true {
			if req.IsV == true {
				//从大V用户基本信息redis获取并存
				var vUser model.User
				if err := redis.GetVUserInfo(req.UserId, &vUser); err != nil {
					zap.L().Error(errGetVUserInfo, zap.Error(err))
					return status.Error(codes.Aborted, err.Error())
				}
				if err := redis.PushVFollowUserInfo(vUser, req.UserId, req.LoginUserId); err != nil {
					zap.L().Error(errorPushVFollowUserInfo, zap.Error(err))
					return status.Error(codes.Aborted, err.Error())
				}
				return nil
			}
			if req.IsActive == true {
				//从活跃用户基本信息redis获取并存
				var ActiveUser model.User
				if err := redis.GetActiveUserInfo(req.UserId, &ActiveUser); err != nil {
					zap.L().Error(errGetActiveUserInfo, zap.Error(err))
					return status.Error(codes.Aborted, err.Error())
				}
				if err := redis.PushVFollowUserInfo(ActiveUser, req.UserId, req.LoginUserId); err != nil {
					zap.L().Error(errorPushVFollowUserInfo, zap.Error(err))
					return status.Error(codes.Aborted, err.Error())
				}
				return nil
			}
			//从mysql取基本信息并存
			var simpleUser model.User
			if err := mysql.GetUserInfo(req.UserId, &simpleUser); err != nil {
				zap.L().Error(errorGetUserInfoFailed, zap.Error(err))
				return status.Error(codes.Aborted, err.Error())
			}
			if err := redis.PushVFollowUserInfo(simpleUser, req.UserId, req.LoginUserId); err != nil {
				zap.L().Error(errorPushVFollowUserInfo, zap.Error(err))
				return status.Error(codes.Aborted, err.Error())
			}
			return nil
		}
		//存活跃用户
		if req.LoginIsActive == true {
			if req.IsV == true {
				//从大V用户基本信息redis获取并存
				var vUser model.User
				if err := redis.GetVUserInfo(req.UserId, &vUser); err != nil {
					zap.L().Error(errGetVUserInfo, zap.Error(err))
					return status.Error(codes.Aborted, err.Error())
				}
				if err := redis.PushActiveFollowUserInfo(vUser, req.UserId, req.LoginUserId); err != nil {
					zap.L().Error(errorPushActiveFollowUserInfo, zap.Error(err))
					return status.Error(codes.Aborted, err.Error())
				}
				return nil
			}
			if req.IsActive == true {
				//从活跃用户基本信息redis获取并存
				var ActiveUser model.User
				if err := redis.GetActiveUserInfo(req.UserId, &ActiveUser); err != nil {
					zap.L().Error(errGetActiveUserInfo, zap.Error(err))
					return status.Error(codes.Aborted, err.Error())
				}
				if err := redis.PushActiveFollowUserInfo(ActiveUser, req.UserId, req.LoginUserId); err != nil {
					zap.L().Error(errorPushActiveFollowUserInfo, zap.Error(err))
					return status.Error(codes.Aborted, err.Error())
				}
				return nil
			}
			//从mysql取基本信息并存
			var simpleUser model.User
			if err := mysql.GetUserInfo(req.UserId, &simpleUser); err != nil {
				zap.L().Error(errorGetUserInfoFailed, zap.Error(err))
				return status.Error(codes.Aborted, err.Error())
			}
			if err := redis.PushActiveFollowUserInfo(simpleUser, req.UserId, req.LoginUserId); err != nil {
				zap.L().Error(errorPushActiveFollowUserInfo, zap.Error(err))
				return status.Error(codes.Aborted, err.Error())
			}
			return nil
		}
	}
	if req.IsFollowFollower == 2 {
		//存大V或活跃用户的粉丝信息
		if req.IsV == true {
			var vUser model.User
			if req.LoginIsV == true {
				if err := redis.GetVUserInfo(req.LoginUserId, &vUser); err != nil {
					zap.L().Error(errGetVUserInfo, zap.Error(err))
					return status.Error(codes.Aborted, err.Error())
				}
			}
			if req.LoginIsActive == true {
				if err := redis.GetActiveUserInfo(req.LoginUserId, &vUser); err != nil {
					zap.L().Error(errGetActiveUserInfo, zap.Error(err))
					return status.Error(codes.Aborted, err.Error())
				}
			}
			if req.LoginIsV == false && req.LoginIsActive == false {
				if err := mysql.GetUserInfo(req.LoginUserId, &vUser); err != nil {
					zap.L().Error(errorGetUserInfoFailed, zap.Error(err))
					return status.Error(codes.Aborted, err.Error())
				}
			}
			if err := redis.PushVFollowerUserInfo(vUser, req.UserId, req.LoginUserId); err != nil {
				zap.L().Error(errorPushVFollowerUserInfo, zap.Error(err))
				return status.Error(codes.Aborted, err.Error())
			}
			return nil
		}
		if req.IsActive == true {
			var vUser model.User
			if req.LoginIsV == true {
				if err := redis.GetVUserInfo(req.LoginUserId, &vUser); err != nil {
					zap.L().Error(errGetVUserInfo, zap.Error(err))
					return status.Error(codes.Aborted, err.Error())
				}
			}
			if req.LoginIsActive == true {
				if err := redis.GetActiveUserInfo(req.LoginUserId, &vUser); err != nil {
					zap.L().Error(errGetActiveUserInfo, zap.Error(err))
					return status.Error(codes.Aborted, err.Error())
				}
			}
			if req.LoginIsV == false && req.LoginIsActive == false {
				if err := mysql.GetUserInfo(req.LoginUserId, &vUser); err != nil {
					zap.L().Error(errorGetUserInfoFailed, zap.Error(err))
					return status.Error(codes.Aborted, err.Error())
				}
			}
			if err := redis.PushActiveFollowerUserInfo(vUser, req.UserId, req.LoginUserId); err != nil {
				zap.L().Error(errorPushActiveFollowerUserInfo, zap.Error(err))
				return status.Error(codes.Aborted, err.Error())
			}
			return nil
		}
	}
	return status.Error(codes.Aborted, errorRpc)
}

func PushVActiveFollowerUserinfoRevert(req *request.DouyinVActiveFollowFollowerUserinfoRequest) error {
	if req.IsFollowFollower == 1 {
		if req.LoginIsV == true {
			if err := redis.DeleteVFollowUserInfo(req.UserId, req.LoginUserId); err != nil {
				zap.L().Error(errDeleteVFollowUserInfo, zap.Error(err))
				return status.Error(codes.Aborted, err.Error())
			}
		}
		if req.LoginIsActive == true {
			if err := redis.DeleteActiveFollowUserInfo(req.UserId, req.LoginUserId); err != nil {
				zap.L().Error(errDeleteActiveFollowUserInfo, zap.Error(err))
				return status.Error(codes.Aborted, err.Error())
			}
		}
	}
	if req.IsFollowFollower == 2 {
		if req.IsV == true {
			if err := redis.DeleteVFollowerUserInfo(req.UserId, req.LoginUserId); err != nil {
				zap.L().Error(errDeleteVFollowerUserInfo, zap.Error(err))
				return status.Error(codes.Aborted, err.Error())
			}
		}
		if req.IsActive == true {
			if err := redis.DeleteActiveFollowerUserInfo(req.UserId, req.LoginUserId); err != nil {
				zap.L().Error(errDeleteActiveFollowerUserInfo, zap.Error(err))
				return status.Error(codes.Aborted, err.Error())
			}
		}
	}
	return nil
}

func DeleteVActiveFollowerUserinfo(req *request.DouyinDeleteVActiveFollowFollowerUserinfoRequest) error {
	if req.IsFollowFollower == 1 {
		if req.LoginIsV == true {
			if err := redis.DeleteVFollowUserInfo(req.UserId, req.LoginUserId); err != nil {
				zap.L().Error(errDeleteVFollowUserInfo, zap.Error(err))
				return status.Error(codes.Aborted, err.Error())
			}
		}
		if req.LoginIsActive == true {
			if err := redis.DeleteActiveFollowUserInfo(req.UserId, req.LoginUserId); err != nil {
				zap.L().Error(errDeleteActiveFollowUserInfo, zap.Error(err))
				return status.Error(codes.Aborted, err.Error())
			}
		}
	}
	if req.IsFollowFollower == 2 {
		if req.IsV == true {
			if err := redis.DeleteVFollowerUserInfo(req.UserId, req.LoginUserId); err != nil {
				zap.L().Error(errDeleteVFollowerUserInfo, zap.Error(err))
				return status.Error(codes.Aborted, err.Error())
			}
		}
		if req.IsActive == true {
			if err := redis.DeleteActiveFollowerUserInfo(req.UserId, req.LoginUserId); err != nil {
				zap.L().Error(errDeleteActiveFollowerUserInfo, zap.Error(err))
				return status.Error(codes.Aborted, err.Error())
			}
		}
	}
	return nil
}
