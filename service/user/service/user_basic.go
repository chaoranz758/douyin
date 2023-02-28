package service

import (
	"context"
	request3 "douyin/proto/favorite/request"
	request1 "douyin/proto/follow/request"
	"douyin/proto/user/request"
	request2 "douyin/proto/video/request"
	"douyin/service/user/dao/mysql"
	"douyin/service/user/dao/redis"
	"douyin/service/user/initialize/grpc_client"
	"douyin/service/user/model"
	"douyin/service/user/util"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

var (
	signature = "抖音用户-英雄联盟玩家"
	avatars   = []string{"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/05d90598a5b55e43cbbb121fc6c3965f.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/1234608dc27c3bb29597f59a4902b741.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/1a9299331799026df2880b0bc1f243d6.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/23a45fe6dc1de67cb0a674654aaaf5d9.jpg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/3047f670ed9800dcb669d9a5a6cda40c.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/44377c8508de5ee8ad10b3f834f46986.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/4768f4e295f3be6e1c250bb1299cdd48.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/51c49387facb88f75bd099f67c9f5f20.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/68d95f81183160d1ca9889652a47d590.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/6dd30df30ab4a48298e54df4d377f95a.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/7632955e68c1a32c8cda324fb9e78b6d.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/81c6e867b8ac907c71c0e2607d01d1b6.jpg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/84544c55badaab660b6b7da804192b62.jpg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/871d0d62506da4ea22584010b0ad2717%20%281%29.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/871d0d62506da4ea22584010b0ad2717.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/87e04c40b3769c7045003092cbcad91c.jpg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/91a15a77d6d5cf508fbbf54ce80339a3.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/9217b61c46f6337b37498b0ba4b86275.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/9bb1fa0d78d02a967e0487ae78925706.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/a16211480af2090761af3e4a29b982db.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/afdaeffa215941611df76006c3a576c1.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/b2cf30f906995736726b65ae71e6cec3.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/bc9c045917ea33125f45b3a4444a0d15.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/bd76e1723f1d0ee240b5f6eeec46622c.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/ca8da9a28451c50a292f1881498596d4.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/cf6d4ce4db9b420cdf4e3048daf71162.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/d79f8fec8873d3fc70a34c725b4d1ecd.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/d7b08f04a942900baac2011f616c3e64.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/f6febeaebc57e4e1f02e67822ed10d34.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%A4%B4%E5%83%8F/%E5%A4%B4%E5%83%8F/fe592260d93bde3d3b73534e651fc85a.jpeg"}
	backgroundImages = []string{"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/0690d39902272986e59097d536145dd7.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/1144605cbc005d4996af9286d7bdc903.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/12485fc02389d49ec9c3df413a3078c8.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/2993917d934382dbd04cbeadfa12b2c4.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/2fc1adcb8dd76d9b7483c2d32730bd71%20%281%29.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/2fc1adcb8dd76d9b7483c2d32730bd71.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/350624206a4856d0fd02b99c0b05f7f9.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/37f7e223e9002532aab844209ed6de07.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/47ea8c41d8698d40394d0369f7ef2347.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/48230ce65f9d5185ce9d87bf11cd4fa7.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/4a3e7fe59a60f48e684cd3bf8bef2feb.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/4b4a2717f3e50a28eb6f937d5788903b.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/4d4e1c3c84da926eaa1d8c839d182fee.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/52c34f362b2b0b7fcc5c10cd1e9a44ab.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/5bd8e335bf8c09c15bb459bc5c163b46.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/5e87652e20697a5fb3054121bc4be529.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/66800ee541dbe13e444ea6bcdb5431b3.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/83ba47a06a95b70460bada086b9edc42.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/8cdac93f2445f129767b0c712975d3d9.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/9117d0077a146cf78eae9a531cf112ce.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/91cc67398b48dce6e59cbe6f5ea7dc26.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/9d188a22e89d01c2918044bacd15ad66.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/a0ce7d66e83646f6d6be505092af1514.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/b6e31bb4b163fc775a5a2bb5d6dd087e.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/c0502fd29b722c688bd4f1196921dbf2.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/c6a26e818fe3f9ba3951fb31225071af.jpg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/cd894589bc2c954e06eb420ca53aa8b7.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/d4e60212878a07a3701e21789fc57de4.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/d9beebf5b2c382c4edfb2e653aa10edc.jpeg",
		"https://simpledouyin.oss-cn-qingdao.aliyuncs.com/%E5%B0%81%E9%9D%A2%E5%9B%BE/%E5%B0%81%E9%9D%A2%E5%9B%BE/ef72faa3f87912e2136934a0f1f97a04.jpeg"}
)

func UserRegister(req *request.DouyinUserRegisterRequest) (int64, string, error) {
	userID := util.GenID()
	password, err := util.ScryptPw(req.Password)
	if err != nil {
		zap.L().Error(errorEncryptPassword, zap.Error(err))
		return 0, "", err
	}
	token, err := util.GenToken(req.Username, userID)
	if err != nil {
		zap.L().Error(errorGenToken, zap.Error(err))
		return 0, "", err
	}
	rand1 := rand.Intn(len(avatars))
	rand2 := rand.Intn(len(backgroundImages))
	var u = model.User{
		UserID:          userID,
		UserName:        req.Username,
		Password:        password,
		Avatar:          avatars[rand1],
		BackgroundImage: backgroundImages[rand2],
		Signature:       signature + fmt.Sprintf("%v", time.Now().UnixMilli()),
	}
	if err := mysql.UserRegister(&u); err != nil {
		zap.L().Error(errorCreateUser, zap.Error(err))
		return 0, "", err
	}
	return userID, token, nil
}

func UserLogin(req *request.DouyinUserLoginRequest) (int64, string, error) {
	password, err := util.ScryptPw(req.Password)
	if err != nil {
		zap.L().Error(errorEncryptPassword, zap.Error(err))
		return 0, "", err
	}
	var u model.User
	if err := mysql.UserLogin(&u, req); err != nil {
		zap.L().Error(errorUserLogin, zap.Error(err))
		return 0, "", err
	}
	if password != u.Password {
		zap.L().Error(errorPassword, zap.Error(err))
		return 0, "", errors.New(errorPassword)
	}
	token, err := util.GenToken(req.Username, u.UserID)
	if err != nil {
		zap.L().Error(errorGenToken, zap.Error(err))
		return 0, "", err
	}
	if err := redis.UserLoginCount(u.UserID); err != nil {
		zap.L().Error(errorUserLoginCountAdd1, zap.Error(err))
		return 0, "", err
	}
	return u.UserID, token, nil
}

func GetUserInfo(req *request.DouyinUserRequest) (*model.User, bool, int64, int64, int64, error) {
	var u model.User
	//读取用户粉丝数和关注数
	res, err := grpc_client.FollowClient.GetFollowFollower(context.Background(), &request1.DouyinFollowFollowerCountRequest{
		UserId: req.UserId,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, false, 0, 0, 0, err
		}
		if res.Code == 2 {
			zap.L().Error(errorGetFollowFollower, zap.Error(err))
			return nil, false, 0, 0, 0, err
		}
	}
	//读取用户作品数量
	resGetUserPublishVideoCount, err := grpc_client.VideoClient.GetUserPublishVideoCount(context.Background(), &request2.DouyinGetUserPublishVideoCountRequest{
		UserId: req.UserId,
	})
	if err != nil {
		if resGetUserPublishVideoCount == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, false, 0, 0, 0, err
		}
		if resGetUserPublishVideoCount.Code == 2 {
			zap.L().Error(errorGetUserPublishVideoCount, zap.Error(err))
			return nil, false, 0, 0, 0, err
		}
	}
	var users []int64
	users = append(users, req.UserId)
	//读取用户点赞视频数量和获赞总数
	resGetUserFavoritedCount, err := grpc_client.FavoriteClient.GetUserFavoritedCount(context.Background(), &request3.DouyinGetUserFavoritedCountRequest{
		UserId: users,
	})
	if err != nil {
		if resGetUserFavoritedCount == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, false, 0, 0, 0, err
		}
		if resGetUserFavoritedCount.Code == 2 {
			zap.L().Error(errorGetUserFavoritedCount, zap.Error(err))
			return nil, false, 0, 0, 0, err
		}
	}
	//判断请求用户是否为大V或活跃用户
	isV, isActive, err := redis.IsVActive(req.UserId)
	if err != nil {
		zap.L().Error(errorExeVActiveSet, zap.Error(err))
		return nil, false, 0, 0, 0, err
	}
	//1.既不是大V也不是活跃用户
	if isV == false && isActive == false {
		//判断当前用户与请求用户是否是同一个人
		if req.UserId == req.LoginUserId {
			//从用户表中读取用户信息
			if err := mysql.GetUserInfo(req.LoginUserId, &u); err != nil {
				zap.L().Error(errorGetUserInfoFailed, zap.Error(err))
				return nil, false, 0, 0, 0, err
			}
			u.Follow = res.FollowCount
			u.Follower = res.FollowerCount
			return &u, true, resGetUserFavoritedCount.FavoritedCount[0], resGetUserPublishVideoCount.Count, resGetUserFavoritedCount.FavoriteCount[0], nil
		}
		res1, err := grpc_client.FollowClient.GetFollowInfo(context.Background(), &request1.DouyinGetFollowRequest{
			UserId:   req.LoginUserId,
			ToUserId: req.UserId,
		})
		if err != nil {
			if res1 == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return nil, false, 0, 0, 0, err
			}
			if res1.Code == 2 {
				zap.L().Error(errorGetFollowInfo, zap.Error(err))
				return nil, false, 0, 0, 0, err
			}
		}
		//从用户表中读取用户信息
		if err := mysql.GetUserInfo(req.UserId, &u); err != nil {
			zap.L().Error(errorGetUserInfoFailed, zap.Error(err))
			return nil, false, 0, 0, 0, err
		}
		u.Follow = res.FollowCount
		u.Follower = res.FollowerCount
		return &u, res1.IsFollow, resGetUserFavoritedCount.FavoritedCount[0], resGetUserPublishVideoCount.Count, resGetUserFavoritedCount.FavoriteCount[0], nil
	}
	//判断当前用户与请求用户是否是同一个人
	//是同一个人
	if req.UserId == req.LoginUserId {
		//如果请求用户是大V，说明当前用户也是大V，可以从记录大V用户信息的redis中查
		if isV == true {
			if err := redis.GetVUserInfo(req.UserId, &u); err != nil {
				zap.L().Error(errGetVUserInfo, zap.Error(err))
				return nil, false, 0, 0, 0, err
			}
			u.Follow = res.FollowCount
			u.Follower = res.FollowerCount
			return &u, true, resGetUserFavoritedCount.FavoritedCount[0], resGetUserPublishVideoCount.Count, resGetUserFavoritedCount.FavoriteCount[0], nil
		}
		if err := redis.GetActiveUserInfo(req.UserId, &u); err != nil {
			zap.L().Error(errGetActiveUserInfo, zap.Error(err))
			return nil, false, 0, 0, 0, err
		}
		u.Follow = res.FollowCount
		u.Follower = res.FollowerCount
		return &u, true, resGetUserFavoritedCount.FavoritedCount[0], resGetUserPublishVideoCount.Count, resGetUserFavoritedCount.FavoriteCount[0], nil
	}
	//不是同一个人
	if isV == true {
		isFollow, err := redis.GetVFollowerInfo(req.UserId, req.LoginUserId)
		if err != nil {
			zap.L().Error(errorGetVFollowerInfo, zap.Error(err))
			return nil, false, 0, 0, 0, err
		}
		if err := redis.GetVUserInfo(req.UserId, &u); err != nil {
			zap.L().Error(errGetVUserInfo, zap.Error(err))
			return nil, false, 0, 0, 0, err
		}
		u.Follow = res.FollowCount
		u.Follower = res.FollowerCount
		return &u, isFollow, resGetUserFavoritedCount.FavoritedCount[0], resGetUserPublishVideoCount.Count, resGetUserFavoritedCount.FavoriteCount[0], nil
	}
	isFollow, err := redis.GetActiveFollowerInfo(req.UserId, req.LoginUserId)
	if err != nil {
		zap.L().Error(errorGetActiveFollowerInfo, zap.Error(err))
		return nil, false, 0, 0, 0, err
	}
	if err := redis.GetActiveUserInfo(req.UserId, &u); err != nil {
		zap.L().Error(errGetActiveUserInfo, zap.Error(err))
		return nil, false, 0, 0, 0, err
	}
	u.Follow = res.FollowCount
	u.Follower = res.FollowerCount
	return &u, isFollow, resGetUserFavoritedCount.FavoritedCount[0], resGetUserPublishVideoCount.Count, resGetUserFavoritedCount.FavoriteCount[0], nil
}
