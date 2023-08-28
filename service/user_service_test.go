package service

import (
	"encoding/json"
	"testing"
)

func TestLogin(t *testing.T) {
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

func TestUserService_QuerySelfInfo(t *testing.T) {
	testService := NewUserServiceInstance()
	res, err := testService.QuerySelfInfo(1)
	if err != nil {
		t.Error(err)
		return
	}
	resJson, _ := json.Marshal(res)
	t.Log(string(resJson))
	info, err := testService.QuerySelfInfo(2)
	if err != nil {
		t.Error(err)
		return
	}
	infoJson, _ := json.Marshal(info)
	t.Log(string(infoJson))
}
