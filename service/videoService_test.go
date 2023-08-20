package service

import (
	"fmt"
	"log"
	"simple_douyin/middleware/database"
	"simple_douyin/middleware/redis"
	"testing"
	"time"
)

// 获取某位用户的视频信息list
func TestVideoServiceImpl_PublishList(t *testing.T) {
	redis.InitRedis()
	videoList, err := videoServiceImp.PublishList(1)
	if err != nil {
		log.Default()
	}
	fmt.Println(videoList)
}

// 返回倒序视频流 GET
func TestVideoServiceImpl_Feed(t *testing.T) {
	database.Init()
	redis.InitRedis()

	videoList, nextTime, err := videoServiceImp.Feed(time.Now(), 1)
	if err != nil {
		log.Default()
	}
	fmt.Println(nextTime)
	fmt.Println(videoList)
}

// GetVideoCnt  GET
func TestVideoServiceImpl_GetVideoCnt(t *testing.T) {
	database.Init()
	redis.InitRedis()
	videoCnt, err := videoServiceImp.GetVideoCnt(1)
	if err != nil {
		log.Default()
	}
	fmt.Println("videoCnt:", videoCnt)
}

// TODO 将视频上传到oss并保存到数据库中
//func TestVideoServiceImpl_SaveVideoInfo(t *testing.T) {
//	database.Init()
//	redis.InitRedis()
//	videoInfo := Video{
//		Video: dao.Video{
//			Id:         1,
//			UserInfoId: 2,
//			Title:      "Sample Video",
//			PlayUrl:    "https://example.com/sample_video",
//			CoverUrl:   "https://example.com/sample_cover",
//			CreatedAt:  time.Now(),
//			UpdatedAt:  time.Now(),
//		},
//		Author: User{
//			Id:   2,
//			Name: "John Doe",
//		},
//		FavoriteCount: 12,
//		CommentCount:  5,
//		IsFavorite:    1,
//	}
//	err := videoServiceImp.Publish(videoInfo)
//	if err != nil {
//		log.Default()
//	}
//	fmt.Println("videoInfo:", videoInfo)
//}
