package service

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"sync"

	"simple_douyin/config"
	"simple_douyin/dao"
	"simple_douyin/middleware/redis"
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

type LoginInfo struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type Credential struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

func (u *UserService) Login(info LoginInfo) (*Credential, error) {
	log.Printf("info: %v", info)
	user, err := u.userDao.GetUserByName(info.UserName)
	if err != nil {
		return nil, errors.New("user not exists")
	}
	success := util.ValidatePassword(user.Password, info.Password)
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

func (u *UserService) Register(info LoginInfo) (*Credential, error) {
	// 1. 从数据库中查询用户是否存在
	user, err := u.userDao.GetUserByName(info.UserName)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	// 2. 对用户输入的密码进行加密
	password, err := util.EncryptPassword(info.Password)
	if err != nil {
		return nil, errors.New("encrypt password failed")
	}
	// 3. 将用户信息插入到数据库中
	user = &dao.User{
		Username:        info.UserName,
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

func (*UserService) QueryUserInfo(userId int64, tokenUserId int64) (*UserInfo, error) {
	return NewQueryUserInfoFlow(userId, tokenUserId).Do()
}

func (*UserService) QuerySelfInfo(userId int64) (*UserInfo, error) {
	return NewQueryUserInfoFlow(userId, userId).Do()
}

func NewQueryUserInfoFlow(userId int64, tokenUserId int64) *QueryUserInfoFlow {
	return &QueryUserInfoFlow{
		userId:      userId,
		tokenUserId: tokenUserId,
	}
}

type QueryUserInfoFlow struct {
	userId      int64
	tokenUserId int64
	userInfo    *UserInfo

	user           *dao.User
	followCount    int64
	followerCount  int64
	isFollow       bool
	totalFavorited int64
	workCount      int64
	favoriteCount  int64
}

func (f *QueryUserInfoFlow) Do() (*UserInfo, error) {
	if err := f.prepareInfo(); err != nil {
		return nil, err
	}
	if err := f.packageInfo(); err != nil {
		return nil, err
	}
	return f.userInfo, nil
}

func (f *QueryUserInfoFlow) prepareInfo() error {
	var wg sync.WaitGroup
	wg.Add(3)
	errChan := make(chan error, 2)
	go func() {
		defer wg.Done()
		userDao := dao.NewUserDaoInstance()
		// 1. 先从redis中获取用户信息
		if user, err := userDao.GetUserFromRedisById(f.userId); err == nil {
			f.user = user
			return
		}
		// todo 日志记录一下为什么从redis中获取失败
		// 2. 如果redis中没有，则从数据库中获取
		user, err := userDao.GetUserById(f.userId)
		if err != nil {
			errChan <- err
			return
		}
		// 3. 将用户信息写入到redis中，即使写入失败，也不影响后续的流程
		if err = userDao.SetUserToRedis(user); err != nil {
			// todo 日志记录一下为什么写入redis失败，并且需要设置一个过期时间，或者直接就不过期
		}
		f.user = user
	}()
	go func() {
		defer wg.Done()
		// 获取用户的关注数、粉丝数、是否关注
		followService := NewFSIInstance()
		followCount, err := followService.GetFollowerCnt(f.userId)
		if err == nil {
			f.followerCount = followCount
		}
		// todo 日志记录一下为什么获取关注数失败
		followingCount, err := followService.GetFollowingCnt(f.userId)
		if err == nil {
			f.followCount = followingCount
		}
		// todo 日志记录一下为什么获取粉丝数失败
		isFollow, err := followService.CheckIsFollowing(f.userId, f.tokenUserId)
		if err != nil {
			errChan <- err
			return
		}
		f.isFollow = isFollow
	}()
	go func() {
		defer wg.Done()
		f.totalFavorited = 0
		f.favoriteCount = 0
		videoService := GetVideoServiceInstance()
		workCnt, err := videoService.GetVideoCnt(f.userId)
		if err != nil {
			// todo 日志记录一下为什么获取作品数失败
		}
		f.workCount = workCnt
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
	f.userInfo = &UserInfo{
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

func (u *UserService) GetUserFromRedisByUserId(userId int64) (dao.User, error) {
	UIdUClients := redis.Clients.UserId_UserR
	key := config.UserId_User_KEY_PREFIX + strconv.FormatInt(userId, 10)
	user, err := redis.GetKeysAndUpdateExpiration(UIdUClients, key)
	if err != nil {
		return dao.User{}, errors.New("get user from redis failed")
	}
	User, ok := user.(string)
	if !ok {
		return dao.User{}, errors.New("type assertion failed")
	}
	//反序列化
	var UserDao dao.User
	err = json.Unmarshal([]byte(User), &UserDao)
	if err != nil {
		return dao.User{}, errors.New("unmarshal failed")
	}
	return UserDao, nil
}
