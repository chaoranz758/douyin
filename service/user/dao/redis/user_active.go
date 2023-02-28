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

const (
	temporarySet = "temporary"
)

func JudgeIsInActiveSet(userId int64) (bool, error) {
	keyActive := getKey(KeyActiveSet)
	return rdb.SIsMember(context.Background(), keyActive, userId).Result()
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
	_, err := rdb.SInterStore(context.Background(), temporarySet, keys...).Result()
	if err != nil {
		return nil, err
	}
	keyActive := getKey(KeyActiveSet)
	result, err := rdb.SDiff(context.Background(), temporarySet, keyActive).Result()
	if err != nil {
		return nil, err
	}
	err = rdb.Unlink(context.Background(), temporarySet).Err()
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
				zap.L().Warn(warnKey)
				continue
			}
			if err != nil && err != redis.TxFailedErr {
				zap.L().Error(errorPushActiveSetAndString, zap.Error(err))
				return err
			}
			break
		}
		return nil
	}, keys...); err != nil {
		zap.L().Error(errorWatchKey, zap.Error(err))
		return err
	}
	return nil
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

func DeleteActiveFollowUserInfo(userID, loginUserID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	keyActiveFollowInfoPart := KeyActiveFollowInfo + loginUserID1
	keyActiveFollowInfo := getKey(keyActiveFollowInfoPart)
	return rdb.HDel(context.Background(), keyActiveFollowInfo, userID1).Err()
}

func DeleteActiveFollowerUserInfo(userID, loginUserID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	loginUserID1 := strconv.FormatInt(loginUserID, 10)
	keyActiveFollowerInfoPart := KeyActiveFollowerInfo + userID1
	keyActiveFollowerInfo := getKey(keyActiveFollowerInfoPart)
	return rdb.HDel(context.Background(), keyActiveFollowerInfo, loginUserID1).Err()
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
				zap.L().Warn(warnKey)
				continue
			}
			if err != nil && err != redis.TxFailedErr {
				zap.L().Error(errorPushActiveFollowInfoInit, zap.Error(err))
				return err
			}
			break
		}
		return nil
	}, keyFollow); err != nil {
		zap.L().Error(errorWatchKey, zap.Error(err))
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
				zap.L().Warn(warnKey)
				continue
			}
			if err != nil && err != redis.TxFailedErr {
				zap.L().Error(errorPushActiveFollowerInfoInit, zap.Error(err))
				return err
			}
			break
		}
		return nil
	}, keyFollow); err != nil {
		zap.L().Error(errorWatchKey, zap.Error(err))
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
				zap.L().Warn(warnKey)
				continue
			}
			if err != nil && err != redis.TxFailedErr {
				zap.L().Error(errorDeleteActiveAllInfo, zap.Error(err))
				return err
			}
			break
		}
		return nil
	}, keyActive, keyUserInfo, keyFollowInfo, keyFollowerInfo); err != nil {
		zap.L().Error(errorWatchKey, zap.Error(err))
		return err
	}
	return nil
}
