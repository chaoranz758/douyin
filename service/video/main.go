package main

import (
	"douyin/service/video/client"
	"douyin/service/video/config"
	"douyin/service/video/dao/mysql"
	"douyin/service/video/dao/redis"
	log1 "douyin/service/video/log"
	"douyin/service/video/pkg/consul"
	"douyin/service/video/pkg/snowflake"
	"douyin/service/video/server"
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
	if err := consul.Init(fmt.Sprintf("%s:%d", config.Config.ConsulServer.Ip, config.Config.ConsulServer.Port)); err != nil {
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
	//6.初始化雪花ID生成
	if err := snowflake.Init(config.Config.SnowFlake.StartTime, config.Config.SnowFlake.MachineID); err != nil {
		zap.L().Error("snowflake generate id failed...", zap.Error(err))
		log.Fatal("snowflake generate id failed...")
	}
	zap.L().Debug("snowflake generate id success...")
	//7.向consul注册服务
	if err := consul.RegisterService(); err != nil {
		zap.L().Error("register service to consul failed...", zap.Error(err))
		log.Fatal("register service to consul failed...")
	}
	zap.L().Debug("register service to consul success...")
	//8.启动grpc客户端
	if err := client.InitFavorite(); err != nil {
		zap.L().Error("init favorite rpc client failed...", zap.Error(err))
		log.Fatal("init favorite rpc client failed...")
	}
	if err := client.InitComment(); err != nil {
		zap.L().Error("init comment rpc client failed...", zap.Error(err))
		log.Fatal("init comment rpc client failed...")
	}
	if err := client.InitUser(); err != nil {
		zap.L().Error("init user rpc client failed...", zap.Error(err))
		log.Fatal("init user rpc client failed...")
	}
	if err := client.InitUserDtm(); err != nil {
		zap.L().Error("init user dtm rpc client failed...", zap.Error(err))
		log.Fatal("init user dtm rpc client failed...")
	}
	//9.开启grpc服务
	if err := server.InitVideo(); err != nil {
		zap.L().Error("init rpc service failed...", zap.Error(err))
		log.Fatal("init rpc service failed...")
	}
	zap.L().Debug("init rpc service success...")
}
