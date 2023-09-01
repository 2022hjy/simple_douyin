package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"net/http"
	"simple_douyin/config"
	"simple_douyin/dao"
	"strconv"
	"time"
)

func VideoRespondWithError(c *gin.Context, statusCode int32, errMsg string) {
	c.JSON(http.StatusOK, Response{StatusCode: statusCode, StatusMsg: errMsg})
}

type FeedResponse struct {
	Response
	VideoList []VideoResponse `json:"video_list"`
	NextTime  int64           `json:"next_time"`
}

// Feed 不限制登录状态，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
func Feed(c *gin.Context) {
	latestTime := c.Query("latest_time")
	log.Println("latestTime:", latestTime)
	log.Println("返回视频的最新投稿时间戳:", latestTime)
	var convTime time.Time
	if latestTime != "0" {
		t, _ := strconv.ParseInt(latestTime, 10, 64)
		if t > math.MaxInt32 {
			convTime = time.Now()
		} else {
			convTime = time.Unix(t, 0)
		}
	} else {
		convTime = time.Now()
	}
	userId := c.GetInt64("token_user_id")
	log.Printf("开始调用 feed 功能")
	convTime = time.Now()
	log.Printf("convTime:%v\n", convTime)
	plainVideos, nextTime, err := videoService.Feed(convTime)
	log.Println("调用 feed 功能结束")

	douyinVideos := make([]VideoResponse, 0, config.VideoInitNumPerRefresh)
	log.Println("在 videoService 的 Feed，转换前plainVideos:", plainVideos)
	douyinVideos, err = getRespVideos(plainVideos, userId)
	log.Println("转换后douyinVideos:", douyinVideos)
	if err != nil {
		MessageRespondWithError(c, -1, "Feed Error: "+err.Error())
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "Feed Success!"},
		VideoList: douyinVideos,
		NextTime:  nextTime.Unix(),
	})
}

// getRespVideos dao.video --> FeedResponse
func getRespVideos(plainVideos []dao.Video, userId int64) ([]VideoResponse, error) {
	var douyinVideos []VideoResponse
	for _, video := range plainVideos {
		response, err := ConvertDBVideoToResponse(video, userId)
		if err != nil {
			log.Println("getRespVideos出现问题:", err)
			return []VideoResponse{}, nil
		}
		douyinVideos = append(douyinVideos, response)
	}
	return douyinVideos, nil
}
