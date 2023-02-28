package handler

import (
	"douyin/web/model/request"
	"douyin/web/service"
	"douyin/web/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	sendMessageSuccess = "send message success"
	getMessageSuccess  = "get message success"
)

func SendMessageHandle(c *gin.Context) {
	//1.接收参数
	var sendMessageRequest request.SendMessageRequest
	if err := c.ShouldBindQuery(&sendMessageRequest); err != nil {
		zap.L().Error(getMsg(CodeInvalidParam), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeInvalidParam,
			"status_msg":  getMsg(CodeInvalidParam),
		})
		c.Abort()
		return
	}
	//token校验
	mc, err := util.ParseToken(sendMessageRequest.Token)
	if err != nil {
		zap.L().Error(getMsg(CodeInvalidToken), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeInvalidToken,
			"status_msg":  getMsg(CodeInvalidToken),
		})
		c.Abort()
		return
	}
	//2.调用rpc
	if err := service.SendMessage(sendMessageRequest, mc.UserID); err != nil {
		zap.L().Error(getMsg(CodeServerBusy), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeServerBusy,
			"status_msg":  getMsg(CodeServerBusy),
		})
		c.Abort()
		return
	}
	zap.L().Info(sendMessageSuccess)
	//3.返回结果
	c.JSON(http.StatusOK, gin.H{
		"status_code": CodeSuccess,
		"status_msg":  getMsg(CodeSuccess),
	})
}

func GetMessageHandle(c *gin.Context) {
	//1.接收参数
	var getMessageRequest request.GetMessageRequest
	if err := c.ShouldBindQuery(&getMessageRequest); err != nil {
		zap.L().Error(getMsg(CodeInvalidParam), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": "1",
			"status_msg":  getMsg(CodeInvalidParam),
		})
		c.Abort()
		return
	}
	//token校验
	mc, err := util.ParseToken(getMessageRequest.Token)
	if err != nil {
		zap.L().Error(getMsg(CodeInvalidToken), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": "7",
			"status_msg":  getMsg(CodeInvalidToken),
		})
		c.Abort()
		return
	}
	//2.调用rpc
	data, err := service.GetMessage(getMessageRequest, mc.UserID)
	if err != nil {
		zap.L().Error(getMsg(CodeServerBusy), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": "5",
			"status_msg":  getMsg(CodeServerBusy),
		})
		c.Abort()
		return
	}
	zap.L().Info(getMessageSuccess)
	//3.返回结果
	c.JSON(http.StatusOK, gin.H{
		"status_code":  CodeSuccess,
		"status_msg":   getMsg(CodeSuccess),
		"message_list": data,
	})
}
