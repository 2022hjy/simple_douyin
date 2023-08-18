package mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"simple_douyin/config"
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

	// 建立连接
	conn, err := amqp.Dial(config.MqUrl)
	if err != nil {
		t.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	// 创建通道
	ch, err := conn.Channel()
	if err != nil {
		t.Fatal("Failed to open a channel:", err)
	}
	defer ch.Close()

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
		err = ch.Publish(
			"events",       // exchange
			msg.routingKey, // routing key
			false,          // mandatory
			false,          // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg.body),
			})
		if err != nil {
			t.Fatal("Failed to publish a message:", err)
		}
	}

	// 给系统一些时间来处理消息
	time.Sleep(1 * time.Second)

}
