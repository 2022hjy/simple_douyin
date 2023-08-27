package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type CommentDao struct {
	Id        int64     `gorm:"column:id"`           //评论id
	UserId    int64     `gorm:"column:user_info_id"` //评论用户id
	VideoId   int64     `gorm:"column:video_id"`     //视频id
	Content   string    `gorm:"column:content"`      //评论内容
	CreatedAt time.Time `gorm:"column:created_at"`
}

// 删除的操作逻辑：物理删除

// TableName 修改表名映射
func (CommentDao) TableName() string {
	return "comment"
}

func handleDBError(result *gorm.DB, action string) error {
	if result.Error != nil {
		errMsg := fmt.Sprintf("%s failed: %s", action, result.Error.Error())
		log.Println(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func AddComment(comment CommentDao) (CommentDao, error) {
	result := Db.Model(CommentDao{}).Create(&comment)
	return comment, handleDBError(result, "Insert comment")
}

func DeleteComment(commentId int64) error {
	return Db.Transaction(func(tx *gorm.DB) error {
		result := tx.Model(CommentDao{}).Where("id = ?", commentId).Delete(&CommentDao{})
		return handleDBError(result, "Delete comment")
	})
}

// deprecated function
//func GetCommentId(userId int64, videoId int64, content string) (int, error) {
//	var comment CommentDao
//	result := Db.Where("user_id = ? and video_id = ? and content = ?", userId, videoId, content).
//		First(&comment)
//	if result.Error != nil {
//		return 0, result.Error
//	}
//	return int(comment.Id), nil
//}

func GetCommentList(videoId int64) ([]CommentDao, error) {
	var commentList []CommentDao
	result := Db.Model(CommentDao{}).Where("video_id = ?", videoId).
		Order("created_at desc").Find(&commentList)
	return commentList, handleDBError(result, "Get comment list")
}

func GetCommentCnt(videoId int64) (int64, error) {
	var count int64
	result := Db.Model(CommentDao{}).Where("video_id = ?", videoId).
		Count(&count)
	return count, handleDBError(result, "Get comment count")
}

// TODO: 需要将user的信息也返回给前端
func GetUserFromCommentId(commentId int64) (int64, error) {
	var comment CommentDao
	result := Db.Model(CommentDao{}).Where("id = ?", commentId).First(&comment)
	if result.Error != nil {
		return 0, result.Error
	}
	return comment.UserId, nil
}
