package dao

import (
    "simple_douyin/middleware/database"
    "simple_douyin/model"
)

var Db = database.Db

func InsertUser(user *model.User) bool {
    if err := Db.Create(user).Error; err != nil {
        return false
    }
    return true
}

func GetUserByName(name string) (model.User, error) {
    user := model.User{}
    if err := Db.Where("username = ?", name).First(&user).Error; err != nil {
        return user, err
    }
    return user, nil
}

func GetUserById(id int64) (model.User, error) {
    user := model.User{}
    if err := Db.Where("user_id = ?", id).First(&user).Error; err != nil {
        return user, err
    }
    return user, nil
}

func GetUsersByIds(ids []int64) ([]model.User, error) {
    users := make([]model.User, 0, len(ids))
    if err := Db.Where("user_id in ?", ids).Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}
