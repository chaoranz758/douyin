package rocketmq

import (
	"context"
	"douyin/service/message/dao/mysql"
	"douyin/service/message/dao/redis"
	"douyin/service/message/model"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
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

func InitCustomer1() error {
	c, err := rocketmq.NewPushConsumer(
		// 指定 Group 可以实现消费者负载均衡进行消费，并且保证他们的Topic+Tag要一样。
		// 如果同一个 GroupID 下的不同消费者实例，订阅了不同的 Topic+Tag 将导致在对Topic 的消费队列进行负载均衡的时候产生不正确的结果，最终导致消息丢失。(官方源码设计)
		consumer.WithGroupName("messageGroup1"),
		consumer.WithNameServer([]string{"192.168.182.137:9876"}),
	)
	if err != nil {
		return err
	}
	err = c.Subscribe("messageTopic1", consumer.MessageSelector{}, messageCustomer1CallBack)
	if err != nil {
		return err
	}
	err = c.Start()
	if err != nil {
		return err
	}
	return nil
}

func messageCustomer1CallBack(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
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
