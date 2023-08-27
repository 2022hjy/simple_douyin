package dao

import (
	"fmt"
	"sync"

	redisv9 "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"simple_douyin/config"
	"simple_douyin/middleware/database"
	"simple_douyin/middleware/redis"
)

var Db = database.Db

type User struct {
	UserId          int64  `json:"id" gorm:"primaryKey;autoIncrement:true" redis:"user_id"`
	Username        string `json:"name" gorm:"unique;not null" redis:"username"`
	Password        string `json:"-" redis:"-"` // 不返回给前端
	Avatar          string `json:"avatar" gorm:"default:''" redis:"avatar"`
	BackgroundImage string `json:"background_image" gorm:"default:''" redis:"background_image"`
	Signature       string `json:"signature" gorm:"default:''" redis:"signature"`
}

func (User) TableName() string {
	return "user"
}

type UserDao struct {
	db          *gorm.DB
	redisClient *redisv9.Client
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
			userDao.redisClient = redis.Clients.UserId_UserR
		})
	return userDao
}

func (u *UserDao) InsertUser(user *User) bool {
	if err := u.db.Create(user).Error; err != nil {
		return false
	}
	return true
}

func (u *UserDao) GetUserByName(name string) (*User, error) {
	var user User
	if err := u.db.Where("username = ?", name).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserDao) GetUserById(id int64) (*User, error) {
	var user User
	if err := u.db.Where("user_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserDao) GetUsersByIds(ids []int64) ([]*User, error) {
	var users []*User
	if err := u.db.Where("user_id in ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserDao) UpdateUser(user *User) bool {
	if err := u.db.Save(user).Error; err != nil {
		return false
	}
	return true
}

func (u *UserDao) GetUserFromRedisById(userId int64) (*User, error) {
	key := fmt.Sprintf("%s:%d", config.UserId_User_KEY_PREFIX, userId)
	var user User
	if err := redis.GetHash(u.redisClient, key, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserDao) SetUserToRedis(user *User) error {
	key := fmt.Sprintf("%s:%d", config.UserId_User_KEY_PREFIX, user.UserId)
	if err := redis.SetHash(u.redisClient, key, user); err != nil {
		return err
	}
	return nil
}
