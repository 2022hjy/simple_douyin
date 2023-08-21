package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_douyin/service"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []UserResponse `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
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
			StatusCode: 200,
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
		UserList: []User{DemoUser},
	})
}

// FriendList all users have same friend list
func FriendList(c *gin.Context) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}
