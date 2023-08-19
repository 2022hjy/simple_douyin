package service

import (
	"encoding/json"
	"errors"
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	"log"
	"simple_douyin/dao"
	"simple_douyin/middleware/database"
	"simple_douyin/middleware/redis"
	"strconv"
	"sync"
	"time"
)

type FavoriteService struct{}

var (
	favoriteService     *FavoriteService
	favoriteServiceOnce sync.Once
)

func init() {
	favoriteService = GetFavoriteServiceInstance()
}

func GetFavoriteServiceInstance() *FavoriteService {
	favoriteServiceOnce.Do(func() {
		favoriteService = &FavoriteService{}
	})
	return favoriteService
}

func (f *FavoriteService) Favorite(userId int64, videoId int64) (favorite dao.FavoriteDao, err error) {
	favorite = dao.FavoriteDao{
		UserId:    userId,
		VideoId:   videoId,
		CreatedAt: time.Unix(time.Now().Unix(), 0),
		UpdatedAt: time.Unix(time.Now().Unix(), 0),
	}
	// Store in database
	AddRes, err := dao.AddFavorite(favorite)
	if err != nil {
		log.Fatalf("add favorite failed, err:%v\n", err)
		return
	}
	// Store in redis
	FavoriteId := AddRes.Id
	updateFavoriteRedis(videoId, FavoriteId, favorite)
	log.Printf("successfully save FavoriteDao in redis, favoriteId:%v\n", FavoriteId)
	return favorite, nil
}

func (f *FavoriteService) Unfavorite(userId int64, videoId int64) error {
	err := database.Db.Transaction(func(tx *gorm.DB) error {
		err := dao.DeleteFavorite(userId, videoId)
		if err != nil {
			log.Printf("delete favorite failed, err:%v\n", err)
			return errors.New("delete favorite in MySQL failed")
		}

		err = deleteFavoriteRedis(videoId, userId)
		if err != nil {
			log.Printf("delete favorite in redis failed, err:%v\n", err)
			return errors.New("delete favorite in redis failed")
		}

		log.Printf("successfully delete FavoriteDao in redis, userId:%v, videoId:%v\n", userId, videoId)
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func updateFavoriteRedis(videoId int64, favoriteId int64, favorite dao.FavoriteDao) {
	// 将 videoId 和 favoriteId 转换为字符串
	vId := strconv.FormatInt(videoId, 10)
	fId := strconv.FormatInt(favoriteId, 10)

	// 将 favorite 对象序列化为 JSON
	favoriteJson, err := json.Marshal(favorite)
	if err != nil {
		log.Fatalf("无法将 favoriteDao 序列化为 JSON，err:%v\n", err)
		return
	}

	// 1. 在 Video_FavoriteId Redis 集合中设置 favoriteId
	VFidClient := redis.Clients.Video_FavoriteIdR
	if VFidClient == nil {
		log.Fatalf("redis 客户端为空")
		return
	}
	err = redis.SetValueWithRandomExp(VFidClient, vId, fId)
	if err != nil {
		log.Fatalf("设置 redis 失败，err:%v\n", err)
		return
	}

	// 2. 在 FavoriteId_Favorite Redis 散列中设置 favorite 对象
	FFidClient := redis.Clients.FavoriteId_FavoriteR
	if FFidClient == nil {
		log.Fatalf("redis 客户端为空")
		return
	}
	err = redis.SetValueWithRandomExp(FFidClient, fId, string(favoriteJson))
	if err != nil {
		log.Fatalf("设置 redis 失败，err:%v\n", err)
		return
	}
}

func deleteFavoriteRedis(videoId int64, favoriteId int64) error {
	// 将 videoId 和 favoriteId 转换为字符串
	vId := strconv.FormatInt(videoId, 10)
	fId := strconv.FormatInt(favoriteId, 10)

	// 1. 从 Video_FavoriteId Redis 集合中移除 favoriteId
	VFidClient := redis.Clients.Video_FavoriteIdR
	if VFidClient == nil {
		log.Fatalf("redis 客户端为空")
		return errors.New("redis 客户端为空")
	}

	// 创建一个分布式互斥锁
	pool := goredis.NewPool(VFidClient)
	rs := redsync.New(pool)
	mutexName := "lock:deleteFavoriteRedis:" + vId + ":" + fId
	mutex := rs.NewMutex(mutexName, redsync.WithExpiry(8*time.Second))
	if err := mutex.Lock(); err != nil {
		log.Printf("无法获取锁，err:%v\n", err)
		return err
	}
	defer mutex.Unlock()

	// 从 Video_FavoriteId Redis 集合中删除 favoriteId
	err := VFidClient.SRem(vId, fId).Err()
	if err != nil {
		log.Printf("删除 redis 失败，err:%v\n", err)
		return err
	}

	// 2. 从 FavoriteId_Favorite Redis 散列中删除 favorite 对象
	FFidClient := redis.Clients.FavoriteId_FavoriteR
	if FFidClient == nil {
		log.Fatalf("redis 客户端为空")
		return errors.New("redis 客户端为空")
	}
	err = redis.DeleteKey(FFidClient, fId)
	if err != nil {
		log.Printf("删除 redis 失败，err:%v\n", err)
		return err
	}
	log.Printf("成功在 redis 中删除 favorite，favoriteId:%v\n", favoriteId)
	return nil
}

func (f *FavoriteService) GetFavoriteList(videoId int64) ([]dao.FavoriteDao, error) {
	favoriteList, err := f.getFavoriteListFromRedis(videoId)
	if err != nil {
		return nil, err
	}

	if len(favoriteList) > 0 {
		return favoriteList, nil
	}

	favoriteList, err = f.getFavoriteListFromDB(videoId)
	if err != nil {
		return nil, err
	}

	return favoriteList, nil
}

func (f *FavoriteService) getFavoriteListFromRedis(videoId int64) ([]dao.FavoriteDao, error) {
	var favoriteList []dao.FavoriteDao
	VidFidR := redis.Clients.Video_FavoriteIdR
	videoIdToStr := strconv.FormatInt(videoId, 10)
	favoriteIdListInterface, err := redis.GetKeysAndUpdateExpiration(VidFidR, videoIdToStr)
	if err != nil {
		log.Println("read redis vId failed", err)
		return nil, err
	}

	favoriteIdStringList, ok := favoriteIdListInterface.(map[string]string)
	if !ok {
		log.Println("failed to assert type to map[string]string")
		return nil, errors.New("type assertion failed")
	}

	FIdFR := redis.Clients.FavoriteId_FavoriteR
	for _, favoriteIdStr := range favoriteIdStringList {
		favoriteInterface, err := redis.GetKeysAndUpdateExpiration(FIdFR, favoriteIdStr)
		if err != nil {
			log.Println("read redis fId failed", err)
			continue
		}

		favoriteString, ok := favoriteInterface.(string)
		if !ok {
			log.Println("failed to assert type to string")
			continue
		}

		var favorite dao.FavoriteDao
		err = json.Unmarshal([]byte(favoriteString), &favorite)
		if err != nil {
			log.Println("unmarshal failed", err)
			continue
		}

		favoriteList = append(favoriteList, favorite)
	}
	return favoriteList, nil
}

func (f *FavoriteService) getFavoriteListFromDB(videoId int64) ([]dao.FavoriteDao, error) {
	favoriteList, err := dao.GetFavoriteList(videoId)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return favoriteList, nil
}
