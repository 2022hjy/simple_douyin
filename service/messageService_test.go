package service

import (
	"fmt"
	"log"
	"simple_douyin/middleware/database"
	"simple_douyin/middleware/redis"
	"testing"
)

// 测试发送消息  GET
func TestMessageServiceImpl_SendMessage(t *testing.T) {
	redis.InitRedis()
	database.Init()
	err := messageServiceImpl.SendMessage(2, 26, "226-1")
	if err == nil {
		log.Println("SendMessage Service 正常")
	}
}

// 测试获取聊天记录 GET
func TestMessageServiceImpl_MessageChat(t *testing.T) {
	database.Init()
	chat, _ := messageServiceImpl.MessageChat(2, 26)
	for _, msg := range chat {
		log.Println(fmt.Sprintf("%+v", msg))
	}
}

// 测试获取最新一条聊天记录
func TestMessageServiceImpl_LatestMessage(t *testing.T) {
	redis.InitRedis()
	database.Init()
	latestMessage, _ := messageServiceImpl.LatestMessage(2, 26)
	log.Println(fmt.Sprintf("%+v", latestMessage))
}
