package mq

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"simple_douyin/config"
	"syscall"
)

var (
	// 所有的队列变量
	CommentMQ amqp.Queue
	LikeMQ    amqp.Queue
	FollowMQ  amqp.Queue
)

func InitMq() {
	// 建立连接
	conn, err := amqp.Dial(config.MqUrl)
	failOnError(err, "Failed to connect to RabbitMQ")

	// 创建通道
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	// 声明交换器
	exchangeName := "events"
	err = ch.ExchangeDeclare(
		exchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// 声明队列
	queueName := "CommentMQ"
	CommentMQ, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 声明队列
	queueName = "LikeMQ"
	LikeMQ, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 声明队列
	queueName = "FollowMQ"
	FollowMQ, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// 绑定LikeMQ队列到交换器，并设置路由键
	routingKeys := []string{"comment_add", "comment_remove"}
	for _, key := range routingKeys {
		err = ch.QueueBind(
			CommentMQ.Name, // queue name
			key,            // routing key
			exchangeName,   // exchange
			false,          // no-wait
			nil,            // arguments
		)
		failOnError(err, "Failed to bind a queue")
	}

	// 绑定LikeMQ队列到交换器，并设置路由键
	routingKeys = []string{"like_add", "like_remove"}
	for _, key := range routingKeys {
		err = ch.QueueBind(
			LikeMQ.Name,  // queue name
			key,          // routing key
			exchangeName, // exchange
			false,        // no-wait
			nil,          // arguments
		)
		failOnError(err, "Failed to bind a queue")
	}

	// 绑定FollowMQ队列到交换器，并设置路由键
	routingKeys = []string{"follow_add", "follow_remove"}
	for _, key := range routingKeys {
		err = ch.QueueBind(
			FollowMQ.Name, // queue name
			key,           // routing key
			exchangeName,  // exchange
			false,         // no-wait
			nil,           // arguments
		)
		failOnError(err, "Failed to bind a queue")
	}

	// 处理CommentMQ队列中的消息
	go handleMessages(CommentMQ.Name, ch, handleCommentMessage)
	// 处理LikeMQ队列中的消息
	go handleMessages(LikeMQ.Name, ch, handleLikeMessage)
	// 处理FollowMQ队列中的消息
	go handleMessages(FollowMQ.Name, ch, handleFollowMessage)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	closeResources(ch, conn)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func closeResources(ch *amqp.Channel, conn *amqp.Connection) {
	ch.Close()
	conn.Close()
}
