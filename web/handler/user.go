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
	userRegisterSuccess = "用户注册成功"
	userLoginSuccess    = "用户登录成功"
	getUserInfoSuccess  = "获取用户信息成功"
)

func UserLoginHandle(c *gin.Context) {
	//1.接收参数
	var userLoginRequest request.UserLoginRequest
	if err := c.ShouldBindQuery(&userLoginRequest); err != nil {
		zap.L().Error(getMsg(CodeInvalidParam), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeInvalidParam,
			"status_msg":  getMsg(CodeInvalidParam),
			"user_id":     0,
			"token":       "",
		})
		c.Abort()
		return
	}
	//2.调用rpc
	userID, token, err := service.UserLogin(userLoginRequest)
	if err != nil {
		zap.L().Error(getMsg(CodeServerBusy), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeServerBusy,
			"status_msg":  getMsg(CodeServerBusy),
			"user_id":     0,
			"token":       "",
		})
		c.Abort()
		return
	}
	zap.L().Info(userLoginSuccess)
	//3.返回结果
	c.JSON(http.StatusOK, gin.H{
		"status_code": CodeSuccess,
		"status_msg":  getMsg(CodeSuccess),
		"user_id":     userID,
		"token":       token,
	})
}

func UserRegisterHandle(c *gin.Context) {
	//1.接收参数
	var userRegisterRequest request.UserRegisterRequest
	if err := c.ShouldBindQuery(&userRegisterRequest); err != nil {
		zap.L().Error(getMsg(CodeInvalidParam), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeInvalidParam,
			"status_msg":  getMsg(CodeInvalidParam),
			"user_id":     0,
			"token":       "",
		})
		c.Abort()
		return
	}
	//2.调用rpc
	userID, token, err := service.UserRegister(userRegisterRequest)
	if err != nil {
		zap.L().Error(getMsg(CodeServerBusy), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeServerBusy,
			"status_msg":  getMsg(CodeServerBusy),
			"user_id":     0,
			"token":       "",
		})
		c.Abort()
		return
	}
	zap.L().Info(userRegisterSuccess)
	//3.返回结果
	c.JSON(http.StatusOK, gin.H{
		"status_code": CodeSuccess,
		"status_msg":  getMsg(CodeSuccess),
		"user_id":     userID,
		"token":       token,
	})
}

func GetUserInfoHandle(c *gin.Context) {
	//接收参数
	var getUserInfoRequest request.GetUserInfoRequest
	if err := c.ShouldBindQuery(&getUserInfoRequest); err != nil {
		zap.L().Error(getMsg(CodeInvalidParam), zap.Error(err))
		responseError(c, CodeInvalidParam, "user")
		c.Abort()
		return
	}
	//token校验
	mc, err := jwt.ParseToken(getUserInfoRequest.Token)
	if err != nil {
		zap.L().Error(getMsg(CodeInvalidToken), zap.Error(err))
		responseError(c, CodeInvalidToken, "user")
		c.Abort()
		return
	}
	//2.调用rpc
	data, err := service.GetUserInfo(getUserInfoRequest, mc.UserID)
	if err != nil {
		zap.L().Error(getMsg(CodeServerBusy), zap.Error(err))
		responseError(c, CodeServerBusy, "user")
		c.Abort()
		return
	}
	zap.L().Info(getUserInfoSuccess)
	//3.返回结果
	responseSuccess(c, CodeSuccess, data, "user")
}
