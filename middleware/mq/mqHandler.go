package mq

import (
	"github.com/streadway/amqp"
	"log"
)

// MessageHandler 定义消息处理函数
func handleMessages(queueName string, ch *amqp.Channel, handler func(string, string)) {
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		routingKey := d.RoutingKey
		body := string(d.Body)
		handler(routingKey, body)
	}
}

// 各自处理消息的具体函数的实现
func handleCommentMessage(routingKey, body string) {
	switch routingKey {
	case "comment_add":
		AddComment(body)
	case "comment_remove":
		DeleteComment(body)
	default:
		log.Println("Unknown routing key:", routingKey)
	}
}

func handleLikeMessage(routingKey, body string) {
	switch routingKey {
	case "favorite_add":
		AddLike(body)
	case "favorite_remove":
		RemoveLike(body)
	default:
		log.Println("Unknown routing key:", routingKey)
	}
}

func handleFollowMessage(routingKey, body string) {
	switch routingKey {
	case "follow_add":
		AddFollow(body)
	case "follow_remove":
		RemoveFollow(body)
	default:
		log.Println("Unknown routing key:", routingKey)
	}
}
