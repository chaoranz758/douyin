package rocketmq

import (
	"douyin/service/message/initialize/config"
	rocketmq1 "douyin/service/message/middleware/rocketmq"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
)

const (
	costumer1 = "messageGroup1"
	topic1    = "messageTopic1"
)

func InitCustomer1() error {
	c, err := rocketmq.NewPushConsumer(
		// 指定 Group 可以实现消费者负载均衡进行消费，并且保证他们的Topic+Tag要一样。
		// 如果同一个 GroupID 下的不同消费者实例，订阅了不同的 Topic+Tag 将导致在对Topic 的消费队列进行负载均衡的时候产生不正确的结果，最终导致消息丢失。(官方源码设计)
		consumer.WithGroupName(costumer1),
		consumer.WithNameServer([]string{config.Config.Rocketmq.Address}),
	)
	if err != nil {
		return err
	}
	err = c.Subscribe(topic1, consumer.MessageSelector{}, rocketmq1.MessageCustomer1CallBack)
	if err != nil {
		return err
	}
	err = c.Start()
	if err != nil {
		return err
	}
	return nil
}
