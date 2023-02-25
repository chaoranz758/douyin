package server

import (
	"douyin/proto/comment/api"
	"douyin/service/comment/handler"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

func InitComment() error {
	listen, err := net.Listen("tcp", ":9080")
	if err != nil {
		zap.L().Error("server tcp port start failed", zap.Error(err))
		return err
	}
	p := grpc.NewServer()
	//向GRPC注册健康检查服务
	healthCheck := health.NewServer()
	healthpb.RegisterHealthServer(p, healthCheck)
	//向GRPC注册用户服务
	api.RegisterCommentServer(p, &handler.Comment{})
	err = p.Serve(listen)
	if err != nil {
		zap.L().Error("server grpc start failed", zap.Error(err))
		return err
	}
	return nil
}