package handler

import (
	"context"
	"douyin/web/initialize/rocketmq"
	"douyin/web/model/request"
	"douyin/web/service"
	"douyin/web/util"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	publishVideoSuccess    = "publish video success"
	getPublishVideoSuccess = "get publish video success"
	getVideoListSuccess    = "get video list success"
	errorSendMessage       = "生产者消息发送失败"
	successSendMessage     = "生产者消息发送成功"
)

const (
	baseUrl       = "C:\\Users\\27155\\Desktop\\douyin_videos\\"
	baseUrlALiYun = "https://simpledouyin.oss-cn-qingdao.aliyuncs.com/"
	houzhui       = ".mp4"
	topic         = "uploadVideoTopic2"
)

type ProducerMessage1 struct {
	VideoURLLocal string `json:"videoURLLocal"`
	VideoName     string `json:"videoName"`
	ImageName     string `json:"imageName"`
}

func PublishVideoHandle(c *gin.Context) {
	//1.接收参数
	var publishVideoRequest request.PublishVideoRequest
	if err := c.ShouldBind(&publishVideoRequest); err != nil {
		zap.L().Error(getMsg(CodeInvalidParam), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeInvalidParam,
			"status_msg":  getMsg(CodeInvalidParam),
		})
		c.Abort()
		return
	}
	//token校验
	mc, err := util.ParseToken(publishVideoRequest.Token)
	if err != nil {
		zap.L().Error(getMsg(CodeInvalidToken), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeInvalidToken,
			"status_msg":  getMsg(CodeInvalidToken),
		})
		c.Abort()
		return
	}
	file, _ := c.FormFile("data")
	videoURLLocal := baseUrl + publishVideoRequest.Title + houzhui
	if err = c.SaveUploadedFile(file, videoURLLocal); err != nil {
		zap.L().Error(getMsg(CodeSaveFileFailed), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeSaveFileFailed,
			"status_msg":  getMsg(CodeSaveFileFailed),
		})
		c.Abort()
		return
	}
	//走消息队列
	var producerMessage1 = ProducerMessage1{
		VideoURLLocal: videoURLLocal,
		VideoName:     file.Filename,
		ImageName:     file.Filename,
	}
	data, _ := json.Marshal(producerMessage1)
	msg := &primitive.Message{
		Topic: topic,
		Body:  data,
	}
	_, err = rocketmq.Producer1.SendSync(context.Background(), msg)
	if err != nil {
		zap.L().Error(errorSendMessage, zap.Error(err))
		zap.L().Error(getMsg(CodeProducerSendMessageFailed), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeProducerSendMessageFailed,
			"status_msg":  getMsg(CodeProducerSendMessageFailed),
		})
		c.Abort()
		return
	}
	zap.L().Info(successSendMessage)
	videoName1 := fmt.Sprintf("videos/%s.mp4", file.Filename)
	videoURL := baseUrlALiYun + videoName1
	imageName1 := fmt.Sprintf("images/%s.jpg", file.Filename)
	imageURL := baseUrlALiYun + imageName1
	//2.调用rpc
	if err = service.PublishVideo(publishVideoRequest, mc.UserID, videoURL, imageURL); err != nil {
		zap.L().Error(getMsg(CodeServerBusy), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeServerBusy,
			"status_msg":  getMsg(CodeServerBusy),
		})
		c.Abort()
		return
	}
	zap.L().Info(publishVideoSuccess)
	//3.返回结果
	c.JSON(http.StatusOK, gin.H{
		"status_code": CodeSuccess,
		"status_msg":  getMsg(CodeSuccess),
	})
}

func GetPublishVideoHandle(c *gin.Context) {
	//接收参数
	var getPublishVideoRequest request.GetPublishVideoRequest
	if err := c.ShouldBindQuery(&getPublishVideoRequest); err != nil {
		zap.L().Error(getMsg(CodeInvalidParam), zap.Error(err))
		responseError(c, CodeInvalidParam, "video_list")
		c.Abort()
		return
	}
	if getPublishVideoRequest.Token == "" {
		responseError(c, CodeNeedLogin, "video_list")
		c.Abort()
		return
	}
	//token校验
	mc, err := util.ParseToken(getPublishVideoRequest.Token)
	if err != nil {
		zap.L().Error(getMsg(CodeInvalidToken), zap.Error(err))
		responseError(c, CodeInvalidToken, "video_list")
		c.Abort()
		return
	}
	//2.调用rpc
	data, err := service.GetPublishVideo(getPublishVideoRequest, mc.UserID)
	if err != nil {
		zap.L().Error(getMsg(CodeServerBusy), zap.Error(err))
		responseError(c, CodeServerBusy, "video_list")
		c.Abort()
		return
	}
	zap.L().Info(getPublishVideoSuccess)
	//3.返回结果
	responseSuccess(c, CodeSuccess, data, "video_list")
}

func GetVideoListHandle(c *gin.Context) {
	//1.接收参数
	var getVideoListRequest request.GetVideoListRequest
	if err := c.ShouldBindQuery(&getVideoListRequest); err != nil {
		zap.L().Error(getMsg(CodeInvalidParam), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeInvalidParam,
			"status_msg":  getMsg(CodeInvalidParam),
			"next_time":   0,
			"video_list":  nil,
		})
		c.Abort()
		return
	}
	var userId int64 = 0
	if getVideoListRequest.Token != "" {
		//token校验
		mc, err := util.ParseToken(getVideoListRequest.Token)
		if err != nil {
			zap.L().Error(getMsg(CodeInvalidToken), zap.Error(err))
			c.JSON(http.StatusOK, gin.H{
				"status_code": CodeInvalidToken,
				"status_msg":  getMsg(CodeInvalidToken),
				"next_time":   0,
				"video_list":  nil,
			})
			c.Abort()
			return
		}
		userId = mc.UserID
	}
	//2.调用rpc
	videoList, nextTime, err := service.GetVideoList(getVideoListRequest, userId)
	if err != nil {
		zap.L().Error(getMsg(CodeServerBusy), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeServerBusy,
			"status_msg":  getMsg(CodeServerBusy),
			"next_time":   0,
			"video_list":  nil,
		})
		c.Abort()
		return
	}
	zap.L().Info(getVideoListSuccess)
	//3.返回结果
	c.JSON(http.StatusOK, gin.H{
		"status_code": CodeSuccess,
		"status_msg":  getMsg(CodeSuccess),
		"next_time":   nextTime,
		"video_list":  videoList,
	})

}
