package server

import (
	"douyin/proto/user/api"
	"douyin/service/user/handler"
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

func InitUser() error {
	listen, err := net.Listen("tcp", ":9084")
	if err != nil {
		zap.L().Error(errorNetListen, zap.Error(err))
		return err
	}
	p := grpc.NewServer()
	//向GRPC注册健康检查服务
	healthCheck := health.NewServer()
	healthpb.RegisterHealthServer(p, healthCheck)
	//向GRPC注册用户服务
	api.RegisterUserServer(p, &handler.User{})
	err = p.Serve(listen)
	if err != nil {
		zap.L().Error(errorGrpcStart, zap.Error(err))
		return err
	}
	return nil
}
