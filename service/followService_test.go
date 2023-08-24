package service

import (
	"fmt"
	"log"
	"simple_douyin/middleware/database"
	"simple_douyin/middleware/mq"
	"simple_douyin/middleware/redis"
	"testing"
)

func TestFollowServiceImp_GetFollowings(t *testing.T) {
	database.Init()
	redis.InitRedis()
	followings, err := followServiceImp.GetFollowings(1)

	if err != nil {
		log.Default()
	}
	fmt.Println(followings)
}

func TestFollowServiceImp_GetFollowers(t *testing.T) {
	database.Init()
	redis.InitRedis()
	followers, err := followServiceImp.GetFollowers(2)

	if err != nil {
		log.Default()
	}
	fmt.Println(followers)
}

func TestFollowServiceImp_GetFollowingCnt(t *testing.T) {
	database.Init()
	redis.InitRedis()
	userIdCnt, err := followServiceImp.GetFollowingCnt(7)
	if err != nil {
		log.Default()
	}
	fmt.Println(userIdCnt)
}

func TestFollowServiceImp_GetFollowerCnt(t *testing.T) {
	database.Init()
	redis.InitRedis()
	userIdCnt, err := followServiceImp.GetFollowerCnt(1)
	if err != nil {
		log.Default()
	}
	fmt.Println(userIdCnt)
}

func TestFollowServiceImp_CheckIsFollowing(t *testing.T) {
	redis.InitRedis()
	database.Init()
	var err error
	result, err := followServiceImp.CheckIsFollowing(1, 5)
	if err != nil {
		log.Default()
	}
	fmt.Println(result)
}

func TestFollowServiceImp_FollowAction(t *testing.T) {
	redis.InitRedis()
	database.Init()
	mq.InitMq()
	result, err := followServiceImp.FollowAction(1, 4)
	if err != nil {
		log.Default()
	}
	fmt.Println(result)
}
