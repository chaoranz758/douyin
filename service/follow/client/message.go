package client

import (
	"douyin/proto/message/api"
	"douyin/service/follow/config"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var MessageClient api.MessageClient

func InitMessage() error {
	p := grpc.NewServer()
	//向GRPC注册健康检查服务
	healthCheck := health.NewServer()
	healthpb.RegisterHealthServer(p, healthCheck)
	client, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?healthy=true",
			config.Config.ConsulServer.Ip,
			config.Config.ConsulServer.Port,
			config.Config.RequestGRPCServer.MessageService.Name),
		// 指定round_robin策略
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	MessageClient = api.NewMessageClient(client)
	return nil
}
