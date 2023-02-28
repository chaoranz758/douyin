package job

import (
	"douyin/service/message/middleware/job"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

const errorPushActiveSet = "push active set failed"

const spec = "@every 10m"

var c *cron.Cron

func InitCron() error {
	c = cron.New()
	_, err := c.AddFunc(spec, func() {
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
