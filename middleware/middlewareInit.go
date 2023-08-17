package middleware

import (
	"github.com/dvwright/xss-mw"
	"github.com/gin-gonic/gin"
	"simple_douyin/log"
	"simple_douyin/middleware/ratelimit"
)

func InitMiddleware(apiRouter *gin.RouterGroup) { // 请将 "YOUR_ROUTER_TYPE" 替换为您的 apiRouter 类型
	// 初始化 logrus
	log.LogrusInit()

	//初始化 xss
	xssMiddleware := xss.XssMw{}
	apiRouter.Use(xssMiddleware.RemoveXss())

	//初始化 Redis
	redis.InitPool(0)

	// 初始化限流器
	rateControlHandler, err := ratelimit.RateControl()
	if err != nil {
		log.Log.Fatalf("Error setting up rate control middleware: %v", err)
	}

	// 检查rateControlHandler是否正确初始化
	if rateControlHandler == nil {
		log.Log.Fatalf("Error initializing rate control handler")
	}

	// 将速率限制中间件添加到路由器
	apiRouter.Use(rateControlHandler)
}
