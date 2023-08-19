package dao

import (
	"log"
	"simple_douyin/config"
	"simple_douyin/middleware/database"
	"time"
)

type Message struct {
	Id         int64     `json:"id" `
	UserId     int64     `json:"user_id" `
	ReceiverId int64     `json:"receiver_id"`
	ActionType int64     `json:"action_type"`
	MsgContent string    `json:"msg_content" `
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" `
}

// APIMessage 返回提取的消息（仅需 Message 的部分字段）
type APIMessage struct {
	Id         int64     `json:"id"`
	MsgContent int64     `json:"msg_content"`
	CreatedAt  time.Time `json:"create_time"`
}

func (Message) TableName() string {
	return "message"
}

// SaveMessage  保存消息
func SaveMessage(msg Message) error {
	result := database.Db.Save(&msg)
	if result.Error != nil {
		log.Println("数据库保存消息失败！", result.Error)
		return result.Error
	}
	return nil
}

// SendMessage  fromUserId 发送消息 content 给 toUserId
func SendMessage(id int64, fromUserId int64, toUserId int64, content string, actionType int64) error {
	var message Message
	message.Id = id
	message.UserId = fromUserId
	message.ReceiverId = toUserId
	message.ActionType = actionType
	message.MsgContent = content
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()
	return SaveMessage(message)
}

// MessageChat 当前登录用户和其他指定用户的聊天记录
func MessageChat(loginUserId int64, targetUserId int64, latestTime time.Time) ([]Message, error) {
	messages := make([]Message, 0, config.VIDEO_NUM_PER_REFRESH)
	result := database.Db.Where("((user_id = ? AND receiver_id = ?) OR (user_id = ? AND receiver_id = ?)) AND created_at BETWEEN ? AND ?",
		loginUserId, targetUserId, targetUserId, loginUserId, latestTime, time.Now()).
		Order("created_at ASC").
		Find(&messages)
	if result.RowsAffected == 0 {
		return messages, nil
	}
	if result.Error != nil {
		log.Println("获取聊天记录失败！")
		return nil, result.Error
	}
	return messages, nil
}

// LatestMessage  返回 loginUserId 和 targetUserId 最近的一条聊天记录
func LatestMessage(loginUserId int64, targetUserId int64) (Message, error) {
	var message Message
	query := database.Db.Where("user_id = ? AND receiver_id = ?", loginUserId, targetUserId).
		Or("user_id = ? AND receiver_id = ?", targetUserId, loginUserId).
		Order("created_at DESC").Limit(1)
	result := query.Take(&message)
	if result.Error != nil {
		log.Println("获取最近一条聊天记录失败！")
		return Message{}, result.Error
	}
	return message, nil
}