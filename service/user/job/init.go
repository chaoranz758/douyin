package job

import (
	"douyin/service/user/service"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

const (
	errorPushActiveSetAndString = "push v or active set and string"
)

var c *cron.Cron

func InitCron() error {
	c = cron.New()
	_, err := c.AddFunc("@every 10m", func() {
		if err := service.PushActiveSetAndString(); err != nil {
			zap.L().Error(errorPushActiveSetAndString, zap.Error(err))
		}
	})
	if err != nil {
		return err
	}
	//开始执行任务
	c.Start()
	return nil
}
