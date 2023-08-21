package service

import (
	"simple_douyin/config"
	"simple_douyin/middleware/redis"
	"simple_douyin/model"
	"strconv"
)

func (u *UserService) GetUserFromRedisByUserId(userId int64) (*model.User, error) {
	UIdUClients := redis.Clients.UserId_UserR
	key := config.UserId_User_KEY_PREFIX + strconv.FormatInt(userId, 10)
	user, err := redis.GetKeysAndUpdateExpiration(UIdUClients, key)
	if err != nil {
		return nil, newError(ErrRedisGet, err)
	}
	User, ok := user.(model.User)
	if !ok {
		return nil,
	}

}
