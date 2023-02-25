package redis

import (
	"context"
	"douyin/service/user/model"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

const (
	errorUserLoginCountAdd1        = "user login count redis add 1 failed"
	errorSetExpireTime             = "set expire time failed"
	errorKeyNotExist               = "key not exist"
	errorAddUserLoginSet           = "add user to user login set failed"
	errorExeVActiveSet             = "exe if request is V or active from redis set failed"
	errorVInfoNotExist             = "v User Information from redis not exist"
	errGetVUserInfo                = "get V User Information from redis failed"
	errJsonUnmarshal               = "json unmarshal failed"
	errorActiveInfoNotExist        = "在执行业务逻辑的时候活跃用户升级成了大V,在活跃用户中找不到相关信息"
	errorGetActiveUserInfo         = "get active User Information from redis failed"
	errorGetVFollowerInfo          = "get v follower information failed"
	errorGetActiveFollowerInfo     = "get active follower information failed"
	errorAddUserVideoSet           = "add user video set failed"
	errorAddUserFavoriteVideoSet   = "add user favorite video set failed"
	errorGetVFollowerInfoList      = "get v follower information list failed"
	errorGetVFollowInfoList        = "get v follow information list failed"
	errorGetActiveFollowInfoList   = "get active follow information list failed"
	errorGetActiveFollowerInfoList = "get active follower information list failed"
	errorJsonMarshal               = "json marshal failed"
	errorJsonUnMarshal             = "json unmarshal failed"
)

var (
	errorKeyNotExistV = errors.New("key not exist")
)

func UserLoginCount(userID int64) error {
	//用户登录不需要考虑并发
	userID1 := strconv.FormatInt(userID, 10)
	keyPart := KeyUserLoginCount + userID1
	key := getKey(keyPart)
	count, err := rdb.Incr(context.Background(), key).Result()
	if err != nil {
		zap.L().Error(errorUserLoginCountAdd1, zap.Error(err))
		return err
	}
	b, err := rdb.Expire(context.Background(), key, time.Hour*72).Result()
	if err != nil {
		zap.L().Error(errorSetExpireTime, zap.Error(err))
		return err
	}
	if b == false {
		zap.L().Error(errorKeyNotExist, zap.Error(errorKeyNotExistV))
		return errorKeyNotExistV
	}
	if count == 5 {
		keySet := getKey(KeyUserLoginSet)
		_, err := rdb.SAdd(context.Background(), keySet, userID1).Result()
		if err != nil {
			zap.L().Error(errorAddUserLoginSet, zap.Error(err))
			return err
		}
	}
	return nil
}

func JudgeIsInActiveSet(userId int64) (bool, error) {
	keyActive := getKey(KeyActiveSet)
	return rdb.SIsMember(context.Background(), keyActive, userId).Result()
}

func IsVActive(userID int64) (bool, bool, error) {
	keyV := getKey(KeyVSet)
	keyActive := getKey(KeyActiveSet)
	pipe := rdb.Pipeline()
	boolCmdV := pipe.SIsMember(context.Background(), keyV, userID)
	boolCmdActive := pipe.SIsMember(context.Background(), keyActive, userID)
	_, err := pipe.Exec(context.Background())
	if err != nil {
		zap.L().Error(errorExeVActiveSet, zap.Error(err))
		return false, false, err
	}
	return boolCmdV.Val(), boolCmdActive.Val(), nil
}

func IsVActiveList(userID []int64) ([]bool, []bool, error) {
	keyV := getKey(KeyVSet)
	keyActive := getKey(KeyActiveSet)
	var boolCmdVList []*redis.BoolCmd
	var boolCmdActiveList []*redis.BoolCmd
	pipe := rdb.Pipeline()
	for i := 0; i < len(userID); i++ {
		boolCmdV := pipe.SIsMember(context.Background(), keyV, userID[i])
		boolCmdVList = append(boolCmdVList, boolCmdV)
		boolCmdActive := pipe.SIsMember(context.Background(), keyActive, userID[i])
		boolCmdActiveList = append(boolCmdActiveList, boolCmdActive)
	}
	_, err := pipe.Exec(context.Background())
	if err != nil {
		zap.L().Error(errorExeVActiveSet, zap.Error(err))
		return nil, nil, err
	}
	var isVList []bool
	var isActiveList []bool
	for i := 0; i < len(userID); i++ {
		isV := boolCmdVList[i].Val()
		isActive := boolCmdActiveList[i].Val()
		isVList = append(isVList, isV)
		isActiveList = append(isActiveList, isActive)
	}
	return isVList, isActiveList, nil
}

func GetVUserInfo(userID int64, u *model.User) error {
	userID1 := strconv.FormatInt(userID, 10)
	keyV := KeyVUserInfo + userID1
	keyVs := getKey(keyV)
	value, err := rdb.Get(context.Background(), keyVs).Result()
	if err != nil {
		if err == redis.Nil {
			zap.L().Error(errorVInfoNotExist, zap.Error(err))
			return err
		}
		zap.L().Error(errGetVUserInfo, zap.Error(err))
		return err
	}
	if err := json.Unmarshal([]byte(value), u); err != nil {
		zap.L().Error(errJsonUnmarshal, zap.Error(err))
		return err
	}
	return nil
}

func GetVUserInfoList(userID []int64) ([]model.User, error) {
	var keys []string
	for i := 0; i < len(userID); i++ {
		userID1 := strconv.FormatInt(userID[i], 10)
		keyV := KeyVUserInfo + userID1
		keyVs := getKey(keyV)
		keys = append(keys, keyVs)
	}
	value, err := rdb.MGet(context.Background(), keys...).Result()
	if err != nil {
		if err == redis.Nil {
			zap.L().Error(errorVInfoNotExist, zap.Error(err))
			return nil, err
		}
		zap.L().Error(errGetVUserInfo, zap.Error(err))
		return nil, err
	}
	if len(value) == 0 {
		zap.L().Error("在大V用户的集合中没有找到相应用户信息")
		return nil, err
	}
	var us []model.User
	for i := 0; i < len(value); i++ {
		v, ok := value[i].(string)
		if !ok {
			return nil, errors.New("类型转换失败")
		}
		var u model.User
		if err := json.Unmarshal([]byte(v), &u); err != nil {
			zap.L().Error(errJsonUnmarshal, zap.Error(err))
			return nil, err
		}
		us = append(us, u)
	}
	return us, nil
}

func GetActiveUserInfo(userID int64, u *model.User) error {
	userID1 := strconv.FormatInt(userID, 10)
	keyActive := KeyActiveUserInfo + userID1
	key := getKey(keyActive)
	value, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			zap.L().Info(errorActiveInfoNotExist, zap.Error(err))
			return err
		}
		zap.L().Error(errorGetActiveUserInfo, zap.Error(err))
		return err
	}
	if err := json.Unmarshal([]byte(value), u); err != nil {
		zap.L().Error(errJsonUnmarshal, zap.Error(err))
		return err
	}
	return nil
}

func GetActiveUserInfoList(userID []int64) ([]model.User, error) {
	var keys []string
	for i := 0; i < len(userID); i++ {
		userID1 := strconv.FormatInt(userID[i], 10)
		keyV := KeyActiveUserInfo + userID1
		keyVs := getKey(keyV)
		keys = append(keys, keyVs)
	}
	value, err := rdb.MGet(context.Background(), keys...).Result()
	if err != nil {
		if err == redis.Nil {
			zap.L().Error(errorActiveInfoNotExist, zap.Error(err))
			return nil, err
		}
		zap.L().Error(errorGetActiveUserInfo, zap.Error(err))
		return nil, err
	}
	//fmt.Printf("%v\n", value)
	var us []model.User
	for i := 0; i < len(value); i++ {
		v, ok := value[i].(string)
		if !ok {
			return nil, errors.New("类型转换失败")
		}
		var u model.User
		if err := json.Unmarshal([]byte(v), &u); err != nil {
			zap.L().Error(errJsonUnmarshal, zap.Error(err))
			return nil, err
		}
		us = append(us, u)
	}
	return us, nil
}

func GetVFollowerInfo(userID, loginUserID int64) (bool, error) {
	userID1 := strconv.FormatInt(userID, 10)
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	keyVFollowerInfoPart := KeyVFollowerInfo + userID1
	keyVFollowerInfo := getKey(keyVFollowerInfoPart)
	b, err := rdb.HExists(context.Background(), keyVFollowerInfo, loginUserID1).Result()
	if err != nil {
		zap.L().Error(errorGetVFollowerInfo, zap.Error(err))
		return false, err
	}
	return b, nil
}

func GetVFollowerInfoList(userID []int64, loginUserID int64) ([]bool, error) {
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	var bs []bool
	pipe := rdb.Pipeline()
	for i := 0; i < len(userID); i++ {
		userID1 := strconv.FormatInt(userID[i], 10)
		keyVFollowerInfoPart := KeyVFollowerInfo + userID1
		keyVFollowerInfo := getKey(keyVFollowerInfoPart)
		pipe.HExists(context.Background(), keyVFollowerInfo, loginUserID1)
	}
	cmds, err := pipe.Exec(context.Background())
	if err != nil {
		zap.L().Error(errorGetVFollowerInfoList, zap.Error(err))
		return nil, err
	}
	for _, value := range cmds {
		boolCmd := value.(*redis.BoolCmd).Val()
		bs = append(bs, boolCmd)
	}
	return bs, nil
}

func GetVFollowInfoList(userID []int64, loginUserID int64) ([]bool, error) {
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	var bs []bool
	pipe := rdb.Pipeline()
	for i := 0; i < len(userID); i++ {
		userID1 := strconv.FormatInt(userID[i], 10)
		keyVFollowInfoPart := KeyVFollowInfo + userID1
		keyVFollowInfo := getKey(keyVFollowInfoPart)
		pipe.HExists(context.Background(), keyVFollowInfo, loginUserID1)
	}
	cmds, err := pipe.Exec(context.Background())
	if err != nil {
		zap.L().Error(errorGetVFollowInfoList, zap.Error(err))
		return nil, err
	}
	for _, value := range cmds {
		boolCmd := value.(*redis.BoolCmd).Val()
		bs = append(bs, boolCmd)
	}
	return bs, nil
}

func GetActiveFollowInfoList(userID []int64, loginUserID int64) ([]bool, error) {
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	var bs []bool
	pipe := rdb.Pipeline()
	for i := 0; i < len(userID); i++ {
		userID1 := strconv.FormatInt(userID[i], 10)
		keyActiveFollowInfoPart := KeyActiveFollowInfo + userID1
		keyActiveFollowInfo := getKey(keyActiveFollowInfoPart)
		pipe.HExists(context.Background(), keyActiveFollowInfo, loginUserID1)
	}
	cmds, err := pipe.Exec(context.Background())
	if err != nil {
		zap.L().Error(errorGetActiveFollowInfoList, zap.Error(err))
		return nil, err
	}
	for _, value := range cmds {
		boolCmd := value.(*redis.BoolCmd).Val()
		bs = append(bs, boolCmd)
	}
	return bs, nil
}

func GetActiveFollowerInfo(userID, loginUserID int64) (bool, error) {
	userID1 := strconv.FormatInt(userID, 10)
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	keyActiveFollowerInfoPart := KeyActiveFollowerInfo + userID1
	keyActiveFollowerInfo := getKey(keyActiveFollowerInfoPart)
	b, err := rdb.HExists(context.Background(), keyActiveFollowerInfo, loginUserID1).Result()
	if err != nil {
		zap.L().Error(errorGetActiveFollowerInfo, zap.Error(err))
		return false, err
	}
	return b, nil
}

func GetActiveFollowerInfoList(userID []int64, loginUserID int64) ([]bool, error) {
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	var bs []bool
	pipe := rdb.Pipeline()
	for i := 0; i < len(userID); i++ {
		userID1 := strconv.FormatInt(userID[i], 10)
		keyVFollowerInfoPart := KeyActiveFollowerInfo + userID1
		keyVFollowerInfo := getKey(keyVFollowerInfoPart)
		pipe.HExists(context.Background(), keyVFollowerInfo, loginUserID1)
	}
	cmds, err := pipe.Exec(context.Background())
	if err != nil {
		zap.L().Error(errorGetActiveFollowerInfoList, zap.Error(err))
		return nil, err
	}
	for _, value := range cmds {
		boolCmd := value.(*redis.BoolCmd).Val()
		bs = append(bs, boolCmd)
	}
	return bs, nil
}

func AddUserVideoCountSet(userID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	keySet := getKey(KeyUserVideoSet)
	_, err := rdb.SAdd(context.Background(), keySet, userID1).Result()
	if err != nil {
		zap.L().Error(errorAddUserVideoSet, zap.Error(err))
		return status.Error(codes.Aborted, errorAddUserVideoSet)
	}
	return nil
}

func AddUserFavoriteVideoCountSet(userID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	keySet := getKey(KeyUserFavoriteVideoSet)
	_, err := rdb.SAdd(context.Background(), keySet, userID1).Result()
	if err != nil {
		zap.L().Error(errorAddUserFavoriteVideoSet, zap.Error(err))
		return status.Error(codes.Aborted, err.Error())
	}
	return nil
}

func PushVSet(userID int64) (bool, error) {
	userID1 := strconv.FormatInt(userID, 10)
	keySet := getKey(KeyVSet)
	b, err := rdb.SIsMember(context.Background(), keySet, userID1).Result()
	if err != nil {
		return false, err
	}
	return b, err
}

func PushVString(userID int64, user model.User) error {
	userID1 := strconv.FormatInt(userID, 10)
	keyV := KeyVUserInfo + userID1
	keyV1 := getKey(keyV)
	keySet := getKey(KeyVSet)
	data, err := json.Marshal(user)
	if err != nil {
		zap.L().Error(errorJsonMarshal, zap.Error(err))
		return err
	}
	if err := rdb.Watch(context.Background(), func(tx *redis.Tx) error {
		for {
			_, err := tx.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
				pipe.SAdd(context.Background(), keySet, userID)
				pipe.Set(context.Background(), keyV1, string(data), 0)
				return nil
			})
			if err == redis.TxFailedErr {
				zap.L().Warn("redis监听的key正在被其它线程执行")
				continue
			}
			if err != nil && err != redis.TxFailedErr {
				zap.L().Error("exec push v string failed", zap.Error(err))
				return err
			}
			break
		}
		return nil
	}, keyV1, keySet); err != nil {
		zap.L().Error("redis watch key failed", zap.Error(err))
		return err
	}
	return nil
}

func GetActiveSet() ([]string, error) {
	keyLoginSet := getKey(KeyUserLoginSet)
	keyVideoSet := getKey(KeyUserVideoSet)
	keyFavoriteVideoSet := getKey(KeyUserFavoriteVideoSet)
	keyFollowSet := getKey(KeyUserFollowCountSet)
	keyFollowerSet := getKey(KeyUserFollowerCountSet)
	var keys = []string{
		keyLoginSet,
		keyVideoSet,
		keyFavoriteVideoSet,
		keyFollowSet,
		keyFollowerSet,
	}
	_, err := rdb.SInterStore(context.Background(), "temporary", keys...).Result()
	if err != nil {
		return nil, err
	}
	keyActive := getKey(KeyActiveSet)
	result, err := rdb.SDiff(context.Background(), "temporary", keyActive).Result()
	if err != nil {
		return nil, err
	}
	err = rdb.Unlink(context.Background(), "temporary").Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func PushActiveSetAndString(users []model.User, idList []int64) error {
	key1 := getKey(KeyActiveSet)
	var datas []string
	var idList1 []string
	var keys []string
	for i := 0; i < len(idList); i++ {
		userId := strconv.FormatInt(idList[i], 10)
		key2Part := KeyActiveUserInfo + userId
		key2 := getKey(key2Part)
		datas = append(datas, key2)
		keys = append(keys, key2)
		data, err := json.Marshal(users[i])
		if err != nil {
			zap.L().Error(errorJsonMarshal, zap.Error(err))
			return err
		}
		datas = append(datas, string(data))
		idList1 = append(idList1, userId)
	}
	keys = append(keys, key1)
	if err := rdb.Watch(context.Background(), func(tx *redis.Tx) error {
		for {
			_, err := tx.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
				pipe.SAdd(context.Background(), key1, idList1)
				pipe.MSet(context.Background(), datas)
				return nil
			})
			if err == redis.TxFailedErr {
				zap.L().Warn("redis监听的key正在被其它线程执行")
				continue
			}
			if err != nil && err != redis.TxFailedErr {
				zap.L().Error("exec push active set and string failed", zap.Error(err))
				return err
			}
			break
		}
		return nil
	}, keys...); err != nil {
		zap.L().Error("redis watch key failed", zap.Error(err))
		return err
	}
	return nil
}

func AddUserFollowUserCountSet(userID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	keySet := getKey(KeyUserFollowCountSet)
	_, err := rdb.SAdd(context.Background(), keySet, userID1).Result()
	if err != nil {
		zap.L().Error(errorAddUserFavoriteVideoSet, zap.Error(err))
		return status.Error(codes.Aborted, err.Error())
	}
	return nil
}

func DeleteUserFollowUserCountSet(userID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	keySet := getKey(KeyUserFollowCountSet)
	if err := rdb.SRem(context.Background(), keySet, userID1).Err(); err != nil {
		zap.L().Error("delete user follow user count set failed")
		return status.New(codes.Aborted, "delete user follow user count set failed").Err()
	}
	return nil
}

func AddUserFollowerUserCountSet(userID int64) error {
	println(userID)
	userID1 := strconv.FormatInt(userID, 10)
	keySet := getKey(KeyUserFollowerCountSet)
	_, err := rdb.SAdd(context.Background(), keySet, userID1).Result()
	if err != nil {
		zap.L().Error(errorAddUserFavoriteVideoSet, zap.Error(err))
		return status.Error(codes.Aborted, err.Error())
	}
	return nil
}

func DeleteUserFollowerUserCountSet(userID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	keySet := getKey(KeyUserFollowerCountSet)
	if err := rdb.SRem(context.Background(), keySet, userID1).Err(); err != nil {
		zap.L().Error("delete user follower user count set failed")
		return status.New(codes.Aborted, "delete user follower user count set failed").Err()
	}
	return nil
}

func PushVFollowUserInfo(u model.User, userID, loginUserID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	keyVFollowInfoPart := KeyVFollowInfo + loginUserID1
	keyVFollowInfo := getKey(keyVFollowInfoPart)
	data, err := json.Marshal(u)
	if err != nil {
		zap.L().Error(errorJsonMarshal, zap.Error(err))
		return err
	}
	return rdb.HSet(context.Background(), keyVFollowInfo, userID1, data).Err()
}

func PushActiveFollowUserInfo(u model.User, userID, loginUserID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	keyActiveFollowInfoPart := KeyActiveFollowInfo + loginUserID1
	keyActiveFollowInfo := getKey(keyActiveFollowInfoPart)
	data, err := json.Marshal(u)
	if err != nil {
		zap.L().Error(errorJsonMarshal, zap.Error(err))
		return err
	}
	return rdb.HSet(context.Background(), keyActiveFollowInfo, userID1, data).Err()
}

func PushVFollowerUserInfo(u model.User, userID, loginUserID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	keyVFollowerInfoPart := KeyVFollowerInfo + userID1
	keyVFollowerInfo := getKey(keyVFollowerInfoPart)
	data, err := json.Marshal(u)
	if err != nil {
		zap.L().Error(errorJsonMarshal, zap.Error(err))
		return err
	}
	return rdb.HSet(context.Background(), keyVFollowerInfo, loginUserID1, data).Err()
}

func PushActiveFollowerUserInfo(u model.User, userID, loginUserID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	keyActiveFollowerInfoPart := KeyActiveFollowerInfo + userID1
	keyActiveFollowerInfo := getKey(keyActiveFollowerInfoPart)
	data, err := json.Marshal(u)
	if err != nil {
		zap.L().Error(errorJsonMarshal, zap.Error(err))
		return err
	}
	return rdb.HSet(context.Background(), keyActiveFollowerInfo, loginUserID1, data).Err()
}

func DeleteVFollowUserInfo(userID, loginUserID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	keyVFollowInfoPart := KeyVFollowInfo + loginUserID1
	keyVFollowInfo := getKey(keyVFollowInfoPart)
	return rdb.HDel(context.Background(), keyVFollowInfo, userID1).Err()
}

func DeleteActiveFollowUserInfo(userID, loginUserID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	keyActiveFollowInfoPart := KeyActiveFollowInfo + loginUserID1
	keyActiveFollowInfo := getKey(keyActiveFollowInfoPart)
	return rdb.HDel(context.Background(), keyActiveFollowInfo, userID1).Err()
}

func DeleteVFollowerUserInfo(userID, loginUserID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	keyVFollowerInfoPart := KeyVFollowerInfo + userID1
	keyVFollowerInfo := getKey(keyVFollowerInfoPart)
	return rdb.HDel(context.Background(), keyVFollowerInfo, loginUserID1).Err()
}

func DeleteActiveFollowerUserInfo(userID, loginUserID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	keyActiveFollowerInfoPart := KeyActiveFollowerInfo + userID1
	keyActiveFollowerInfo := getKey(keyActiveFollowerInfoPart)
	return rdb.HDel(context.Background(), keyActiveFollowerInfo, loginUserID1).Err()
}

func GetVFollowUserInfo(userID int64) ([]model.User, error) {
	userID1 := strconv.FormatInt(userID, 10)
	keyVFollowInfoPart := KeyVFollowInfo + userID1
	keyVFollowInfo := getKey(keyVFollowInfoPart)
	datas, err := rdb.HVals(context.Background(), keyVFollowInfo).Result()
	if err != nil {
		return nil, err
	}
	if len(datas) == 0 {
		return nil, nil
	}
	var users []model.User
	for i := 0; i < len(datas); i++ {
		var user model.User
		if err = json.Unmarshal([]byte(datas[i]), &user); err != nil {
			zap.L().Error(errorJsonUnMarshal, zap.Error(err))
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetActiveFollowUserInfo(userID int64) ([]model.User, error) {
	userID1 := strconv.FormatInt(userID, 10)
	keyActiveFollowInfoPart := KeyActiveFollowInfo + userID1
	keyActiveFollowInfo := getKey(keyActiveFollowInfoPart)
	datas, err := rdb.HVals(context.Background(), keyActiveFollowInfo).Result()
	if err != nil {
		return nil, err
	}
	if len(datas) == 0 {
		return nil, nil
	}
	var users []model.User
	for i := 0; i < len(datas); i++ {
		var user model.User
		if err = json.Unmarshal([]byte(datas[i]), &user); err != nil {
			zap.L().Error(errorJsonUnMarshal, zap.Error(err))
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetVFollowerUserInfo(userID int64) ([]model.User, error) {
	userID1 := strconv.FormatInt(userID, 10)
	keyVFollowerInfoPart := KeyVFollowerInfo + userID1
	keyVFollowerInfo := getKey(keyVFollowerInfoPart)
	datas, err := rdb.HVals(context.Background(), keyVFollowerInfo).Result()
	if err != nil {
		return nil, err
	}
	if len(datas) == 0 {
		return nil, nil
	}
	var users []model.User
	for i := 0; i < len(datas); i++ {
		var user model.User
		if err = json.Unmarshal([]byte(datas[i]), &user); err != nil {
			zap.L().Error(errorJsonUnMarshal, zap.Error(err))
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetActiveFollowerUserInfo(userID int64) ([]model.User, error) {
	userID1 := strconv.FormatInt(userID, 10)
	keyActiveFollowerInfoPart := KeyActiveFollowerInfo + userID1
	keyActiveFollowerInfo := getKey(keyActiveFollowerInfoPart)
	datas, err := rdb.HVals(context.Background(), keyActiveFollowerInfo).Result()
	if err != nil {
		return nil, err
	}
	if len(datas) == 0 {
		return nil, nil
	}
	var users []model.User
	for i := 0; i < len(datas); i++ {
		var user model.User
		if err = json.Unmarshal([]byte(datas[i]), &user); err != nil {
			zap.L().Error(errorJsonUnMarshal, zap.Error(err))
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func PushVFollowFollowerInfoInit(usersFollow, usersFollower []model.User, userId int64) error {
	var datasFollow []string
	var datasFollower []string
	for i := 0; i < len(usersFollow); i++ {
		data, err := json.Marshal(usersFollow[i])
		if err != nil {
			zap.L().Error(errorJsonMarshal, zap.Error(err))
			return err
		}
		userID1 := strconv.FormatInt(usersFollow[i].UserID, 10)
		datasFollow = append(datasFollow, userID1)
		datasFollow = append(datasFollow, string(data))
	}
	for i := 0; i < len(usersFollower); i++ {
		data, err := json.Marshal(usersFollower[i])
		if err != nil {
			zap.L().Error(errorJsonMarshal, zap.Error(err))
			return err
		}
		userID1 := strconv.FormatInt(usersFollower[i].UserID, 10)
		datasFollower = append(datasFollower, userID1)
		datasFollower = append(datasFollower, string(data))
	}
	userID1 := strconv.FormatInt(userId, 10)
	keyFollowPart := KeyVFollowInfo + userID1
	keyFollow := getKey(keyFollowPart)
	keyFollowerPart := KeyVFollowerInfo + userID1
	keyFollower := getKey(keyFollowerPart)
	if err := rdb.Watch(context.Background(), func(tx *redis.Tx) error {
		for {
			_, err := tx.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
				if len(usersFollow) != 0 {
					pipe.HMSet(context.Background(), keyFollow, datasFollow)
				}
				pipe.HMSet(context.Background(), keyFollower, datasFollower)
				return nil
			})
			if err == redis.TxFailedErr {
				zap.L().Warn("redis监听的key正在被其它线程执行")
				continue
			}
			if err != nil && err != redis.TxFailedErr {
				zap.L().Error("exec push v follow follower information failed", zap.Error(err))
				return err
			}
			break
		}
		return nil
	}, keyFollow, keyFollower); err != nil {
		zap.L().Error("redis watch key failed", zap.Error(err))
		return err
	}
	return nil
}

func PushActiveFollowInfoInit(us []model.User, userId int64) error {
	var datasFollow []string
	for i := 0; i < len(us); i++ {
		data, err := json.Marshal(us[i])
		if err != nil {
			zap.L().Error(errorJsonMarshal, zap.Error(err))
			return err
		}
		userID1 := strconv.FormatInt(us[i].UserID, 10)
		datasFollow = append(datasFollow, userID1)
		datasFollow = append(datasFollow, string(data))
	}
	userID1 := strconv.FormatInt(userId, 10)
	keyFollowPart := KeyActiveFollowInfo + userID1
	keyFollow := getKey(keyFollowPart)
	if err := rdb.Watch(context.Background(), func(tx *redis.Tx) error {
		for {
			_, err := tx.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
				pipe.HMSet(context.Background(), keyFollow, datasFollow)
				return nil
			})
			if err == redis.TxFailedErr {
				zap.L().Warn("redis监听的key正在被其它线程执行")
				continue
			}
			if err != nil && err != redis.TxFailedErr {
				zap.L().Error("exec push active follow information failed", zap.Error(err))
				return err
			}
			break
		}
		return nil
	}, keyFollow); err != nil {
		zap.L().Error("redis watch key failed", zap.Error(err))
		return err
	}
	return nil
}

func PushActiveFollowerInfoInit(us []model.User, userId int64) error {
	var datasFollow []string
	for i := 0; i < len(us); i++ {
		data, err := json.Marshal(us[i])
		if err != nil {
			zap.L().Error(errorJsonMarshal, zap.Error(err))
			return err
		}
		userID1 := strconv.FormatInt(us[i].UserID, 10)
		datasFollow = append(datasFollow, userID1)
		datasFollow = append(datasFollow, string(data))
	}
	userID1 := strconv.FormatInt(userId, 10)
	keyFollowPart := KeyActiveFollowerInfo + userID1
	keyFollow := getKey(keyFollowPart)
	if err := rdb.Watch(context.Background(), func(tx *redis.Tx) error {
		for {
			_, err := tx.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
				pipe.HMSet(context.Background(), keyFollow, datasFollow)
				return nil
			})
			if err == redis.TxFailedErr {
				zap.L().Warn("redis监听的key正在被其它线程执行")
				continue
			}
			if err != nil && err != redis.TxFailedErr {
				zap.L().Error("exec push active follower information failed", zap.Error(err))
				return err
			}
			break
		}
		return nil
	}, keyFollow); err != nil {
		zap.L().Error("redis watch key failed", zap.Error(err))
		return err
	}
	return nil
}

func DeleteActiveAllInfo(userId int64) error {
	//1.移除activeset
	keyActive := getKey(KeyActiveSet)
	//2.清楚用户信息 string
	userID1 := strconv.FormatInt(userId, 10)
	key := KeyActiveUserInfo + userID1
	keyUserInfo := getKey(key)
	//3.清楚粉丝信息
	keyActiveFollowerInfoPart := KeyActiveFollowerInfo + userID1
	keyFollowerInfo := getKey(keyActiveFollowerInfoPart)
	//4.清楚关注信息
	keyActiveFollowInfoPart := KeyActiveFollowInfo + userID1
	keyFollowInfo := getKey(keyActiveFollowInfoPart)
	if err := rdb.Watch(context.Background(), func(tx *redis.Tx) error {
		for {
			_, err := tx.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
				pipe.SRem(context.Background(), keyActive, userId)
				pipe.Del(context.Background(), keyUserInfo, keyFollowInfo, keyFollowerInfo)
				return nil
			})
			if err == redis.TxFailedErr {
				zap.L().Warn("redis监听的key正在被其它线程执行")
				continue
			}
			if err != nil && err != redis.TxFailedErr {
				zap.L().Error("exec deleteActiveAllInfo failed", zap.Error(err))
				return err
			}
			break
		}
		return nil
	}, keyActive, keyUserInfo, keyFollowInfo, keyFollowerInfo); err != nil {
		zap.L().Error("redis watch key failed", zap.Error(err))
		return err
	}
	return nil
}

func DeletedUserVideoCountSet(userID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	keySet := getKey(KeyUserVideoSet)
	if err := rdb.SRem(context.Background(), keySet, userID1).Err(); err != nil {
		zap.L().Error("delete user video count set failed")
		return status.New(codes.Aborted, "delete user video count set failed").Err()
	}
	return nil
}

func DeleteUserFavoriteVideoCountSet(userID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	keySet := getKey(KeyUserFavoriteVideoSet)
	if err := rdb.SRem(context.Background(), keySet, userID1).Err(); err != nil {
		zap.L().Error("delete user video count set failed")
		return status.New(codes.Aborted, "delete user video count set failed").Err()
	}
	return nil
}
