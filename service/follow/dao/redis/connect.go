package redis

import (
	"context"
	"douyin/service/follow/initialize/config"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

var (
	rdb *redis.Client
)

// 初始化连接

func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Config.Redis.Host, config.Config.Redis.Port),
		Password: config.Config.Redis.Password, // no password set
		DB:       config.Config.Redis.Db,       // use default DB
		PoolSize: config.Config.Redis.PoolSize, // 连接池大小
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		zap.L().Error(errorConnectRedis, zap.Error(err))
		return err
	}
	return nil
}
