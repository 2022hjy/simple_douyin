package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"simple_douyin/service"
	"strconv"
)

// VideoResponse 返回给 Controller 层的 VideoResponse 结构体
type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type FavoriteActionResponse struct {
	Response
}

type GetFavouriteListResponse struct {
	Response
	VideoList []service.Video `json:"video_list"`
}

func init() {
	commentService = service.GetCommentServiceInstance()
}

// FavoriteAction 处理点赞操作
func FavoriteAction(c *gin.Context, service service.FavoriteService) {
	videoID, err := strconv.ParseInt(c.Query("video_id"), 10, 64)

	if err != nil {
		log.Printf("解析 video_id 失败：%v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{StatusCode: -1, StatusMsg: "无效的 video_id"})
		return
	}

	if err := service.FavoriteAction(c.GetInt64("userId"), videoID); err != nil {
		log.Printf("用户 %d 对视频 %d 的点赞操作失败：%v", c.GetInt64("userId"), videoID, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{StatusCode: -1, StatusMsg: "收藏操作失败"})
		return
	}

	log.Printf("用户 %d 成功点赞了视频 %d", c.GetInt64("userId"), videoID)
	c.JSON(http.StatusOK, FavoriteActionResponse{Response{StatusCode: 0, StatusMsg: "收藏操作成功"}})
}

// FavoriteList 获取收藏列表
func FavoriteList(c *gin.Context, service service.FavoriteService) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		log.Printf("解析 user_id 失败：%v", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{StatusCode: -1, StatusMsg: "无效的 user_id"})
		return
	}

	videoList, err := service.GetFavoriteList(userId)
	if err != nil {
		log.Printf("获取用户 %d 的收藏列表失败：%v", userId, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{StatusCode: -1, StatusMsg: "获取收藏列表失败"})
		return
	}

	log.Printf("成功获取了用户 %d 的收藏列表", userId)
	c.JSON(http.StatusOK, GetFavouriteListResponse{Response: Response{StatusCode: 0, StatusMsg: "获取收藏列表成功"}, VideoList: videoList})
}
