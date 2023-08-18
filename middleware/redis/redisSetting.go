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

// MRLock (MutexRedisLock）keyAccessMapLock 用于保护 keyAccessMap
// 注意：MRLock应该是一个指针类型，而不是一个值类型。因为您需要的是一个全局的互斥锁，而不是一个值的副本。
var MRLock = &sync.Mutex{}

type RedisClients struct {
	Test           *redis.Client
	VCid           *redis.Client
	CVid           *redis.Client
	CIdComment     *redis.Client
	UVid           *redis.Client
	VUid           *redis.Client
	UserFollowers  *redis.Client
	UserFollowings *redis.Client
	UserFriends    *redis.Client
}

var Clients *RedisClients

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
		VCid: redis.NewClient(&redis.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       1,
		}),
		CVid: redis.NewClient(&redis.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       2,
		}),
		CIdComment: redis.NewClient(&redis.Options{
			Addr:     ProdRedisAddr,
			Password: ProRedisPwd,
			DB:       3,
		}),
		UVid: redis.NewClient(&redis.Options{
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

// SetValue 设置 Redis 键值对
func SetValue(client *redis.Client, key string, value interface{}) error {
	if client == nil {
		return errors.New("client is nil")
	}
	// 设置 2 min 过期，如果 expiration 为 0 表示永不过期
	return client.Set(Ctx, key, value, 2*time.Minute).Err()
}

/*
思路：一开始去 set 数值在 redis 中的时候，我们采取随机数种子随机产生一个短时间的过期的时间数值。
当这个这个数值（键值对）被采取了 get 操作的时候，进行一个缓存的过期时长的一个升级。
当被访问一定次数的时候，缓存的过期时长将升级成永不过期。
同时，如果出现了缓存穿透的情况，可以采取缓存空数据的方式，将空数据缓存一段时间，防止缓存穿透。在代码GetKeyAndUpdateExpiration内部也已经实现
*/

// SetValueWithRandomExp 设置 Redis 键值对，过期时间随机
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

// GetKeyAndUpdateExpiration 获取 Redis 中的值，并更新（膨胀）过期时间
func GetKeyAndUpdateExpiration(client *redis.Client, key string) (string, error) {
	val, err := client.Get(Ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// 缓存穿透，将空数据缓存一段时间，防止缓存穿透
			SetValueWithRandomExp(client, key, "")
			//SetValueWithRandomExp(client, key, "nil")
			return "", nil
		}
		return "", err
	}

	MRLock.Lock()
	keyAccessMap[key]++
	count := keyAccessMap[key]
	MRLock.Unlock()

	if count >= config.THREASHOLD {
		client.Persist(Ctx, key) // 删除过期时间，使键永不过期
	} else {
		// 增加过期时间
		client.Expire(Ctx, key, 2*time.Minute)
	}
	return val, nil
}

// GetValue 获取 Redis 中的值（常规）
func GetValue(client *redis.Client, key string) (string, error) {
	if client == nil {
		return "", errors.New("client is nil")
	}
	return client.Get(Ctx, key).Result()
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

	// 通过 GetKeyAndUpdateExpiration 获取 Redis 中的值，并更新（膨胀）过期时间
	value, err = GetKeyAndUpdateExpiration(Clients.Test, key)
	if err != nil {
		log.Printf("Error getting value: %v", err)
	} else {
		log.Printf("Value is: %s", value)
	}
}
