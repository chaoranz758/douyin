package main

import (
	"douyin/web/initialize/config"
	consul2 "douyin/web/initialize/consul"
	"douyin/web/initialize/grpc_client"
	log1 "douyin/web/initialize/log"
	"douyin/web/initialize/oss"
	"douyin/web/initialize/rocketmq"
	"douyin/web/router"
	"douyin/web/util"
	"fmt"
	"go.uber.org/zap"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("need config file.eg: config/config.json")
		return
	}
	//1.加载配置
	if err := config.Init(os.Args[1]); err != nil {
		log.Fatal("init settings failed, err:", err)
	}
	//2.初始化日志
	if err := log1.Init(config.Config.DouYinWeb.Mode); err != nil {
		log.Fatal("init logger failed, err:", err)
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")
	//3.初始化consul连接
	if err := consul2.Init(fmt.Sprintf("%s:%d", config.Config.ConsulServer.Ip, config.Config.ConsulServer.Port)); err != nil {
		zap.L().Error("init consul connection failed...", zap.Error(err))
		log.Fatal("init consul connection failed...")
	}
	zap.L().Debug("init consul connection success...")
	//4.启动GRPC客户端
	if err := grpc.InitUserClient(); err != nil {
		zap.L().Error("init user rpc client failed...", zap.Error(err))
		log.Fatal("init user rpc client failed...")
	}
	zap.L().Debug("init user rpc client success...")
	if err := grpc.InitVideoClient(); err != nil {
		zap.L().Error("init video rpc client failed...", zap.Error(err))
		log.Fatal("init video rpc client failed...")
	}
	zap.L().Debug("init video rpc client success...")
	if err := grpc.InitCommentClient(); err != nil {
		zap.L().Error("init comment rpc client failed...", zap.Error(err))
		log.Fatal("init comment rpc client failed...")
	}
	zap.L().Debug("init comment rpc client success...")
	if err := grpc.InitFavoriteClient(); err != nil {
		zap.L().Error("init favorite rpc client failed...", zap.Error(err))
		log.Fatal("init favorite rpc client failed...")
	}
	zap.L().Debug("init favorite rpc client success...")
	if err := grpc.InitFollowClient(); err != nil {
		zap.L().Error("init follow rpc client failed...", zap.Error(err))
		log.Fatal("init follow rpc client failed...")
	}
	zap.L().Debug("init follow rpc client success...")
	if err := grpc.InitMessageClient(); err != nil {
		zap.L().Error("init message rpc client failed...", zap.Error(err))
		log.Fatal("init message rpc client failed...")
	}
	zap.L().Debug("init message rpc client success...")
	//5.向Consul注册服务
	if err := consul2.RegisterService(); err != nil {
		zap.L().Error("register service to consul failed...", zap.Error(err))
		log.Fatal("register service to consul failed...")
	}
	zap.L().Debug("register service to consul success...")
	//6.启动阿里云oss上传客户端
	if err := oss.InitOSS(); err != nil {
		zap.L().Error("init oss failed", zap.Error(err))
		log.Fatal("init oss failed")
	}
	zap.L().Debug("init oss success")
	//7.开启消息队列生产者、消费者
	if err := rocketmq.InitProducer1(); err != nil {
		zap.L().Error("init producer1 failed...", zap.Error(err))
		log.Fatal("init producer1 failed, err:", err)
	}
	zap.L().Debug("init producer1 success...")
	if err := rocketmq.InitCustomer1(); err != nil {
		zap.L().Error("init customer1 failed...", zap.Error(err))
		log.Fatal("init customer1 failed, err:", err)
	}
	zap.L().Debug("init customer1 success...")
	//8.开启路由
	r := router.Init()
	//9.启动服务（优雅关机）
	util.Spin(r)
}
