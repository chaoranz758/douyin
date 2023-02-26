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

func InitUser() error {
	listen, err := net.Listen("tcp", ":9084")
	if err != nil {
		zap.L().Error("server tcp port start failed", zap.Error(err))
		return err
	}
	p := grpc.NewServer()
	//向GRPC注册健康检查服务
	healthCheck := health.NewServer()
	healthpb.RegisterHealthServer(p, healthCheck)
	//向GRPC注册用户服务
	api.RegisterUserServer(p, &handler.User{})
	zap.L().Debug("init rpc service success...")
	//workflow.InitGrpc(DtmServer, "192.168.182.1:9084", p)
	err = p.Serve(listen)
	if err != nil {
		zap.L().Error("server grpc start failed", zap.Error(err))
		return err
	}
	return nil
}
