package redis

import (
	"context"
	"douyin/service/user/model"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
)

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
		zap.L().Error(errorNoUserInfoInVSet)
		return nil, err
	}
	var us []model.User
	for i := 0; i < len(value); i++ {
		v, ok := value[i].(string)
		if !ok {
			return nil, errors.New(errorTypeTurnFailed)
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
				zap.L().Warn(warnKey)
				continue
			}
			if err != nil && err != redis.TxFailedErr {
				zap.L().Error(errorPushVString, zap.Error(err))
				return err
			}
			break
		}
		return nil
	}, keyV1, keySet); err != nil {
		zap.L().Error(errorWatchKey, zap.Error(err))
		return err
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

func DeleteVFollowUserInfo(userID, loginUserID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	keyVFollowInfoPart := KeyVFollowInfo + loginUserID1
	keyVFollowInfo := getKey(keyVFollowInfoPart)
	return rdb.HDel(context.Background(), keyVFollowInfo, userID1).Err()
}

func DeleteVFollowerUserInfo(userID, loginUserID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	keyVFollowerInfoPart := KeyVFollowerInfo + userID1
	keyVFollowerInfo := getKey(keyVFollowerInfoPart)
	return rdb.HDel(context.Background(), keyVFollowerInfo, loginUserID1).Err()
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
				zap.L().Warn(warnKey)
				continue
			}
			if err != nil && err != redis.TxFailedErr {
				zap.L().Error(errorPushVFollowFollowerInfoInit, zap.Error(err))
				return err
			}
			break
		}
		return nil
	}, keyFollow, keyFollower); err != nil {
		zap.L().Error(errorWatchKey, zap.Error(err))
		return err
	}
	return nil
}
