package jwt

import (
	"errors"
	"net/http"

	"github.com/RaymondCode/simple-demo/controller"
	"github.com/gin-gonic/gin"
	"simple_douyin/util"
)

// Auth 鉴权中间件，token存储在query中
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果请求为Get请求，则从query中获取token
		// 如果请求为Post请求，则从body中获取token
		token := c.Query("token")
		if len(token) == 0 {
			token = c.PostForm("token")
		}
		claims, err := authCheckToken(token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, controller.Response{
				StatusCode: http.StatusUnauthorized,
				StatusMsg:  err.Error(),
			})
		} else {
			c.Set("token_user_id", claims.ID)
			c.Next()
		}
	}
}

func authCheckToken(token string) (*util.Claims, error) {
	// 没携带token，返回错误
	if len(token) == 0 {
		//return nil, error(nil)
		return nil, errors.New("token is missing")
	}
	return util.ParseToken(token)
}

// AuthWithoutLogin 未登录情况，若携带token,解析用户id放入context;如果没有携带，则将用户id默认为0
func AuthWithoutLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果请求为Get请求，则从query中获取token
		// 如果请求为Post请求，则从body中获取token
		token := c.Query("token")
		if len(token) == 0 {
			token = c.PostForm("token")
		}
		var userId int64
		claims, err := authCheckToken(token)
		if err == nil {
			userId = claims.ID
		}
		c.Set("token_user_id", userId)
		c.Next()
	}
}
