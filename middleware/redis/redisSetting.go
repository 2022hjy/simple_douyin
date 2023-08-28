package redis

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"

	redisv9 "github.com/redis/go-redis/v9"
	"simple_douyin/config"
)

// Ctx Path: middleware/redis/redisv9.go
var Ctx = context.Background()

// NilError redisv9.Nil 的别名
var NilError = redisv9.Nil

// keyAccessMap 用于记录 key 的访问次数
var keyAccessMap = make(map[string]int)

var Clients *RedisClients

// MRLock (MutexRedisLock）keyAccessMapLock 用于保护 keyAccessMap
// 注意：MRLock应该是一个指针类型，而不是一个值类型。因为您需要的是一个全局的互斥锁，而不是一个值的副本。
//var MRLock = &sync.Mutex{}

// RedisClients redis 客户端
// 命名规则：[Key]_[Val]R
type RedisClients struct {
	Test               *redisv9.Client //0
	VideoId_CommentIdR *redisv9.Client //1
	CommentId_CommentR *redisv9.Client //3
	//F : favorite
	UserId_FavoriteVideoIdR  *redisv9.Client //4
	VideoId_FavoritebUserIdR *redisv9.Client // 5 todo 用视频ID作为键，存储的是点赞这个视频的所有用户ID。

	VideoId_VideoR *redisv9.Client
	//点赞数 和 被点赞数
	//获赞数：value 对应 total_favorited（前端返回值
	//此方法已经 deprecated 了，不再使用
	//只是为了给大家提个醒，现在的使用的思路是使用上面的UserId_FavoriteVideoIdR的 redis
	//：直接通过获取的关联的 id 集合，从而去获得对应的点赞的视频数目（len 方法去获取长度）
	//UserId_FavoritedNumR *redisv9.Client

	//（给别人的）点赞数：value 对应 favorite_count （前端返回值
	UserId_FavoriteNumR *redisv9.Client
	UserId_UserR        *redisv9.Client
	UserId_FollowersR   *redisv9.Client
	UserId_FollowingsR  *redisv9.Client
	UserId_FriendsR     *redisv9.Client
	UserAllId_MessageR  *redisv9.Client
}

//    获取用户的所有点赞视频：只需查询 UserId_FavoriteVideoIdR 集合。
//    获取点赞某视频的所有用户：只需查询 RdbVUid 集合。
//    检查用户是否点赞了某个视频：查询 UserId_FavoriteVideoIdR 集合中用户ID对应的集合是否包含该视频ID。
//    检查某个视频是否被用户点赞：查询 RdbVUid 集合中视频ID对应的集合是否包含该用户ID。

const (
	ProdRedisAddr = config.RedisAddr
	ProRedisPwd   = config.RedisPwd
)

func init() {
	InitRedis()
}

// InitRedis 初始化 Redis 连接
func InitRedis() {
	Clients = &RedisClients{
		Test: redisv9.NewClient(&redisv9.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       0,
		}),
		VideoId_CommentIdR: redisv9.NewClient(&redisv9.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       1,
		}),
		//卧槽，改 bug 改了一天，原来是这里的问题，我把这里的 DB 改成了 2是被注释掉的，因此在表的监控软件里面应该使用的是 3
		//CVid: redisv9.NewClient(&redisv9.Options{
		//	Addr:     ProdRedisAddr,
		//	Password: ProRedisPwd,
		//	DB:       2,
		//}),
		CommentId_CommentR: redisv9.NewClient(&redisv9.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       3,
		}),
		UserId_FavoriteVideoIdR: redisv9.NewClient(&redisv9.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       4,
		}),
		UserId_UserR: redisv9.NewClient(&redisv9.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       5,
		}),
		UserId_FollowingsR: redisv9.NewClient(&redisv9.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       11,
		}),
		UserId_FollowersR: redisv9.NewClient(&redisv9.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       12,
		}),
		UserId_FriendsR: redisv9.NewClient(&redisv9.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       13,
		}),
		VideoId_VideoR: redisv9.NewClient(&redisv9.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       14,
		}),
		UserAllId_MessageR: redisv9.NewClient(&redisv9.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       15,
		}),
	}
}

/*
思路：一开始去 set 数值在 redis 中的时候，我们采取随机数种子随机产生一个短时间的过期的时间数值。
当这个这个数值（键值对）被采取了 get 操作的时候，进行一个缓存的过期时长的一个升级。
当被访问一定次数的时候，缓存的过期时长将升级成永不过期。
同时，如果出现了缓存穿透的情况，可以采取缓存空数据的方式，将空数据缓存一段时间，防止缓存穿透。在代码GetKeyAndUpdateExpiration内部也已经实现
*/

// SetValueWithRandomExp 设置 Redis 集合的批量添加，采取过期时间随机
func SetValueWithRandomExp(client *redisv9.Client, key string, value interface{}) error {
	log.Println("正在调用SetValueWithRandomExp方法进行 Redis 的信息的存储，key:", key, "value:", value)
	if client == nil {
		return errors.New("client is nil")
	}
	rand.Seed(time.Now().UnixNano())
	//exp := time.Duration(rand.Intn(20)) * time.Minute

	SetErr := client.Set(Ctx, key, value, 0).Err()
	if SetErr != nil {
		log.Fatalf("redis set error: %v\n", SetErr)
	}
	return SetErr
}

func PushValueToListR(client *redisv9.Client, key string, value interface{}) error {
	log.Println("正在调用PushValueToListR 方法进行 Redis 的信息的存储，key:", key, "value:", value)
	if client == nil {
		return errors.New("client is nil")
	}
	//rand.Seed(time.Now().UnixNano())
	//exp := time.Duration(rand.Intn(20)) * time.Minute

	_, SetErr := client.RPush(Ctx, key, value).Result()
	if SetErr != nil {
		log.Printf("redis RPush set error: %v\n\n", SetErr)
		return errors.New("redis set error: " + SetErr.Error())
	}
	return SetErr
}

// SetHashWithExpiration 设置 Redis 哈希表，并设置过期时间
func SetHashWithExpiration(client *redisv9.Client, key string, data map[string]interface{}, exp time.Duration) error {
	SetErr := SetHash(client, key, data)
	if SetErr != nil {
		return SetErr
	}
	return SetExpiration(client, key, exp)
}

// SetExpiration 设置 Redis 键值对的过期时间
func SetExpiration(client *redisv9.Client, key string, exp time.Duration) error {
	if client == nil {
		return errors.New("client is nil")
	}
	ExpErr := client.Expire(Ctx, key, exp).Err()
	if ExpErr != nil {
		log.Fatalf("redis set error: %v\n", ExpErr)
	}
	return ExpErr
}

// GetHash 获取 Redis Hash结构
// 参数client 为 redis 客户端，key 为键，data 为存储数据的结构体的指针
func GetHash(client *redisv9.Client, key string, dataPtr interface{}) error {
	if client == nil {
		return errors.New("client is nil")
	}
	// 1. 首先判断键是否存在
	isExist, err := client.Exists(Ctx, key).Result()
	if err != nil {
		return err
	}
	if isExist == 0 {
		return errors.New("key does not exist")
	}
	if err := client.HGetAll(Ctx, key).Scan(dataPtr); err != nil {
		return err
	}
	return nil
}

// SetHash 设置 Redis 哈希表
func SetHash(client *redisv9.Client, key string, data interface{}) error {
	if client == nil {
		return errors.New("client is nil")
	}
	if SetErr := client.HSet(Ctx, key, data).Err(); SetErr != nil {
		log.Fatalf("redis set hash error: %v\n", SetErr)
		return SetErr
	}
	return nil
}

type HotData struct {
	ExpireTime time.Time
	Data       interface{}
}

func SetHashWithLogic(client *redisv9.Client, key string, data interface{}, exp time.Duration) error {
	if client == nil {
		return errors.New("client is nil")
	}
	hotData := &HotData{
		ExpireTime: time.Now().Add(exp),
		Data:       data,
	}
	if err := SetHash(client, key, hotData); err != nil {
		return err
	}
	return nil
}

func GetHashWithLogic(client *redisv9.Client, key string, dataPtr interface{}) error {
	if client == nil {
		return errors.New("client is nil")
	}
	hotData := &HotData{}
	if err := GetHash(client, key, hotData); err != nil {
		return err
	}
	if hotData.ExpireTime.Before(time.Now()) {
		return errors.New("data expired")
	}
	dataPtr = hotData.Data
	return nil
}

/*

似乎没有必要，因为在单个数值添加也没有关系，只有在批量添加的时候才有必要，
addErr := client.SAdd(Ctx, key, value, exp).Err()
	if addErr != nil {
		log.Fatalf("redis set error: %v\n", addErr)
		return addErr
	}
*/

// Deprecate： GetKeyAndUpdateExpiration 获取 Redis 中的值，并更新（膨胀）过期时间
//func GetKeyAndUpdateExpiration(client *redisv9.Client, key string) (interface{}, error) {
//	val, err := client.Get(Ctx, key).Result()
//	if err != nil {
//		if err == redisv9.Nil {
//			// 缓存穿透，将空数据缓存一段时间，防止缓存穿透
//			SetValueWithRandomExp(client, key, "")
//			//SetValueWithRandomExp(client, key, "nil")
//			return "", nil
//		}
//		return "", err
//	}
//
//	MRLock.Lock()
//	keyAccessMap[key]++
//	count := keyAccessMap[key]
//	MRLock.Unlock()
//
//	if count >= config.THREASHOLD {
//		client.Persist(Ctx, key) // 删除过期时间，使键永不过期
//	} else {
//		// 增加过期时间
//		client.Expire(Ctx, key, 2*time.Minute)
//	}
//	return val, nil
//}

// GetKeysAndUpdateExpiration 获取与给定模式匹配的所有Redis中的值，并更新（膨胀）它们的过期时间，支持返回一对一，一对多的键值对
// 该函数返回一个interface，因此需要在使用的时候使用类型断言
func GetKeysAndUpdateExpiration(client *redisv9.Client, key string) (interface{}, error) {
	keys, err := client.Keys(Ctx, key).Result()
	log.Println("在GetKeysAndUpdateExpiration内部获得的keys:", keys)
	if err != nil {
		return nil, err
	}

	// 没有匹配的键
	if len(keys) == 0 {
		return nil, nil
	}

	values, err := client.MGet(Ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	results := make(map[string]string)
	for i, key := range keys {
		log.Printf("正在处理的键值对：key == %v value == %v ", key, values[i])
		val, ok := values[i].(string)
		if !ok {
			//MRLock.Unlock()
			return nil, errors.New("value is not of type string")
		}
		results[key] = val
		keyAccessMap[key]++
		count := keyAccessMap[key]
		if count >= config.THREASHOLD {
			client.Persist(Ctx, key) // 删除过期时间，使键永不过期
		} else {
			// 增加过期时间
			client.Expire(Ctx, key, 2*time.Minute)
		}
	}

	// 如果只有一个键值对，直接返回键值对
	if len(results) == 1 {
		for _, value := range results {
			return value, nil
		}
	}
	return results, nil
}

// DeleteKey 删除 Redis 中的键值对，线程安全
func DeleteKey(client *redisv9.Client, key string) error {
	if client == nil {
		return errors.New("client is nil")
	}
	//MRLock.Lock()
	delete(keyAccessMap, key)
	//MRLock.Unlock()
	return client.Del(Ctx, key).Err()
}

// IsKeyExist 判断 Redis 中是否存在某个键
// 返回值：bool，error。如果存在，返回 true，nil；
// 如果不存在，返回 false，nil；如果出错，返回 false，error
func IsKeyExist(client *redisv9.Client, key string) (bool, error) {
	if client == nil {
		return false, errors.New("client is nil")
	}

	isExist, err := client.Exists(Ctx, key).Result()
	if err != nil {
		return false, err
	}

	return isExist == 1, nil
}

// CountElements 获取与给定键关联的 Redis 集合的元素数量
func CountElements(client *redisv9.Client, key string) (int64, error) {
	if client == nil {
		return 0, errors.New("客户端为空")
	}

	count, err := client.SCard(Ctx, key).Result()
	if err != nil {
		log.Printf("获取 Redis 集合元素数量失败: %v\n", err)
		return 0, err
	}
	return count, nil
}

// GetValue 获取 Redis 中的值（常规）
func GetValue(client *redisv9.Client, key string) (string, error) {
	if client == nil {
		return "", errors.New("client is nil")
	}
	return client.Get(Ctx, key).Result()
}

// SetValue 设置 Redis 键值对
func SetValue(client *redisv9.Client, key string, value interface{}) error {
	if client == nil {
		return errors.New("client is nil")
	}
	// 设置 2 min 过期，如果 expiration 为 0 表示永不过期
	return client.Set(Ctx, key, value, 2*time.Minute).Err()
}
