package redis

import (
	"context"
	"errors"
	"strconv"
)

func AddVideoCommentCount(videoId int64) error {
	videoId1 := strconv.FormatInt(videoId, 10)
	keyPart := KeyVideoCommentCount + videoId1
	key := getKey(keyPart)
	if err := rdb.Incr(context.Background(), key).Err(); err != nil {
		return err
	}
	return nil
}

func SubVideoCommentCount(videoId int64) error {
	videoId1 := strconv.FormatInt(videoId, 10)
	keyPart := KeyVideoCommentCount + videoId1
	key := getKey(keyPart)
	if err := rdb.Decr(context.Background(), key).Err(); err != nil {
		return err
	}
	return nil
}

func GetCommentCount(videoIdList []int64) ([]int64, error) {
	var keys []string
	for i := 0; i < len(videoIdList); i++ {
		videoId1 := strconv.FormatInt(videoIdList[i], 10)
		keyPart := KeyVideoCommentCount + videoId1
		key := getKey(keyPart)
		keys = append(keys, key)
	}
	result, err := rdb.MGet(context.Background(), keys...).Result()
	if err != nil {
		return nil, err
	}
	var counts []int64
	for i := 0; i < len(result); i++ {
		if result[i] == nil {
			counts = append(counts, 0)
		} else {
			count, ok := result[i].(string)
			count1, _ := strconv.ParseInt(count, 10, 64)
			if !ok {
				return nil, errors.New(errorTypeTurnFailed)
			}
			counts = append(counts, count1)
		}
	}
	return counts, nil
}
