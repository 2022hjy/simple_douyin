package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type Comment struct {
	Id        int64     //评论id
	UserId    int64     //评论用户id
	VideoId   int64     //视频id
	Content   string    //评论内容
	CreatedAt time.Time //评论发布的日期mm-dd
	UpdatedAt time.Time
}

// 删除的操作逻辑：物理删除

// TableName 修改表名映射
func (Comment) TableName() string {
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

func AddComment(comment Comment) (Comment, error) {
	result := Db.Model(Comment{}).Create(&comment)
	return comment, handleDBError(result, "Insert comment")
}

func DeleteComment(commentId int64) error {
	result := Db.Where("id = ?", commentId).Delete(&Comment{})
	return handleDBError(result, "Delete comment")
}

func GetCommentList(videoId int64) ([]Comment, error) {
	var commentList []Comment
	result := Db.Model(Comment{}).Where("video_id = ?", videoId).
		Order("created_at desc").Find(&commentList)
	return commentList, handleDBError(result, "Get comment list")
}

func GetCommentCnt(videoId int64) (int64, error) {
	var count int64
	result := Db.Model(Comment{}).Where("video_id = ?", videoId).
		Count(&count)
	return count, handleDBError(result, "Get comment count")
}
