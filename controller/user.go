package controller

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "simple_douyin/model"
    "simple_douyin/service"
)

// todo 能不能把检验参数和返回结果的代码抽象出来

func Register(c *gin.Context) {
    var registerRequest model.RegisterRequest

    // 1. 检验参数，如果不符合要求，返回错误信息
    if err := c.ShouldBindJSON(&registerRequest); err != nil {
        c.JSON(http.StatusBadRequest,
            model.ErrorResponse(model.ErrorCode, "Invalid parameter"))
        return
    }
    c.JSON(http.StatusOK, service.Register(registerRequest))
}

func Login(c *gin.Context) {
    var loginRequest model.LoginRequest
    // 1. 检验参数，如果不符合要求，返回错误信息
    if err := c.ShouldBindJSON(&loginRequest); err != nil {
        c.JSON(http.StatusBadRequest,
            model.ErrorResponse(model.ErrorCode, "Invalid parameter"))
        return
    }
    // 2. 调用service层的Login方法，返回结果
    c.JSON(http.StatusOK, service.Login(loginRequest))
}

func UserInfo(c *gin.Context) {
    // 从context中获取userId并转为int64类型
    userId := c.GetInt64("userId")
    // todo 是否需要加入检验userId的代码
    c.JSON(http.StatusOK, service.UserInfo(userId))
}
