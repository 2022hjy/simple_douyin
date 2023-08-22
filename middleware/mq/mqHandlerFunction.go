package mq

import (
	"encoding/json"
	"fmt"
	"log"
	"simple_douyin/dao"
)

// MessageHandler 定义消息处理函数
// 按照注释的要求去定义消息处理函数：return的值是一个函数，此函数是在需要处理的 service 层内部的函数
// 例如：return func(msg string) { service.AddComment(...) }
// 可以选择一个返回值，也可以选择两个返回值，第一个返回值是一个函数，第二个返回值是一个 error
type MessageHandler func(string)

func AddComment(msg string) {
	comment := dao.CommentDao{}
	err := json.Unmarshal([]byte(msg), &comment)
	if err != nil {
		log.Println("Failed to unmarshal comment:", err)
		return
	}

	_, err = dao.AddComment(comment)
	if err != nil {
		log.Println("Failed to add comment to the database:", err)
	}
}

func DeleteComment(msg string) {
	var commentId int64
	err := json.Unmarshal([]byte(msg), &commentId)
	if err != nil {
		log.Println("Failed to unmarshal commentId:", err)
		return
	}

	err = dao.DeleteComment(commentId)
	if err != nil {
		log.Println("Failed to delete comment from the database:", err)
	}
}

func AddLike(body string) {
	favorite := dao.FavoriteDao{}
	err := json.Unmarshal([]byte(body), &favorite)
	if err != nil {
		log.Println("Failed to unmarshal favorite:", err)
		return
	}
	err = dao.InsertFavoriteInfo(favorite)
	if err != nil {
		log.Println("Failed to add like:", err)
	}
}

func RemoveLike(body string) {
	favorite := dao.FavoriteDao{}
	err := json.Unmarshal([]byte(body), &favorite)
	if err != nil {
		log.Println("Failed to unmarshal favorite:", err)
		return
	}
	err = dao.DeleteFavoriteInfo(favorite.UserId, favorite.VideoId)
	if err != nil {
		log.Println("Failed to remove like:", err)
	}
}

func AddFollow(msg string) {
	fmt.Println("Adding follow:", msg)
}

func RemoveFollow(msg string) {
	fmt.Println("Removing follow:", msg)
}
