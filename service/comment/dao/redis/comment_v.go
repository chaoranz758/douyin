package redis

import (
	"context"
	"douyin/service/comment/model"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

func PushCommentInfo(comment model.Commit, authorId int64) error {
	authorId1 := strconv.FormatInt(authorId, 10)
	videoId1 := strconv.FormatInt(comment.VideoID, 10)
	commentId1 := strconv.FormatInt(comment.CommitID, 10)
	keyPart := fmt.Sprintf("%s:%s", authorId1, videoId1)
	keyPart1 := KeyVVideoComment + keyPart
	key := getKey(keyPart1)
	data, err := json.Marshal(comment)
	if err != nil {
		zap.L().Error(errorJsonMarshal, zap.Error(err))
		return err
	}
	if err := rdb.HSet(context.Background(), key, commentId1, data).Err(); err != nil {
		return err
	}
	return nil
}

func DeleteCommentInfo(videoId, authorId, commentId int64) error {
	authorId1 := strconv.FormatInt(authorId, 10)
	videoId1 := strconv.FormatInt(videoId, 10)
	commentId1 := strconv.FormatInt(commentId, 10)
	keyPart := fmt.Sprintf("%s:%s", authorId1, videoId1)
	keyPart1 := KeyVVideoComment + keyPart
	key := getKey(keyPart1)
	if err := rdb.HDel(context.Background(), key, commentId1).Err(); err != nil {
		return err
	}
	return nil
}

func GetCommentInfo(videoId, authorId int64) ([]string, error) {
	authorId1 := strconv.FormatInt(authorId, 10)
	videoId1 := strconv.FormatInt(videoId, 10)
	keyPart := fmt.Sprintf("%s:%s", authorId1, videoId1)
	keyPart1 := KeyVVideoComment + keyPart
	key := getKey(keyPart1)
	return rdb.HVals(context.Background(), key).Result()
}

func PushVCommentBasicInfoInit(videoId, authorId int64, cs []model.Commit) error {
	authorId1 := strconv.FormatInt(authorId, 10)
	videoId1 := strconv.FormatInt(videoId, 10)
	keyPart := fmt.Sprintf("%s:%s", authorId1, videoId1)
	keyPart1 := KeyVVideoComment + keyPart
	key := getKey(keyPart1)
	var keys []string
	for i := 0; i < len(cs); i++ {
		data, err := json.Marshal(cs[i])
		if err != nil {
			zap.L().Error(errorJsonMarshal, zap.Error(err))
			return err
		}
		commentId1 := strconv.FormatInt(cs[i].CommitID, 10)
		keys = append(keys, commentId1)
		keys = append(keys, string(data))
	}
	return rdb.HMSet(context.Background(), key, keys).Err()
}
