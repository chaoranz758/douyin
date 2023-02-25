package main

import (
	"douyin/service/user/client"
	"douyin/service/user/config"
	"douyin/service/user/dao/mysql"
	"douyin/service/user/dao/redis"
	"douyin/service/user/job"
	log1 "douyin/service/user/log"
	"douyin/service/user/pkg/consul"
	"douyin/service/user/pkg/snowflake"
	"douyin/service/user/server"
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
	if err := client.InitFollow(); err != nil {
		zap.L().Error("init follow rpc client failed...", zap.Error(err))
		log.Fatal("init follow rpc client failed...")
	}
	zap.L().Debug("init follow rpc client success...")
	if err := client.InitVideo(); err != nil {
		zap.L().Error("init video rpc client failed...", zap.Error(err))
		log.Fatal("init video rpc client failed...")
	}
	zap.L().Debug("init video rpc client success...")
	if err := client.InitFavorite(); err != nil {
		zap.L().Error("init favorite rpc client failed...", zap.Error(err))
		log.Fatal("init favorite rpc client failed...")
	}
	zap.L().Debug("init favorite rpc client success...")
	//9.开启定时任务
	if err := job.InitCron(); err != nil {
		zap.L().Error("init cron job failed...", zap.Error(err))
		log.Fatal("init cron job failed, err:", err)
	}
	zap.L().Debug("init cron job success...")
	//10.开启grpc服务
	if err := server.InitUser(); err != nil {
		zap.L().Error("init rpc service failed...", zap.Error(err))
		log.Fatal("init rpc service failed...")
	}
}
