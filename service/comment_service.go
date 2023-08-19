package service

import (
	"encoding/json"
	"log"
	"simple_douyin/dao"
	"simple_douyin/middleware/redis"
	"strconv"
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

func (c *CommentService) Comment(userId int64, videoId int64, content string) (comment dao.Comment, err error) {
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
		CommentId := AddRes.Id
		updateCommentRedis(videoId, CommentId, comment)
		log.Printf("successfully save Comment in redis, commentId:%v\n", CommentId)
	}()
	//返回给前端的数据
	return comment, nil
}

func (c *CommentService) DeleteComment(commentId int64) {
	//删除数据库中的评论，删除redis中的评论
	//操作逻辑：先更新数据库，再更新redis
	// TODO 采取消息队列的方式，异步使用 rabbitmq删除数据库中的评论
	err := dao.DeleteComment(commentId)
	if err != nil {
		log.Fatalf("delete comment failed, err:%v\n", err)
		return
	}

	//删除redis中的评论 - 分割线

	go func() {
		//1.删除 videoId -> commentId
		VCidClient := redis.Clients.Video_CommentIdR
		if VCidClient == nil {
			log.Fatalf("redis client is nil")
			return
		}
		vId := strconv.FormatInt(videoId, 10)
		err := redis.DeleteKey(VCidClient, vId)
		if err != nil {
			log.Fatalf("delete redis failed, err:%v\n", err)
			return
		}
		//2.删除 commentId -> comment
		CCidClient := redis.Clients.CommentId_CommentR
		if CCidClient == nil {
			log.Fatalf("redis client is nil")
			return
		}
		cId := strconv.FormatInt(commentId, 10)
		err = redis.DeleteValue(CCidClient, cId)
		if err != nil {
			log.Fatalf("delete redis failed, err:%v\n", err)
			return
		}
		log.Printf("successfully delete comment in redis, commentId:%v\n", commentId)
	}()
}

func (c *CommentService) GetCommentList() {

}

func (c *CommentService) GetCommentNum() {

}

func updateCommentRedis(videoId int64, commentId int64, comment dao.Comment) {
	//采取缓存的单向添加（添加两部分的缓存）
	//1. videoId -> commentId
	//2. commentId -> comment（序列化成 json）

	//发现之所以一直无法正常导入的原因是出现了包的名称的冲突的问题，导入了 go-redis 包，和 middleware 内部的 redis 包起了冲突，不知道具体导入哪一个
	VCidClient := redis.Clients.Video_CommentIdR
	if VCidClient == nil {
		log.Fatalf("redis client is nil")
		return
	}
	vId := strconv.FormatInt(videoId, 10)
	//1.key == videoId, value == commentId
	err := redis.SetValueWithRandomExp(VCidClient, vId, commentId)
	if err != nil {
		log.Fatalf("set redis failed, err:%v\n", err)
		return
	}

	//2.key == commentId, value == comment
	CCIdClient := redis.Clients.CommentId_CommentR
	if CCIdClient == nil {
		log.Fatalf("redis client is nil")
		return
	}
	cId := strconv.FormatInt(commentId, 10)
	comment_json, serializationErr := json.Marshal(comment)
	if serializationErr != nil {
		log.Fatalf("jsonfiy comment failed, err:%v\n", jsonfiyErr)
		return
	}
	err = redis.SetValueWithRandomExp(CCIdClient, cId, string(comment_json))
	if err != nil {
		log.Fatalf("set redis failed, err:%v\n", err)
		return
	}
}
