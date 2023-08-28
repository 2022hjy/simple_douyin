package controller

import "simple_douyin/service"

var (
	commentService  = service.GetCommentServiceInstance()
	favoriteService = service.GetFavoriteServiceInstance()
	videoService    = service.GetVideoServiceInstance()
	messageService  = service.GetMessageServiceInstance()
	userService     = service.NewUserServiceInstance()
)
