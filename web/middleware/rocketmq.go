package middleware

import (
	"context"
	"douyin/web/util"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
)

const (
	errorUploadFile    = "upload file failed"
	errorJsonUnmarshal = "json unmarshal failed"
)

type ProducerMessage1 struct {
	VideoURLLocal string `json:"videoURLLocal"`
	VideoName     string `json:"videoName"`
	ImageName     string `json:"imageName"`
}

func UploadVideoCallBack(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	var producer1Message ProducerMessage1
	if err := json.Unmarshal(msgs[0].Body, &producer1Message); err != nil {
		zap.L().Error(errorJsonUnmarshal, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	err := util.Upload(producer1Message.VideoURLLocal, producer1Message.VideoName, producer1Message.ImageName)
	if err != nil {
		zap.L().Error(errorUploadFile, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	return consumer.ConsumeSuccess, nil
}
