package service

import (
	"encoding/json"
	"errors"
	"github.com/go-redsync/redsync/v4"
	"log"
	"simple_douyin/controller"
	"simple_douyin/dao"
	"simple_douyin/middleware/mq"
	"simple_douyin/middleware/redis"
	"simple_douyin/util"
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
		UpdatedAt: time.Unix(time.Now().Unix(), 0),
	}

	// 判断用户是否已经点赞过该视频
	isFavorited, err := dao.IsVideoFavoritedByUser(userId, videoId)

	go syncLikeRedis(userId, videoId, 1) // 更新redis

	go func() {
		// 使用消息队列异步更新数据库
		if isFavorited { // 用户之前没有点赞，所以现在执行点赞操作
			// 将 FavoriteDao 对象序列化为 JSON 字符串
			favoriteJson, err := json.Marshal(favorite)
			if err != nil {
				log.Println("Failed to marshal favorite:", err)
				return
			}
			// 发送消息到消息队列
			mq.SendMessage(mq.FAVORITE_ADD, string(favoriteJson))
		} else { // 用户之前已经点赞，所以现在执行取消点赞操作
			// 将 FavoriteDao 对象序列化为 JSON 字符串
			favoriteJson, err := json.Marshal(favorite)
			if err != nil {
				log.Println("Failed to marshal favorite:", err)
				return
			}
			// 发送消息到消息队列
			mq.SendMessage(mq.FAVORITE_REMOVE, string(favoriteJson))
		}
	}()

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
func (f *FavoriteService) GetFavoriteList(UserId int64) ([]dao.FavoriteDao, error) {
	favoriteList, err := f.getFavoriteListFromRedis(UserId)
	if err != nil {
		return nil, err
	}
	if len(favoriteList) > 0 {
		return favoriteList, nil
	}

	favoriteList, err = f.getFavoriteListFromDB(UserId)
	if err != nil {
		return nil, err
	}
	return favoriteList, nil
}

//todo 视频模块
func (f *FavoriteService) getFavoriteListFromRedis(UserId int64) ([]controller.VideoResponse, error) {
	var favoriteList []dao.FavoriteDao

	// 获取 UserId_FavoriteId Redis 集合客户端
	//todo: 直接通过获取的关联的 id 集合，从而去获得对应的点赞的视频数目（len 方法去获取长度）
	UFidClient := redis.Clients.UserId_FavoriteVideoIdR
	userIdToStr := strconv.FormatInt(UserId, 10)

	// 从 UserId_FavoriteId Redis 集合中获取 FavoriteId 列表
	isExist, err := redis.IsKeyExist(UFidClient, userIdToStr)
	if err != nil {
		log.Println("判断 Redis 中是否存在 uId 失败", err)
		return nil, err
	}
	if !isExist {
		log.Println("Redis 中不存在 uId")
		//return nil, errors.New("redis 中不存在 uId")
		// 从数据库中获取，并写入 Redis（不能直接返回）

	}
	// 如果 id 组存在的情况
	videoIdList, err := UFidClient.SMembers(redis.Ctx,userIdToStr).Result()
	if err != nil {
		log.Println("从 Redis 中获取 videoId 失败", err)
		return nil, err
	}

	// 从 VideoId_Video Redis 集合中获取 video 对象
	// todo:实现一个 videoId 到 video 对象的映射，再从 VideoId_Video Redis 集合获得 video 的所有信息（序列化为 json 格式的字符串，取出的时候再反序列化）
	VidClient := redis.Clients.VideoId_VideoR
	VideoList := make([]controller.VideoResponse, 0)
	for _, vId := range videoIdList {
		videoStr, err := redis.GetKeysAndUpdateExpiration(VidClient, vId)
		if err != nil {
			log.Println("从 Redis 中获取 video 失败", err)
			return nil, err
		}
		var video dao.Video
		videoByte,ok := videoStr.(string)
		if !ok {
			log.Println("类型断言失败")
			return nil, errors.New("类型断言失败")
		}
		err = json.Unmarshal([]byte(videoByte), &video)
		if err != nil {
			log.Println("反序列化 video 失败", err)
			return nil, err
		}

		//userId -> commentId, commentId -> comment
		//将 dao.Video 转换为 controller.VideoResponse，再补充  controller.UserResponse 再添加到 VideoList 中
		//Redis库内部存放的统一都是 dao 类型的对象，比如dao.comment 所以需要转换为 controller 类型的对象

		//1.获取User信息：先以参数传递的形式获得UserId
		//1.1 从 UserId_UserR Redis 集合中获取 User 对象=> key:UserId: + UserId
		//1.2 将 User 对象转换为 controller.UserResponse
		//	1.2.1 获取 User 对象的 Id，通过 Redis 获取 User 对象的信息，拿出来的是 json 格式的字符串，再反序列化为 User 对象
		//	1.2.2 将 User 对象转换为 controller.UserResponse
		//		补充：存放在 Redis 中的 User 对象是
		//		1.2.3  补充额外的属性：点赞数，粉丝数，关注数，是否关注，总共获赞数，作品数
		//		1.2.4  其中，（点赞模块）获赞数 total_favorited，点赞数 favorite_count，   	需要从 Redis 中获取
		//		      	 Redis 的 Client：UserId_FavoritedNumR，
		//				 Redis 的 Client：UserId_FavoriteNumR
		//			采取逻辑过期的方式，每次点赞，都会自动更新一个固定的时间 5 分钟，如果 5 分钟内没有点赞，就会过期，需要重新计算

		//					  （视频模块）作品数work_count 									需要从 Redis 中获取
		//		      Redis 的Client：UserId_VideoNumR，
		//		1.2.5  粉丝数 follower_count，关注数 follow_count,是否关注 IsFollow 需要从数据库中获取
		//		1.2.6  得到 controller.UserResponse，拿一个变量进行接受
		//
		//2.获取Video信息：
		//	2.1 从 VideoId_Video Redis 集合中获取 Video 对象
		//	2.2 进行反序列化，得到 dao.Video 对象
		//	2.3 调用 ConvertDBVideoToResponse 将 dao.Video 对象转换为 controller.VideoResponse
		//	2.4 返回 controller.VideoResponse 的集合

	}
	return favoriteList, nil
}

// GetVideoResponseList 获取用户的视频列表
func getVideoResponseList(userId int64, video dao.Video) ([]controller.VideoResponse, error) {
	// 存放结果的集合
	var videoResponseList []controller.VideoResponse
	// 1. 获取User信息
	// 1.1 从 UserId_UserR Redis 集合中获取 User 对象
	userKey := "UserId:" + string(userId)
	userStr, err := redisClient.Get(ctx, userKey).Result()
	if err != nil {
		return nil, err
	}

	// 1.2 将 User 对象转换为 controller.UserResponse
	var user dao.UserDAO
	err = json.Unmarshal([]byte(userStr), &user)
	if err != nil {
		return nil, err
	}

	// 1.2.3 补充额外的属性
	favoritedKey := "UserId_FavoritedNumR:" + string(userId)
	favoriteKey := "UserId_FavoriteNumR:" + string(userId)
	videoNumKey := "UserId_VideoNumR:" + string(userId)
	totalFavorited, err := redisClient.Get(ctx, favoritedKey).Int64()
	if err != nil {
		return nil, err
	}
	favoriteCount, err := redisClient.Get(ctx, favoriteKey).Int64()
	if err != nil {
		return nil, err
	}
	workCount, err := redisClient.Get(ctx, videoNumKey).Int64()
	if err != nil {
		return nil, err
	}
	// 1.2.5 从数据库获取粉丝数，关注数和是否关注
	followerCount, followCount, isFollow, err := getFollowInfo(userId)
	if err != nil {
		return nil, err
	}

	userResponse := util.ConvertDBUserToResponse(user, favoriteCount, followCount, followerCount, isFollow, string(totalFavorited), workCount)

	// 2. 获取Video信息
	videoList, err := getVideosByUserId(userId)
	if err != nil {
		return nil, err
	}

	for _, v := range videoList {
		// 2.1 从 VideoId_Video Redis 集合中获取 Video 对象
		videoKey := "VideoId_Video:" + string(v.Id)
		videoStr, err := redisClient.Get(ctx, videoKey).Result()
		if err != nil {
			return nil, err
		}

		// 2.2 进行反序列化，得到 dao.Video 对象
		var video dao.Video
		err = json.Unmarshal([]byte(videoStr), &video)
		if err != nil {
			return nil, err
		}

		// 2.3 调用 ConvertDBVideoToResponse 将 dao.Video 对象转换为 controller.VideoResponse
		videoResponse := util.ConvertDBVideoToResponse(video, userResponse)
		videoResponseList = append(videoResponseList, videoResponse)
	}

	// 2.4 返回 controller.VideoResponse 的集合
	return videoResponseList, nil
}

func getFollowInfo(userId int64) (int64, int64, bool, error) {
	// Implement this function to get the follower count, follow count, and follow status from the database.
	return 0, 0, false, nil
}

func getVideosByUserId(userId int64) ([]dao.Video, error) {
	// Implement this function to get the videos of a user.
	return nil, nil
}



	//事实上，感觉这里不需要关心并发问题，因为这里只是读取，不涉及到写入，所以不需要加锁
	//// 类型断言
	//favoriteIdStringList, ok := favoriteIdListInterface.(map[string]string)
	//if !ok {
	//	log.Println("类型断言失败：无法转换为 map[string]string")
	//	return nil, errors.New("类型断言失败")
	//}
	//// 再通过 favoriteId 获取 favorite 对象
	//FFidClient := redis.Clients.FavoriteId_FavoriteR
	//for _, favoriteIdStr := range favoriteIdStringList {
	//
	//}

}

// ImportVideoIdsFromDb 从数据库内部获取数据
func ImportVideoIdsFromDb(userId int64, videoIds []int64) error {
	userIdStr := strconv.FormatInt(userId, 10)
	for _, videoId := range videoIds {
		redis.RdbUVid.SAdd(redis.Ctx, userIdStr, videoId)
	}
	// 设置过期时间，为数据不一致情况兜底
	redis.RdbUVid.Expire(redis.Ctx, userIdStr, config.ExpireTime)
	return nil
}

// 点赞/取消时同步更新 redis 中的数据
func syncLikeRedis(userId int64, videoId int64, actionType bool) error {
	userIdStr := strconv.FormatInt(userId, 10)
	videoIdStr := strconv.FormatInt(videoId, 10)
	lockName := "lock:syncLikeRedis:" + userIdStr + ":" + videoIdStr
	mutex := rs.NewMutex(lockName, redsync.WithExpiry(10*time.Second))
	err := mutex.Lock()
	if err != nil {
		return errors.New("无法获取锁")
	}
	defer mutex.Unlock()
	switch actionType {
	case true:
		// 点赞
		redis.RdbUVid.SAdd(redis.Ctx, userIdStr, videoId)
		redis.RdbVUid.SAdd(redis.Ctx, videoIdStr, userId)
	case false:
		// 取消点赞
		redis.RdbUVid.SRem(redis.Ctx, userIdStr, videoId)
		redis.RdbVUid.SRem(redis.Ctx, videoIdStr, userId)
	default:
		log.Println("syncLikeRedis 传入的 ActionType 参数错误")
	}
	return nil
}

// ConvertDBVideoToResponse 转换数据库视频结构体到前端返回结构体
//todo dao.User 的内容不全，需要补充
