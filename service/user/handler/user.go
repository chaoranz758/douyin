package handler

import (
	"context"
	"douyin/proto/user/api"
	"douyin/proto/user/request"
	"douyin/proto/user/response"
	"douyin/service/user/service"
	"go.uber.org/zap"
)

const (
	errorCreateUser                         = "create user failed"
	errorUserLogin                          = "user login failed"
	errorGetUserInfo                        = "get user information failed"
	errorExeVActiveSet                      = "exe if request is V or active from redis set failed"
	errorGetUserInfoList                    = "get user information list failed"
	errorAddUserVideoCountSet               = "add user video count set failed"
	errorAddUserFavoriteVideoCountSet       = "add user favorite video count set failed"
	errorPushVActiveFollowerUserinfo        = "push v or active follower user information failed"
	errorPushVActiveFollowerUserinfoRevert  = "push v or active follower user information revert failed"
	errorDeleteVActiveFollowerUserinfo      = "delete v or active follower user information failed"
	errorGetVActiveFollowerUserinfo         = "get v or active follower information failed"
	errorPushVSet                           = "push v set failed"
	errorAddUserFollowUserCountSet          = "add user follow user count failed"
	errorAddUserFollowerUserCountSet        = "add user follower user count failed"
	errorAddUserVideoCountSetRevert         = "add user video count set revert failed"
	errorAddUserFavoriteVideoCountSetRevert = "add user favorite video count set revert failed"
	errorAddUserFollowUserCountSetRevert    = "add user follow user count set revert failed"
	errorAddUserFollowerUserCountSetRevert  = "add user follower user count set revert failed"
)

var u1 response.User

type User struct {
	api.UnimplementedUserServer
}

func (user *User) UserRegister(ctx context.Context, req *request.DouyinUserRegisterRequest) (*response.DouyinUserRegisterResponse, error) {
	userID, token, err := service.UserRegister(req)
	if err != nil {
		zap.L().Error(errorCreateUser, zap.Error(err))
		return &response.DouyinUserRegisterResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinUserRegisterResponse{
		Code:   1,
		UserId: userID,
		Token:  token,
	}, nil
}

func (user *User) UserLogin(ctx context.Context, req *request.DouyinUserLoginRequest) (*response.DouyinUserLoginResponse, error) {
	userID, token, err := service.UserLogin(req)
	if err != nil {
		zap.L().Error(errorUserLogin, zap.Error(err))
		return &response.DouyinUserLoginResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinUserLoginResponse{
		Code:   1,
		UserId: userID,
		Token:  token,
	}, nil
}

func (user *User) GetUserInfo(ctx context.Context, req *request.DouyinUserRequest) (*response.DouyinUserResponse, error) {
	u, b, count1, count2, count3, err := service.GetUserInfo(req)
	if err != nil {
		zap.L().Error(errorGetUserInfo, zap.Error(err))
		return &response.DouyinUserResponse{
			User: nil,
			Code: 2,
		}, err
	}
	u1.Id = u.UserID
	u1.Name = u.UserName
	u1.FollowCount = u.Follow
	u1.FollowerCount = u.Follower
	u1.IsFollow = b
	u1.Avatar = u.Avatar
	u1.BackgroundImage = u.BackgroundImage
	u1.Signature = u.Signature
	u1.TotalFavorited = count1
	u1.WorkCount = count2
	u1.FavoriteCount = count3
	return &response.DouyinUserResponse{
		User: &u1,
		Code: 1,
	}, nil
}

func (user *User) GetUserInfoList(ctx context.Context, req *request.DouyinUserListRequest) (*response.DouyinUserListResponse, error) {
	data, err := service.GetUserInfoList(req)
	if err != nil {
		zap.L().Error(errorGetUserInfoList, zap.Error(err))
		return &response.DouyinUserListResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinUserListResponse{
		Code: 1,
		User: data,
	}, nil
}

func (user *User) UserIsInfluencerActiver(ctx context.Context, req *request.DouyinUserIsInfluencerActiverRequest) (*response.DouyinUserIsInfluencerActiverResponse, error) {
	isV, isActive, err := service.UserIsInfluencerActiver(req)
	if err != nil {
		zap.L().Error(errorExeVActiveSet, zap.Error(err))
		return &response.DouyinUserIsInfluencerActiverResponse{
			IsInfluencer: false,
			IsActiver:    false,
		}, err
	}
	return &response.DouyinUserIsInfluencerActiverResponse{
		IsInfluencer: isV,
		IsActiver:    isActive,
	}, nil
}

func (user *User) AddUserVideoCountSet(ctx context.Context, req *request.DouyinUserVideoCountSetRequest) (*response.DouyinUserVideoCountSetResponse, error) {
	if err := service.AddUserVideoCountSet(req); err != nil {
		zap.L().Error(errorAddUserVideoCountSet, zap.Error(err))
		return &response.DouyinUserVideoCountSetResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinUserVideoCountSetResponse{
		Code: 1,
	}, nil
}

func (user *User) AddUserFavoriteVideoCountSet(ctx context.Context, req *request.DouyinUserVideoCountSetRequest) (*response.DouyinUserVideoCountSetResponse, error) {
	if err := service.AddUserFavoriteVideoCountSet(req); err != nil {
		zap.L().Error(errorAddUserFavoriteVideoCountSet, zap.Error(err))
		return &response.DouyinUserVideoCountSetResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinUserVideoCountSetResponse{
		Code: 1,
	}, nil
}

func (user *User) PushVActiveFollowerUserinfo(ctx context.Context, req *request.DouyinVActiveFollowFollowerUserinfoRequest) (*response.DouyinVActiveFollowFollowerUserinfoResponse, error) {
	if err := service.PushVActiveFollowerUserinfo(req); err != nil {
		zap.L().Error(errorPushVActiveFollowerUserinfo, zap.Error(err))
		return &response.DouyinVActiveFollowFollowerUserinfoResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinVActiveFollowFollowerUserinfoResponse{
		Code: 1,
	}, nil
}

func (user *User) PushVActiveFollowerUserinfoRevert(ctx context.Context, req *request.DouyinVActiveFollowFollowerUserinfoRequest) (*response.DouyinVActiveFollowFollowerUserinfoResponse, error) {
	if err := service.PushVActiveFollowerUserinfoRevert(req); err != nil {
		zap.L().Error(errorPushVActiveFollowerUserinfoRevert, zap.Error(err))
		return &response.DouyinVActiveFollowFollowerUserinfoResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinVActiveFollowFollowerUserinfoResponse{
		Code: 1,
	}, nil
}

func (user *User) DeleteVActiveFollowerUserinfo(ctx context.Context, req *request.DouyinDeleteVActiveFollowFollowerUserinfoRequest) (*response.DouyinDeleteVActiveFollowFollowerUserinfoResponse, error) {
	if err := service.DeleteVActiveFollowerUserinfo(req); err != nil {
		zap.L().Error(errorDeleteVActiveFollowerUserinfo, zap.Error(err))
		return &response.DouyinDeleteVActiveFollowFollowerUserinfoResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinDeleteVActiveFollowFollowerUserinfoResponse{
		Code: 1,
	}, nil
}

func (user *User) GetVActiveFollowerUserinfo(ctx context.Context, req *request.DouyinGetVActiveFollowFollowerUserinfoRequest) (*response.DouyinGetVActiveFollowFollowerUserinfoResponse, error) {
	data, err := service.GetVActiveFollowerUserinfo(req)
	if err != nil {
		zap.L().Error(errorGetVActiveFollowerUserinfo, zap.Error(err))
		return &response.DouyinGetVActiveFollowFollowerUserinfoResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinGetVActiveFollowFollowerUserinfoResponse{
		User: data,
		Code: 1,
	}, nil
}

func (user *User) PushVUserRelativeInfoInit(ctx context.Context, req *request.DouyinPushVSetRequest) (*response.DouyinPushVSetResponse, error) {
	isActive, b, err := service.PushVSet(req)
	if err != nil {
		zap.L().Error(errorPushVSet, zap.Error(err))
		return &response.DouyinPushVSetResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinPushVSetResponse{
		IsActive: isActive,
		IsExist:  b,
		Code:     1,
	}, nil
}

func (user *User) AddUserFollowUserCountSet(ctx context.Context, req *request.DouyinUserFollowCountSetRequest) (*response.DouyinUserFollowCountSetResponse, error) {
	if err := service.AddUserFollowUserCountSet(req); err != nil {
		zap.L().Error(errorAddUserFollowUserCountSet, zap.Error(err))
		return &response.DouyinUserFollowCountSetResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinUserFollowCountSetResponse{
		Code: 1,
	}, nil
}

func (user *User) AddUserFollowUserCountSetRevert(ctx context.Context, req *request.DouyinUserFollowCountSetRequest) (*response.DouyinUserFollowCountSetResponse, error) {
	if err := service.AddUserFollowUserCountSetRevert(req); err != nil {
		zap.L().Error(errorAddUserFollowUserCountSetRevert, zap.Error(err))
		return &response.DouyinUserFollowCountSetResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinUserFollowCountSetResponse{
		Code: 1,
	}, nil

}

func (user *User) AddUserFollowerUserCountSet(ctx context.Context, req *request.DouyinUserFollowerCountSetRequest) (*response.DouyinUserFollowerCountSetResponse, error) {
	if err := service.AddUserFollowerUserCountSet(req); err != nil {
		zap.L().Error(errorAddUserFollowerUserCountSet, zap.Error(err))
		return &response.DouyinUserFollowerCountSetResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinUserFollowerCountSetResponse{
		Code: 1,
	}, nil
}

func (user *User) AddUserFollowerUserCountSetRevert(ctx context.Context, req *request.DouyinUserFollowerCountSetRequest) (*response.DouyinUserFollowerCountSetResponse, error) {
	if err := service.AddUserFollowerUserCountSetRevert(req); err != nil {
		zap.L().Error(errorAddUserFollowerUserCountSetRevert, zap.Error(err))
		return &response.DouyinUserFollowerCountSetResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinUserFollowerCountSetResponse{
		Code: 1,
	}, nil
}

func (user *User) AddUserVideoCountSetRevert(ctx context.Context, req *request.DouyinUserVideoCountSetRequest) (*response.DouyinUserVideoCountSetResponse, error) {
	if err := service.AddUserVideoCountSetRevert(req); err != nil {
		zap.L().Error(errorAddUserVideoCountSetRevert, zap.Error(err))
		return &response.DouyinUserVideoCountSetResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinUserVideoCountSetResponse{
		Code: 1,
	}, nil
}

func (user *User) AddUserFavoriteVideoCountSetRevert(ctx context.Context, req *request.DouyinUserVideoCountSetRequest) (*response.DouyinUserVideoCountSetResponse, error) {
	if err := service.AddUserFavoriteVideoCountSetRevert(req); err != nil {
		zap.L().Error(errorAddUserFavoriteVideoCountSetRevert, zap.Error(err))
		return &response.DouyinUserVideoCountSetResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinUserVideoCountSetResponse{
		Code: 1,
	}, nil
}
