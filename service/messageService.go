package service

import (
	"simple_douyin/dao"
)

// Message CreateTime字段类型为time.Time，但是在响应中需要返回时间戳，所以需要定义一个新的结构体
type Message struct {
	Id         int64  `json:"id"`           // 消息ID
	FromUserID int64  `json:"from_user_id"` // 发送消息的用户ID
	ToUserID   int64  `json:"to_user_id"`   // 接收消息的用户ID
	Content    string `json:"content"`      // 消息内容
	CreateTime int64  `json:"create_time"`  // 创建时间
}

// LatestMessage 提供给用户好友列表接口的最新一条聊天信息, msgType 消息类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
type LatestMessage struct {
	Message string `json:"message"`
	MsgType int64  `json:"msg_type"`
}

type MessageService interface {

	/*
		业务所需服务接口
	*/
	// SendMessage 发送消息服务
	SendMessage(id int64, fromUserId int64, toUserId int64, content string, actionType int64) error
	// MessageChat 聊天记录服务
	MessageChat(loginUserId int64, targetUserId int64) (dao.Message, error)

	/*
		对外提供服务接口
	*/
	// LatestMessage 返回最近的一条聊天记录  -- followerServiceImpl
	LatestMessage(loginUserId int64, targetUserId int64) (LatestMessage, error)
}
