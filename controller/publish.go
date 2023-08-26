package controller

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"simple_douyin/config"
	"simple_douyin/service"
)

type VideoListResponse struct {
	Response
	VideoList []VideoResponse `json:"video_list"`
}

// Publish 投稿视频
func Publish(c *gin.Context) {
	userId := c.GetInt64("token_user_id")
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	title := c.PostForm("title")
	log.Printf("视频 title: %v\n", title)
	videoService := service.GetVideoServiceInstance()
	// 从 token 中获取 userId
	//err = videoService.Publish(data, title, userId)
	err = videoService.Publish(c, data, title, userId)
	log.Printf("在 controller 视频 title: %v\n", title)
	//log.Println("视频 data：", data)
	log.Println("视频 userId：", userId)
	if err != nil {
		log.Println(err.Error())
		log.Println("上传文件失败")
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  fmt.Sprintf("《%s》视频上传成功", title),
	})
}

// PublishList 用户的视频发布列表，直接列出用户所有投稿过的视频
func PublishList(c *gin.Context) {
	// 获取到的userId，被访问的用户id
	reqUserId := c.Query("user_id")
	userId, _ := strconv.ParseInt(reqUserId, 10, 64)
	log.Println("获取到用户 Id：", userId)
	token := c.Query("token")
	log.Println("获取到用户 token：", token)
	videoService := service.GetVideoServiceInstance()
	plainVideos, err := videoService.PublishList(userId)
	videos := make([]VideoResponse, 0, config.VideoInitNum)
	videos, err = getRespVideos(plainVideos, userId)
	if err != nil {
		log.Println("getRespVideos:", err)
	}

	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 1, StatusMsg: "获取用户视频发布列表失败!"},
			VideoList: nil,
		})
		return
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "获取用户发布的视频列表成功！",
		},
		VideoList: videos,
	})
}
