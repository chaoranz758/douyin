package rocketmq

import (
	"douyin/service/message/initialize/config"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

const (
	producer1  = "messageProducer1"
	retryCount = 5
)

var (
	Producer1 rocketmq.Producer
)

func InitProducer1() error {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{config.Config.Rocketmq.Address}),
		producer.WithGroupName(producer1),
		producer.WithRetry(retryCount),
	)
	if err != nil {
		return err
	}
	err = p.Start()
	if err != nil {
		return err
	}
	Producer1 = p
	return nil
}
