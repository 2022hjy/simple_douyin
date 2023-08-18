package service

import (
	"log"
	"simple_douyin/dao"
	"sync"
	"time"
)

type CommentService struct {
}

var (
	commentService     *CommentService
	commentServiceOnce sync.Once
)

func init() {
	commentService = GetCommentServiceInstance()
}

func GetCommentServiceInstance() *CommentService {
	commentServiceOnce.Do(func() {
		commentService = &CommentService{
			//需要嵌套其他service时，使用&CommentService{&OtherService{}}
		}
	})
	return commentService
}

func (c *CommentService) Comment(userId int64, videoId int64, content string) (comment Comment, err error) {
	comment = dao.Comment{
		UserId:    userId,
		VideoId:   videoId,
		Content:   content,
		CreatedAt: time.Unix(time.Now().Unix(), 0),
		UpdatedAt: time.Unix(time.Now().Unix(), 0),
	}
	//存入数据库
	AddRes, err := dao.AddComment(comment)
	if err != nil {
		log.Fatalf("add comment failed, err:%v\n", err)
		return
	}
	//存入redis
	go func() {

	}()
}

func (c *CommentService) DeleteComment() {

}

func (c *CommentService) GetCommentList() {

}

func (c *CommentService) GetCommentNum() {

}

func updateCommentRedis(videoId int64, commentId int64, comment Comment) {

}
