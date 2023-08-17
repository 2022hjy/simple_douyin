package ratelimit

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"sync"
)

const (
	// RATE 每秒生成10个令牌
	RATE = 10
	// CAPACITY 令牌桶的容量
	CAPACITY = 10
)

var buckets = make(map[string]*ratelimit.Bucket)

var MLock sync.Mutex

func RateControl() (gin.HandlerFunc, error) {
	return func(c *gin.Context) {
		CurrentIp := c.ClientIP()

		MLock.Lock()
		if buckets[CurrentIp] == nil {
			buckets[CurrentIp] = ratelimit.NewBucketWithRate(RATE, CAPACITY)
		}
		MLock.Unlock()

		rl := buckets[CurrentIp]
		if rl.TakeAvailable(1) < 1 {
			c.String(http.StatusForbidden, "请勿频繁访问")
			c.Abort()
			return
		}
		c.Next()
	}, nil
}
