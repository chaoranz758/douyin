package server

import (
	"douyin/proto/video/api"
	"douyin/service/video/handler"
	"douyin/service/video/initialize/config"
	"fmt"
	"github.com/dtm-labs/dtmgrpc/workflow"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

const (
	errorNetListen = "server tcp port start failed"
	errorGrpcStart = "server grpc start failed"
)

func InitVideo() error {
	p := grpc.NewServer()
	workflow.InitGrpc(config.Config.Dtm.Address, fmt.Sprintf("%s:%d", config.Config.ConsulRegister.IP, config.Config.ConsulRegister.Port), p)
	listen, err := net.Listen(config.Config.DouYinService.Protocol, fmt.Sprintf(":%d", config.Config.DouYinService.Port))
	if err != nil {
		zap.L().Error(errorNetListen, zap.Error(err))
		return err
	}
	//向GRPC注册健康检查服务
	healthCheck := health.NewServer()
	healthpb.RegisterHealthServer(p, healthCheck)
	//向GRPC注册用户服务
	api.RegisterVideoServer(p, &handler.Video{})
	err = p.Serve(listen)
	if err != nil {
		zap.L().Error(errorGrpcStart, zap.Error(err))
		return err
	}
	return nil
}
