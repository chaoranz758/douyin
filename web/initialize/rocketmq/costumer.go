package rocketmq

import (
	"douyin/web/initialize/config"
	"douyin/web/middleware"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
)

const (
	group = "uploadVideoGroup"
	topic = "uploadVideoTopic2"
)

func InitCustomer1() error {
	c, err := rocketmq.NewPushConsumer(
		// 指定 Group 可以实现消费者负载均衡进行消费，并且保证他们的Topic+Tag要一样。
		// 如果同一个 GroupID 下的不同消费者实例，订阅了不同的 Topic+Tag 将导致在对Topic 的消费队列进行负载均衡的时候产生不正确的结果，最终导致消息丢失。(官方源码设计)
		consumer.WithGroupName(group),
		consumer.WithNameServer([]string{config.Config.Rocketmq.Address}),
	)
	if err != nil {
		return err
	}
	err = c.Subscribe(topic, consumer.MessageSelector{}, middleware.UploadVideoCallBack)
	if err != nil {
		return err
	}
	err = c.Start()
	if err != nil {
		return err
	}
	return nil
}
