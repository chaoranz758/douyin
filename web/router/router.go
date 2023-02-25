package router

import (
	"douyin/web/handler"
	"douyin/web/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init() *gin.Engine {
	r := gin.New()
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	r.Use(middleware.Cors())
	//r.Use(middleware.CSRFMiddle())
	//r.Use(middleware.Ja3())
	//健康检查接口
	r.GET("/checkHealth", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"success": "success",
		})
	})
	//用户相关
	//用户注册
	r.POST("/douyin/user/register/", handler.UserRegisterHandle)
	//用户登录
	r.POST("/douyin/user/login/", handler.UserLoginHandle)
	//用户信息
	r.GET("/douyin/user/", handler.GetUserInfoHandle)

	//视频相关
	//投稿接口
	r.POST("/douyin/publish/action/", handler.PublishVideoHandle)
	//发布列表
	r.GET("/douyin/publish/list/", handler.GetPublishVideoHandle)
	//视频流接口
	r.GET("/douyin/feed/", handler.GetVideoListHandle)

	//视频点赞相关
	//赞操作
	r.POST("/douyin/favorite/action/", handler.FavoriteVideoHandle)
	//喜欢列表
	r.GET("/douyin/favorite/list/", handler.GetFavoriteVideoListHandle)

	//评论相关
	//评论操作
	r.POST("/douyin/comment/action/", handler.CommentVideoHandle)
	//评论列表
	r.GET("/douyin/comment/list/", handler.GetCommentVideoListHandle)

	//关注相关
	r.POST("/douyin/relation/action/", handler.FollowUserHandle)
	//关注列表
	r.GET("/douyin/relation/follow/list/", handler.GetFollowListHandle)
	//粉丝列表
	r.GET("/douyin/relation/follower/list/", handler.GetFollowerListHandle)
	//好友列表
	r.GET("/douyin/relation/friend/list/", handler.GetFriendListHandle)

	//消息相关
	//发送消息
	r.POST("/douyin/message/action/", handler.SendMessageHandle)
	//聊天记录
	r.GET("/douyin/message/chat/", handler.GetMessageHandle)
	return r
}
