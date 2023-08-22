package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"net/http"
	"simple_douyin/service"
	"strconv"
	"time"
)

func VideoRespondWithError(c *gin.Context, statusCode int32, errMsg string) {
	c.JSON(http.StatusOK, Response{StatusCode: statusCode, StatusMsg: errMsg})
}

var (
	videoService *service.VideoServiceImpl
)

func init() {
	videoService = service.GetVideoServiceInstance()
}

type FeedResponse struct {
	Response
	VideoList []VideoResponse `json:"video_list,omitempty"`
	NextTime  int64           `json:"next_time,omitempty"`
}

// Feed 不限制登录状态，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
func Feed(c *gin.Context) {
	latestTime := c.Query("latest_time")
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
	userId := c.GetInt64("userId")
	douyinVideos, nextTime, err := videoService.Feed(convTime, userId)
	if err != nil {
		MessageRespondWithError(c, -1, "Feed Error: "+err.Error())
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "Feed Success!"},
		VideoList: douyinVideos,
		NextTime:  nextTime.Unix(),
	})
}
