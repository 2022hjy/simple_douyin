package dao

import (
	"log"
	"simple_douyin/config"
	"simple_douyin/middleware/database"
	"time"
)

type Message struct {
	Id         int64     `json:"id"`           // 消息ID
	FromUserID int64     `json:"from_user_id"` // 发送消息的用户ID
	ToUserID   int64     `json:"to_user_id"`   // 接收消息的用户ID
	Content    string    `json:"content"`      // 消息内容
	CreateTime time.Time `json:"create_time"`  // 创建时间
}

func (Message) TableName() string {
	return "message"
}

// SendMessage  fromUserId 发送消息给 toUserId
func SendMessage(message Message) (Message, error) {
	result := database.Db.Model(Message{}).Create(&message)
	if result.Error != nil {
		log.Println("SendMessage失败：", result.Error.Error())
		return Message{}, result.Error
	}
	return message, nil
}

// MessageChat 当前登录用户和其他指定用户的聊天记录
// TODO 是否要考虑上次最新消息时间，还是直接获取全部聊天记录
func MessageChat(loginUserId int64, targetUserId int64) ([]Message, error) {
	messages := make([]Message, 0, config.MessageInitNum)
	result := database.Db.Where("(from_user_id = ? and to_user_id = ?) or (from_user_id = ? and to_user_id = ?)",
		loginUserId, targetUserId, targetUserId, loginUserId).
		Order("create_time ASC").
		Find(&messages)
	if result.RowsAffected == 0 {
		return messages, nil
	}
	if result.Error != nil {
		log.Println("获取MessageChat失败：", result.Error.Error())
		return nil, result.Error
	}
	return messages, nil
}

// LatestMessage  返回 loginUserId 和 targetUserId 最近的一条聊天记录
func LatestMessage(loginUserId int64, targetUserId int64) (Message, error) {
	var message Message
	query := database.Db.Where("from_user_id = ? AND to_user_id = ?", loginUserId, targetUserId).
		Or("to_user_id = ? AND from_user_id = ?", targetUserId, loginUserId).
		Order("create_time DESC").
		Limit(1)
	result := query.Take(&message)
	if result.Error != nil {
		log.Println("获取最近一条聊天记录失败:", result.Error.Error())
		return Message{}, result.Error
	}
	return message, nil
}
