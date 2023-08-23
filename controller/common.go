package controller

import (
	"simple_douyin/dao"
	"simple_douyin/service"
	"strconv"
	"time"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

func Success() Response {
	return Response{
		StatusCode: 0,
		StatusMsg:  "success",
	}
}

// VideoResponse data 内部返回给前端的结构体
type VideoResponse struct {
	Id            int64        `json:"id,omitempty"`
	Author        UserResponse `json:"author"`
	PlayUrl       string       `json:"play_url,omitempty"`
	CoverUrl      string       `json:"cover_url,omitempty"`
	FavoriteCount int64        `json:"favorite_count,omitempty"`
	CommentCount  int64        `json:"comment_count,omitempty"`
	IsFavorite    bool         `json:"is_favorite,omitempty"`
	Title         string       `json:"title,omitempty"`
}

// UserResponse  data 返回给前端的结构体
type UserResponse struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	Signature       string `json:"signature"`        // 个人简介
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow,omitempty"`
	FavoriteCount   int64  `json:"favorite_count"`  // 喜欢数
	TotalFavorited  string `json:"total_favorited"` // 获赞数量
	WorkCount       int64  `json:"work_count"`      // 作品数
}

type FriendUser struct {
	User
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
	Avatar        string `json:"avatar"`            //头像
	Message       string `json:"message,omitempty"` //聊天记录
	MsgType       int64  `json:"msg_type"`          //消息类型
}

type CommentResponse struct {
	Id         int64        `json:"id,omitempty"`
	User       UserResponse `json:"user"`
	Content    string       `json:"content,omitempty"`
	CreateDate string       `json:"create_date,omitempty"`
}

type User struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	Signature       string `json:"signature"`        // 个人简介
	Password        string `json:"-"`                // 不返回给前端
}

type MessageSave struct {
	Id         int64  `json:"id,omitempty"`
	FromUserId int64  `json:"from_user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

// ConvertDBUserToResponse 将数据库中的用户信息转换为响应的用户信息
// Finish this function
//func ConvertDBUserToResponse(dbUser dao.User, FavoriteCount int64, FollowCount int64, FollowerCount int64, IsFollow bool, TotalFavorited string, WorkCount int64) controller.UserResponse {
//	return controller.UserResponse{
//		Avatar:          dbUser.Avatar,
//		BackgroundImage: dbUser.BackgroundImage,
//		Id:              dbUser.UserId,
//		Name:            dbUser.Username,
//		Signature:       dbUser.Signature,
//		FavoriteCount:   FavoriteCount,
//		FollowCount:     FollowCount,
//		FollowerCount:   FollowerCount,
//		IsFollow:        IsFollow,
//		TotalFavorited:  TotalFavorited,
//		WorkCount:       WorkCount,
//	}
//}

// ConvertDBVideoToResponse 将数据库中的视频信息转换为响应的视频信息
// tips: 这里的User是controller.UserResponse类型，不是dao.UserDAO类型！！！
func ConvertDBVideoToResponse(dbVideo dao.Video, tokenId int64) (VideoResponse, error) {
	userService := service.NewUserServiceInstance()

	// 使用 QueryUserInfo 获取视频作者的信息
	userInfo, err := userService.QueryUserInfo(dbVideo.UserInfoId, tokenId) // 假设 tokenUserId 为0
	if err != nil {
		return VideoResponse{}, err
	}

	// 将 UserInfo 转换为 UserResponse
	userResponse := UserResponse{
		Id:              userInfo.User.UserId,
		Name:            userInfo.User.Username,
		FollowCount:     userInfo.FollowCount,
		FollowerCount:   userInfo.FollowerCount,
		IsFollow:        userInfo.IsFollow,
		Avatar:          userInfo.User.Avatar,
		BackgroundImage: userInfo.User.BackgroundImage,
		Signature:       userInfo.User.Signature,
		TotalFavorited:  strconv.FormatInt(userInfo.TotalFavorited, 10), // 将 int64 转换为 string
		WorkCount:       userInfo.WorkCount,
		FavoriteCount:   userInfo.FavoriteCount,
	}

	var isFavorite bool
	if dbVideo.IsFavorite == 1 {
		isFavorite = true
	} else {
		isFavorite = false
	}

	return VideoResponse{
		Id:            dbVideo.Id,
		Author:        userResponse,
		PlayUrl:       dbVideo.PlayUrl,
		CoverUrl:      dbVideo.CoverUrl,
		FavoriteCount: dbVideo.FavoriteCount,
		CommentCount:  dbVideo.CommentCount,
		IsFavorite:    isFavorite,
		Title:         dbVideo.Title,
	}, nil
}

//func ConvertDBCommentToResponse(comment dao.CommentDao, UserResponse controller.UserResponse) controller.CommentResponse {
//	// Convert time.Time to string
//	createDate := comment.CreatedAt.Format(time.RFC3339)
//	return controller.CommentResponse{
//		Id:         comment.Id,
//		User:       UserResponse,
//		Content:    comment.Content,
//		CreateDate: createDate,
//	}
//}

func ConvertDBCommentToResponse(comment dao.CommentDao, tokenId int64) CommentResponse {
	userService := service.NewUserServiceInstance()

	// 使用 QueryUserInfo 获取用户信息
	userInfo, err := userService.QueryUserInfo(comment.UserId, tokenId) // 假设 tokenUserId 为0
	if err != nil {
		// 处理错误或记录日志
	}

	// 将 UserInfo 转换为 UserResponse
	userResponse := UserResponse{
		Id:              userInfo.User.UserId,
		Name:            userInfo.User.Username,
		FollowCount:     userInfo.FollowCount,
		FollowerCount:   userInfo.FollowerCount,
		IsFollow:        userInfo.IsFollow,
		Avatar:          userInfo.User.Avatar,
		BackgroundImage: userInfo.User.BackgroundImage,
		Signature:       userInfo.User.Signature,
		TotalFavorited:  strconv.FormatInt(userInfo.TotalFavorited, 10),
		WorkCount:       userInfo.WorkCount,
		FavoriteCount:   userInfo.FavoriteCount,
	}

	// Convert time.Time to string
	createDate := comment.CreatedAt.Format(time.RFC3339)
	return CommentResponse{
		Id:         comment.Id,
		User:       userResponse,
		Content:    comment.Content,
		CreateDate: createDate,
	}
}
