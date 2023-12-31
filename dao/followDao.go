package dao

import (
	"log"
	"simple_douyin/config"
	"simple_douyin/middleware/database"
	"sync"
	"time"
)

type Userfollow struct {
	User

	FollowCount   int64 `json:"follow_count"`
	FollowerCount int64 `json:"follower_count"`
	IsFollow      bool  `json:"is_follow"`
}
type FriendUser struct {
	Userfollow
	Avatar  string `json:"avatar"`
	Message string `json:"message"`
	MsgType int64  `json:"msg_type"`
}
type Follow struct {
	Id          int64  `gorm:"column:id"`           //关系ID
	UserId      int64  `gorm:"column:user_id"`      //用户ID
	FollowingId int64  `gorm:"column:following_id"` //关注的用户ID
	Followed    int8   `gorm:"column:is_followed"`  //是否已关注，1表示已关注，0表示未关注
	CreatedAt   string `gorm:"column:created_at"`   //记录创建时间
	UpdatedAt   string `gorm:"column:updated_at"`   //记录更新时间
}

func (Follow) TableName() string {
	return "relation"
}

type FollowDao struct {
}

var (
	followDao  *FollowDao
	followOnce sync.Once
)

// NewFollowDaoInstance 生成并返回followDao的单例对象。
func NewFollowDaoInstance() *FollowDao {
	followOnce.Do(
		func() {
			followDao = &FollowDao{}
		})
	return followDao
}

// GetNameByUserId 在user表中根据id查询用户姓名
func (*FollowDao) GetNameByUserId(userId int64) (string, error) {
	var name string

	err := database.Db.Table("user").Where("user_id = ?", userId).Pluck("username", &name).Error

	if nil != err {
		log.Println(err.Error())
		return "", err
	}

	return name, nil
}

// GetAvatarByUserId 在user表中根据id查询用户姓名
func (*FollowDao) GetAvatarByUserId(userId int64) (string, error) {
	var avatar string

	err := database.Db.Table("user").Where("user_id = ?", userId).Pluck("avatar", &avatar).Error

	if nil != err {
		log.Println(err.Error())
		return "", err
	}

	return avatar, nil
}

// CheckFollowRelation 给定当前用户和目标用户id，查看曾经是否有关注关系。
func (*FollowDao) CheckFollowRelation(userId int64, targetId int64) (*Follow, error) {
	// 用于存储查出来的关注关系。
	follow := Follow{}
	// 查询是否存在记录

	err := database.Db.
		Where("user_id = ?", userId).
		Where("following_id = ?", targetId).
		//Where("is_followed = ? or is_followed = ?", 0, 1).
		Take(&follow).Error
	// 当查询出现错误时，日志打印err msg，并return err.
	if nil != err {
		// 当没查到记录报错时，不当做错误处理。
		if "record not found" == err.Error() {
			return nil, nil
		}
		log.Println(err.Error())
		return nil, err
	}
	// 正常情况，返回取到的关系和空err.
	return &follow, nil
}

// InsertFollowRelation 给定用户和目标对象id，插入其关注关系。
func (*FollowDao) InsertFollowRelation(userId int64, targetId int64) (bool, error) {
	// 生成需要插入的关系结构体。
	follow := Follow{
		UserId:      userId,
		FollowingId: targetId,
		Followed:    1,
		CreatedAt:   time.Now().Format(config.GO_STARTER_TIME),
	}
	// 插入用户与目标用户的关注记录
	err := database.Db.Select("UserId", "FollowingId", "Followed", "CreatedAt").Create(&follow).Error
	// 插入失败，返回err.
	if nil != err {
		log.Println(err.Error())
		return false, err
	}
	// 插入成功
	return true, nil
}

// UpdateFollowRelation 给定用户和目标用户的id，更新他们的关系为取消关注或再次关注。
func (*FollowDao) UpdateFollowRelation(userId int64, targetId int64, followed int8) (bool, error) {
	// 更新用户与目标用户的关注记录（正在关注或者取消关注）
	err := database.Db.Model(Follow{}).
		Where("user_id = ?", userId).
		Where("following_id = ?", targetId).
		Update("is_followed", followed).Error
	// 更新失败，返回错误。
	if nil != err {
		// 更新失败，打印错误日志。
		log.Println(err.Error())
		return false, err
	}
	// 更新成功。
	return true, nil
}

// FindFollowRelation 给定当前用户和目标用户id，查询relation表是否存在关注关系
func (*FollowDao) FindFollowRelation(userId int64, targetId int64) (bool, error) {
	// follow变量用于后续存储数据库查出来的用户关系。
	follow := Follow{}
	//当查询出现错误时，日志打印err msg，并return err.
	if err := database.Db.
		Where("user_id = ?", userId).
		Where("following_id = ?", targetId).
		Where("is_followed = ?", 1).
		Take(&follow).Error; nil != err {
		// 当没查到数据时，gorm也会报错。
		if "record not found" == err.Error() {
			return false, nil
		}
		log.Println(err.Error())
		return false, err
	}
	//正常情况，返回取到的值和空err.
	return true, nil
}

// GetFollowingsInfo 返回当前用户正在关注的用户信息列表，包括当前用户正在关注的用户ID列表和正在关注的用户总数
func (*FollowDao) GetFollowingsInfo(userId int64) ([]int64, int64, error) {

	var followingCnt int64
	var followingId []int64

	// user_id -> following_id
	result := database.Db.Model(&Follow{}).Where("user_id = ?", userId).Where("is_followed = ?", 1).Pluck("following_id", &followingId)
	followingCnt = result.RowsAffected

	if nil != result.Error {
		log.Println(result.Error.Error())
		return nil, 0, result.Error
	}

	return followingId, followingCnt, nil

}

// GetFollowersInfo 返回当前用户的粉丝用户信息列表，包括当前用户的粉丝用户ID列表和粉丝总数
func (*FollowDao) GetFollowersInfo(userId int64) ([]int64, int64, error) {

	var followerCnt int64
	var followerId []int64

	// following_id -> user_id
	result := database.Db.Model(&Follow{}).Where("following_id = ?", userId).Where("is_followed = ?", 1).Pluck("user_id", &followerId)
	followerCnt = result.RowsAffected

	if nil != result.Error {
		log.Println(result.Error.Error())
		return nil, 0, result.Error
	}

	return followerId, followerCnt, nil
}

func (*FollowDao) GetFriendsInfo(userId int64) ([]int64, int64, error) {

	friendId, friendCnt, err := followDao.GetFollowingsInfo(userId)

	if nil != err {
		log.Println(err.Error())
		return nil, -1, err
	}

	for i := 0; int64(i) < friendCnt; i++ {
		// 判断每一个登陆用户的关注用户是否关注了登陆用户，没关注就从集合里面剔除
		if flag, err1 := followDao.FindFollowRelation(friendId[i], userId); !flag {
			if err1 != nil {
				return nil, -1, err1
			}
			friendId = append(friendId[:i], friendId[i+1:]...)
			friendCnt--
			i--
		}

	}
	return friendId, friendCnt, nil

}
