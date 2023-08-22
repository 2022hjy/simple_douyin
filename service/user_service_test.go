package service

import (
	"encoding/json"
	"testing"

	"simple_douyin/middleware/database"
	"simple_douyin/middleware/redis"
)

func TestLogin(t *testing.T) {
	database.Init()
	req := LoginInfo{
		UserName: "guest",
		Password: "guest",
	}
	testService := NewUserServiceInstance()
	res, err := testService.Login(req)
	if err != nil {
		t.Error(err)
		return
	}
	// 将结构体转为json字符串
	resJson, _ := json.Marshal(res)
	t.Log(string(resJson))
}

func TestRegister(t *testing.T) {

}

func TestGetUserInfo(t *testing.T) {
	database.Init()
	redis.InitRedis()

}
