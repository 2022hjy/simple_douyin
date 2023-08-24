package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"simple_douyin/service"
)

var (
	userService = service.NewUserServiceInstance()
)

type LoginResponse struct {
	Response
	*service.Credential
}

type InfoResponse struct {
	Response
	*service.UserInfo
}

func Register(c *gin.Context) {
	var loginInfo service.LoginInfo

	// 从查询字符串中获取参数
	loginInfo.UserName = c.Query("username")
	loginInfo.Password = c.Query("password")

	// 1. 检验参数，如果不符合要求，返回错误信息
	if loginInfo.UserName == "" || loginInfo.Password == "" {
		c.JSON(http.StatusBadRequest,
			Error(http.StatusBadRequest, "Invalid parameter1"))
		return
	}

	credential, err := userService.Register(loginInfo)
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
	var loginInfo service.LoginInfo
	// 1. 检验参数，如果不符合要求，返回错误信息
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest,
			Error(http.StatusBadRequest, "Invalid parameter"))
		return
	}
	// 2. 调用service层的Login方法，返回结果
	credential, err := userService.Login(loginInfo)
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
