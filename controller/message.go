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

func MessageRespondWithError(c *gin.Context, statusCode int, errMsg string) {
	c.JSON(http.StatusOK, Response{StatusCode: int32(statusCode), StatusMsg: errMsg})
}

// MessageAction 发送消息
func MessageAction(c *gin.Context) {
	FromUserID := c.GetInt64("token_user_id")
	ToUserID, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)

	// 前缀树过滤敏感词
	//TODO 要刷新才会显示过滤后的消息bug
	if filteredContentValue, exists := c.Get("filteredContent"); exists {
		filteredContent := filteredContentValue.(string) // 类型断言
		log.Println("过滤后的filteredContent:", filteredContent)
		err = messageService.SendMessage(FromUserID, ToUserID, filteredContent)
	} else {
		content := c.Query("content")
		err = messageService.SendMessage(FromUserID, ToUserID, content)
	}

	if err != nil {
		MessageRespondWithError(c, -1, "MessageAction Error: "+err.Error())
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "MessageAction Success!"})
}

func MessageChat(c *gin.Context) {
	FromUserID := c.GetInt64("token_user_id")
	ToUserID, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		MessageRespondWithError(c, -1, "MessageChat Error: "+err.Error())
		return
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
	if err != nil {
		c.JSON(http.StatusOK, ChatResponse{
			Response: Response{StatusCode: -1, StatusMsg: "MessageChat Error"},
		})
		return
	} else {
		c.JSON(http.StatusOK, ChatResponse{
			Response:    Response{StatusCode: 0, StatusMsg: "MessageChat Success!"},
			MessageList: messages})
	}
}
