package mq

import (
	"fmt"
	"testing"
	"time"
)

func TestMq(t *testing.T) {
	// 初始化消息队列系统
	go InitMq()

	//log.Log.Info("TestMq")
	fmt.Print("TestMq")

	// 给系统一些时间来初始化
	time.Sleep(1 * time.Second)

	// 定义测试消息
	messages := []struct {
		routingKey string
		body       string
	}{
		{"comment_add", "Test add comment"},
		{"comment_remove", "Test remove comment"},
		{"like_add", "Test add like"},
		{"like_remove", "Test remove like"},
		{"follow_add", "Test add follow"},
		{"follow_remove", "Test remove follow"},
	}

	// 向交换器发送测试消息
	for _, msg := range messages {
		SendMessage(msg.routingKey, msg.body)
	}

	// 给系统一些时间来处理消息
	time.Sleep(1 * time.Second)

}
