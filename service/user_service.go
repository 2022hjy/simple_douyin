package service

import (
	"errors"
	"sync"

	redisv9 "github.com/redis/go-redis/v9"
	"simple_douyin/controller"
	"simple_douyin/dao"
	"simple_douyin/model"
	"simple_douyin/util"
)

type UserService struct {
	userDao         *dao.UserDao
	redisUserFollow *redisv9.Client
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
			//userService.redisUserFollow = redis.Clients.UserFollowings
		})
	return userService
}

type Credential struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

func (u *UserService) Login(request controller.LoginRequest) (*Credential, error) {
	user, err := u.userDao.GetUserByName(request.UserName)
	if err != nil {
		return nil, errors.New("user not exists")
	}
	success := util.ValidatePassword(user.Password, request.Password)
	if !success {
		return nil, errors.New("password is wrong")
	}
	token, err := util.GenerateToken(user.UserId, user.Username)
	if err != nil {
		return nil, errors.New("generate token failed")
	}
	return &Credential{
		Token:  token,
		UserId: user.UserId,
	}, nil
}

func (u *UserService) Register(request controller.RegisterRequest) (*Credential, error) {
	// 1. 从数据库中查询用户是否存在
	user, err := u.userDao.GetUserByName(request.UserName)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	// 2. 对用户输入的密码进行加密
	password, err := util.EncryptPassword(request.Password)
	if err != nil {
		return nil, errors.New("encrypt password failed")
	}
	// 3. 将用户信息插入到数据库中
	user = &dao.User{
		Username:        request.UserName,
		Password:        password,
		Avatar:          "",
		BackgroundImage: "",
		Signature:       "",
	}
	ret := u.userDao.InsertUser(user)
	if !ret {
		return nil, errors.New("insert user failed")
	}
	token, err := util.GenerateToken(user.UserId, user.Username)
	if err != nil {
		return nil, errors.New("generate token failed")
	}
	return &Credential{
		Token:  token,
		UserId: user.UserId,
	}, nil
}

type UserInfo struct {
	*dao.User

	FollowCount   int64 `json:"follow_count,omitempty"`
	FollowerCount int64 `json:"follower_count,omitempty"`
	IsFollow      bool  `json:"is_follow,omitempty"`

	TotalFavorited int64 `json:"total_favorited,omitempty"`
	WorkCount      int64 `json:"work_count,omitempty"`
	FavoriteCount  int64 `json:"favorite_count,omitempty"`
}

func QueryUserInfo(userId int64) (*UserInfo, error) {
	return NewQueryUserInfoFlow(userId).Do()
}

func NewQueryUserInfoFlow(userId int64) *QueryUserInfoFlow {
	return &QueryUserInfoFlow{
		userId: userId,
	}
}

type QueryUserInfoFlow struct {
	userId   int64
	userInfo *UserInfo

	user           *dao.User
	followCount    int64
	followerCount  int64
	isFollow       bool
	totalFavorited int64
	workCount      int64
	favoriteCount  int64
}

func (f *QueryUserInfoFlow) Do() (*UserInfo, error) {
	err := f.prepareInfo()
	if err != nil {
		return nil, err
	}
	err = f.packageInfo()
	if err != nil {
		return nil, err
	}
	return f.userInfo, nil
}

func (f *QueryUserInfoFlow) prepareInfo() error {
	var wg sync.WaitGroup
	wg.Add(3)
	errChan := make(chan error, 7)
	go func() {
		defer wg.Done()
		// 1. 先从redis中获取用户信息
		// 2. 如果redis中没有，则从数据库中获取
		user, err := dao.NewUserDaoInstance().GetUserById(f.userId)
		if err != nil {
			errChan <- err
			return
		}
		f.user = user
	}()
	go func() {
		defer wg.Done()
		// 获取用户的关注数、粉丝数、是否关注

	}()
	go func() {
		defer wg.Done()
		f.isFollow = false
	}()
	go func() {
		defer wg.Done()
		f.totalFavorited = 0
	}()
	go func() {
		defer wg.Done()
		videoCnt, err := dao.GetVideoCnt(f.userId)
		if err != nil {
			errChan <- err
			return
		}
		f.workCount = videoCnt
	}()
	go func() {
		defer wg.Done()
		f.favoriteCount = 0
	}()
	wg.Wait()
	select {
	case err := <-errChan:
		return err
	default:
	}
	return nil
}

func (f *QueryUserInfoFlow) packageInfo() error {
	f.userInfo = &model.UserInfo{
		User:           f.user,
		FollowerCount:  f.followerCount,
		FollowCount:    f.followCount,
		IsFollow:       f.isFollow,
		TotalFavorited: f.totalFavorited,
		WorkCount:      f.workCount,
		FavoriteCount:  f.favoriteCount,
	}
	return nil
}
