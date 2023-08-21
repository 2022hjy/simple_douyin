package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"simple_douyin/service"
)

var (
	userService = service.NewUserServiceInstance()
)

type LoginRequest struct {
	UserName string
	Password string
}

type RegisterRequest struct {
	UserName string
	Password string
}

type LoginResponse struct {
	Response
	*service.Credential
}

type InfoResponse struct {
	Response
	*service.UserInfo
}

func Register(c *gin.Context) {
	var registerRequest RegisterRequest

	// 1. 检验参数，如果不符合要求，返回错误信息
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest,
			Error(http.StatusBadRequest, "Invalid parameter"))
		return
	}
	credential, err := userService.Register(registerRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			Error(http.StatusBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, LoginResponse{
		Response:   Success(),
		Credential: credential,
	})
}

func Login(c *gin.Context) {
	var loginRequest LoginRequest
	// 1. 检验参数，如果不符合要求，返回错误信息
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest,
			Error(http.StatusBadRequest, "Invalid parameter"))
		return
	}
	// 2. 调用service层的Login方法，返回结果
	credential, err := userService.Login(loginRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			Error(http.StatusBadRequest, err.Error()))
		return
	}
	// 3. 返回结果
	c.JSON(http.StatusOK, LoginResponse{
		Response:   Success(),
		Credential: credential,
	})
}

func UserInfo(c *gin.Context) {
	// 从context中获取userId并转为int64类型
	userId := c.GetInt64("userId")
	userInfo, err := userService.QuerySelfInfo(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			Error(http.StatusBadRequest, err.Error()))
		return
	}
	c.JSON(http.StatusOK, InfoResponse{
		Response: Success(),
		UserInfo: userInfo,
	})
}
