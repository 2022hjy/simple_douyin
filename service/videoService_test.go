package service

import (
	"fmt"
	"log"
	"simple_douyin/dao"
	"simple_douyin/middleware/database"
	"simple_douyin/middleware/redis"
	"testing"
	"time"
)

// 获取某位用户的视频信息list
func TestVideoServiceImpl_PublishList(t *testing.T) {
	redis.InitRedis()
	videoList, err := videoServiceImpl.PublishList(226)
	if err != nil {
		log.Default()
	}
	fmt.Println(videoList)
}

// 返回倒序视频流 GET
func TestVideoServiceImpl_Feed(t *testing.T) {
	database.Init()
	redis.InitRedis()

	videoList, nextTime, err := videoServiceImpl.Feed(time.Now())
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
	videoCnt, err := videoServiceImpl.GetVideoCnt(226)
	if err != nil {
		log.Default()
	}
	fmt.Println("videoCnt:", videoCnt)
}

//// 将视频上传到oss 并保存到数据库
//func TestPublish(t *testing.T) {
//	database.Init()
//
//	filePath := "videotest.mp4"
//	videocontent, err := ioutil.ReadFile(filePath)
//	if err != nil {
//		log.Fatal("Error reading file:", err)
//		return
//	}
//
//	fileHeader := &multipart.FileHeader{
//		Filename: "videotest.mp4",
//		Size:     int64(len(videocontent)),
//	}
//	log.Println(fileHeader.Filename)
//
//	title := "test title"
//	userId := int64(1)
//	videoService := &VideoServiceImpl{}
//
//	err = videoService.Publish(fileHeader, title, userId)
//	log.Println(err)
//}

func TestUploadVideo(t *testing.T) {
	database.Init()

	title := "test title"
	userId := int64(1)

	err := dao.UploadVideo("videoname", userId, title)
	log.Println(err)
}
