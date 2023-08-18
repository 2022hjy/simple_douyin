package service

import (
    "encoding/json"
    "testing"

    "simple_douyin/model"
)

func TestLogin(t *testing.T) {
    req := model.LoginRequest{
        UserName: "guest",
        Password: "guest",
    }
    res := Login(req)
    // 将结构体转为json字符串
    resJson, _ := json.Marshal(res)
    t.Log(string(resJson))
}

func TestRegister(t *testing.T) {
    req := model.RegisterRequest{
        UserName: "guest",
        Password: "guest",
    }
    res := Register(req)
    resJson, _ := json.Marshal(res)
    t.Log(string(resJson))
}

func TestUserInfo(t *testing.T) {
    res := UserInfo(7)
    resJson, _ := json.Marshal(res)
    t.Log(string(resJson))
}
