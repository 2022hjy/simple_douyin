package mq

import (
	"github.com/streadway/amqp"
	"log"
	"simple_douyin/config"
)

// 路由键常量
const (
	COMMENT_ADD     = "comment_add"
	COMMENT_REMOVE  = "comment_remove"
	FAVORITE_ADD    = "favorite_add"
	FAVORITE_REMOVE = "favorite_remove"
	FOLLOW_ADD      = "follow_add"
	FOLLOW_REMOVE   = "follow_remove"
)

var (
	// 全局的通道和交换器名变量
	ch           *amqp.Channel
	exchangeName string

	// 所有的队列变量
	CommentMQ amqp.Queue
	LikeMQ    amqp.Queue
	FollowMQ  amqp.Queue
)

func init() {
	InitMq()
}

func InitMq() {
	// 建立连接
	conn, err := amqp.Dial(config.MqUrl)
	failOnError(err, "Failed to connect to RabbitMQ")

	// 创建通道
	var Cerr error
	ch, Cerr = conn.Channel()
	failOnError(Cerr, "Failed to open a channel")

	// 声明交换器
	exchangeName = "events"
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

	// 绑定CommentMQ队列到交换器，并设置路由键
	commentRoutingKeys := []string{COMMENT_ADD, COMMENT_REMOVE}
	for _, key := range commentRoutingKeys {
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
	likeRoutingKeys := []string{FAVORITE_ADD, FAVORITE_REMOVE}
	for _, key := range likeRoutingKeys {
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
	followRoutingKeys := []string{FOLLOW_ADD, FOLLOW_REMOVE}
	for _, key := range followRoutingKeys {
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

	//sig := make(chan os.Signal, 1)
	//signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	//<-sig

	closeResources(ch, conn)
}

// SendMessage 用于发送消息到交换器, routingKey为路由键，body为消息内容, 交换器名为events
// 交换机将会根据路由键将消息发送到对应的队列中，无须指定队列名
// 在这里加一个初始化检查
func SendMessage(routingKey string, string2 string) {
	InitMq()
	log.Println("我在 SendMessage 中，准备发送信息")
	if ch == nil {
		log.Println("Channel is nil. Trying to initialize again.")
		if ch == nil {
			log.Fatal("Failed to re-initialize channel.")
			return
		}
	}
	log.Println("发送消息到交换器：", exchangeName, "，路由键：", routingKey)
	err := ch.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(string2),
		})
	if err != nil {
		log.Println("我在 SendMessage 的错误处理中")
		log.Println("Failed to publish a message:", err)
	}
}

func closeResources(ch *amqp.Channel, conn *amqp.Connection) {
	ch.Close()
	conn.Close()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Print("我在 failOnError 中")
		log.Printf("%s: %s", msg, err)
	}
}

//
//func establishConnection(url string, maxRetries int) error {
//	var err error
//	var conn = *new(*amqp.Connection)
//	attempts := 0
//	for attempts < maxRetries {
//		conn, err = amqp.Dial(url)
//		if err == nil {
//			return nil
//		}
//		attempts++
//		log.Printf("Failed to connect to RabbitMQ (attempt %d/%d). Retrying in 5 seconds...", attempts, maxRetries)
//		time.Sleep(5 * time.Second)
//	}
//	return err
//}
//
//func monitorConnection(conn *amqp.Connection, url string) {
//	for {
//		reason, ok := <-conn.NotifyClose(make(chan *amqp.Error))
//		if !ok {
//			log.Println("Connection closed. Reconnecting...")
//			err := establishConnection(url)
//			if err != nil {
//				log.Fatal("Failed to reconnect:", err)
//			}
//		} else {
//			log.Printf("Connection closed due to: %s. Reconnecting...", reason)
//			err := establishConnection(url)
//			if err != nil {
//				log.Fatal("Failed to reconnect:", err)
//			}
//		}
//	}
//}
//
//// 定时进行健康检查
//func healthCheck() {
//	ticker := time.NewTicker(10 * time.Second)
//	defer ticker.Stop()
//
//	for {
//		select {
//		case <-ticker.C:
//			if conn == nil || conn.IsClosed() {
//				log.Println("Connection is not available. Trying to reconnect.")
//				_ = establishConnection(config.MqUrl)
//			}
//
//			if ch == nil || ch.IsClosed() {
//				log.Println("Channel is not available. Trying to reinitialize.")
//				initChannelResources()
//			}
//		}
//	}
//}
