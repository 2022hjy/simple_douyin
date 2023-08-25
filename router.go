package main

import (
	//"github.com/RaymondCode/simple-demo/controller"
	"github.com/gin-gonic/gin"
	"simple_douyin/controller"
	"simple_douyin/middleware/corsUtils"
	"simple_douyin/middleware/jwt"
)

func InitRouter(apiRouter *gin.RouterGroup) *gin.RouterGroup {
	// public directory is used to serve static resources
	r := gin.Default()
	r.Static("/static", "./public")

	apiRouter.GET("/feed/", jwt.AuthWithoutLogin(), controller.Feed)

	// 用户相关路由
	rUser := apiRouter.Group("/user")
	{
		rUser.POST("/register/", controller.Register)
		rUser.POST("/login/", controller.Login)
		rUser.GET("/", jwt.Auth(), controller.UserInfo)
	}

	// 互动相关路由
	rInteraction := apiRouter.Group("/relation")
	{
		rInteraction.POST("/action/", controller.RelationAction)
		rInteraction.GET("/follow/list/", controller.Followlist)
		rInteraction.GET("/follower/list/", controller.FollowerList)
		rInteraction.GET("/friend/list/", controller.FriendList)
	}

	// 点赞相关路由
	rFavorite := apiRouter.Group("/favorite")
	{
		rFavorite.POST("/action/", jwt.Auth(), controller.FavoriteAction)
		rFavorite.GET("/list/", jwt.AuthWithoutLogin(), controller.FavoriteList)
	}

	// 视频相关路由
	rVideo := apiRouter.Group("/publish")
	{
		rVideo.POST("/action/", controller.Publish)
		rVideo.GET("/list/", jwt.AuthWithoutLogin(), controller.PublishList)
	}

	// 评论相关路由
	rComment := apiRouter.Group("/comment")
	{
		rComment.POST("/action/", controller.CommentAction)
		rComment.GET("/list/", controller.CommentList)
	}

	// 私信相关路由
	rMessage := apiRouter.Group("/message")
	{
		rMessage.POST("/action/", controller.MessageAction)
		rMessage.GET("/chat/", controller.MessageChat)
	}

	// 允许跨域
	apiRouter.Use(corsUtils.AllowAllCORS())
	return apiRouter
}
