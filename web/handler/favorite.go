package handler

import (
	"douyin/web/model/request"
	"douyin/web/pkg/jwt"
	"douyin/web/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const (
	favoriteVideoSuccess        = "favorite video success"
	getFavoriteVideoListSuccess = "get favorite video list success"
)

func FavoriteVideoHandle(c *gin.Context) {
	//1.接收参数
	var favoriteVideoRequest request.FavoriteVideoRequest
	if err := c.ShouldBindQuery(&favoriteVideoRequest); err != nil {
		zap.L().Error(getMsg(CodeInvalidParam), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeInvalidParam,
			"status_msg":  getMsg(CodeInvalidParam),
		})
		c.Abort()
		return
	}
	//token校验
	mc, err := jwt.ParseToken(favoriteVideoRequest.Token)
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
	if err := service.FavoriteVideo(favoriteVideoRequest, mc.UserID); err != nil {
		if err.Error() == "rpc error: code = Unknown desc = 点赞关系数据没有插入进favorite表但没有报错" {
			zap.L().Error(getMsg(CodeFavoriteRepeated), zap.Error(err))
			c.JSON(http.StatusOK, gin.H{
				"status_code": CodeFavoriteRepeated,
				"status_msg":  getMsg(CodeFavoriteRepeated),
			})
			c.Abort()
			return
		}
		zap.L().Error(getMsg(CodeServerBusy), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeServerBusy,
			"status_msg":  getMsg(CodeServerBusy),
		})
		c.Abort()
		return
	}
	zap.L().Info(favoriteVideoSuccess)
	//3.返回结果
	c.JSON(http.StatusOK, gin.H{
		"status_code": CodeSuccess,
		"status_msg":  getMsg(CodeSuccess),
	})
}

func GetFavoriteVideoListHandle(c *gin.Context) {
	//接收参数
	var getFavoriteVideoListRequest request.GetFavoriteVideoListRequest
	if err := c.ShouldBindQuery(&getFavoriteVideoListRequest); err != nil {
		zap.L().Error(getMsg(CodeInvalidParam), zap.Error(err))
		responseError(c, CodeInvalidParam, "video_list")
		c.Abort()
		return
	}
	if getFavoriteVideoListRequest.Token == "" {
		responseError(c, CodeNeedLogin, "video_list")
		c.Abort()
		return
	}
	//token校验
	mc, err := jwt.ParseToken(getFavoriteVideoListRequest.Token)
	if err != nil {
		zap.L().Error(getMsg(CodeInvalidToken), zap.Error(err))
		responseError(c, CodeInvalidToken, "video_list")
		c.Abort()
		return
	}
	//2.调用rpc
	data, err := service.GetFavoriteVideoList(getFavoriteVideoListRequest, mc.UserID)
	if err != nil {
		zap.L().Error(getMsg(CodeServerBusy), zap.Error(err))
		responseError(c, CodeServerBusy, "video_list")
		c.Abort()
		return
	}
	zap.L().Info(getFavoriteVideoListSuccess)
	//3.返回结果
	responseSuccess(c, CodeSuccess, data, "video_list")
}
