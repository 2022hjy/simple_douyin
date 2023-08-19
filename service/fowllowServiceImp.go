package service

import (
	"math/rand"
	"simple_douyin/dao"
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
