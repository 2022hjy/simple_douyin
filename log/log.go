package log

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)
}

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		passTime := time.Since(startTime)
		log.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"requestUri": c.Request.RequestURI,
			"status":     c.Writer.Status(),
			"passTime":   passTime,
			}).Info("request info")
		}
	}
}
