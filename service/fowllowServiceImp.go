package service

import (
	"log"
	"math/rand"
	"simple_douyin/dao"
	"simple_douyin/middleware/redis"
	"strconv"
	"sync"
	"time"
)

// FollowServiceImp 该结构体继承FollowService接口。
type FollowServiceImp struct {
	//MessageService
	FollowService
	*UserService
}

var (
	followServiceImp  *FollowServiceImp //controller层通过该实例变量调用service的所有业务方法。
	followServiceOnce sync.Once         //限定该service对象为单例，节约内存。
)

func CacheTimeGenerator() time.Duration {
	// 先设置随机数 - 这里比较重要
	rand.Seed(time.Now().Unix())
	// 再设置缓存时间
	// 10 + [0~20) 分钟的随机时间
	return time.Duration((10 + rand.Int63n(20)) * int64(time.Minute))
}

func convertToInt64Array(strArr []string) ([]int64, error) {
	int64Arr := make([]int64, len(strArr))
	for i, str := range strArr {
		int64Val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		int64Arr[i] = int64Val
	}
	return int64Arr, nil
}

// NewFSIInstance 生成并返回FollowServiceImp结构体单例变量。
func NewFSIInstance() *FollowServiceImp {
	followServiceOnce.Do(
		func() {
			followServiceImp = &FollowServiceImp{
				UserService: NewUserServiceInstance(),
			}
		})
	return followServiceImp
}

//-------------------------------------API IMPLEMENT--------------------------------------------

/*
	关注业务
*/

// FollowAction 关注操作的业务
func (followService *FollowServiceImp) FollowAction(userId int64, targetId int64) (bool, error) {
	followDao := dao.NewFollowDaoInstance()
	follow, err := followDao.FindEverFollowing(userId, targetId)
	// 获取关注的消息队列
	if nil != err {
		return false, err
	}
	if nil != follow {
		//发送消息队列

	}
	return true, nil
}

func (followService *FollowServiceImp) AddToRDBWhenFollow(userId int64, targetId int64) {
	followDao := dao.NewFollowDaoInstance()
	client := redis.Clients.UserFollowings
	// 尝试给following数据库追加user关注target的记录
	keyCnt1, err1 := client.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()
	if err1 != nil {
		log.Println(err1.Error())
	}

	// 只判定键是否不存在，若不存在即从数据库导入
	if keyCnt1 <= 0 {
		userFollowingsId, _, err := followDao.GetFollowingsInfo(userId)
		if err != nil {
			log.Println(err.Error())
			return
		}
		ImportToRDBFollowing(userId, userFollowingsId)
	}
	// 数据库导入到redis结束后追加记录
	client.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), targetId)

	// 尝试给follower数据库追加target的粉丝有user的记录
	keyCnt2, err2 := client.Exists(redis.Ctx, strconv.FormatInt(targetId, 10)).Result()

	if err2 != nil {
		log.Println(err2.Error())
	}

	if keyCnt2 <= 0 {
		//获取target的粉丝，直接刷新，关注时刷新target的粉丝
		userFollowersId, _, err := followDao.GetFollowersInfo(targetId)
		if err != nil {
			log.Println(err.Error())
			return
		}
		ImportToRDBFollower(targetId, userFollowersId)
	}

	client.SAdd(redis.Ctx, strconv.FormatInt(targetId, 10), userId)

	// 进行好友的判定，本接口实现user对target的关注，若此时target也关注了user，进行friend数据库的记录追加
	// user的好友有target，target的好友有user
	if flag, _ := followService.CheckIsFollowing(targetId, userId); flag {
		// 尝试给friend数据库追加user的好友有target的记录
		keyCnt3, err3 := client.Exists(redis.Ctx, strconv.FormatInt(userId, 10)).Result()

		if err3 != nil {
			log.Println(err3.Error())
		}
		if keyCnt3 <= 0 {
			userFriendsId1, _, err := followDao.GetFriendsInfo(userId)
			if err != nil {
				log.Println(err)
				return
			}
			ImportToRDBFriend(userId, userFriendsId1)
		}

		client.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), targetId)

		// 尝试给friend数据库追加target的好友有user的记录
		keyCnt4, err4 := client.Exists(redis.Ctx, strconv.FormatInt(targetId, 10)).Result()

		if err4 != nil {
			log.Println(err4.Error())
		}
		if keyCnt4 <= 0 {
			//获取target的粉丝，直接刷新，关注时刷新target的粉丝
			userFriendsId2, _, err := followDao.GetFriendsInfo(targetId)
			if err != nil {
				log.Println(err)
				return
			}
			ImportToRDBFriend(targetId, userFriendsId2)
		}

		client.SAdd(redis.Ctx, strconv.FormatInt(targetId, 10), userId)
	}
}

/*
	提供目标用户id和对应的id列表导入到redis中的方法，一般用在更新失效键的逻辑中
*/

// ImportToRDBFollowing 将登陆用户的关注id列表导入到following数据库中
func ImportToRDBFollowing(userId int64, ids []int64) {
	client := redis.Clients.UserFollowings
	// 将传入的userId及其关注用户id更新至redis中
	for _, id := range ids {
		client.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), int(id))
	}

	client.Expire(redis.Ctx, strconv.FormatInt(userId, 10), CacheTimeGenerator())
}

func ImportToRDBFriend(userId int64, ids []int64) {
	client := redis.Clients.UserFollowings
	// 将传入的userId及其好友id更新至redis中
	for _, id := range ids {
		client.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), int(id))
	}

	client.Expire(redis.Ctx, strconv.FormatInt(userId, 10), CacheTimeGenerator())
}

// ImportToRDBFollower 将登陆用户的关注id列表导入到follower数据库中
func ImportToRDBFollower(userId int64, ids []int64) {
	client := redis.Clients.UserFollowings
	// 将传入的userId及其粉丝id更新至redis中
	for _, id := range ids {
		client.SAdd(redis.Ctx, strconv.FormatInt(userId, 10), int(id))
	}

	client.Expire(redis.Ctx, strconv.FormatInt(userId, 10), CacheTimeGenerator())
}
