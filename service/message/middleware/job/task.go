package job

import (
	"douyin/service/message/dao/mysql"
	"douyin/service/message/dao/redis"
	"douyin/service/message/model"
	"go.uber.org/zap"
	"strconv"
)

const (
	errorPushActiveSet      = "push active set failed"
	errorGetMessageByUserId = "get message by user id failed"
	errorCreateMessage      = "create message failed"
)

const (
	success = "执行定时任务成功"
)

func PushActiveSet() error {
	result, err := redis.PushActiveSet()
	if err != nil {
		zap.L().Error(errorPushActiveSet, zap.Error(err))
		return err
	}
	if len(result) == 0 {
		//zap.L().Info("执行了定时任务但没有要更新的活跃用户")
		return nil
	}
	var userId []int64
	if len(result) != 0 {
		for i := 0; i < len(result); i++ {
			userId1, _ := strconv.ParseInt(result[i], 10, 64)
			ms := make([]model.Message, 0)
			if err = mysql.GetMessageByUserId(&ms, userId1); err != nil {
				zap.L().Error(errorGetMessageByUserId)
				return err
			}
			for j := 0; j < len(ms); j++ {
				var producerMessage = redis.ProducerMessage1{
					MessageId:  ms[j].MessageID,
					UserId:     ms[j].UserID,
					ToUserId:   ms[j].ToUserID,
					Content:    ms[j].Content,
					CreateTime: ms[j].CreatedAt,
				}
				if err = redis.CreateMessage(producerMessage); err != nil {
					zap.L().Error(errorCreateMessage)
					return err
				}
			}
			userId = append(userId, userId1)
		}
		if err = redis.PushActiveSet1(userId); err != nil {
			zap.L().Error(errorPushActiveSet)
			return err
		}
	}
	zap.L().Info(success)
	return nil
}
