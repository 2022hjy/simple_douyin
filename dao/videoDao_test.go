package dao

import (
	"log"
	"simple_douyin/config"
	"simple_douyin/middleware/database"
	"testing"
	"time"
)

// 测试保存视频到数据库  GET
func TestSaveVideo(t *testing.T) {
	database.Init()

	video := Video{
		Id:            2,
		UserInfoId:    226,
		Title:         "测试视频1",
		PlayUrl:       "https://www.baidu.com",
		CoverUrl:      "https://www.baidu.com",
		IsFavorite:    1,
		FavoriteCount: 23,
		CommentCount:  0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	err := SaveVideo(video)
	if err != nil {
		log.Println(err)
	}
	log.Println("保存视频成功！")
}

// 根据视频 Id 获取视频信息  GET
func TestGetVideoByVideoId(t *testing.T) {
	database.Init()

	// 测试获取单个视频
	video, err := GetVideoByVideoId(1)
	if err != nil {
		t.Errorf("获取视频失败：%v", err)
	} else {
		log.Println("单个视频：", video)
	}

	// 测试获取多个视频
	videoIds := []int64{1, 2, 3, 4}
	for _, videoId := range videoIds {
		video, err := GetVideoByVideoId(videoId)
		if err != nil {
			t.Errorf("获取视频失败：%v", err)
		} else {
			log.Println("视频 ID", videoId, "：", video)
		}
	}
}

// 测试根据用户 Id 获取该用户已发布的所有视频  GET
func TestGetVideosByUserId(t *testing.T) {
	database.Init()

	res, err := GetVideosByUserId(1)
	if err == nil {
		for _, re := range res {
			log.Println(re)
		}
	}
}

// 根据视频 Id 列表获取视频信息列表  GET
func TestGetVideoListById(t *testing.T) {
	database.Init()
	videoList, err := GetVideoListById([]int64{1, 2, 3, 4})
	if err == nil {
		log.Println(len(videoList))
	}
	for _, video := range videoList {
		log.Println(video)
	}
}

// 测试按投稿时间倒序的视频列表 GET
func TestGetVideosByLatestTime(t *testing.T) {
	database.Init()
	// 时区修正
	mockTime, _ := time.ParseInLocation(config.GO_STARTER_TIME, "2024-01-29 21:20:04", time.Local)
	log.Println(mockTime)
	res, err := GetVideosByLatestTime(mockTime)
	if err == nil {
		for _, re := range res {
			log.Println(re)
		}
	}
}

// 根据userId获取作品数量  GET
func TestGetVideoCnt(t *testing.T) {
	database.Init()
	count, err := GetVideoCnt(1)
	if err == nil {
		log.Println(count)
	}
}

// 测试上传视频 GET
func TestUploadVideo(t *testing.T) {
	database.Init()
	err := UploadVideo("VID_2023_1_29", 1, "测试视频1")
	if err != nil {
		return
	}
}
