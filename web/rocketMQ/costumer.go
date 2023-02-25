package rocketMQ

import (
	"context"
	"douyin/web/util"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
)

const (
	errorUploadFile = "upload file failed"
)

type ProducerMessage1 struct {
	VideoURLLocal string `json:"videoURLLocal"`
	VideoName     string `json:"videoName"`
	ImageName     string `json:"imageName"`
}

func InitCustomer1() error {
	c, err := rocketmq.NewPushConsumer(
		// 指定 Group 可以实现消费者负载均衡进行消费，并且保证他们的Topic+Tag要一样。
		// 如果同一个 GroupID 下的不同消费者实例，订阅了不同的 Topic+Tag 将导致在对Topic 的消费队列进行负载均衡的时候产生不正确的结果，最终导致消息丢失。(官方源码设计)
		consumer.WithGroupName("uploadVideoGroup"),
		consumer.WithNameServer([]string{"192.168.182.137:9876"}),
	)
	if err != nil {
		return err
	}
	err = c.Subscribe("uploadVideoTopic2", consumer.MessageSelector{}, uploadVideoCallBack)
	if err != nil {
		return err
	}
	err = c.Start()
	if err != nil {
		return err
	}
	return nil
}

func uploadVideoCallBack(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	var producer1Message ProducerMessage1
	if err := json.Unmarshal(msgs[0].Body, &producer1Message); err != nil {
		zap.L().Error("json解析失败", zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	err := util.Upload(producer1Message.VideoURLLocal, producer1Message.VideoName, producer1Message.ImageName)
	if err != nil {
		zap.L().Error(errorUploadFile, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	return consumer.ConsumeSuccess, nil
}
