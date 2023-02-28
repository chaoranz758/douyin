package redis

import "errors"

const (
	errorUserLoginCountAdd1             = "user login count redis add 1 failed"
	errorSetExpireTime                  = "set expire time failed"
	errorKeyNotExist                    = "key not exist"
	errorAddUserLoginSet                = "add user to user login set failed"
	errorExeVActiveSet                  = "exe if request is V or active from redis set failed"
	errorVInfoNotExist                  = "v User Information from redis not exist"
	errGetVUserInfo                     = "get V User Information from redis failed"
	errJsonUnmarshal                    = "json unmarshal failed"
	errorActiveInfoNotExist             = "在执行业务逻辑的时候活跃用户升级成了大V,在活跃用户中找不到相关信息"
	errorGetActiveUserInfo              = "get active User Information from redis failed"
	errorGetVFollowerInfo               = "get v follower information failed"
	errorGetActiveFollowerInfo          = "get active follower information failed"
	errorAddUserVideoSet                = "add user video set failed"
	errorAddUserFavoriteVideoSet        = "add user favorite video set failed"
	errorGetVFollowerInfoList           = "get v follower information list failed"
	errorGetVFollowInfoList             = "get v follow information list failed"
	errorGetActiveFollowInfoList        = "get active follow information list failed"
	errorGetActiveFollowerInfoList      = "get active follower information list failed"
	errorJsonMarshal                    = "json marshal failed"
	errorJsonUnMarshal                  = "json unmarshal failed"
	errorTypeTurnFailed                 = "type turn failed"
	errorDeleteUserVideoSet             = "delete user video set failed"
	errorDeleteUserFavoriteVideoSet     = "delete user favorite video set failed"
	errorAddUserFollowerUserCountSet    = "add user follower user count set failed"
	errorDeleteUserFollowerUserCountSet = "delete user follower user count set failed"
	errorAddUserFollowUserCountSet      = "add user follow user count set failed"
	errorDeleteUserFollowUserCountSet   = "delete user follow user count set failed"
	errorNoUserInfoInVSet               = "not found user information in v set"
	errorPushVString                    = "push v string failed"
	warnKey                             = "redis监听的key正在被其它线程执行"
	errorWatchKey                       = "redis watch key failed"
	errorPushVFollowFollowerInfoInit    = "push v follow follower information init phase failed"
	errorDeleteActiveAllInfo            = "delete active all information failed"
	errorPushActiveFollowerInfoInit     = "push active follower information init phase failed"
	errorPushActiveFollowInfoInit       = "push active follow information init phase failed"
	errorPushActiveSetAndString         = "push active set and string failed"
	errorConnectRedis                   = "connect redis failed"
)

var (
	errorKeyNotExistV = errors.New("key not exist")
)
