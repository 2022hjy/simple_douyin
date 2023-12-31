package controller

import (
	"errors"
	"log"
	"simple_douyin/dao"
	"simple_douyin/service"
	"time"
)

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
		return VideoResponse{999, UserResponse{}, "999", "999", 666, 666, false, "WRONG GET MESSAGE"}, errors.New(" userInfo 的获取出问题")
	}

	log.Printf("在 controller 的 convertVideo 里面，userInfo 的内容是：", userInfo)
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
		TotalFavorited:  userInfo.TotalFavorited, // 将 int64 转换为 string
		WorkCount:       userInfo.WorkCount,
		FavoriteCount:   userInfo.FavoriteCount,
	}

	log.Println("在 controller 的 convertVideo 里面，userResponse 的内容是：", userResponse)

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
		log.Println("在 controller 的 convertComment 里面，userInfo 的获取出问题")
		return CommentResponse{}
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
		TotalFavorited:  userInfo.TotalFavorited, // 将 int64 转换为 string
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
