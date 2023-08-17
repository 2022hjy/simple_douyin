// log.go
package log

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// Log 是 logrus 的实例
var Log = logrus.New()

// LogrusInit 初始化logrus
func LogrusInit() {
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.DebugLevel)
}

// Logging 日志中间件
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		passTime := time.Since(startTime)
		Log.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"requestUri": c.Request.RequestURI,
			"status":     c.Writer.Status(),
			"passTime":   passTime,
		}).Info("request info")
	}
}
