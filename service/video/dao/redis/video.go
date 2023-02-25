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

const (
	errorJsonMarshal                 = "json marshal failed"
	errorAddUserInfoToZSet           = "add user information to zset failed"
	errorGetVVideoInfo               = "get v video information from redis failed"
	errorGetActiveVideoInfo          = "get active video information from redis failed"
	errorComputeZSetCount            = "compute zset count failed"
	errorSetKeyVVideoID              = "set key v video id failed"
	errorSetKeyActiveVideoID         = "set key active video id failed"
	errorGetVVideoInfoByVideoId      = "get v video information by video id failed"
	errorGetActiveVideoInfoByVideoId = "get active video information by video id failed"
	errorJsonUnMarshal               = "json unmarshal failed"
	errorZRevRangeByScore            = "z rev range by score failed"
	errorZCard                       = "zCard failed"
)

func JudgeFeedIsFirst(t string) (int64, error) {
	key := getKey(KeyVideoInfoZSet)
	return rdb.ZCount(context.Background(), key, t, t).Result()
}

func PublishVUserVideo(v *model.Video) error {
	userID1 := strconv.FormatInt(v.UserID, 10)
	keyPart := KeyVPublishVideoInfo + userID1
	key := getKey(keyPart)
	data, err := json.Marshal(v)
	if err != nil {
		zap.L().Error(errorJsonMarshal, zap.Error(err))
		return err
	}
	if err := rdb.HSet(context.Background(), key, v.VideoID, data).Err(); err != nil {
		return err
	}
	key1 := getKey(KeyVVideoID)
	if err := rdb.HSet(context.Background(), key1, v.VideoID, v.UserID).Err(); err != nil {
		zap.L().Error(errorSetKeyVVideoID, zap.Error(err))
		return err
	}
	return nil
}

func DeleteVUserVideo(userId, videoId int64) error {
	userID1 := strconv.FormatInt(userId, 10)
	videoId1 := strconv.FormatInt(videoId, 10)
	keyPart := KeyVPublishVideoInfo + userID1
	key := getKey(keyPart)
	if err := rdb.HDel(context.Background(), key, videoId1).Err(); err != nil {
		return err
	}
	key1 := getKey(KeyVVideoID)
	if err := rdb.HDel(context.Background(), key1, videoId1).Err(); err != nil {
		zap.L().Error(errorSetKeyVVideoID, zap.Error(err))
		return err
	}
	return nil
}

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
	//TODO
	//后续添加watch
	//返回队列长度
	length, err := rdb.ZCard(context.Background(), key).Result()
	if err != nil {
		zap.L().Error(errorComputeZSetCount, zap.Error(err))
		return err
	}
	//return status.New(codes.Aborted, errorComputeZSetCount).Err()
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
		zap.L().Info("系统里一个视频都没有")
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

func GetVVideoInfo(userID int64) ([]string, error) {
	userID1 := strconv.FormatInt(userID, 10)
	keyPart := KeyVPublishVideoInfo + userID1
	key := getKey(keyPart)
	m, err := rdb.HVals(context.Background(), key).Result()
	if err != nil {
		zap.L().Error(errorGetVVideoInfo, zap.Error(err))
		return nil, err
	}
	return m, nil
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

func GetVVideoInfoByVideoId(userID int64, videoId int64) (string, error) {
	userID1 := strconv.FormatInt(userID, 10)
	videoId1 := strconv.FormatInt(videoId, 10)
	keyPart := KeyVPublishVideoInfo + userID1
	key := getKey(keyPart)
	m, err := rdb.HGet(context.Background(), key, videoId1).Result()
	if err != nil {
		if err == redis.Nil {
			zap.L().Error("没有从redis中读到", zap.Error(err))
			return "", err
		}
		zap.L().Error(errorGetVVideoInfoByVideoId, zap.Error(err))
		return "", err
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
			zap.L().Error("没有从redis中读到", zap.Error(err))
			return "", err
		}
		zap.L().Error(errorGetActiveVideoInfoByVideoId, zap.Error(err))
		return "", err
	}
	return m, nil
}

func JudgeVideoAuthor(videoId int64) (string, bool, bool, error) {
	videoId1 := strconv.FormatInt(videoId, 10)
	key1 := getKey(KeyVVideoID)
	key2 := getKey(KeyActiveVideoID)
	authorId, err := rdb.HGet(context.Background(), key1, videoId1).Result()
	if err != nil {
		if err == redis.Nil {
			zap.L().Info("视频不属于大V用户")
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
			zap.L().Info("视频不属于活跃用户")
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

func PushVFavoriteVideoInfo(data string, userId, videoId int64) error {
	videoId1 := strconv.FormatInt(videoId, 10)
	userId1 := strconv.FormatInt(userId, 10)
	keyPart := KeyVFavoriteVideo + userId1
	key := getKey(keyPart)
	return rdb.HSet(context.Background(), key, videoId1, data).Err()
}

func PushActiveFavoriteVideoInfo(data string, userId, videoId int64) error {
	videoId1 := strconv.FormatInt(videoId, 10)
	userId1 := strconv.FormatInt(userId, 10)
	keyPart := KeyActiveFavoriteVideo + userId1
	key := getKey(keyPart)
	return rdb.HSet(context.Background(), key, videoId1, data).Err()
}

func DeleteVFavoriteVideoInfo(userId, videoId int64) error {
	videoId1 := strconv.FormatInt(videoId, 10)
	userId1 := strconv.FormatInt(userId, 10)
	keyPart := KeyVFavoriteVideo + userId1
	key := getKey(keyPart)
	return rdb.HDel(context.Background(), key, videoId1).Err()
}

func DeleteActiveFavoriteVideoInfo(userId, videoId int64) error {
	videoId1 := strconv.FormatInt(videoId, 10)
	userId1 := strconv.FormatInt(userId, 10)
	keyPart := KeyActiveFavoriteVideo + userId1
	key := getKey(keyPart)
	return rdb.HDel(context.Background(), key, videoId1).Err()
}

func GetVFavoriteVideoInfo(userId int64) ([]model.Video, error) {
	userId1 := strconv.FormatInt(userId, 10)
	keyPart := KeyVFavoriteVideo + userId1
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

func PushVBasicInfoInit(vs1, vs2 []model.Video, videoIdList1, videoIdList2 []int64, userId int64) error {
	key1 := getKey(KeyVVideoID) //hash video user
	userID1 := strconv.FormatInt(userId, 10)
	keyPart := KeyVPublishVideoInfo + userID1
	key2 := getKey(keyPart) //hash video data
	keyPart1 := KeyVFavoriteVideo + userID1
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
				zap.L().Warn("redis监听的key正在被其它线程执行")
				continue
			}
			if err != nil && err != redis.TxFailedErr {
				zap.L().Error("exec push v basic information init failed", zap.Error(err))
				return err
			}
			break
		}
		return nil
	}, key1, key2, key3); err != nil {
		zap.L().Error("redis watch key failed", zap.Error(err))
		return err
	}
	return nil
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
				zap.L().Warn("redis监听的key正在被其它线程执行")
				continue
			}
			if err != nil && err != redis.TxFailedErr {
				zap.L().Error("exec push active basic information failed", zap.Error(err))
				return err
			}
			break
		}
		return nil
	}, key1, key2, key3); err != nil {
		zap.L().Error("redis watch key failed", zap.Error(err))
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
				zap.L().Warn("redis监听的key正在被其它线程执行")
				continue
			}
			if err != nil && err != redis.TxFailedErr {
				zap.L().Error("exec delete active video information failed", zap.Error(err))
				return err
			}
			break
		}
		return nil
	}, key1, key2, key3); err != nil {
		zap.L().Error("redis watch key failed", zap.Error(err))
		return err
	}
	return nil
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
				return nil, errors.New("类型转换失败")
			}
			result1 = append(result1, r11)
		}
	}
	return result1, nil
}
