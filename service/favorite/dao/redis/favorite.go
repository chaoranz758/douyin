package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

const (
	errorAddFavoriteCount  = "add favorite count failed"
	errorSubFavoriteCount  = "sub favorite count failed"
	MGetVideoFavoriteCount = "get video favorite count failed"
)

func AddFavoriteCount(userID, videoID, authorId int64) (int64, error) {
	userID1 := strconv.FormatInt(userID, 10)
	videoID1 := strconv.FormatInt(videoID, 10)
	authorId1 := strconv.FormatInt(authorId, 10)
	keyPart := KeyUserFavoriteVideoCount + userID1
	key := getKey(keyPart)
	keyPart1 := KeyVideoFavoriteCount + videoID1
	key1 := getKey(keyPart1)
	keyPart2 := KeyUserFavoritedTotalCount + authorId1
	key2 := getKey(keyPart2)
	pipe := rdb.Pipeline()
	result := pipe.Incr(context.Background(), key)
	pipe.Incr(context.Background(), key1)
	pipe.Incr(context.Background(), key2)
	_, err := pipe.Exec(context.Background())
	if err != nil {
		zap.L().Error(errorAddFavoriteCount, zap.Error(err))
		return 0, err
	}
	return result.Val(), nil
}

func GetFavoriteCount(userId int64) (int64, error) {
	userID1 := strconv.FormatInt(userId, 10)
	keyPart := KeyUserFavoriteVideoCount + userID1
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

func SubFavoriteCount(userID, videoID, authorId int64) error {
	userID1 := strconv.FormatInt(userID, 10)
	videoID1 := strconv.FormatInt(videoID, 10)
	authorId1 := strconv.FormatInt(authorId, 10)
	keyPart := KeyUserFavoriteVideoCount + userID1
	key := getKey(keyPart)
	keyPart1 := KeyVideoFavoriteCount + videoID1
	key1 := getKey(keyPart1)
	keyPart2 := KeyUserFavoritedTotalCount + authorId1
	key2 := getKey(keyPart2)
	pipe := rdb.Pipeline()
	pipe.Decr(context.Background(), key)
	pipe.Decr(context.Background(), key1)
	pipe.Decr(context.Background(), key2)
	_, err := pipe.Exec(context.Background())
	if err != nil {
		zap.L().Error(errorSubFavoriteCount, zap.Error(err))
		return err
	}
	return nil
}

func GetVideoFavoriteCount(idList []int64) ([]int64, error) {
	var key []string
	for i := 0; i < len(idList); i++ {
		videoID := strconv.FormatInt(idList[i], 10)
		keyPart1 := KeyVideoFavoriteCount + videoID
		key1 := getKey(keyPart1)
		key = append(key, key1)
	}
	result, err := rdb.MGet(context.Background(), key...).Result()
	if err != nil {
		zap.L().Error(MGetVideoFavoriteCount, zap.Error(err))
		return nil, err
	}
	//if result[0] == nil && len(result) == 1 {
	//	data := make([]int64, len(idList), len(idList))
	//	return data, nil
	//}
	var result1 []int64
	for i := 0; i < len(result); i++ {
		if result[i] == nil {
			result1 = append(result1, 0)
		} else {
			x, ok := result[i].(string)
			x1, _ := strconv.ParseInt(x, 10, 64)
			if !ok {
				return nil, errors.New("类型转换失败")
			}
			result1 = append(result1, x1)
		}
	}
	return result1, nil
}

func GetUserFavoritedCount(userId []int64) ([]int64, []int64, error) {
	var keys1 []string
	var keys2 []string
	for i := 0; i < len(userId); i++ {
		userID1 := strconv.FormatInt(userId[i], 10)
		keyPart1 := KeyUserFavoriteVideoCount + userID1
		key1 := getKey(keyPart1)
		keys1 = append(keys1, key1)
		keyPart2 := KeyUserFavoritedTotalCount + userID1
		key2 := getKey(keyPart2)
		keys2 = append(keys2, key2)
	}
	pipe := rdb.Pipeline()
	cmd1 := pipe.MGet(context.Background(), keys1...)
	cmd2 := pipe.MGet(context.Background(), keys2...)
	_, err := pipe.Exec(context.Background())
	if err != nil && err != redis.Nil {
		return nil, nil, err
	}
	var result1 []int64
	var result2 []int64
	for i := 0; i < len(cmd1.Val()); i++ {
		if cmd1.Val()[i] == nil {
			result1 = append(result1, 0)
		} else {
			r1, ok := cmd1.Val()[i].(string)
			r11, _ := strconv.ParseInt(r1, 10, 64)
			if !ok {
				return nil, nil, errors.New("类型转换失败")
			}
			result1 = append(result1, r11)
		}
	}
	for i := 0; i < len(cmd2.Val()); i++ {
		if cmd2.Val()[i] == nil {
			result2 = append(result2, 0)
		} else {
			r2, ok := cmd2.Val()[i].(string)
			r21, _ := strconv.ParseInt(r2, 10, 64)
			if !ok {
				return nil, nil, errors.New("类型转换失败")
			}
			result2 = append(result2, r21)
		}
	}
	return result1, result2, nil
}
