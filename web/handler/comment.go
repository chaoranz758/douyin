package handler

import (
	"douyin/web/model/request"
	"douyin/web/service"
	"douyin/web/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	commentVideoSuccess        = "comment video success"
	getCommentVideoListSuccess = "get comment video list success"
)

func CommentVideoHandle(c *gin.Context) {
	//1.接收参数
	var commentVideoRequest request.CommentVideoRequest
	if err := c.ShouldBindQuery(&commentVideoRequest); err != nil {
		zap.L().Error(getMsg(CodeInvalidParam), zap.Error(err))
		responseError(c, CodeInvalidParam, "comment")
		c.Abort()
		return
	}
	fmt.Printf("%v\n", commentVideoRequest)
	//token校验
	mc, err := util.ParseToken(commentVideoRequest.Token)
	if err != nil {
		zap.L().Error(getMsg(CodeInvalidToken), zap.Error(err))
		responseError(c, CodeInvalidToken, "comment")
		c.Abort()
		return
	}
	//2.调用rpc
	data, err := service.CommentVideo(commentVideoRequest, mc.UserID)
	if err != nil {
		zap.L().Error(getMsg(CodeServerBusy), zap.Error(err))
		responseError(c, CodeServerBusy, "comment")
		c.Abort()
		return
	}
	zap.L().Info(commentVideoSuccess)
	//3.返回结果
	responseSuccess(c, CodeSuccess, data, "comment")
}

func GetCommentVideoListHandle(c *gin.Context) {
	//1.接收参数
	var getCommentVideoListRequest request.GetCommentVideoListRequest
	if err := c.ShouldBindQuery(&getCommentVideoListRequest); err != nil {
		zap.L().Error(getMsg(CodeInvalidParam), zap.Error(err))
		responseError(c, CodeInvalidParam, "comment_list")
		c.Abort()
		return
	}
	//token校验
	mc, err := util.ParseToken(getCommentVideoListRequest.Token)
	if err != nil {
		zap.L().Error(getMsg(CodeInvalidToken), zap.Error(err))
		responseError(c, CodeInvalidToken, "comment_list")
		c.Abort()
		return
	}
	//2.调用rpc
	data, err := service.GetCommentVideoList(getCommentVideoListRequest, mc.UserID)
	if err != nil {
		zap.L().Error(getMsg(CodeServerBusy), zap.Error(err))
		responseError(c, CodeServerBusy, "comment_list")
		c.Abort()
		return
	}
	zap.L().Info(getCommentVideoListSuccess)
	//3.返回结果
	responseSuccess(c, CodeSuccess, data, "comment_list")

}
