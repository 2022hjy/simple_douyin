package service

import (
	// 处理多部分表单数据，通常在处理 HTTP 请求中的文件上传时使用
	"mime/multipart"
	"simple_douyin/dao"
	"time"
)

// Video 返回给 Controller 层的 Video 结构体
type Video struct {
	dao.Video
	Author dao.User `json:"author"`
}

type VideoService interface {
	// Publish 将传入的视频流保存到 OSS 中，并在数据库中添加记录
	Publish(data *multipart.FileHeader, title string, userId int64) error

	// Feed 通过传入时间，当前用户的id，返回对应的返回视频流，以及视频流中最早的视频投稿时间
	Feed(latestTime time.Time, userId int64) ([]Video, time.Time, error)

	// PublishList 查询用户 userId 所发布的所有视频
	PublishList(userId int64) ([]Video, error)

	// GetVideoCnt 根据用户id查询用户的作品数
	GetVideoCnt(userId int64) (int64, error)

	// GetVideoListById 查询videoId列表的视频信息
	//GetVideoListById(videoIdList []int64, userId int64) ([]Video, error)
}
