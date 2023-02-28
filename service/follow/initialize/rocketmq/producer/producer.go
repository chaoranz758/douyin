package producer

import (
	"douyin/service/follow/initialize/config"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

const (
	producer1  = "followProducer1"
	producer2  = "followProducer2"
	producer3  = "followProducer3"
	retryCount = 5
)

var (
	Producer1 rocketmq.Producer
	Producer2 rocketmq.Producer
	Producer3 rocketmq.Producer
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

func InitProducer2() error {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{config.Config.Rocketmq.Address}),
		producer.WithGroupName(producer2),
		producer.WithRetry(retryCount),
	)
	if err != nil {
		return err
	}
	err = p.Start()
	if err != nil {
		return err
	}
	Producer2 = p
	return nil
}

func InitProducer3() error {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{config.Config.Rocketmq.Address}),
		producer.WithGroupName(producer3),
		producer.WithRetry(retryCount),
	)
	if err != nil {
		return err
	}
	err = p.Start()
	if err != nil {
		return err
	}
	Producer3 = p
	return nil
}
