package util

import (
	"simple_douyin/controller"
	"simple_douyin/dao"
	"time"
)

func convertToResponse(comment dao.CommentDao) controller.Comment {
	// Convert time.Time to string
	createDate := comment.CreatedAt.Format(time.RFC3339)
	User, _ := dao.GetUserById(comment.UserId)
	return controller.Comment{
		Id: comment.Id,
		//User:		User,
		//TODO 需要将user的信息也返回给前端，等待用户数据库完成后再补充
		Content:    comment.Content,
		CreateDate: createDate,
	}
}
