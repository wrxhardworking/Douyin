package main

import (
	"douyin/handle"
	"douyin/middleware"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
)

func RouterInit(r *gin.Engine) {
	// public directory is used to serve static resources
	//静态文件的映射
	r.Static("/static", "./public")
	apiRouter := r.Group("/douyin")
	//基本api
	// basic apis
	apiRouter.GET("/feed/", handle.Feed)
	apiRouter.GET("/user/", middleware.UserAuth(), handle.UserInfo)
	//中间件检查账号密码是否合法
	apiRouter.POST("/user/register/", middleware.Check(), handle.Register)
	apiRouter.POST("/user/login/", handle.Login)
	//加入中间件限制整个post请求的大小，这里设置的是3MB Limit size of POST requests for Gin framework
	apiRouter.POST("/publish/action/", middleware.UserAuth(), limits.RequestSizeLimiter(4<<20), handle.Publish)
	apiRouter.GET("/publish/list/", middleware.UserAuth(), handle.PublishList)
	//  extra apis - I
	apiRouter.POST("/favorite/action/", middleware.UserAuth(), handle.FavoriteAction)
	apiRouter.GET("/favorite/list/", middleware.UserAuth(), handle.FavoriteList)
	apiRouter.POST("/comment/action/", middleware.UserAuth(), handle.CommentAction)
	apiRouter.GET("/comment/list/", middleware.UserAuth(), handle.CommentList)
	//  jextra apis - II
	apiRouter.POST("/relation/action/", middleware.UserAuth(), handle.RelationAction)
	apiRouter.GET("/relation/follow/list/", middleware.UserAuth(), handle.FollowList)
	apiRouter.GET("/relation/follower/list/", middleware.UserAuth(), handle.FollowerList)
	apiRouter.GET("/relation/friend/list/", middleware.UserAuth(), handle.FriendList)
	apiRouter.GET("/message/chat/", middleware.UserAuth(), handle.MessageChat)
	apiRouter.POST("/message/action/", middleware.UserAuth(), handle.MessageAction)

}
