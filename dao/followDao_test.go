package dao

import (
	"fmt"
	"log"
	"simple_douyin/middleware/database"
	"testing"
)

func TestFollowDao_InsertFollowRelation(t *testing.T) {
	followDao.InsertFollowRelation(6, 3)
}

func TestFollowDao_FindRelation(t *testing.T) {
	database.Init()
	follow, err := followDao.FindEverFollowing(3, 1)
	if err == nil {
		log.Default()
	}
	fmt.Println(follow.Followed)
	fmt.Print(follow)
}

func TestFollowDao_UpdateFollowRelation(t *testing.T) {
	// followDao.UpdateFollowRelation(2, 3, 1)

}

func TestFollowDao_GetFollowingsInfo(t *testing.T) {
	database.Init()
	followingsID, followingsCnt, err := followDao.GetFollowingsInfo(1)

	if err != nil {
		log.Default()
	}

	fmt.Println(followingsID)
	fmt.Println(followingsCnt)

}

func TestFollowDao_GetUserName(t *testing.T) {
	database.Init()
	name, err := followDao.GetUserName(1)
	if err != nil {
		log.Default()
	}
	fmt.Println(name)
}

func TestFollowDao_GetFriendsInfo(t *testing.T) {
	database.Init()
	friendId, friendCnt, _ := followDao.GetFriendsInfo(1)

	fmt.Println(friendId)
	fmt.Println(friendCnt)

}
