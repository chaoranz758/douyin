package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

const (
	userLoginCount = 5
)

func UserLoginCount(userID int64) error {
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
	if count == userLoginCount {
		keySet := getKey(KeyUserLoginSet)
		_, err := rdb.SAdd(context.Background(), keySet, userID1).Result()
		if err != nil {
			zap.L().Error(errorAddUserLoginSet, zap.Error(err))
			return err
		}
	}
	return nil
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

func DeletedUserVideoCountSet(userID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	keySet := getKey(KeyUserVideoSet)
	if err := rdb.SRem(context.Background(), keySet, userID1).Err(); err != nil {
		zap.L().Error(errorDeleteUserVideoSet)
		return status.New(codes.Aborted, errorDeleteUserVideoSet).Err()
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

func DeleteUserFavoriteVideoCountSet(userID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	keySet := getKey(KeyUserFavoriteVideoSet)
	if err := rdb.SRem(context.Background(), keySet, userID1).Err(); err != nil {
		zap.L().Error(errorDeleteUserFavoriteVideoSet)
		return status.New(codes.Aborted, errorDeleteUserFavoriteVideoSet).Err()
	}
	return nil
}

func AddUserFollowUserCountSet(userID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	keySet := getKey(KeyUserFollowCountSet)
	_, err := rdb.SAdd(context.Background(), keySet, userID1).Result()
	if err != nil {
		zap.L().Error(errorAddUserFollowUserCountSet, zap.Error(err))
		return status.Error(codes.Aborted, err.Error())
	}
	return nil
}

func DeleteUserFollowUserCountSet(userID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	keySet := getKey(KeyUserFollowCountSet)
	if err := rdb.SRem(context.Background(), keySet, userID1).Err(); err != nil {
		zap.L().Error(errorDeleteUserFollowUserCountSet)
		return status.New(codes.Aborted, errorDeleteUserFollowUserCountSet).Err()
	}
	return nil
}

func AddUserFollowerUserCountSet(userID int64) error {
	println(userID)
	userID1 := strconv.FormatInt(userID, 10)
	keySet := getKey(KeyUserFollowerCountSet)
	_, err := rdb.SAdd(context.Background(), keySet, userID1).Result()
	if err != nil {
		zap.L().Error(errorAddUserFollowerUserCountSet, zap.Error(err))
		return status.Error(codes.Aborted, err.Error())
	}
	return nil
}

func DeleteUserFollowerUserCountSet(userID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	keySet := getKey(KeyUserFollowerCountSet)
	if err := rdb.SRem(context.Background(), keySet, userID1).Err(); err != nil {
		zap.L().Error(errorDeleteUserFollowerUserCountSet)
		return status.New(codes.Aborted, errorDeleteUserFollowerUserCountSet).Err()
	}
	return nil
}
