package service

import (
	"encoding/json"
	"testing"

	"simple_douyin/middleware/database"
	"simple_douyin/model"
)

func TestLogin(t *testing.T) {
	database.Init()
	req := model.LoginRequest{
		UserName: "guest",
		Password: "guest",
	}
	service := NewUserServiceInstance()
	res := service.Login(req)
	// 将结构体转为json字符串
	resJson, _ := json.Marshal(res)
	t.Log(string(resJson))
}

func TestRegister(t *testing.T) {
	database.Init()
	req := model.RegisterRequest{
		UserName: "guest",
		Password: "guest",
	}
	service := NewUserServiceInstance()
	res := service.Register(req)
	resJson, _ := json.Marshal(res)
	t.Log(string(resJson))
}

func TestUserInfo(t *testing.T) {
	database.Init()
	service := NewUserServiceInstance()
	res := service.UserInfo(7)
	resJson, _ := json.Marshal(res)
	t.Log(string(resJson))
}
