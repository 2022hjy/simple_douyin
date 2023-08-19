package service

import (
	"fmt"
	"log"
	"simple_douyin/middleware/database"
	"testing"
	"time"
)

// 测试发送消息  GET
func TestMessageServiceImpl_SendMessage(t *testing.T) {
	database.Init()
	err := messageServiceImpl.SendMessage(10, 2, 8, "2 号用户发送消息给 8 号用户", 1)
	if err == nil {
		log.Println("SendMessage Service 正常")
	}
}

// 测试获取聊天记录 GET
func TestMessageServiceImpl_MessageChat(t *testing.T) {
	database.Init()
	lastTime := time.Date(2023, time.August, 10, 10, 0, 0, 0, time.UTC)
	chat, _ := messageServiceImpl.MessageChat(8, 2, lastTime)
	for _, msg := range chat {
		log.Println(fmt.Sprintf("%+v", msg))
	}
}

// 测试获取最新一条聊天记录 GET
func TestMessageServiceImpl_LatestMessage(t *testing.T) {
	database.Init()
	message, _ := messageServiceImpl.LatestMessage(2, 8)
	log.Println(fmt.Sprintf("%+v", message))
}
