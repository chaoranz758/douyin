package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	CodeSuccess = 0 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeNeedLogin
	CodeInvalidToken
	CodeUploadFileFailed
	CodeSaveFileFailed
	CodeProducerSendMessageFailed
	CodeFollowFailed
	CodeFavoriteRepeated
	CodeFollowRepeated
	CodeCrawlerRequest
)

var codeMsgMap = map[int]string{
	CodeSuccess:                   "success",
	CodeInvalidParam:              "请求参数错误",
	CodeUserExist:                 "用户名已存在",
	CodeUserNotExist:              "用户名不存在",
	CodeInvalidPassword:           "用户名或密码错误",
	CodeServerBusy:                "服务繁忙",
	CodeNeedLogin:                 "请先登录",
	CodeInvalidToken:              "无效的token",
	CodeUploadFileFailed:          "上传视频文件失败",
	CodeSaveFileFailed:            "保存视频文件失败",
	CodeProducerSendMessageFailed: "生产者发送消息失败",
	CodeFollowFailed:              "不能对自己进行关系操作",
	CodeFavoriteRepeated:          "请不要重复点赞",
	CodeFollowRepeated:            "请不要重复关注",
	CodeCrawlerRequest:            "恶意爬虫请求",
}

func getMsg(code int) string {
	return codeMsgMap[code]
}

//把响应部分的函数单独封装起来

func responseError(c *gin.Context, code int, dataName string) {
	c.JSON(http.StatusOK, gin.H{
		"status_code": code,
		"status_msg":  getMsg(code),
		dataName:      nil,
	})
}

func responseSuccess(c *gin.Context, code int, data interface{}, dataName string) {
	c.JSON(http.StatusOK, gin.H{
		"status_code": code,
		"status_msg":  getMsg(code),
		dataName:      data,
	})
}
