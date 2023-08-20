package service

import (
	"encoding/json"
	"errors"
	"github.com/go-redsync/redsync/v4"
	"log"
	"simple_douyin/controller"
	"simple_douyin/dao"
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
	favorite = dao.FavoriteDao{
		UserId:    userId,
		VideoId:   videoId,
		CreatedAt: time.Unix(time.Now().Unix(), 0),
		UpdatedAt: time.Unix(time.Now().Unix(), 0),
	}

	isFavorited, err := dao.IsVideoFavoritedByUser(userId, videoId)

	go syncLikeRedis(userId, videoId, isFavorited) // 更新redis

	go func() {
		//使用消息队列异步更新数据库
		if isFavorited == STATUS_NOT_LIKE_BEFORE { // 用户之前没有点赞，所以现在执行点赞操作
			//todo: 插入一条点赞的新记录
		} else { // 用户之前已经点赞，所以现在执行取消点赞操作
			//todo: 删除一条点赞的记录
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

//
//func updateFavoriteRedis(videoId int64, favoriteId int64, favorite dao.FavoriteDao) {
//	// 将 videoId 和 favoriteId 转换为字符串
//	vId := strconv.FormatInt(videoId, 10)
//	fId := strconv.FormatInt(favoriteId, 10)
//
//	// 将 favorite 对象序列化为 JSON
//	favoriteJson, err := json.Marshal(favorite)
//	if err != nil {
//		log.Fatalf("无法将 favoriteDao 序列化为 JSON，err:%v\n", err)
//		return
//	}
//
//	// 1. 在 Video_FavoriteId Redis 集合中设置 favoriteId
//	VFidClient := redis.Clients.Video_FavoriteIdR
//	if VFidClient == nil {
//		log.Fatalf("redis 客户端为空")
//		return
//	}
//	err = redis.SetValueWithRandomExp(VFidClient, vId, fId)
//	if err != nil {
//		log.Fatalf("设置 redis 失败，err:%v\n", err)
//		return
//	}
//
//	// 2. 在 FavoriteId_Favorite Redis 散列中设置 favorite 对象
//	FFidClient := redis.Clients.FavoriteId_FavoriteR
//	if FFidClient == nil {
//		log.Fatalf("redis 客户端为空")
//		return
//	}
//	err = redis.SetValueWithRandomExp(FFidClient, fId, string(favoriteJson))
//	if err != nil {
//		log.Fatalf("设置 redis 失败，err:%v\n", err)
//		return
//	}
//}
//
//func deleteFavoriteRedis(videoId int64, favoriteId int64) error {
//	// 将 videoId 和 favoriteId 转换为字符串
//	vId := strconv.FormatInt(videoId, 10)
//	fId := strconv.FormatInt(favoriteId, 10)
//
//	// 1. 从 Video_FavoriteId Redis 集合中移除 favoriteId
//	VFidClient := redis.Clients.Video_FavoriteIdR
//	if VFidClient == nil {
//		log.Fatalf("redis 客户端为空")
//		return errors.New("redis 客户端为空")
//	}
//
//	// 创建一个分布式互斥锁
//	pool := goredis.NewPool(VFidClient)
//	rs := redsync.New(pool)
//	mutexName := "lock:deleteFavoriteRedis:" + vId + ":" + fId
//	mutex := rs.NewMutex(mutexName, redsync.WithExpiry(8*time.Second))
//	if err := mutex.Lock(); err != nil {
//		log.Printf("无法获取锁，err:%v\n", err)
//		return err
//	}
//	defer mutex.Unlock()
//
//	// 从 Video_FavoriteId Redis 集合中删除 favoriteId
//	err := VFidClient.SRem(vId, fId).Err()
//	if err != nil {
//		log.Printf("删除 redis 失败，err:%v\n", err)
//		return err
//	}
//
//	// 2. 从 FavoriteId_Favorite Redis 散列中删除 favorite 对象
//	FFidClient := redis.Clients.FavoriteId_FavoriteR
//	if FFidClient == nil {
//		log.Fatalf("redis 客户端为空")
//		return errors.New("redis 客户端为空")
//	}
//	err = redis.DeleteKey(FFidClient, fId)
//	if err != nil {
//		log.Printf("删除 redis 失败，err:%v\n", err)
//		return err
//	}
//	log.Printf("成功在 redis 中删除 favorite，favoriteId:%v\n", favoriteId)
//	return nil
//}

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
	UFidClient := redis.Clients.UserId_FVideoIdR
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
	//// 将 videoId 转换为 int64 类型
	//var videoIdInt64 []int64
	//for i , v := range videoId {
	//	videoIdInt64[i], err = strconv.ParseInt(v, 10, 64)
	//	if err != nil {
	//		log.Println("转换 videoId 失败", err)
	//		return nil, err
	//	}
	//}

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

		//将 dao.Video 转换为 controller.VideoResponse，再添加到 VideoList 中
		//Redis库内部存放的统一都是 dao 类型的对象，比如dao.comment 所以需要转换为 controller 类型的对象
		//
		//1.获取User信息：先以参数传递的形式获得UserId
		//1.1 从 UserId_User Redis 集合中获取 User 对象
		//1.2 将 User 对象转换为 controller.UserResponse
		//	1.2.1 获取 User 对象的 Id，通过 Redis 获取 User 对象的信息，拿出来的是 json 格式的字符串，再反序列化为 User 对象
		//	1.2.2 将 User 对象转换为 controller.UserResponse
		//		补充：存放在 Redis 中的 User 对象是
		//		1.2.3  补充额外的属性：点赞数，粉丝数，关注数，是否关注，总共获赞数，作品数
		//		1.2.4  其中，（点赞模块）获赞数 total_favorited，点赞数favorite_count，   	需要从 Redis 中获取
		//		      Redis 的Client：UserId_FavoritedNumR，
		//				 Redis 的 Client：UserId_FavoriteNumR
		//			采取逻辑过期的方式，每次点赞，都会自动更新一个固定的时间 5 分钟，如果 5 分钟内没有点赞，就会过期，需要重新计算
		//
		//					  （视频模块）作品数work_count 									需要从 Redis 中获取
		//		      Redis 的Client：UserId_VideoNumR，
		//		1.2.5  粉丝数 follower_count，关注数 follow_count,是否关注IsFollow 需要从数据库中获取
		//		1.2.6  得到 controller.UserResponse，拿一个变量进行接受
		//
		//2.获取Video信息：
		//	2.1 从 VideoId_Video Redis 集合中获取 Video 对象
		//	2.2 进行反序列化，得到 dao.Video 对象
		//	2.3 调用ConvertDBVideoToResponse 将 dao.Video 对象转换为 controller.VideoResponse
		//	2.4 返回 controller.VideoResponse 的集合
		////





	}
	return favoriteList, nil
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
func syncLikeRedis(userId int64, videoId int64, actionType int) error {
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
	case ACTION_UPDATE_LIKE:
		// 点赞
		redis.RdbUVid.SAdd(redis.Ctx, userIdStr, videoId)
		redis.RdbVUid.SAdd(redis.Ctx, videoIdStr, userId)
	case ACTION_CANCEL_LIKE:
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






func ConvertDBVideoToResponse(dbVideo dao.Video, author dao.User) controller.VideoResponse {
	return controller.VideoResponse{
		Id:            dbVideo.Id,
		User:          author,
		PlayUrl:       dbVideo.PlayUrl,
		CoverUrl:      dbVideo.CoverUrl,
		FavoriteCount: int64(dbVideo.FavoriteCount),
		CommentCount:  int64(dbVideo.CommentCount),
		IsFavorite:    true,
		Title:         dbVideo.Title,
	}
}