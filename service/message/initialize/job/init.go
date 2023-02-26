package job

import (
	"douyin/service/message/middleware/job"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

const errorPushActiveSet = "push active set failed"

var c *cron.Cron

func InitCron() error {
	c = cron.New()
	_, err := c.AddFunc("@every 10m", func() {
		if err := job.PushActiveSet(); err != nil {
			zap.L().Error(errorPushActiveSet, zap.Error(err))
		}
	})
	if err != nil {
		return err
	}
	//开始执行任务
	c.Start()
	return nil
}
