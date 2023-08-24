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
		token := c.Query("token")

		claims, err := AuthCheckToken(token)
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

func AuthCheckToken(token string) (*util.Claims, error) {
	// 没携带token，返回错误
	if len(token) == 0 {
		return nil, errors.New("token is empty")
	}
	return util.ParseToken(token)
}

// AuthWithoutLogin 未登录情况，若携带token,解析用户id放入context;如果没有携带，则将用户id默认为0
func AuthWithoutLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		var userId int64
		claims, err := AuthCheckToken(token)
		if err != nil {
			userId = 0
		} else {
			userId = claims.ID
		}
		c.Set("token_user_id", userId)
		c.Next()
	}
}

// AuthFromBody token存储在body中
func AuthFromBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.PostForm("token")
		claims, err := AuthCheckToken(token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, controller.Response{
				StatusCode: http.StatusUnauthorized,
				StatusMsg:  "Unauthorized",
			})
		} else {
			c.Set("token_user_id", claims.ID)
			c.Next()
		}
	}
}
