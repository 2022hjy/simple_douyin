package dao

import (
	"testing"

	"simple_douyin/middleware/database"
	"simple_douyin/model"
)

func TestUserDao_UpdateUser(t *testing.T) {
	database.Init()
	userDao := NewUserDaoInstance()
	user := &model.User{
		Username: "test",
		Password: "test",
	}
	userDao.InsertUser(user)
	user, err := userDao.GetUserById(user.UserId)
	if err != nil {
		t.Error(err)
	}
	user.Password = "test2"
	userDao.UpdateUser(user)
	user, err = userDao.GetUserById(user.UserId)
	t.Log(user)
}
