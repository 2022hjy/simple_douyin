package main

import (
	//"github.com/RaymondCode/simple-demo/controller"
	"github.com/gin-gonic/gin"
	"log"
	"simple_douyin/controller"
	"simple_douyin/middleware/corsUtils"
	"simple_douyin/middleware/jwt"
	"simple_douyin/middleware/word_filter"
)

func InitRouter(apiRouter *gin.RouterGroup) *gin.RouterGroup {
	// public directory is used to serve static resources
	r := gin.Default()
	r.Static("/static", "./public")

	WordFilter, err := word_filter.NewWordFilterMiddleware("middleware/word_filter/sensitive_words.txt")
	if err != nil {
		log.Printf("Error setting up word filter middleware: %v\n", err)
	}

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
		rInteraction.POST("/action/", jwt.Auth(), controller.RelationAction)
		rInteraction.GET("/follow/list/", jwt.Auth(), controller.Followlist)
		rInteraction.GET("/follower/list/", jwt.Auth(), controller.FollowerList)
		rInteraction.GET("/friend/list/", jwt.Auth(), controller.FriendList)
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
		rVideo.POST("/action/", jwt.Auth(), controller.Publish)
		rVideo.GET("/list/", jwt.AuthWithoutLogin(), controller.PublishList)
	}

	// 评论相关路由
	rComment := apiRouter.Group("/comment")
	{
		rComment.POST("/action/", jwt.Auth(), WordFilter, controller.CommentAction)
		rComment.GET("/list/", jwt.Auth(), controller.CommentList)
	}

	// 私信相关路由
	rMessage := apiRouter.Group("/message")
	{
		rMessage.POST("/action/", jwt.Auth(), controller.MessageAction)
		rMessage.GET("/chat/", jwt.Auth(), controller.MessageChat)
	}

	// 允许跨域
	apiRouter.Use(corsUtils.AllowAllCORS())
	log.Printf("Init Router success ！")
	return apiRouter
}
