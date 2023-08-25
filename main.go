package main

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"simple_douyin/middleware"
	"simple_douyin/middleware/database"
)

func main() {
	go service.RunMessageServer()

	r := gin.Default()

	ApiRouter := r.Group("/douyin")

	middleware.InitMiddleware(ApiRouter)

	database.Init()

	InitRouter(ApiRouter)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
