package service

import (
	"context"
	request3 "douyin/proto/favorite/request"
	request1 "douyin/proto/follow/request"
	response1 "douyin/proto/follow/response"
	"douyin/proto/user/request"
	"douyin/proto/user/response"
	request2 "douyin/proto/video/request"
	"douyin/service/user/client"
	"douyin/service/user/dao/mysql"
	"douyin/service/user/dao/redis"
	"douyin/service/user/model"
	"douyin/service/user/pkg/jwt"
	"douyin/service/user/pkg/snowflake"
	"douyin/service/user/util"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math/rand"
	"time"
)

const (
	errorEncryptPassword             = "encrypt password failed"
	errorCreateUser                  = "create user failed"
	errorGenToken                    = "generate token failed"
	errorUserLogin                   = "user login failed"
	errorPassword                    = "password compare failed"
	errorUserLoginCountAdd1          = "user login count redis add 1 failed"
	errorConnectToGRPCServer         = "connect to grpc server failed"
	errorGetFollowFollower           = "get user's follow follower count failed"
	errorGetFollowFollowerList       = "get user's follow follower list count failed"
	errorExeVActiveSet               = "exe if request is V or active from redis set failed"
	errorGetUserInfoFailed           = "get user information failed"
	errorGetFollowInfo               = "get user follow relation from mysql failed"
	errGetVUserInfo                  = "get V User Information from redis failed"
	errGetActiveUserInfo             = "get active User Information from redis failed"
	errorGetVFollowerInfo            = "get v follower information failed"
	errorGetActiveFollowerInfo       = "get active follower information failed"
	errorGetUserInfoListFailed       = "get user information list failed"
	errorGetFollowInfoList           = "get follow information list failed"
	errGetVUserInfoList              = "get v user information list failed"
	errGetActiveUserInfoList         = "get active user information list failed"
	errorGetVFollowerInfoList        = "get v follower information list failed"
	errorGetActiveFollowerInfoList   = "get active follower information list failed"
	errorPushVFollowUserInfo         = "push v follow user information failed"
	errorPushActiveFollowUserInfo    = "push active follow user information failed"
	errorPushVFollowerUserInfo       = "push v follower user information failed"
	errorPushActiveFollowerUserInfo  = "push active follower user information failed"
	errDeleteVFollowUserInfo         = "delete v follow user information failed"
	errDeleteActiveFollowUserInfo    = "delete active follow user information failed"
	errDeleteVFollowerUserInfo       = "delete v follower user information failed"
	errDeleteActiveFollowerUserInfo  = "delete active follower user information failed"
	errorGetVFollowerUserInfo        = "get v follower user information failed"
	errorGetActiveFollowerUserInfo   = "get active follower user information failed"
	errorGetVFollowUserInfo          = "get v follower user information failed"
	errorPushVSet                    = "push v set failed"
	errorPushVString                 = "push v string failed"
	errorGetFollowFollowerIdList     = "get follow follower id list failed"
	errorPushVFollowFollowerInfoInit = "push v follow follower information init failed"
	errorDeleteActiveAllInfo         = "delete active all information failed"
	errorGetUserPublishVideoCount    = "get user publish video count failed"
	errorGetUserFavoritedCount       = "get user favorited count failed"
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
	userID := snowflake.GenID()
	password, err := util.ScryptPw(req.Password)
	if err != nil {
		zap.L().Error(errorEncryptPassword, zap.Error(err))
		return 0, "", err
	}
	token, err := jwt.GenToken(req.Username, userID)
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
	zap.L().Info(fmt.Sprintf("用户%s注册成功", req.Username))
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
	token, err := jwt.GenToken(req.Username, u.UserID)
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
	res, err := client.FollowClient.GetFollowFollower(context.Background(), &request1.DouyinFollowFollowerCountRequest{
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
	resGetUserPublishVideoCount, err := client.VideoClient.GetUserPublishVideoCount(context.Background(), &request2.DouyinGetUserPublishVideoCountRequest{
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
	resGetUserFavoritedCount, err := client.FavoriteClient.GetUserFavoritedCount(context.Background(), &request3.DouyinGetUserFavoritedCountRequest{
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
		res1, err := client.FollowClient.GetFollowInfo(context.Background(), &request1.DouyinGetFollowRequest{
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

func GetUserInfoList(req *request.DouyinUserListRequest) ([]*response.User, error) {
	if len(req.UserId) == 0 {
		return nil, nil
	}
	result := make([]*response.User, 0, len(req.UserId))
	//批量读取用户粉丝数和关注数
	res, err := client.FollowClient.GetFollowFollowerList(context.Background(), &request1.DouyinFollowFollowerListCountRequest{
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
	resGetUserPublishVideoCount, err := client.VideoClient.GetUserPublishVideoCountList(context.Background(), &request2.DouyinGetUserPublishVideoCountListRequest{
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
	resGetUserFavoritedCount, err := client.FavoriteClient.GetUserFavoritedCount(context.Background(), &request3.DouyinGetUserFavoritedCountRequest{
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
			res1, err = client.FollowClient.GetFollowInfoList(context.Background(), &request1.DouyinGetFollowListRequest{
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

func AddUserVideoCountSet(req *request.DouyinUserVideoCountSetRequest) error {
	return redis.AddUserVideoCountSet(req.UserId)
}

func AddUserFavoriteVideoCountSet(req *request.DouyinUserVideoCountSetRequest) error {
	return redis.AddUserFavoriteVideoCountSet(req.UserId)
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
	return status.Error(codes.Aborted, "错误的rpc调用输入")
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
		res, err := client.FollowClient.GetFollowFollowerList(context.Background(), &request1.DouyinFollowFollowerListCountRequest{
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
		resGetUserPublishVideoCount, err := client.VideoClient.GetUserPublishVideoCountList(context.Background(), &request2.DouyinGetUserPublishVideoCountListRequest{
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
		resGetUserFavoritedCount, err := client.FavoriteClient.GetUserFavoritedCount(context.Background(), &request3.DouyinGetUserFavoritedCountRequest{
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
		res1, err := client.FollowClient.GetFollowInfoList(context.Background(), &request1.DouyinGetFollowListRequest{
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
		res, err := client.FollowClient.GetFollowFollowerList(context.Background(), &request1.DouyinFollowFollowerListCountRequest{
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
		resGetUserPublishVideoCount, err := client.VideoClient.GetUserPublishVideoCountList(context.Background(), &request2.DouyinGetUserPublishVideoCountListRequest{
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
		resGetUserFavoritedCount, err := client.FavoriteClient.GetUserFavoritedCount(context.Background(), &request3.DouyinGetUserFavoritedCountRequest{
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
		res1, err := client.FollowClient.GetFollowInfoList(context.Background(), &request1.DouyinGetFollowListRequest{
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
	return nil, errors.New("错误的rpc调用输入")
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

func AddUserVideoCountSetRevert(req *request.DouyinUserVideoCountSetRequest) error {
	return redis.DeletedUserVideoCountSet(req.UserId)
}

func AddUserFavoriteVideoCountSetRevert(req *request.DouyinUserVideoCountSetRequest) error {
	return redis.DeleteUserFavoriteVideoCountSet(req.UserId)
}
