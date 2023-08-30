package dao

import (
	"log"
	"simple_douyin/middleware/database"
	"testing"
	"time"
)

// 测试保存消息  GET
func TestSaveMessage(t *testing.T) {
	database.Init() // 初始化数据库连接
	message := Message{
		FromUserID: 1,
		ToUserID:   2,
		Content:    "测试消息",
		CreateTime: time.Now(),
	}
	resultMessage, err := SendMessage(message)
	if err == nil {
		log.Println("SaveMessage 测试成功！")
	}
	log.Println(resultMessage)
}

// TestGetUserBasicInfoByName 测试获取聊天记录 GET
func TestMessageChat(t *testing.T) {
	database.Init()
	loginUserId := 1
	targetUserId := 2
	LastTime := time.Now()
	messages, err := MessageChat(int64(loginUserId), int64(targetUserId), LastTime)
	if err != nil {
		log.Println("MessageChat 测试失败")
	}
	for _, msg := range messages {
		log.Println(msg)
	}
}

// TestGetUserBasicInfoByName 测试 获取最新一条消息 GET
func TestLatestMessage(t *testing.T) {
	database.Init()
	loginUserId := 1
	targetUserId := 2
	message, err := LatestMessage(int64(loginUserId), int64(targetUserId))
	if err != nil {
		log.Println("LatestMessage 测试失败")
	}
	log.Println(message)
}
