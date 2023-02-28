package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
)

func AddUserFollowFollowerCount(loginUserId, userId int64) (int64, int64, error) {
	loginUserId1 := strconv.FormatInt(loginUserId, 10)
	userID1 := strconv.FormatInt(userId, 10)
	keyPart1 := KeyUserFollowCount + loginUserId1
	key1 := getKey(keyPart1)
	keyPart2 := KeyUserFollowerCount + userID1
	key2 := getKey(keyPart2)
	pipe := rdb.Pipeline()
	intCmd1 := pipe.Incr(context.Background(), key1)
	intCmd2 := pipe.Incr(context.Background(), key2)
	_, err := pipe.Exec(context.Background())
	if err != nil {
		return 0, 0, err
	}
	return intCmd1.Val(), intCmd2.Val(), nil
}
func SubUserFollowFollowerCount(loginUserId, userId int64) error {
	loginUserId1 := strconv.FormatInt(loginUserId, 10)
	userID1 := strconv.FormatInt(userId, 10)
	keyPart1 := KeyUserFollowCount + loginUserId1
	key1 := getKey(keyPart1)
	keyPart2 := KeyUserFollowerCount + userID1
	key2 := getKey(keyPart2)
	pipe := rdb.Pipeline()
	pipe.Decr(context.Background(), key1)
	pipe.Decr(context.Background(), key2)
	_, err := pipe.Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func GetUserFollowFollowerCount(userId int64) (int64, int64, error) {
	userID1 := strconv.FormatInt(userId, 10)
	keyPart1 := KeyUserFollowCount + userID1
	key1 := getKey(keyPart1)
	keyPart2 := KeyUserFollowerCount + userID1
	key2 := getKey(keyPart2)
	pipe := rdb.Pipeline()
	cmd1 := pipe.Get(context.Background(), key1)
	cmd2 := pipe.Get(context.Background(), key2)
	_, err := pipe.Exec(context.Background())
	if err != nil && err != redis.Nil {
		return 0, 0, err
	}
	count1 := 0
	count2 := 0
	if cmd1.Val() != "" {
		count1, _ = strconv.Atoi(cmd1.Val())
	}
	if cmd2.Val() != "" {
		count2, _ = strconv.Atoi(cmd2.Val())
	}
	return int64(count1), int64(count2), nil
}

func GetUserFollowFollowerCountInner(userId, toUserId int64) (int64, int64, error) {
	userID1 := strconv.FormatInt(userId, 10)
	toUserId1 := strconv.FormatInt(toUserId, 10)
	keyPart1 := KeyUserFollowCount + userID1
	key1 := getKey(keyPart1)
	keyPart2 := KeyUserFollowerCount + toUserId1
	key2 := getKey(keyPart2)
	pipe := rdb.Pipeline()
	cmd1 := pipe.Get(context.Background(), key1)
	cmd2 := pipe.Get(context.Background(), key2)
	_, err := pipe.Exec(context.Background())
	if err != nil && err != redis.Nil {
		return 0, 0, err
	}
	count1 := 0
	count2 := 0
	if cmd1.Val() != "" {
		count1, _ = strconv.Atoi(cmd1.Val())
	}
	if cmd2.Val() != "" {
		count2, _ = strconv.Atoi(cmd2.Val())
	}
	return int64(count1), int64(count2), nil
}

func GetUserFollowFollowerCountList(userId []int64) ([]int64, []int64, error) {
	var keys1 []string
	var keys2 []string
	for i := 0; i < len(userId); i++ {
		userID1 := strconv.FormatInt(userId[i], 10)
		keyPart1 := KeyUserFollowCount + userID1
		key1 := getKey(keyPart1)
		keys1 = append(keys1, key1)
		keyPart2 := KeyUserFollowerCount + userID1
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
				return nil, nil, errors.New(errorTypeTurnFailed)
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
				return nil, nil, errors.New(errorTypeTurnFailed)
			}
			result2 = append(result2, r21)
		}
	}
	return result1, result2, nil
}
