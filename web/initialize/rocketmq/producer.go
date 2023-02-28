package rocketmq

import (
	"douyin/web/initialize/config"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

const (
	retryCount    = 5
	groupProducer = "uploadProducer1"
)

var (
	Producer1 rocketmq.Producer
)

func InitProducer1() error {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{config.Config.Rocketmq.Address}),
		producer.WithGroupName(groupProducer),
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
