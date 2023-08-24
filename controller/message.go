package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simple_douyin/dao"
	"simple_douyin/service"
)

type ChatResponse struct {
	Response
	MessageList []dao.Message `json:"message_list"`
}

var (
	messageService *service.MessageServiceImpl
)

func init() {
	messageService = service.GetMessageServiceInstance()
}

func MessageRespondWithError(c *gin.Context, statusCode int, errMsg string) {
	c.JSON(http.StatusOK, Response{StatusCode: int32(statusCode), StatusMsg: errMsg})
}

// MessageAction 发送消息
func MessageAction(c *gin.Context) {
	FromUserID := c.GetInt64("from_user_id")
	ToUserID := c.GetInt64("to_user_id")
	Content := c.Query("content")
	err := messageService.SendMessage(FromUserID, ToUserID, Content)
	if err != nil {
		MessageRespondWithError(c, -1, "MessageAction Error: "+err.Error())
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "MessageAction Success!"})
}

// MessageChat 消息列表
func MessageChat(c *gin.Context) {
	FromUserID := c.GetInt64("userId")
	ToUserID := c.GetInt64("to_user_id")
	preMsgTime := c.Query("content")
	log.Println("content", preMsgTime)
	messages, err := messageService.MessageChat(FromUserID, ToUserID)
	log.Println(messages)
	if err != nil {
		c.JSON(http.StatusOK, ChatResponse{
			Response: Response{StatusCode: -1, StatusMsg: "MessageChat Error"},
		})
	} else {
		c.JSON(http.StatusOK, ChatResponse{Response: Response{StatusCode: 0, StatusMsg: "MessageChat Success!"}, MessageList: messages})
	}
}
