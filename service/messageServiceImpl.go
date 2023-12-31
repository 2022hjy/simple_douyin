package service

import (
	"encoding/json"
	"fmt"
	"log"
	"simple_douyin/config"
	"simple_douyin/dao"
	"simple_douyin/middleware/redis"
	"strconv"
	"sync"
	"time"
)

type MessageServiceImpl struct{}

var (
	messageServiceImpl *MessageServiceImpl
	messageServiceOnce sync.Once
)

func GetMessageServiceInstance() *MessageServiceImpl {
	messageServiceOnce.Do(func() {
		messageServiceImpl = &MessageServiceImpl{}
	})
	return messageServiceImpl
}

func (c *MessageServiceImpl) SendMessage(fromUserId int64, toUserId int64, content string) (err error) {
	message := dao.Message{
		FromUserID: fromUserId,
		ToUserID:   toUserId,
		Content:    content,
		CreateTime: time.Unix(time.Now().Unix(), 0),
	}
	LaseMessage, err := dao.SendMessage(message)

	// 在发送消息的时候就存入redis 对fromUserId和toUserId进行排序 保证LastMessage的redis一致性
	if fromUserId > toUserId {
		fromUserId, toUserId = toUserId, fromUserId
	}
	updateLastMessageRedis(fromUserId, toUserId, LaseMessage)

	return nil
}

func (c *MessageServiceImpl) MessageChat(loginUserId int64, targetUserId int64, latestTime time.Time) ([]Message, error) {
	messages := make([]Message, 0, config.MessageInitNum)
	daoMessages, err := dao.MessageChat(loginUserId, targetUserId, latestTime)
	if err != nil {
		log.Println("MessageChat Service出错:", err.Error())
		return nil, err
	}
	for _, tmpMessage := range daoMessages {
		var message Message
		message.Id = tmpMessage.Id
		message.FromUserID = tmpMessage.FromUserID
		message.ToUserID = tmpMessage.ToUserID
		message.Content = tmpMessage.Content
		message.CreateTime = tmpMessage.CreateTime.Unix()
		messages = append(messages, message)
	}
	return messages, nil
}

//======================   LatestMessage   =========================

func (c *MessageServiceImpl) LatestMessage(loginUserId int64, targetUserId int64) (LatestMessage, error) {
	if loginUserId > targetUserId {
		loginUserId, targetUserId = targetUserId, loginUserId
	}
	lastMessage, err := c.getLastMessageFromRedis(loginUserId, targetUserId)
	if err != nil {
		return LatestMessage{}, err
	}

	return lastMessage, nil
}

func (c *MessageServiceImpl) getLastMessageFromDB(loginUserId int64, targetUserId int64) (LatestMessage, error) {
	plainMessage, err := dao.LatestMessage(loginUserId, targetUserId)
	if err != nil {
		log.Println("LatestMessage Service:", err)
		return LatestMessage{}, err
	}
	var latestMessage LatestMessage
	latestMessage.Message = plainMessage.Content
	if plainMessage.FromUserID == loginUserId {
		// 最新一条消息是当前登录用户发送的
		latestMessage.MsgType = 1
	} else {
		latestMessage.MsgType = 0
	}
	return latestMessage, nil
}

func (c *MessageServiceImpl) getLastMessageFromRedis(loginUserId int64, targetUserId int64) (LatestMessage, error) {
	var latestMessage LatestMessage
	uId := fmt.Sprintf("%s%d-%d", config.UserAllId_Message_KEY_PREFIX, loginUserId, targetUserId)
	UMClient := redis.Clients.UserAllId_MessageR
	if UMClient == nil {
		return latestMessage, fmt.Errorf("redis client is nil")
	}

	messageJson, err := redis.GetValue(UMClient, uId)
	if err != nil {
		return latestMessage, fmt.Errorf("get redis value failed: %v", err)
	}

	if unmarshalErr := json.Unmarshal([]byte(messageJson), &latestMessage); unmarshalErr != nil {
		return latestMessage, fmt.Errorf("unmarshal message failed: %v", unmarshalErr)
	}

	_, err = redis.GetKeysAndUpdateExpiration(UMClient, uId)
	if err != nil {
		return latestMessage, fmt.Errorf("update redis expiration failed: %v", err)
	}

	log.Println("redis读取Message成功！")
	return latestMessage, nil
}

// updateLastMessageRedis 更新最新一条消息redis
func updateLastMessageRedis(loginUserId int64, targetUserId int64, message dao.Message) {
	// UserAllId --> message
	UMClient := redis.Clients.UserAllId_MessageR
	if UMClient == nil {
		log.Fatalf("redis client is nil")
		return
	}

	uId := config.UserAllId_Message_KEY_PREFIX + strconv.FormatInt(loginUserId, 10) + "-" + strconv.FormatInt(targetUserId, 10)

	LatestMessage := LatestMessage{
		Message: message.Content,
	}
	if message.FromUserID == loginUserId {
		LatestMessage.MsgType = 1
	} else {
		LatestMessage.MsgType = 0
	}

	messageJson, serializationErr := json.Marshal(LatestMessage)
	if serializationErr != nil {
		log.Fatalf("jsonfiy messag failed, err:%v\n", messageJson)
		return
	}
	err := redis.SetValueWithRandomExp(UMClient, uId, string(messageJson))
	if err != nil {
		log.Fatalf("set redis failed, err:%v\n", err)
		return
	}
	log.Println("redis缓存Message成功！")
}
