package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simple_douyin/service"
	"strconv"
	"time"
)

type ChatResponse struct {
	Response
	MessageList []service.Message `json:"message_list"`
}

// MessageChat 消息列表
func MessageChat(c *gin.Context) {
	loginUserId := c.GetInt64("userId")
	toUserId := c.Query("to_user_id")
	preMsgTime := c.Query("pre_msg_time")
	log.Println("preMsgTime", preMsgTime)
	covPreMsgTime, err := strconv.ParseInt(preMsgTime, 10, 64)
	if err != nil {
		log.Println("preMsgTime 参数错误")
		return
	}
	latestTime := time.Unix(covPreMsgTime, 0)
	targetUserId, err := strconv.ParseInt(toUserId, 10, 64)
	if err != nil {
		log.Println("toUserId 参数错误")
		return
	}
	messageService := service.GetMessageServiceInstance()
	messages, err := messageService.MessageChat(loginUserId, targetUserId, latestTime)
	log.Println(messages)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0, StatusMsg: "获取消息成功"}, MessageList: messages})
	}
}

// MessageAction 发送消息
// 主要负责从 HTTP 请求中获取参数（如 fromUserId、toUserId、content 和 actionType），并根据这些参数协调调用服务层来执行相应的操作。
func MessageAction(c *gin.Context) {
	id := c.GetInt64("id")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")
	actionType := c.Query("action_type")
	loginUserId := c.GetInt64("userId")
	targetUserId, err := strconv.ParseInt(toUserId, 10, 64)
	targetActionType, err1 := strconv.ParseInt(actionType, 10, 64)
	if err != nil || err1 != nil {
		log.Println("toUserId/actionType 参数错误")
		return
	}
	messageService := service.GetMessageServiceInstance()
	err = messageService.SendMessage(id, loginUserId, targetUserId, content, targetActionType)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{StatusCode: 1, StatusMsg: "Send Message 接口错误"})
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0})
}
