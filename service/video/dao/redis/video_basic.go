package redis

import (
	"context"
	"douyin/service/video/model"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

func JudgeFeedIsFirst(t string) (int64, error) {
	key := getKey(KeyVideoInfoZSet)
	return rdb.ZCount(context.Background(), key, t, t).Result()
}

func PublishVideoInfoZSet(v *model.Video, times int64) error {
	key := getKey(KeyVideoInfoZSet)
	data, err := json.Marshal(v)
	if err != nil {
		zap.L().Error(errorJsonMarshal, zap.Error(err))
	}
	_, err = rdb.ZAdd(context.Background(), key, &redis.Z{
		Score:  float64(times),
		Member: string(data),
	}).Result()
	if err != nil {
		zap.L().Error(errorAddUserInfoToZSet, zap.Error(err))
		return err
	}
	//返回队列长度
	length, err := rdb.ZCard(context.Background(), key).Result()
	if err != nil {
		zap.L().Error(errorComputeZSetCount, zap.Error(err))
		return err
	}
	//若队列长度大于300,则裁剪
	if length > 300 {
		rdb.ZRemRangeByRank(context.Background(), key, 0, length-301)
	}
	return nil
}

func GetVideoInfoZSetInitial() ([]string, error) {
	key := getKey(KeyVideoInfoZSet)
	time1 := strconv.FormatInt(time.Now().UnixMilli(), 10)
	results, err := rdb.ZRevRangeByScore(context.Background(), key, &redis.ZRangeBy{
		Max:   time1,
		Count: 30,
	}).Result()
	if err != nil {
		zap.L().Error(errorZRevRangeByScore, zap.Error(err))
		return nil, err
	}
	return results, nil
}

func GetVideoInfoZSet(time string) (int64, []string, error) {
	key := getKey(KeyVideoInfoZSet)
	count, err := rdb.ZCard(context.Background(), key).Result()
	if err != nil {
		zap.L().Error(errorZCard, zap.Error(err))
		return 0, nil, err
	}
	if count == 0 {
		return 0, nil, nil
	}
	results, err := rdb.ZRevRangeByScore(context.Background(), key, &redis.ZRangeBy{
		Max:   time,
		Count: 31,
	}).Result()

	if err != nil {
		zap.L().Error(errorZRevRangeByScore, zap.Error(err))
		return 0, nil, err
	}
	return count, results, nil
}

func IncrUserVideoCount(userID int64) (int64, error) {
	userID1 := strconv.FormatInt(userID, 10)
	keyPart := KeyUserVideoCount + userID1
	key := getKey(keyPart)
	return rdb.Incr(context.Background(), key).Result()
}

func SubUserVideoCount(userID int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	keyPart := KeyUserVideoCount + userID1
	key := getKey(keyPart)
	return rdb.Decr(context.Background(), key).Err()
}

func GetUserVideoCount(userID int64) (int64, error) {
	userID1 := strconv.FormatInt(userID, 10)
	keyPart := KeyUserVideoCount + userID1
	key := getKey(keyPart)
	countString, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, status.Error(codes.Aborted, err.Error())
	}
	count, err := strconv.ParseInt(countString, 10, 64)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func JudgeVideoAuthor(videoId int64) (string, bool, bool, error) {
	videoId1 := strconv.FormatInt(videoId, 10)
	key1 := getKey(KeyVVideoID)
	key2 := getKey(KeyActiveVideoID)
	authorId, err := rdb.HGet(context.Background(), key1, videoId1).Result()
	if err != nil {
		if err == redis.Nil {
			//zap.L().Info("视频不属于大V用户")
		} else {
			return "", false, false, err
		}
	}
	if err == nil {
		return authorId, true, false, nil
	}
	authorId, err = rdb.HGet(context.Background(), key2, videoId1).Result()
	if err != nil {
		if err == redis.Nil {
			//zap.L().Info("视频不属于活跃用户")
			return "", false, false, nil
		} else {
			return "", false, false, err
		}
	}
	if err == nil {
		return authorId, false, true, nil
	}
	return "", false, false, nil
}

func GetUserPublishVideoCount(userId int64) (int64, error) {
	userID1 := strconv.FormatInt(userId, 10)
	keyPart1 := KeyUserVideoCount + userID1
	key1 := getKey(keyPart1)
	count, err := rdb.Get(context.Background(), key1).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	count1 := 0
	count1, _ = strconv.Atoi(count)
	return int64(count1), nil
}

func GetUserPublishVideoCountList(userId []int64) ([]int64, error) {
	var keys1 []string
	for i := 0; i < len(userId); i++ {
		userID1 := strconv.FormatInt(userId[i], 10)
		keyPart1 := KeyUserVideoCount + userID1
		key1 := getKey(keyPart1)
		keys1 = append(keys1, key1)
	}
	count, err := rdb.MGet(context.Background(), keys1...).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	var result1 []int64
	for i := 0; i < len(count); i++ {
		if count[i] == nil {
			result1 = append(result1, 0)
		} else {
			r1, ok := count[i].(string)
			r11, _ := strconv.ParseInt(r1, 10, 64)
			if !ok {
				return nil, errors.New(errorTypeTurnFailed)
			}
			result1 = append(result1, r11)
		}
	}
	return result1, nil
}
