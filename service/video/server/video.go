package server

import (
	"douyin/proto/video/api"
	"douyin/service/video/handler"
	"github.com/dtm-labs/dtmgrpc/workflow"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

var DtmServer = "127.0.0.1:36790"

func InitVideo() error {
	p := grpc.NewServer()
	workflow.InitGrpc(DtmServer, "10.122.238.133:9085", p)
	listen, err := net.Listen("tcp", ":9085")
	if err != nil {
		zap.L().Error("server tcp port start failed", zap.Error(err))
		return err
	}
	//向GRPC注册健康检查服务
	healthCheck := health.NewServer()
	healthpb.RegisterHealthServer(p, healthCheck)
	//向GRPC注册用户服务
	api.RegisterVideoServer(p, &handler.Video{})
	err = p.Serve(listen)
	if err != nil {
		zap.L().Error("server grpc start failed", zap.Error(err))
		return err
	}
	return nil
}
