package main

import (
	"douyin/service/favorite/dao/mysql"
	"douyin/service/favorite/dao/redis"
	"douyin/service/favorite/initialize/config"
	consul2 "douyin/service/favorite/initialize/consul"
	"douyin/service/favorite/initialize/grpc_client"
	log1 "douyin/service/favorite/initialize/log"
	rocketmq2 "douyin/service/favorite/initialize/rocketmq"
	"douyin/service/favorite/initialize/server"
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
		log.Fatal("init settings failed, err: ", err)
	}
	//2.初始化日志
	if err := log1.Init(config.Config.DouYinService.Mode); err != nil {
		log.Fatal("init logger failed, err: ", err)
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")
	//3.初始化consul连接
	if err := consul2.Init(fmt.Sprintf("%s:%d", config.Config.ConsulServer.Ip, config.Config.ConsulServer.Port)); err != nil {
		zap.L().Error("init consul connection failed...", zap.Error(err))
		log.Fatal("init consul connection failed...")
	}
	zap.L().Debug("init consul connection success...")
	//4.初始化mysql连接
	if err := mysql.Init(); err != nil {
		zap.L().Error("init mysql failed...", zap.Error(err))
		log.Fatal("init mysql failed, err:", err)
	}
	zap.L().Debug("init mysql success...")
	//5.初始化redis连接
	if err := redis.Init(); err != nil {
		zap.L().Error("init redis failed...", zap.Error(err))
		log.Fatal("init redis failed, err:", err)
	}
	zap.L().Debug("init redis success...")
	//6.向consul注册服务
	if err := consul2.RegisterService(); err != nil {
		zap.L().Error("register service to consul failed...", zap.Error(err))
		log.Fatal("register service to consul failed...")
	}
	zap.L().Debug("register service to consul success...")
	//7.启动grpc客户端
	if err := grpc_client.InitUser(); err != nil {
		zap.L().Error("init user rpc client failed...", zap.Error(err))
		log.Fatal("init user rpc client failed...")
	}
	zap.L().Debug("init user rpc client success...")
	if err := grpc_client.InitVideo(); err != nil {
		zap.L().Error("init video rpc client failed...", zap.Error(err))
		log.Fatal("init video rpc client failed...")
	}
	zap.L().Debug("init video rpc client success...")
	if err := grpc_client.InitComment(); err != nil {
		zap.L().Error("init comment rpc client failed...", zap.Error(err))
		log.Fatal("init comment rpc client failed...")
	}
	zap.L().Debug("init comment rpc client success...")
	if err := grpc_client.InitUserDtm(); err != nil {
		zap.L().Error("init user dtm rpc client failed...", zap.Error(err))
		log.Fatal("init user dtm rpc client failed...")
	}
	zap.L().Debug("init user dtm rpc client success...")
	if err := grpc_client.InitVideoDtm(); err != nil {
		zap.L().Error("init video dtm rpc client failed...", zap.Error(err))
		log.Fatal("init video dtm rpc client failed...")
	}
	zap.L().Debug("init video dtm rpc client success...")
	//8.启动grpc生产者
	if err := rocketmq2.InitProducer1(); err != nil {
		zap.L().Error("init producer1 failed...", zap.Error(err))
		log.Fatal("init producer1 failed, err:", err)
	}
	zap.L().Debug("init producer1 success...")
	if err := rocketmq2.InitProducer2(); err != nil {
		zap.L().Error("init producer2 failed...", zap.Error(err))
		log.Fatal("init producer2 failed, err:", err)
	}
	zap.L().Debug("init producer2 success...")
	if err := rocketmq2.InitProducer3(); err != nil {
		zap.L().Error("init producer3 failed...", zap.Error(err))
		log.Fatal("init producer3 failed, err:", err)
	}
	zap.L().Debug("init producer3 success...")
	if err := rocketmq2.InitProducer4(); err != nil {
		zap.L().Error("init producer4 failed...", zap.Error(err))
		log.Fatal("init producer4 failed, err:", err)
	}
	zap.L().Debug("init producer4 success...")
	//9.开启grpc消费者
	if err := rocketmq2.InitCustomer1(); err != nil {
		zap.L().Error("init customer1 failed...", zap.Error(err))
		log.Fatal("init customer1 failed, err:", err)
	}
	zap.L().Debug("init customer1 success...")
	if err := rocketmq2.InitCustomer2(); err != nil {
		zap.L().Error("init customer2 failed...", zap.Error(err))
		log.Fatal("init customer2 failed, err:", err)
	}
	zap.L().Debug("init customer2 success...")
	if err := rocketmq2.InitCustomer3(); err != nil {
		zap.L().Error("init customer3 failed...", zap.Error(err))
		log.Fatal("init customer3 failed, err:", err)
	}
	zap.L().Debug("init customer3 success...")
	if err := rocketmq2.InitCustomer4(); err != nil {
		zap.L().Error("init customer4 failed...", zap.Error(err))
		log.Fatal("init customer4 failed, err:", err)
	}
	zap.L().Debug("init customer4 success...")
	//10.开启grpc服务
	if err := server.InitFavorite(); err != nil {
		zap.L().Error("init rpc service failed...", zap.Error(err))
		log.Fatal("init rpc service failed...")
	}
	zap.L().Debug("init rpc service success...")
}
