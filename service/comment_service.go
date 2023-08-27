package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"gorm.io/gorm"
	"log"
	"simple_douyin/config"
	"simple_douyin/dao"
	"simple_douyin/middleware/database"
	"simple_douyin/middleware/redis"
	"strconv"
	"sync"
	"time"
)

type CommentService struct{}

var (
	commentService     *CommentService
	commentServiceOnce sync.Once
)

func init() {
	commentService = GetCommentServiceInstance()
}

func GetCommentServiceInstance() *CommentService {
	commentServiceOnce.Do(func() {
		commentService = &CommentService{}
	})
	return commentService
}

func (c *CommentService) Comment(userId int64, videoId int64, content string) (comment dao.CommentDao, err error) {
	comment = dao.CommentDao{
		UserId:    userId,
		VideoId:   videoId,
		Content:   content,
		CreatedAt: time.Unix(time.Now().Unix(), 0),
	}
	//存入数据库
	// TODO 采取消息队列的方式，异步使用 rabbitmq存入数据库
	AddRes, err := dao.AddComment(comment)
	if err != nil {
		log.Fatalf("add comment failed, err:%v\n", err)
		return
	}
	//存入redis
	CommentId := AddRes.Id
	updateCommentRedis(videoId, CommentId, comment)
	log.Printf("successfully save CommentDao in redis, commentId:%v\n", CommentId)
	//返回给前端的数据
	return comment, nil
}

// DeleteComment 使用分布式锁，防止并发删除；同时使用了 Transaction 事务，保证数据库和缓存的一致性
func (c *CommentService) DeleteComment(videoId int64, commentId int64) error {
	err := database.Db.Transaction(func(tx *gorm.DB) error {
		// 1. 删除数据库中的评论
		// TODO 采取消息队列的方式，异步使用 rabbitmq删除数据库中的评论
		err := dao.DeleteComment(commentId)
		if err != nil {
			log.Printf("delete comment failed, err:%v\n", err)
			return errors.New("delete comment in MySQL failed")
		}

		// 2. 删除Redis中的评论
		err = deleteCommentRedis(videoId, commentId)
		if err != nil {
			// 如果Redis操作失败，回滚数据库事务
			// TODO: 如何处理Redis操作失败的情况
			log.Printf("delete comment in redis failed, err:%v\n", err)
			return errors.New("delete comment in redis failed")
		}

		log.Printf("successfully delete CommentDao in redis, commentId:%v\n", commentId)
		return nil
	})
	// 检查事务是否成功
	if err != nil {
		return err
	}
	return nil
}

func (c *CommentService) GetCommentList(videoId int64) ([]dao.CommentDao, error) {
	commentList, err := c.getCommentListFromRedis(videoId)
	log.Println("现在从redis中获取commentList")
	log.Println("commentList:", commentList)
	if err != nil {
		return nil, err
	}

	if len(commentList) > 0 {
		return commentList, nil
	}

	commentList, err = c.getCommentListFromDB(videoId)
	if err != nil {
		return nil, err
	}

	return commentList, nil
}

func (c *CommentService) getCommentListFromRedis(videoId int64) ([]dao.CommentDao, error) {
	var commentList []dao.CommentDao
	VidCidR := redis.Clients.VideoId_CommentIdR
	videoIdToStr := strconv.FormatInt(videoId, 10)
	commentIdListInterface, err := redis.GetKeysAndUpdateExpiration(VidCidR, videoIdToStr)
	if err != nil {
		// todo 这里出现了问题
		log.Println("read redis vId failed", err)
		return nil, err
	}

	commentIdStringList, ok := commentIdListInterface.(map[string]string)
	if !ok {
		log.Println("failed to assert type to map[string]string")
		return nil, errors.New("type assertion failed")
	}

	CIdCR := redis.Clients.CommentId_CommentR
	for _, commentIdStr := range commentIdStringList {
		commentInterface, err := redis.GetKeysAndUpdateExpiration(CIdCR, commentIdStr)
		if err != nil {
			log.Println("read redis cId failed", err)
			continue
		}

		commentString, ok := commentInterface.(string)
		if !ok {
			log.Println("failed to assert type to string")
			continue
		}

		var comment dao.CommentDao
		err = json.Unmarshal([]byte(commentString), &comment)
		if err != nil {
			log.Println("unmarshal failed", err)
			continue
		}

		commentList = append(commentList, comment)
	}
	return commentList, nil
}

func (c *CommentService) getCommentListFromDB(videoId int64) ([]dao.CommentDao, error) {
	commentList, err := dao.GetCommentList(videoId)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return commentList, nil
}

func updateCommentRedis(videoId int64, commentId int64, comment dao.CommentDao) {
	//采取缓存的单向添加（添加两部分的缓存）
	//1. videoId -> commentId
	//2. commentId -> comment（序列化成 json）

	//发现之所以一直无法正常导入的原因是出现了包的名称的冲突的问题，导入了 go-redis 包，和 middleware 内部的 redis 包起了冲突，不知道具体导入哪一个
	VCidClient := redis.Clients.VideoId_CommentIdR
	if VCidClient == nil {
		log.Fatalf("redis client is nil")
		return
	}
	vId := config.VideoId_CommentId_KEY_PREFIX + strconv.FormatInt(videoId, 10)
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
	cId := config.CommentId_Comment_KEY_PREFIX + strconv.FormatInt(commentId, 10)

	comment_json, serializationErr := json.Marshal(comment)
	if serializationErr != nil {
		log.Fatalf("jsonfiy comment failed, err:%v\n", comment_json)
		return
	}
	err = redis.SetValueWithRandomExp(CCIdClient, cId, string(comment_json))
	if err != nil {
		log.Fatalf("set redis failed, err:%v\n", err)
		return
	}
}

// 牛刀小试，尝试使用分布式锁，防止并发删除
func deleteCommentRedis(videoId int64, commentId int64) error {
	//1.删除 videoId -> commentId
	VCidClient := redis.Clients.VideoId_CommentIdR
	if VCidClient == nil {
		log.Fatalf("redis client is nil")
		return errors.New("redis client is nil")
	}
	vId := config.VideoId_CommentId_KEY_PREFIX + strconv.FormatInt(videoId, 10)
	cId := config.CommentId_Comment_KEY_PREFIX + strconv.FormatInt(commentId, 10)

	// 创建一个红色同步互斥锁
	pool := goredis.NewPool(VCidClient) // or, pool := goredis.NewPool(&redis.Pool{…})
	rs := redsync.New(pool)
	mutexName := "lock:deleteCommentRedis:" + vId + ":" + cId
	mutex := rs.NewMutex(mutexName, redsync.WithExpiry(8*time.Second))
	if err := mutex.Lock(); err != nil {
		log.Printf("Failed to acquire lock, err:%v\n", err)
		return err
	}
	defer mutex.Unlock()

	ctx := context.Background() // 创建一个新的context
	err := VCidClient.SRem(ctx, vId, cId).Err()

	if err != nil {
		log.Printf("delete redis failed, err:%v\n", err)
		return err
	}

	//2.删除 commentId -> comment
	CCidClient := redis.Clients.CommentId_CommentR
	if CCidClient == nil {
		log.Fatalf("redis client is nil")
		return errors.New("redis client is nil")
	}
	err = redis.DeleteKey(CCidClient, cId)
	if err != nil {
		log.Printf("delete redis failed, err:%v\n", err)
		return err
	}
	log.Printf("successfully delete comment in redis, commentId:%v\n", commentId)
	return nil
}
