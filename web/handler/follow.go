package handler

import (
	"douyin/web/model/request"
	"douyin/web/pkg/jwt"
	"douyin/web/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

const (
	followUserSuccess      = "follow user success"
	getFollowListSuccess   = "get follow list success"
	getFollowerListSuccess = "get follower list success"
	getFriendListSuccess   = "get friend list success"
)

func FollowUserHandle(c *gin.Context) {
	//1.接收参数
	var followUserRequest request.FollowUserRequest
	if err := c.ShouldBindQuery(&followUserRequest); err != nil {
		zap.L().Error(getMsg(CodeInvalidParam), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeInvalidParam,
			"status_msg":  getMsg(CodeInvalidParam),
		})
		c.Abort()
		return
	}
	//token校验
	mc, err := jwt.ParseToken(followUserRequest.Token)
	if err != nil {
		zap.L().Error(getMsg(CodeInvalidToken), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeInvalidToken,
			"status_msg":  getMsg(CodeInvalidToken),
		})
		c.Abort()
		return
	}
	//如果是同一个人,返回不能关注自己
	userId1, _ := strconv.ParseInt(followUserRequest.ToUserID, 10, 64)
	if userId1 == mc.UserID {
		zap.L().Error(getMsg(CodeFollowFailed), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status_code": CodeFollowFailed,
			"status_msg":  getMsg(CodeFollowFailed),
		})
		c.Abort()
		return
	}
	//2.调用rpc
	if err := service.FollowUser(followUserRequest, mc.UserID); err != nil {
		if err.Error() == "rpc error: code = Unknown desc = 关注关系数据没有插入进follow表但没有报错" {
			zap.L().Error(getMsg(CodeFollowRepeated), zap.Error(err))
			c.JSON(http.StatusOK, gin.H{
				"status_code": CodeFollowRepeated,
				"status_msg":  getMsg(CodeFollowRepeated),
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
	zap.L().Info(followUserSuccess)
	//3.返回结果
	c.JSON(http.StatusOK, gin.H{
		"status_code": CodeSuccess,
		"status_msg":  getMsg(CodeSuccess),
	})
}

func GetFollowListHandle(c *gin.Context) {
	//1.接收参数
	var getFollowListRequest request.GetFollowListRequest
	if err := c.ShouldBindQuery(&getFollowListRequest); err != nil {
		zap.L().Error(getMsg(CodeInvalidParam), zap.Error(err))
		responseError(c, CodeInvalidParam, "user_list")
		c.Abort()
		return
	}
	//token校验
	mc, err := jwt.ParseToken(getFollowListRequest.Token)
	if err != nil {
		zap.L().Error(getMsg(CodeInvalidToken), zap.Error(err))
		responseError(c, CodeInvalidToken, "user_list")
		c.Abort()
		return
	}
	//2.调用rpc
	data, err := service.GetFollowList(getFollowListRequest, mc.UserID)
	if err != nil {
		zap.L().Error(getMsg(CodeServerBusy), zap.Error(err))
		responseError(c, CodeServerBusy, "user_list")
		c.Abort()
		return
	}
	zap.L().Info(getFollowListSuccess)
	//3.返回结果
	responseSuccess(c, CodeSuccess, data, "user_list")
}

func GetFollowerListHandle(c *gin.Context) {
	//1.接收参数
	var getFollowerListRequest request.GetFollowerListRequest
	if err := c.ShouldBindQuery(&getFollowerListRequest); err != nil {
		zap.L().Error(getMsg(CodeInvalidParam), zap.Error(err))
		responseError(c, CodeInvalidParam, "user_list")
		c.Abort()
		return
	}
	//token校验
	mc, err := jwt.ParseToken(getFollowerListRequest.Token)
	if err != nil {
		zap.L().Error(getMsg(CodeInvalidToken), zap.Error(err))
		responseError(c, CodeInvalidToken, "user_list")
		c.Abort()
		return
	}
	//2.调用rpc
	data, err := service.GetFollowerList(getFollowerListRequest, mc.UserID)
	if err != nil {
		zap.L().Error(getMsg(CodeServerBusy), zap.Error(err))
		responseError(c, CodeServerBusy, "user_list")
		c.Abort()
		return
	}
	zap.L().Info(getFollowerListSuccess)
	//3.返回结果
	responseSuccess(c, CodeSuccess, data, "user_list")

}

func GetFriendListHandle(c *gin.Context) {
	//1.接收参数
	var getFriendListRequest request.GetFriendListRequest
	if err := c.ShouldBindQuery(&getFriendListRequest); err != nil {
		zap.L().Error(getMsg(CodeInvalidParam), zap.Error(err))
		responseError(c, CodeInvalidParam, "user_list")
		c.Abort()
		return
	}
	//token校验
	mc, err := jwt.ParseToken(getFriendListRequest.Token)
	if err != nil {
		zap.L().Error(getMsg(CodeInvalidToken), zap.Error(err))
		responseError(c, CodeInvalidToken, "user_list")
		c.Abort()
		return
	}
	//2.调用rpc
	data, err := service.GetFriendList(getFriendListRequest, mc.UserID)
	if err != nil {
		zap.L().Error(getMsg(CodeServerBusy), zap.Error(err))
		responseError(c, CodeServerBusy, "user_list")
		c.Abort()
		return
	}
	zap.L().Info(getFriendListSuccess)
	//3.返回结果
	responseSuccess(c, CodeSuccess, data, "user_list")
}
