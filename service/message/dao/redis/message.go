package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

const (
	maxInt64     = "9223372036854775807"
	temporarySet = "temporary"
)

type ProducerMessage1 struct {
	MessageId  int64  `json:"messageId"`
	UserId     int64  `json:"userId"`
	ToUserId   int64  `json:"toUserId"`
	Content    string `json:"content"`
	CreateTime int64  `json:"createTime"`
}

func IsMessageActiveUser(loginUserId int64) (bool, error) {
	key := getKey(KeyMessageActiveUser)
	return rdb.SIsMember(context.Background(), key, loginUserId).Result()
}

func CreateMessage(producer1Message ProducerMessage1) error {
	userId1 := strconv.FormatInt(producer1Message.UserId, 10)
	toUserId1 := strconv.FormatInt(producer1Message.ToUserId, 10)
	key1 := KeyMessageActiveUserMessage + userId1 + ":" + toUserId1
	key := getKey(key1)
	data, _ := json.Marshal(producer1Message)
	return rdb.ZAdd(context.Background(), key, &redis.Z{
		Score:  float64(producer1Message.CreateTime),
		Member: data,
	}).Err()
}

func GetMessage(userId, toUserId, latestTime int64, isFirst bool) ([]string, error) {
	userId1 := strconv.FormatInt(userId, 10)
	toUserId1 := strconv.FormatInt(toUserId, 10)
	key1 := KeyMessageActiveUserMessage + userId1 + ":" + toUserId1
	key := getKey(key1)
	if isFirst == true {
		return rdb.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
			Max: fmt.Sprintf("%v", time.Now().UnixMilli()),
		}).Result()
	}
	return rdb.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
		Min: fmt.Sprintf("%v", latestTime),
		//Max: fmt.Sprintf("%v", time.Now().UnixMilli()),
		Max: maxInt64,
	}).Result()
}

func AddUserGetMessage(userId int64) error {
	userId1 := strconv.FormatInt(userId, 10)
	key1 := KeyUserGetMessageCount + userId1
	key := getKey(key1)
	pipe := rdb.Pipeline()
	intCmd := pipe.Incr(context.Background(), key)
	pipe.Expire(context.Background(), key, time.Hour*72)
	_, err := pipe.Exec(context.Background())
	count := intCmd.Val()
	if count == 20 {
		keyGet := getKey(KeyGetCountSet)
		if err = rdb.SAdd(context.Background(), keyGet, userId1).Err(); err != nil {
			return err
		}
	}
	return err
}

func AddUserSendMessage(userId int64) error {
	userId1 := strconv.FormatInt(userId, 10)
	key1 := KeyUserSendMessageCount + userId1
	key := getKey(key1)
	pipe := rdb.Pipeline()
	intCmd := pipe.Incr(context.Background(), key)
	pipe.Expire(context.Background(), key, time.Hour*72)
	_, err := pipe.Exec(context.Background())
	count := intCmd.Val()
	if count == 10 {
		keySend := getKey(KeySendCountSet)
		if err = rdb.SAdd(context.Background(), keySend, userId1).Err(); err != nil {
			return err
		}
	}
	return err
}

func PushActiveSet() ([]string, error) {
	keyGet := getKey(KeyGetCountSet)
	keySend := getKey(KeySendCountSet)
	keySet := getKey(KeyMessageActiveUser)
	var keys = []string{
		keyGet,
		keySend,
	}
	_, err := rdb.SInterStore(context.Background(), temporarySet, keys...).Result()
	if err != nil {
		return nil, err
	}
	result, err := rdb.SDiff(context.Background(), temporarySet, keySet).Result()
	if err != nil {
		return nil, err
	}
	err = rdb.Unlink(context.Background(), temporarySet).Err()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func PushActiveSet1(userId []int64) error {
	keySet := getKey(KeyMessageActiveUser)
	var userIds []string
	for i := 0; i < len(userId); i++ {
		userId1 := strconv.FormatInt(userId[i], 10)
		userIds = append(userIds, userId1)
	}
	return rdb.SAdd(context.Background(), keySet, userIds).Err()
}
