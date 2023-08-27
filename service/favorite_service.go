package service

import (
	"encoding/json"
	"errors"
	"github.com/go-redsync/redsync/v4"
	"log"
	"simple_douyin/config"
	"simple_douyin/dao"
	"simple_douyin/middleware/mq"
	"simple_douyin/middleware/redis"
	"strconv"
	"sync"
	"time"
)

type FavoriteService struct{}

const (
	// ACTION_UPDATE_LIKE 点赞行为
	ACTION_UPDATE_LIKE = 1
	ACTION_CANCEL_LIKE = 2
)

// STATUS_NOT_LIKE_BEFORE 点赞状态
const STATUS_NOT_LIKE_BEFORE = 0
const STATUS_NOT_LIKE = 1

var (
	favoriteService     *FavoriteService
	favoriteServiceOnce sync.Once
	rs                  *redsync.Redsync
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

func (f *FavoriteService) FavoriteAction(userId int64, videoId int64) error {
	// 创建一个 FavoriteDao 对象
	favorite := dao.FavoriteDao{
		UserId:    userId,
		VideoId:   videoId,
		CreatedAt: time.Unix(time.Now().Unix(), 0),
	}

	log.Println("正在执行点赞/取消点赞操作")
	// 判断用户是否已经点赞过该视频
	log.Println("正在判断用户是否已经点赞过该视频")
	//如果还没有点赞，数值是 false，那么就执行点赞操作；如果已经点赞，数值是 true，那么就执行取消点赞操作
	isFavorited, err := dao.IsVideoFavoritedByUser(userId, videoId)

	log.Println("点赞状态：", isFavorited)
	if err != nil {
		log.Println("判断用户是否已经点赞过该视频失败")
		return err
	}

	//go func() {
	log.Println("正在更新redis")
	err = UpdateRedis(userId, videoId, isFavorited)
	if err != nil {
		log.Println("Failed to sync like redis:", err)
		return nil
	}
	//}() // 更新redis

	//go func() {
	// 使用消息队列异步更新数据库
	if isFavorited { // 用户之前没有点赞，所以现在执行点赞操作
		// 将 FavoriteDao 对象序列化为 JSON 字符串
		log.Println("用户之前没有点赞，所以现在执行点赞操作")
		log.Println("准备插入的favorite对象:", favorite)
		favoriteJson, err := json.Marshal(favorite)
		if err != nil {
			log.Println("Failed to marshal favorite:", err)
			return nil
		}
		// 发送消息到消息队列
		mq.SendMessage(mq.FAVORITE_ADD, string(favoriteJson))
	} else { // 用户之前已经点赞，所以现在执行取消点赞操作
		// 将 FavoriteDao 对象序列化为 JSON 字符串
		favoriteJson, err := json.Marshal(favorite)
		if err != nil {
			log.Println("Failed to marshal favorite:", err)
			return nil
		}
		// 发送消息到消息队列
		mq.SendMessage(mq.FAVORITE_REMOVE, string(favoriteJson))
	}

	if err != nil {
		log.Print(err.Error() + " Favorite action failed!")
		return err
	} else {
		log.Print("Favorite action succeed!")
	}
	return nil
}

// GetFavoriteList 获取User 点赞过的视频列表
// 逻辑：先从 UserId_VideoId Redis 集合中获取 videoId，再从 VideoId_Video Redis 集合获得 video 的所有信息（序列化为 json 格式的字符串，取出的时候再反序列化）
func (f *FavoriteService) GetFavoriteList(UserId int64) ([]dao.Video, error) {
	favoriteIdList, err := f.getFavoriteIdListByUserIdFromRedis(UserId)
	if err != nil {
		log.Println("从 Redis 中获取 favoriteIdList 失败", err)
		return nil, err
	}
	// 通过 videoId 获取 video
	var videoList []dao.Video
	for _, videoId := range favoriteIdList {
		video, err := dao.GetVideoByVideoId(videoId)
		if err != nil {
			log.Println("从数据库获取 video 失败", err)
			return nil, err
		}
		videoList = append(videoList, video)
	}
	return videoList, nil
}

func (f *FavoriteService) getFavoriteIdListByUserIdFromRedis(UserId int64) ([]int64, error) {
	// 从 Redis 中获取
	UIdFVIdR := redis.Clients.UserId_FavoriteVideoIdR
	key := config.UserId_FVideoId_KEY_PREFIX + strconv.FormatInt(UserId, 10)
	videoIdList, err := redis.GetKeysAndUpdateExpiration(UIdFVIdR, key)

	VideoIdList, ok := videoIdList.([]int64)
	if !ok {
		log.Println("类型断言失败：无法转换为 []int64")
		return nil, errors.New("类型断言失败")
	}
	if err != nil {
		log.Println("从 Redis 中获取 videoIdList 失败", err)
		return nil, err
	}
	if len(VideoIdList) == 0 {
		log.Printf("用户没有点赞过任何视频")
		log.Println("从数据库获取 videoIdList")
		VideoIdList, _ = dao.GetVideoIdListByUserId(UserId)
		// 将 VideoIdList 存入 Redis
		err = ImportVideoIdsFromDb(UserId, VideoIdList)
	}
	return VideoIdList, nil
}

// todo ：其他部分也需要进行同步列表更新
// ImportVideoIdsFromDb 从数据库内部获取数据
func ImportVideoIdsFromDb(userId int64, videoIds []int64) error {
	userIdStr := config.UserId_FVideoId_KEY_PREFIX + strconv.FormatInt(userId, 10)

	// 转化videoIds为[]interface{}以满足Redis的LPush方法的参数要求
	idsInterface := make([]interface{}, len(videoIds))
	for i, v := range videoIds {
		idsInterface[i] = v
	}
	// 使用RPush将所有videoIds存储为一个Redis List
	_, err := redis.Clients.UserId_FavoriteVideoIdR.RPush(redis.Ctx, userIdStr, idsInterface...).Result()
	if err != nil {
		return errors.New("ImportVideoIdsFromDb failed: " + err.Error())
	}
	return nil
}

//func ImportVideoIdsFromDb(userId int64, videoIds []int64) error {
//	userIdStr := config.UserId_FVideoId_KEY_PREFIX + strconv.FormatInt(userId, 10)
//
//	// 在 8.27 的凌晨的时候发现了这个问题
//	// 如果对相同的一个 id，多次存储了不同的数据，后续添加的数据将会被覆盖
//	// 例如：第一次存储了 1，第二次存储了 2，第三次存储了 3，那么最终存储的结果是 3
//
//	//for _, videoId := range videoIds {
//	//	// 将 videoId 存入 Redis，逐个存入，不使用批量存入
//	//	err := redis.SetValueWithRandomExp(redis.Clients.UserId_FavoriteVideoIdR, userIdStr, videoId)
//	//	if err != nil {
//	//		return errors.New("ImportVideoIdsFromDb failed")
//	//	}
//	//}
//	return nil
//}

// 点赞/取消时同步更新 redis 中的数据
func UpdateRedis(userId int64, videoId int64, isfavorited bool) error {
	userIdStr := strconv.FormatInt(userId, 10)
	userIdStrKey := config.UserId_FVideoId_KEY_PREFIX + userIdStr
	videoIdStr := strconv.FormatInt(videoId, 10)
	videoIdStrKey := config.VideoId_FavoritebUserId_KEY_PREFIX + videoIdStr

	log.Println("正在更新redis，参数分别是：", userIdStrKey, videoIdStrKey, isfavorited)

	// Assuming UIdFVIdR is a Redis client or similar
	var UIdFVIdR = redis.Clients.UserId_FavoriteVideoIdR

	switch isfavorited {
	case false:
		// 点赞
		log.Println("在 service 的 UpdateRedis 中，进行点赞操作")
		err := redis.SetValueWithRandomExp(UIdFVIdR, userIdStrKey, videoIdStr)
		if err != nil {
			return errors.New("UpdateRedis failed in redis: " + err.Error())
		}
		err = redis.SetValueWithRandomExp(UIdFVIdR, videoIdStrKey, userIdStr)
		if err != nil {
			return errors.New("UpdateRedis failed in redis: " + err.Error())
		}
	case true:
		// 取消点赞
		//todo 发现了一个问题，如果原本在 Redis 里面本来就没有对应的点赞的记录，那么取消点赞在删除点赞的信息的时候，就会阻塞在这里
		log.Println("在 service 的 UpdateRedis 中，取消点赞")
		err := redis.DeleteKey(UIdFVIdR, userIdStrKey)
		if err != nil {
			return errors.New("UpdateRedis failed in redis: " + err.Error())
		}
		err = redis.DeleteKey(UIdFVIdR, videoIdStrKey)
		if err != nil {
			return errors.New("UpdateRedis failed in redis: " + err.Error())
		}
	default:
		log.Println("UpdateRedis 传入的 ActionType 参数错误")
	}
	return nil
}

// GettotalFavorited 返回用户点赞的视频数量。
// 逻辑是计算与用户ID关联在Redis集合中的videoIds的数量。
func GettotalFavorited(userId int64) (int, error) {
	UIdFVIdR := redis.Clients.UserId_FavoriteVideoIdR
	key := config.UserId_FVideoId_KEY_PREFIX + strconv.FormatInt(userId, 10)

	// 使用一个函数来获取集合中的元素数量。
	// 假设有一个函数像redis.CountElements或类似的。
	count, err := redis.CountElements(UIdFVIdR, key)
	if err != nil {
		log.Println("获取用户点赞的视频数量失败：", err)
		return 0, err
	}
	return count, nil
}

// GetfavoriteCount 返回点赞该视频的用户数量。
// 逻辑是计算与视频ID关联在Redis集合中的userIds的数量。
func GetfavoriteCount(videoId int64) (int, error) {
	UIdFVIdR := redis.Clients.UserId_FavoriteVideoIdR
	key := config.VideoId_FavoritebUserId_KEY_PREFIX + strconv.FormatInt(videoId, 10)

	// 使用一个函数来获取集合中的元素数量。
	// 假设有一个函数像redis.CountElements或类似的。
	count, err := redis.CountElements(UIdFVIdR, key)
	if err != nil {
		log.Println("获取点赞该视频的用户数量失败：", err)
		return 0, err
	}
	return count, nil
}
