package rocketmq

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var (
	Producer1 rocketmq.Producer
	Producer2 rocketmq.Producer
	Producer3 rocketmq.Producer
)

func InitProducer1() error {
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{"192.168.182.137:9876"}),
		producer.WithGroupName("followProducer1"),
		producer.WithRetry(5),
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
		producer.WithNameServer([]string{"192.168.182.137:9876"}),
		producer.WithGroupName("followProducer2"),
		producer.WithRetry(5),
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
		producer.WithNameServer([]string{"192.168.182.137:9876"}),
		producer.WithGroupName("followProducer3"),
		producer.WithRetry(5),
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
