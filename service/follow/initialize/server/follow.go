package server

import (
	"douyin/proto/follow/api"
	"douyin/service/follow/handler"
	"github.com/dtm-labs/dtmgrpc/workflow"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

var DtmServer = "127.0.0.1:36790"

func InitFollow() error {
	p := grpc.NewServer()
	workflow.InitGrpc(DtmServer, "10.122.238.133:9082", p)
	listen, err := net.Listen("tcp", ":9082")
	if err != nil {
		zap.L().Error("server tcp port start failed", zap.Error(err))
		return err
	}
	//p := grpc.NewServer()
	//向GRPC注册健康检查服务
	healthCheck := health.NewServer()
	healthpb.RegisterHealthServer(p, healthCheck)
	//向GRPC注册用户服务
	api.RegisterFollowServer(p, &handler.Follow{})
	err = p.Serve(listen)
	if err != nil {
		zap.L().Error("server grpc start failed", zap.Error(err))
		return err
	}
	return nil
}
