package rocketmq

import (
	"context"
	"douyin/service/message/dao/mysql"
	"douyin/service/message/dao/redis"
	"douyin/service/message/model"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
)

const (
	errorCreateMessage      = "create message failed"
	errorAddUserSendMessage = "add user send message count failed"
)

type ProducerMessage1 struct {
	MessageId  int64  `json:"messageId"`
	UserId     int64  `json:"userId"`
	ToUserId   int64  `json:"toUserId"`
	Content    string `json:"content"`
	CreateTime int64  `json:"createTime"`
}

func MessageCustomer1CallBack(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	var producer1Message redis.ProducerMessage1
	if err := json.Unmarshal(msgs[0].Body, &producer1Message); err != nil {
		zap.L().Error("json解析失败", zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	if err := redis.CreateMessage(producer1Message); err != nil {
		zap.L().Error(errorCreateMessage, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	//存入mysql
	var m = model.Message{
		MessageID: producer1Message.MessageId,
		UserID:    producer1Message.UserId,
		ToUserID:  producer1Message.ToUserId,
		Content:   producer1Message.Content,
		CreatedAt: producer1Message.CreateTime,
	}
	if err := mysql.CreateMessage(&m); err != nil {
		zap.L().Error(errorCreateMessage, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	if err := redis.AddUserSendMessage(producer1Message.UserId); err != nil {
		zap.L().Error(errorAddUserSendMessage, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	zap.L().Info("消费者1执行成功")
	return consumer.ConsumeSuccess, nil
}
