package dao

import (
	"fmt"
	"log"
	"simple_douyin/middleware/database"
	"testing"
	"time"
)

// 测试保存消息  GET
func TestSaveMessage(t *testing.T) {
	database.Init()

	message := Message{
		Id:         1,
		UserId:     1,
		ReceiverId: 2,
		ActionType: 1,
		MsgContent: "测试消息",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	err := SaveMessage(message)
	if err == nil {
		log.Println("SaveMessage 测试成功！")
	}
}

// TestGetUserBasicInfoById 测试发送消息 GET
func TestSendMessage(t *testing.T) {
	database.Init()
	id := 9
	fromUserId := 1
	toUserId := 7

	err := SendMessage(int64(id), int64(fromUserId), int64(toUserId), fmt.Sprintf("我是 userId=%d,发送消息给 userId=%d", fromUserId, toUserId), 1)
	if err == nil {
		log.Println("SendMessage 测试成功！")
	}
}

// TestGetUserBasicInfoByName 测试获取聊天记录 GET
func TestMessageChat(t *testing.T) {
	database.Init()
	loginUserId := 1
	targetUserId := 7
	lastTime := time.Date(2023, time.August, 10, 10, 0, 0, 0, time.UTC)
	messages, err := MessageChat(int64(loginUserId), int64(targetUserId), lastTime)
	if err != nil {
		log.Println("MessageChat 测试失败")
	}
	for _, msg := range messages {
		log.Println(fmt.Sprintf("%d -> %d: %s (sendTime:%v)", msg.UserId, msg.ReceiverId, msg.MsgContent, msg.CreatedAt))
		log.Println("测试成功！")
		log.Println(msg.MsgContent)
	}
}

// TestGetUserBasicInfoByName 测试 获取最新一条消息 GET
func TestLatestMessage(t *testing.T) {
	database.Init()
	loginUserId := 1
	targetUserId := 7
	message, err := LatestMessage(int64(loginUserId), int64(targetUserId))
	if err != nil {
		log.Println("LatestMessage 测试失败")
	}
	log.Println(fmt.Sprintf("%d -> %d 的最新一条消息记录：%s", message.UserId, message.ReceiverId, message.MsgContent))
}
