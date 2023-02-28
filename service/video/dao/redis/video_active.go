package redis

import (
	"context"
	"douyin/service/video/model"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"strconv"
)

func PublishActiveUserVideo(v *model.Video) error {
	userID1 := strconv.FormatInt(v.UserID, 10)
	keyPart := KeyActivePublishVideoInfo + userID1
	key := getKey(keyPart)
	data, err := json.Marshal(v)
	if err != nil {
		zap.L().Error(errorJsonMarshal, zap.Error(err))
		return err
	}
	if err := rdb.HSet(context.Background(), key, v.VideoID, data).Err(); err != nil {
		return err
	}
	key1 := getKey(KeyActiveVideoID)
	if err := rdb.HSet(context.Background(), key1, v.VideoID, v.UserID).Err(); err != nil {
		zap.L().Error(errorSetKeyActiveVideoID, zap.Error(err))
		return err
	}
	return nil
}

func DeleteActiveUserVideo(userId, videoId int64) error {
	userID1 := strconv.FormatInt(userId, 10)
	videoId1 := strconv.FormatInt(videoId, 10)
	keyPart := KeyActivePublishVideoInfo + userID1
	key := getKey(keyPart)
	if err := rdb.HDel(context.Background(), key, videoId1).Err(); err != nil {
		return err
	}
	key1 := getKey(KeyActiveVideoID)
	if err := rdb.HDel(context.Background(), key1, videoId1).Err(); err != nil {
		zap.L().Error(errorSetKeyActiveVideoID, zap.Error(err))
		return err
	}
	return nil
}

func GetActiveVideoInfo(userID int64) ([]string, error) {
	userID1 := strconv.FormatInt(userID, 10)
	keyPart := KeyActivePublishVideoInfo + userID1
	key := getKey(keyPart)
	m, err := rdb.HVals(context.Background(), key).Result()
	if err != nil {
		zap.L().Error(errorGetActiveVideoInfo, zap.Error(err))
		return nil, err
	}
	return m, nil
}

func GetActiveVideoInfoByVideoId(userID int64, videoId int64) (string, error) {
	userID1 := strconv.FormatInt(userID, 10)
	videoId1 := strconv.FormatInt(videoId, 10)
	keyPart := KeyActivePublishVideoInfo + userID1
	key := getKey(keyPart)
	m, err := rdb.HGet(context.Background(), key, videoId1).Result()
	if err != nil {
		if err == redis.Nil {
			//zap.L().Error("没有从redis中读到", zap.Error(err))
			return "", err
		}
		zap.L().Error(errorGetActiveVideoInfoByVideoId, zap.Error(err))
		return "", err
	}
	return m, nil
}

func PushActiveFavoriteVideoInfo(data string, userId, videoId int64) error {
	videoId1 := strconv.FormatInt(videoId, 10)
	userId1 := strconv.FormatInt(userId, 10)
	keyPart := KeyActiveFavoriteVideo + userId1
	key := getKey(keyPart)
	return rdb.HSet(context.Background(), key, videoId1, data).Err()
}

func DeleteActiveFavoriteVideoInfo(userId, videoId int64) error {
	videoId1 := strconv.FormatInt(videoId, 10)
	userId1 := strconv.FormatInt(userId, 10)
	keyPart := KeyActiveFavoriteVideo + userId1
	key := getKey(keyPart)
	return rdb.HDel(context.Background(), key, videoId1).Err()
}

func GetActiveFavoriteVideoInfo(userId int64) ([]model.Video, error) {
	userId1 := strconv.FormatInt(userId, 10)
	keyPart := KeyActiveFavoriteVideo + userId1
	key := getKey(keyPart)
	results, err := rdb.HVals(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	var videos []model.Video
	for i := 0; i < len(results); i++ {
		var video model.Video
		if err := json.Unmarshal([]byte(results[i]), &video); err != nil {
			zap.L().Error(errorJsonUnMarshal, zap.Error(err))
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}

func PushActiveBasicInfoInit(vs1, vs2 []model.Video, videoIdList1, videoIdList2 []int64, userId int64) error {
	key1 := getKey(KeyActiveVideoID) //hash video user
	userID1 := strconv.FormatInt(userId, 10)
	keyPart := KeyActivePublishVideoInfo + userID1
	key2 := getKey(keyPart) //hash video data
	keyPart1 := KeyActiveFavoriteVideo + userID1
	key3 := getKey(keyPart1) //hash video data
	var data1 []string
	var data2 []string
	var data3 []string
	for i := 0; i < len(vs1); i++ {
		videoID1 := strconv.FormatInt(videoIdList1[i], 10)
		data, err := json.Marshal(vs1[i])
		if err != nil {
			zap.L().Error(errorJsonMarshal, zap.Error(err))
			return err
		}
		data1 = append(data1, videoID1)
		data1 = append(data1, string(data))
		data3 = append(data3, videoID1)
		data3 = append(data3, userID1)
	}
	for i := 0; i < len(vs2); i++ {
		videoID1 := strconv.FormatInt(videoIdList2[i], 10)
		data, err := json.Marshal(vs2[i])
		if err != nil {
			zap.L().Error(errorJsonMarshal, zap.Error(err))
			return err
		}
		data2 = append(data2, videoID1)
		data2 = append(data2, string(data))
	}
	if err := rdb.Watch(context.Background(), func(tx *redis.Tx) error {
		for {
			_, err := tx.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
				if len(vs1) != 0 {
					pipe.HMSet(context.Background(), key2, data1)
					pipe.HMSet(context.Background(), key1, data3)
				}
				if len(vs2) != 0 {
					pipe.HMSet(context.Background(), key3, data2)
				}
				return nil
			})
			if err == redis.TxFailedErr {
				zap.L().Warn(warnKey)
				continue
			}
			if err != nil && err != redis.TxFailedErr {
				zap.L().Error(errorPushActiveBasicInfoInit, zap.Error(err))
				return err
			}
			break
		}
		return nil
	}, key1, key2, key3); err != nil {
		zap.L().Error(errorWatchKey, zap.Error(err))
		return err
	}
	return nil
}

func DeleteActiveVideoInfo(userId int64, videoId []int64) error {
	key1 := getKey(KeyActiveVideoID) //hash video user
	userID1 := strconv.FormatInt(userId, 10)
	keyPart := KeyActivePublishVideoInfo + userID1
	key2 := getKey(keyPart) //hash video data
	keyPart1 := KeyActiveFavoriteVideo + userID1
	key3 := getKey(keyPart1) //hash video data
	var videoIds []string
	for i := 0; i < len(videoId); i++ {
		v := strconv.FormatInt(videoId[i], 10)
		videoIds = append(videoIds, v)
	}
	if err := rdb.Watch(context.Background(), func(tx *redis.Tx) error {
		for {
			_, err := tx.TxPipelined(context.Background(), func(pipe redis.Pipeliner) error {
				pipe.HDel(context.Background(), key1, videoIds...)
				pipe.Del(context.Background(), key2, key3)
				return nil
			})
			if err == redis.TxFailedErr {
				zap.L().Warn(warnKey)
				continue
			}
			if err != nil && err != redis.TxFailedErr {
				zap.L().Error(errorDeleteActiveVideoInfo, zap.Error(err))
				return err
			}
			break
		}
		return nil
	}, key1, key2, key3); err != nil {
		zap.L().Error(errorWatchKey, zap.Error(err))
		return err
	}
	return nil
}
