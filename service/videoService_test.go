package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
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

// 将视频上传到oss 并保存到数据库
func TestPublish(t *testing.T) {
	filename := "videotest.mp4"
	content, err := ioutil.ReadFile(filename)
	videocontent := []byte(content)

	fileHeader := &multipart.FileHeader{
		Filename: "videotest.mp4",
		Size:     int64(len(videocontent)),
	}
	log.Println(fileHeader.Filename)

	title := "test title"
	userId := int64(1)
	videoService := &VideoServiceImpl{}

	err = videoService.Publish(fileHeader, title, userId)
	log.Println(err)
}

//func TestUploadVideo(t *testing.T) {
//	database.Init()
//	err := UploadVideo("VID_2023_1_29", 1, "测试视频1")
//	if err != nil {
//		return
//	}
//}
