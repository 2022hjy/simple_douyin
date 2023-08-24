package controller

import (
	"net/http"
	"strconv"

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

	// 1. 检验参数，如果不符合要求，返回错误信息
	if err := c.ShouldBind(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest,
			Error(http.StatusBadRequest, "Invalid parameter"))
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
	if err := c.ShouldBind(&loginInfo); err != nil {
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
	// 从query参数中获取user_id
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			Error(http.StatusBadRequest, err.Error()))
		return
	}
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
