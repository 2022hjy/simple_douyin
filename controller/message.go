package controller

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"simple_douyin/dao"
)

type ChatResponse struct {
	Response
	MessageList []dao.Message `json:"message_list"`
}

var ()

//func init() {
//	messageService = service.GetMessageServiceInstance()
//}

func MessageRespondWithError(c *gin.Context, statusCode int, errMsg string) {
	c.JSON(http.StatusOK, Response{StatusCode: int32(statusCode), StatusMsg: errMsg})
}

// MessageAction 发送消息
func MessageAction(c *gin.Context) {
	FromUserID := c.GetInt64("token_user_id")
	//ToUserID := c.GetInt64("to_user_id")
	ToUserID, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		MessageRespondWithError(c, -1, "MessageAction Error: "+err.Error())
	}
	Content := c.Query("content")

	log.Println("FromUserID:", FromUserID)
	log.Println("ToUserID:", ToUserID)
	log.Println("Content:", Content)

	if err := messageService.SendMessage(FromUserID, ToUserID, Content); err != nil {
		MessageRespondWithError(c, -1, "MessageAction Error: "+err.Error())
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "MessageAction Success!"})
}

// MessageChat 消息列表
func MessageChat(c *gin.Context) {
	FromUserID := c.GetInt64("token_user_id")
	ToUserID, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		MessageRespondWithError(c, -1, "MessageChat Error: "+err.Error())
	}
	preMsgTime := c.Query("pre_msg_time")
	log.Println("pre_msg_time", preMsgTime)
	covPreMsgTime, err := strconv.ParseInt(preMsgTime, 10, 64)
	if err != nil {
		log.Println("preMsgTime 参数错误")
		return
	}
	latestTime := time.Unix(covPreMsgTime, 0)
	messages, err := messageService.MessageChat(FromUserID, ToUserID, latestTime)
	log.Println(messages)
	if err != nil {
		c.JSON(http.StatusOK, ChatResponse{
			Response: Response{StatusCode: -1, StatusMsg: "MessageChat Error"},
		})
	} else {
		r := ChatResponse{
			Response: Response{
				StatusCode: 0, StatusMsg: "MessageChat Success!"}, MessageList: messages}
		log.Println("r:", r)
		c.JSON(http.StatusOK, r)
	}
}
