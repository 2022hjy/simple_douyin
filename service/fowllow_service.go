package service

import (
	"simple_douyin/dao"
)

type FollowService interface {

	// AddFollowAction 当前用户关注目标用户
	AddFollowAction(userId int64, targetId int64) (bool, error)
	// CancelFollowAction 当前用户取消对目标用户的关注
	CancelFollowAction(userId int64, targetId int64) (bool, error)
	// GetFollowList 获取当前用户的关注列表
	GetFollowList(userId int64) ([]dao.Userfollow, error)
	// GetFollowerList 获取当前用户的粉丝列表
	GetFollowerList(userId int64) ([]dao.Userfollow, error)
	// GetFriendList 获取好友
	GetFriendList(userId int64) ([]dao.FriendUser, error)
	// GetFollowingCnt 根据用户id查询关注数
	GetFollowingCnt(userId int64) (int64, error)
	// GetFollowerCnt 根据用户id查询粉丝数
	GetFollowerCnt(userId int64) (int64, error)
	// CheckIsFollowing 判断当前登录用户是否关注了目标用户
	CheckIsFollowing(userId int64, targetId int64) (bool, error)
}
