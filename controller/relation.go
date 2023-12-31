package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simple_douyin/dao"
	"simple_douyin/service"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []dao.Userfollow `json:"user_list"`
}
type FriendUserListResponse struct {
	Response
	FriendUserList []dao.FriendUser `json:"user_list"`
}

// 关系操作
// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	userId := c.GetInt64("token_user_id")
	log.Println("userid参数为：", userId)
	//userId, err1 := strconv.ParseInt(c.Query("userId"), 10, 64)
	ToUserId := c.Query("to_user_id")
	log.Println("to_user参数分别为：", ToUserId)
	ActionType := c.Query("action_type")
	log.Println("action参数分别为：", ActionType)

	//批量转化
	//userId, err2 := strconv.ParseInt(UserId, 10, 64)
	toUserId, err3 := strconv.ParseInt(ToUserId, 10, 64)
	log.Println("参数分别为：", userId, toUserId)
	actionType, _ := strconv.ParseInt(ActionType, 10, 64)
	log.Println("参数分别为：", userId, toUserId, actionType)
	// 传入参数格式有问题。

	if nil != err3 || actionType < 1 || actionType > 2 {
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
			_, err := fsi.AddFollowAction(userId, toUserId)
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

func Followlist(c *gin.Context) {
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
	followings, err1 := fsi.GetFollowList(userId)
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
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)

	fmt.Println(userId)

	if err != nil {
		fmt.Printf("fail")
		c.JSON(http.StatusOK, UserListResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "请求参数格式错误",
			},
			nil,
		})
		return
	}

	fsi := service.NewFSIInstance()
	followers, err1 := fsi.GetFollowerList(userId)
	if err1 != nil {
		fmt.Printf("fail")
		c.JSON(http.StatusOK, UserListResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "获取粉丝列表失败",
			},
			nil,
		})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response{
			StatusCode: 0,
			StatusMsg:  "获取粉丝列表成功",
		},
		followers,
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)

	fmt.Println(userId)

	if err != nil {
		fmt.Printf("fail")
		c.JSON(http.StatusOK, FriendUserListResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "请求参数格式错误",
			},
			nil,
		})
		return
	}

	fsi := service.NewFSIInstance()
	followers, err1 := fsi.GetFriendList(userId)
	if err1 != nil {
		fmt.Printf("fail")
		c.JSON(http.StatusOK, FriendUserListResponse{
			Response{
				StatusCode: -1,
				StatusMsg:  "获取好友列表失败",
			},
			nil,
		})
		return
	}

	c.JSON(http.StatusOK, FriendUserListResponse{
		Response{
			StatusCode: 0,
			StatusMsg:  "获取好友列表成功",
		},
		followers,
	})
}
