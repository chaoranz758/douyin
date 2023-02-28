package rocketmq

import (
	"douyin/service/favorite/initialize/config"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

const (
	group1     = "favoriteProducer1"
	group2     = "favoriteProducer2"
	group3     = "favoriteProducer3"
	group4     = "favoriteProducer4"
	retryCount = 5
)

var (
	Producer1 rocketmq.Producer
	Producer2 rocketmq.Producer
	Producer3 rocketmq.Producer
	Producer4 rocketmq.Producer
)

func InitProducer1() error {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{config.Config.Rocketmq.Address}),
		producer.WithGroupName(group1),
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
		producer.WithGroupName(group2),
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
		producer.WithGroupName(group3),
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

func InitProducer4() error {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{config.Config.Rocketmq.Address}),
		producer.WithGroupName(group4),
		producer.WithRetry(retryCount),
	)
	if err != nil {
		return err
	}
	err = p.Start()
	if err != nil {
		return err
	}
	Producer4 = p
	return nil
}
