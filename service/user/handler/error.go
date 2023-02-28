package handler

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
