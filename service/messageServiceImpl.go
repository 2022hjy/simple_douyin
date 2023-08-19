package service

import (
	"fmt"
	"log"
	"simple_douyin/config"
	"simple_douyin/dao"
	"time"
)

type MessageServiceImpl struct {
}

// 单例模式简化代码
var messageServiceImpl = &MessageServiceImpl{}

// SendMessage 处理err的，有用到dao.SendMessage的方法
func (messageService *MessageServiceImpl) SendMessage(id int64, fromUserId int64, toUserId int64, content string, actionType int64) error {
	var err error
	switch actionType {
	case 1:
		err = dao.SendMessage(id, fromUserId, toUserId, content, actionType)
	default:
		err = fmt.Errorf("未定义的 actionType=%d", actionType)
		log.Println(err)
	}
	return err
}

func (messageService *MessageServiceImpl) MessageChat(loginUserId int64, targetUserId int64, latestTime time.Time) ([]Message, error) {
	messages := make([]Message, 0, config.VIDEO_INIT_NUM)
	//MessageChat:当前登录用户和其他指定用户的聊天记录
	plainMessages, err := dao.MessageChat(loginUserId, targetUserId, latestTime)
	if err != nil {
		log.Println("MessageChat Service:", err)
		return nil, err
	}
	// 将原始消息数据转换为响应所需的消息格式，并将这些消息添加到 messages 切片中
	err = messageService.getRespMessage(&messages, &plainMessages)
	if err != nil {
		log.Println("getRespMessage:", err)
		return nil, err
	}
	return messages, nil
}

func (messageService *MessageServiceImpl) LatestMessage(loginUserId int64, targetUserId int64) (LatestMessage, error) {
	plainMessage, err := dao.LatestMessage(loginUserId, targetUserId)
	if err != nil {
		log.Println("LatestMessage Service:", err)
		return LatestMessage{}, err
	}
	var latestMessage LatestMessage
	latestMessage.Message = plainMessage.MsgContent
	if plainMessage.UserId == loginUserId {
		// 最新一条消息是当前登录用户发送的
		latestMessage.MsgType = 1
	} else {
		// 最新一条消息是当前好友发送的
		latestMessage.MsgType = 0
	}
	return latestMessage, nil
}

// 返回 message list 接口所需的 Message 结构体
//
//	下一个函数的list类型，获取回复消息
func (messageService *MessageServiceImpl) getRespMessage(messages *[]Message, plainMessages *[]dao.Message) error {
	for _, tmpMessage := range *plainMessages {
		var message Message
		err := messageService.combineMessage(&message, &tmpMessage)
		if err != nil {
			return err
		}
		*messages = append(*messages, message)
	}
	return nil
}

// 把service的message转成dao的message
func (messageService *MessageServiceImpl) combineMessage(message *Message, plainMessage *dao.Message) error {
	message.Id = plainMessage.Id
	message.UserId = plainMessage.UserId
	message.ReceiverId = plainMessage.ReceiverId
	message.MsgContent = plainMessage.MsgContent
	message.CreatedAt = plainMessage.CreatedAt.Unix()
	return nil
}
