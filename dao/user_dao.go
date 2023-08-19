package dao

import (
	"sync"

	"gorm.io/gorm"
	"simple_douyin/middleware/database"
	"simple_douyin/model"
)

var Db = database.Db

type UserDao struct {
	db *gorm.DB
}

var (
	userDao  *UserDao
	userOnce sync.Once
)

func NewUserDaoInstance() *UserDao {
	userOnce.Do(
		func() {
			userDao = &UserDao{}
			userDao.db = database.Db
		})
	return userDao
}

func (u *UserDao) InsertUser(user *model.User) bool {
	if err := u.db.Create(user).Error; err != nil {
		return false
	}
	return true
}

func (u *UserDao) GetUserByName(name string) (model.User, error) {
	user := model.User{}
	if err := u.db.Where("username = ?", name).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (u *UserDao) GetUserById(id int64) (model.User, error) {
	user := model.User{}
	if err := u.db.Where("user_id = ?", id).First(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (u *UserDao) GetUsersByIds(ids []int64) ([]model.User, error) {
	users := make([]model.User, 0, len(ids))
	if err := u.db.Where("user_id in ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
