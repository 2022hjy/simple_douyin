package service

import (
	"sync"

	"simple_douyin/dao"
	"simple_douyin/model"
	"simple_douyin/util"
)

type UserService struct {
	userDao *dao.UserDao
}

var (
	userService     *UserService
	userServiceOnce sync.Once
)

func NewUserServiceInstance() *UserService {
	userServiceOnce.Do(
		func() {
			userService = &UserService{}
			userService.userDao = dao.NewUserDaoInstance()
		})
	return userService
}

func (u *UserService) Login(request model.LoginRequest) *model.LoginResponse {
	user, err := u.userDao.GetUserByName(request.UserName)
	if err != nil {
		return &model.LoginResponse{
			Response: model.ErrorResponse(model.ErrorCode, "User not exists"),
		}
	}
	success := util.ValidatePassword(user.Password, request.Password)
	if !success {
		return &model.LoginResponse{
			Response: model.ErrorResponse(model.ErrorCode, "Wrong password"),
		}
	}
	token, err := util.GenerateToken(user.UserId, user.Username)
	if err != nil {
		return &model.LoginResponse{
			Response: model.ErrorResponse(model.ErrorCode, "Generate token failed"),
		}
	}
	return &model.LoginResponse{
		Response: model.SuccessResponse(),
		Token:    token,
		UserId:   user.UserId,
	}
}

func (u *UserService) Register(request model.RegisterRequest) *model.RegisterResponse {
	// 1. 从数据库中查询用户是否存在
	user, err := u.userDao.GetUserByName(request.UserName)
	if err == nil {
		return &model.RegisterResponse{
			Response: model.ErrorResponse(model.ErrorCode, "User already exists"),
		}
	}
	// 2. 对用户输入的密码进行加密
	password, err := util.EncryptPassword(request.Password)
	if err != nil {
		return &model.RegisterResponse{
			Response: model.ErrorResponse(model.ErrorCode, "Encrypt password failed"),
		}
	}
	// 3. 将用户信息插入到数据库中
	user = model.User{
		Username:        request.UserName,
		Password:        password,
		Avatar:          "",
		BackgroundImage: "",
		Signature:       "",
	}
	ret := u.userDao.InsertUser(&user)
	if !ret {
		return &model.RegisterResponse{
			Response: model.ErrorResponse(model.ErrorCode, "Create user failed"),
		}
	}
	token, err := util.GenerateToken(user.UserId, user.Username)
	if err != nil {
		return &model.RegisterResponse{
			Response: model.ErrorResponse(model.ErrorCode, "Generate token failed"),
		}
	}
	return &model.RegisterResponse{
		Response: model.SuccessResponse(),
		Token:    token,
		UserId:   user.UserId,
	}
}

func (u *UserService) UserInfo(userId int64) *model.UserResponse {
	user, err := u.userDao.GetUserById(userId)
	if err != nil {
		return &model.UserResponse{
			Response: model.ErrorResponse(model.ErrorCode, "User not exists"),
		}
	}
	// todo 获取用户的关注数、粉丝数、是否关注，作品数、获赞数、收藏数

	return &model.UserResponse{
		Response: model.SuccessResponse(),
		User:     user,
	}
}
