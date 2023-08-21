package dao

import (
	"sync"

	"gorm.io/gorm"
	"simple_douyin/middleware/database"
)

var Db = database.Db

type User struct {
	UserId          int64  `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Username        string `json:"name" gorm:"unique;not null"`
	Password        string `json:"-"` // 不返回给前端
	Avatar          string `json:"avatar,omitempty" gorm:"default:''"`
	BackgroundImage string `json:"background_image,omitempty" gorm:"default:''"`
	Signature       string `json:"signature,omitempty" gorm:"default:''"`
}

func (User) TableName() string {
	return "user"
}

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
