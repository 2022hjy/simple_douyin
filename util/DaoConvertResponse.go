package util

import (
	"simple_douyin/controller"
	"simple_douyin/dao"
	"time"
)

// ConvertDBUserToResponse 将数据库中的用户信息转换为响应的用户信息
// Finish this function
func ConvertDBUserToResponse(dbUser dao.UserDAO, FavoriteCount int64, FollowCount int64, FollowerCount int64, IsFollow bool, TotalFavorited string, WorkCount int64) controller.UserResponse {
	return controller.UserResponse{
		Avatar:          dbUser.Avatar,
		BackgroundImage: dbUser.BackgroundImage,
		ID:              dbUser.UserId,
		Name:            dbUser.Username,
		Signature:       dbUser.Signature,
		FavoriteCount:   FavoriteCount,
		FollowCount:     FollowCount,
		FollowerCount:   FollowerCount,
		IsFollow:        IsFollow,
		TotalFavorited:  TotalFavorited,
		WorkCount:       WorkCount,
	}
}

// ConvertDBVideoToResponse 将数据库中的视频信息转换为响应的视频信息
// tips: 这里的User是controller.UserResponse类型，不是dao.UserDAO类型！！！
func ConvertDBVideoToResponse(dbVideo dao.Video, User controller.UserResponse) controller.VideoResponse {
	var isFavorite bool
	if dbVideo.IsFavorite == 1 {
		isFavorite = true
	} else {
		isFavorite = false
	}
	return controller.VideoResponse{
		Id:            dbVideo.Id,
		Author:        User,
		PlayUrl:       dbVideo.PlayUrl,
		CoverUrl:      dbVideo.CoverUrl,
		FavoriteCount: int64(dbVideo.FavoriteCount),
		CommentCount:  int64(dbVideo.CommentCount),
		IsFavorite:    isFavorite,
		Title:         dbVideo.Title,
	}
}

func ConvertDBCommentToResponse(comment dao.CommentDao, UserResponse controller.UserResponse) controller.CommentResponse {
	// Convert time.Time to string
	createDate := comment.CreatedAt.Format(time.RFC3339)
	return controller.CommentResponse{
		Id:         comment.Id,
		User:       UserResponse,
		Content:    comment.Content,
		CreateDate: createDate,
	}
}
