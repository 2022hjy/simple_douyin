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
	log.Println("Entering Comment function...")
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
	log.Println("Entering Comment function...")
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
		// 判断是否在redis中存在
		key := config.CommentId_Comment_KEY_PREFIX + strconv.FormatInt(commentId, 10)
		log.Println("这个评论在redis中的key是：", key)
		isExist, err := redis.IsKeyExist(redis.Clients.CommentId_CommentR, key)
		if err != nil {
			log.Println("Error checking if key exists in Redis:", err)
			return errors.New("check if key exists in Redis failed")
		}
		log.Println("这个评论isExist In redis？", isExist)
		if !isExist {
			err = deleteCommentRedis(videoId, commentId)
			if err != nil {
				// 如果Redis操作失败，回滚数据库事务
				log.Printf("delete comment in redis failed, err:%v\n", err)
				return errors.New("delete comment in redis failed")
			}

			log.Printf("successfully delete CommentDao in redis, commentId:%v\n", commentId)
		}
		return nil
	},
	)
	// 检查事务是否成功
	if err != nil {
		return err
	}
	return nil
}

func (c *CommentService) GetCommentList(videoId int64) ([]dao.CommentDao, error) {
	log.Println("现在从redis中获取commentList")
	commentList, err := c.getCommentListFromRedis(videoId)
	log.Println("从 Redis 里面获取的commentList:", commentList)
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
	videoIdToStr := config.VideoId_CommentId_KEY_PREFIX + strconv.FormatInt(videoId, 10)
	//commentIdListInterface, err := redis.GetKeysAndUpdateExpiration(VidCidR, videoIdToStr)
	//commentIdListInterface, err := VidCidR.SMembers(context.Background(), videoIdToStr).Result()  这是基于 go-redis 里面使用 set 的方式
	values, err := VidCidR.LRange(context.Background(), videoIdToStr, 0, -1).Result()
	if err != nil {
		log.Fatalf("Error getting values from list: %v", err)
	}

	log.Println("我在 comment_service内部，commentIdStringList:", values)

	CIdCR := redis.Clients.CommentId_CommentR //一对一

	for _, commentIdStr := range values {
		commentString, err := CIdCR.Get(context.Background(), config.CommentId_Comment_KEY_PREFIX+commentIdStr).Result()

		if err != nil {
			log.Println("Error getting comment from Redis:", err)
			continue
		}

		if commentString == "" {
			log.Println("Empty comment data for ID:", commentIdStr)
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
	log.Println("Entering updateCommentRedis function...")
	//采取缓存的单向添加（添加两部分的缓存）
	//1. videoId -> commentId个
	//2. commentId -> comment（序列化成 json）
	//发现之所以一直无法正常导入的原因是出现了包的名	称的冲突的问题，导入了 go-redis 包，和 middleware 内部的 redis 包起了冲突，不知道具体导入哪一个
	VCidClient := redis.Clients.VideoId_CommentIdR
	if VCidClient == nil {
		log.Fatalf("redis client is nil")
		return
	}
	//1.key == videoId, value == commentId
	vId := config.VideoId_CommentId_KEY_PREFIX + strconv.FormatInt(videoId, 10)
	err := redis.PushValueToListR(VCidClient, vId, strconv.FormatInt(commentId, 10))
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

	log.Println("cId:", cId)
	log.Println("comment:", comment)

	comment_json, serializationErr := json.Marshal(comment)
	if serializationErr != nil {
		log.Fatalf("jsonfiy comment failed, err:%v\n", serializationErr) // 注意，这里改为了打印serializationErr
		return
	} else {
		log.Println("Comment serialization successful:", string(comment_json)) // 打印序列化后的json
	}

	err = CCIdClient.Set(context.Background(), cId, string(comment_json), 0).Err()
	if err != nil {
		log.Fatalf("set redis failed, err:%v\n", err)
		return
	} else {
		log.Println("Successfully set comment in redis with key:", cId) // 确认成功地在Redis中设置了数据
	}

	// 检查Redis中是否真的设置了这个值
	valueInRedis, err := CCIdClient.Get(context.Background(), cId).Result()
	if err != nil {
		log.Println("Error getting value back from Redis for key:", cId, "Error:", err)
	} else {
		log.Println("Value in Redis for key:", cId, "is:", valueInRedis)
	}
	log.Printf("successfully save CommentDao in redis, commentId:%v\n", commentId)
	log.Println("Entering updateCommentRedis function...")
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
		log.Println("redis client is nil")
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
