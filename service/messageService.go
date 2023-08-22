package service

import (
	"simple_douyin/dao"
)

// LatestMessage 提供给用户好友列表接口的最新一条聊天信息, msgType 消息类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
type LatestMessage struct {
	Message string `json:"message"`
	MsgType int    `json:"msg_type"`
}

type MessageService interface {

	/*
		模块业务所需的服务接口
	*/
	// SendMessage 发送消息服务
	SendMessage(id int64, fromUserId int64, toUserId int64, content string, actionType int64) error
	// MessageChat 聊天记录服务
	MessageChat(loginUserId int64, targetUserId int64) (dao.Message, error)

	/*
		模块对外提供的服务接口
	*/
	// LatestMessage 返回最近的一条聊天记录  -- followerServiceImpl
	LatestMessage(loginUserId int64, targetUserId int64) (LatestMessage, error)
}
