package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_douyin/service"
	"strconv"
)

const (
	ADD_COMMENT    = 1
	DELETE_COMMENT = 2
)

func respondWithError(c *gin.Context, statusCode int, errMsg string) {
	c.JSON(http.StatusOK, Response{StatusCode: statusCode, StatusMsg: errMsg})
}

var (
	commentService *service.CommentService
)

func init() {
	commentService = service.GetCommentServiceInstance()
}

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	userId := c.GetInt64("userId")
	videoId, ConvertVideoErr := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if ConvertVideoErr != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: -1, StatusMsg: "comment videoId json invalid"},
		})
		return
	}

	actionType, convertActionTypeErr := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if convertActionTypeErr != nil || actionType < 1 || actionType > 2 {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: -1, StatusMsg: "comment actionType json invalid"},
		})
		return
	}

	// 评论类型：1-评论 2-删除评论
	switch {
	case actionType == ADD_COMMENT:
		content := c.Query("comment_text")
		commentRes, err := commentService.Comment(userId, videoId, content)
		//评论操作时
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: -1, StatusMsg: "comment failed"},
			})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0,
				StatusMsg: "comment success"},
			Comment: commentRes,
		})
		return

	case actionType == DELETE_COMMENT:
		commentId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: -1,
					StatusMsg: "delete commentId invalid"},
			})
			return
		}
		err = commentService.DeleteComment(commentId)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: -1,
					StatusMsg: err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0,
				StatusMsg: "delete commentId success"},
		})
		return
	}
}

func CommentList(c *gin.Context) {
	userId := c.GetInt64("userId")
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		respondWithError(c, -1, "comment videoId json invalid: "+err.Error())
		return
	}
	commentList, err := commentService.GetCommentList(videoId, userId)
	if err != nil {
		respondWithError(c, -1, err.Error())
		return
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentList,
	})
	return
}
