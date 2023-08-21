package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simple_douyin/service"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []UserResponse `json:"user_list"`
}
type FriendUserListResponse struct {
	Response
	FriendUserList []FriendUser `json:"user_list"`
}

// 关系操作
// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	userId := c.GetInt64("userId")
	//userId, err1 := strconv.ParseInt(c.Query("userId"), 10, 64)
	toUserId, err2 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	actionType, err3 := strconv.ParseInt(c.Query("action_type"), 10, 64)
	//fmt.Println(userId)
	//fmt.Println(toUserId)
	//fmt.Println(actionType)
	// 传入参数格式有问题。
	if nil != err2 || nil != err3 || actionType < 1 || actionType > 2 {
		fmt.Printf("fail")
		c.JSON(http.StatusOK, Response{
			StatusCode: 400,
			StatusMsg:  "请求参数格式错误",
		})
		return
	}
	// 正常处理
	fsi := service.NewFSIInstance()
	switch {
	// 关注
	case 1 == actionType:
		go func() {
			_, err := fsi.FollowAction(userId, toUserId)
			if err != nil {
				log.Println(err)
			}
		}()
	// 取关
	case 2 == actionType:
		go func() {
			_, err := fsi.CancelFollowAction(userId, toUserId)
			if err != nil {
				log.Println(err)
			}
		}()
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "操作成功"})
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)

	fmt.Println(userId)
	if err != nil {
		fmt.Printf("fail")
		c.JSON(http.StatusOK, UserListResponse{
			Response{
				StatusCode: 400,
				StatusMsg:  "请求参数格式错误",
			},
			nil,
		})
		return
	}

	fsi := service.NewFSIInstance()
	followings, err1 := fsi.GetFollowings(userId)
	if err1 != nil {
		fmt.Printf("fail")
		c.JSON(http.StatusOK, UserListResponse{
			Response{
				StatusCode: 500,
				StatusMsg:  "获取关注列表失败",
			},
			nil,
		})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response{
			StatusCode: 0,
			StatusMsg:  "获取关注列表成功",
		},
		followings,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []UserResponse{DemoUser},
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []UserResponse{DemoUser},
	})
}
