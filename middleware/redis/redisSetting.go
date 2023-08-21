package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand"
	"simple_douyin/config"
	"sync"
	"time"
)

// Ctx Path: middleware/redis/redis.go
var Ctx = context.Background()

// NilError redis.Nil 的别名
var NilError = redis.Nil

// keyAccessMap 用于记录 key 的访问次数
var keyAccessMap = make(map[string]int)

var Clients *RedisClients

// MRLock (MutexRedisLock）keyAccessMapLock 用于保护 keyAccessMap
// 注意：MRLock应该是一个指针类型，而不是一个值类型。因为您需要的是一个全局的互斥锁，而不是一个值的副本。
var MRLock = &sync.Mutex{}

// RedisClients redis 客户端
// 命名规则：[Key]_[Val]R
type RedisClients struct {
	Test               *redis.Client
	VideoId_CommentIdR *redis.Client
	CommentId_CommentR *redis.Client
	//F : favorite
	UserId_FVideoIdR *redis.Client
	VideoId_VideoR   *redis.Client
	VUid             *redis.Client
	UserFollowers    *redis.Client
	UserFollowings   *redis.Client
	UserFriends      *redis.Client
}

const (
	ProdRedisAddr = config.RedisAddr
	ProRedisPwd   = config.RedisPwd
)

// InitRedis 初始化 Redis 连接
func InitRedis() {
	Clients = &RedisClients{
		Test: redis.NewClient(&redis.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       0,
		}),
		VideoId_CommentIdR: redis.NewClient(&redis.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       1,
		}),
		//CVid: redis.NewClient(&redis.Options{
		//	Addr:     ProdRedisAddr,
		//	Password: ProRedisPwd,
		//	DB:       2,
		//}),
		CommentId_CommentR: redis.NewClient(&redis.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       3,
		}),
		UserId_FVideoIdR: redis.NewClient(&redis.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       4,
		}),
		VUid: redis.NewClient(&redis.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       5,
		}),
		UserFollowings: redis.NewClient(&redis.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       11,
		}),
		UserFollowers: redis.NewClient(&redis.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       12,
		}),
		UserFriends: redis.NewClient(&redis.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       13,
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
func SetValueWithRandomExp(client *redis.Client, key string, value interface{}) error {
	if client == nil {
		return errors.New("client is nil")
	}
	rand.Seed(time.Now().UnixNano())
	exp := time.Duration(rand.Intn(20)) * time.Minute

	SetErr := client.Set(Ctx, key, value, exp).Err()
	if SetErr != nil {
		log.Fatalf("redis set error: %v\n", SetErr)
	}
	return SetErr
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
//func GetKeyAndUpdateExpiration(client *redis.Client, key string) (interface{}, error) {
//	val, err := client.Get(Ctx, key).Result()
//	if err != nil {
//		if err == redis.Nil {
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
func GetKeysAndUpdateExpiration(client *redis.Client, pattern string) (interface{}, error) {
	keys, err := client.Keys(Ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	// 没有匹配的键
	if len(keys) == 0 {
		return nil, errors.New("no keys match the pattern")
	}

	values, err := client.MGet(Ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	results := make(map[string]string)
	//MRLock.Lock()
	for i, key := range keys {
		val, ok := values[i].(string)
		if !ok {
			MRLock.Unlock()
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
	//MRLock.Unlock()

	// 如果只有一个键值对，直接返回键值对
	if len(results) == 1 {
		for _, value := range results {
			return value, nil
		}
	}
	return results, nil
}

// DeleteKey 删除 Redis 中的键值对，线程安全
func DeleteKey(client *redis.Client, key string) error {
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
func IsKeyExist(client *redis.Client, key string) (bool, error) {
	if client == nil {
		return false, errors.New("client is nil")
	}

	isExist, err := client.Exists(Ctx, key).Result()
	if err != nil {
		return false, err
	}

	return isExist == 1, nil
}

// GetValue 获取 Redis 中的值（常规）
func GetValue(client *redis.Client, key string) (string, error) {
	if client == nil {
		return "", errors.New("client is nil")
	}
	return client.Get(Ctx, key).Result()
}

// SetValue 设置 Redis 键值对
func SetValue(client *redis.Client, key string, value interface{}) error {
	if client == nil {
		return errors.New("client is nil")
	}
	// 设置 2 min 过期，如果 expiration 为 0 表示永不过期
	return client.Set(Ctx, key, value, 2*time.Minute).Err()
}

func testRedis() {
	// 初始化 Redis 连接
	InitRedis()

	// 使用 Redis 客户端
	key := "my-key"
	err := SetValue(Clients.Test, key, "my-value")
	if err != nil {
		log.Printf("Error setting value: %v", err)
	}

	value, err := GetValue(Clients.Test, key)
	if err != nil {
		log.Printf("Error getting value: %v", err)
	} else {
		log.Printf("Value is: %s", value)
	}

	// 通过 SetValueWithRandomExp 设置 Redis 键值对，过期时间随机
	SetValueWithRandomExp(Clients.Test, key, "my-value")

	// 通过 GetKeysAndUpdateExpiration 获取 Redis 中的值，并更新（膨胀）过期时间
	valueInterface, err := GetKeysAndUpdateExpiration(Clients.Test, key)
	if err != nil {
		log.Printf("Error getting value: %v", err)
	} else {
		log.Printf("Value is: %s", value)
	}

	value, ok := valueInterface.(string)
	if !ok {
		log.Printf("Value is not of type string")
	} else {
		log.Printf("Value is: %s", value)
	}
}
