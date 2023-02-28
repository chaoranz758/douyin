package rocketmq

import (
	"douyin/service/favorite/initialize/config"
	rocketmq1 "douyin/service/favorite/middleware/rocketmq"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
)

const (
	groupCostumer1 = "favoriteGroup1"
	groupCostumer2 = "favoriteGroup2"
	groupCostumer3 = "favoriteGroup3"
	groupCostumer4 = "favoriteGroup4"
	topic1         = "favoriteTopic1"
	topic2         = "favoriteTopic2"
	topic3         = "favoriteTopic3"
	topic4         = "favoriteTopic4"
)

func InitCustomer1() error {
	c, err := rocketmq.NewPushConsumer(
		// 指定 Group 可以实现消费者负载均衡进行消费，并且保证他们的Topic+Tag要一样。
		// 如果同一个 GroupID 下的不同消费者实例，订阅了不同的 Topic+Tag 将导致在对Topic 的消费队列进行负载均衡的时候产生不正确的结果，最终导致消息丢失。(官方源码设计)
		consumer.WithGroupName(groupCostumer1),
		consumer.WithNameServer([]string{config.Config.Rocketmq.Address}),
	)
	if err != nil {
		return err
	}
	err = c.Subscribe(topic1, consumer.MessageSelector{}, rocketmq1.FavoriteCustomer1CallBack)
	if err != nil {
		return err
	}
	err = c.Start()
	if err != nil {
		return err
	}
	return nil
}

func InitCustomer2() error {
	c, err := rocketmq.NewPushConsumer(
		// 指定 Group 可以实现消费者负载均衡进行消费，并且保证他们的Topic+Tag要一样。
		// 如果同一个 GroupID 下的不同消费者实例，订阅了不同的 Topic+Tag 将导致在对Topic 的消费队列进行负载均衡的时候产生不正确的结果，最终导致消息丢失。(官方源码设计)
		consumer.WithGroupName(groupCostumer2),
		consumer.WithNameServer([]string{config.Config.Rocketmq.Address}),
	)
	if err != nil {
		return err
	}
	err = c.Subscribe(topic2, consumer.MessageSelector{}, rocketmq1.FavoriteCustomer2CallBack)
	if err != nil {
		return err
	}
	err = c.Start()
	if err != nil {
		return err
	}
	return nil
}

func InitCustomer3() error {
	c, err := rocketmq.NewPushConsumer(
		// 指定 Group 可以实现消费者负载均衡进行消费，并且保证他们的Topic+Tag要一样。
		// 如果同一个 GroupID 下的不同消费者实例，订阅了不同的 Topic+Tag 将导致在对Topic 的消费队列进行负载均衡的时候产生不正确的结果，最终导致消息丢失。(官方源码设计)
		consumer.WithGroupName(groupCostumer3),
		consumer.WithNameServer([]string{config.Config.Rocketmq.Address}),
	)
	if err != nil {
		return err
	}
	err = c.Subscribe(topic3, consumer.MessageSelector{}, rocketmq1.FavoriteCustomer3CallBack)
	if err != nil {
		return err
	}
	err = c.Start()
	if err != nil {
		return err
	}
	return nil
}

func InitCustomer4() error {
	c, err := rocketmq.NewPushConsumer(
		// 指定 Group 可以实现消费者负载均衡进行消费，并且保证他们的Topic+Tag要一样。
		// 如果同一个 GroupID 下的不同消费者实例，订阅了不同的 Topic+Tag 将导致在对Topic 的消费队列进行负载均衡的时候产生不正确的结果，最终导致消息丢失。(官方源码设计)
		consumer.WithGroupName(groupCostumer4),
		consumer.WithNameServer([]string{config.Config.Rocketmq.Address}),
	)
	if err != nil {
		return err
	}
	err = c.Subscribe(topic4, consumer.MessageSelector{}, rocketmq1.FavoriteCustomer4CallBack)
	if err != nil {
		return err
	}
	err = c.Start()
	if err != nil {
		return err
	}
	return nil
}
