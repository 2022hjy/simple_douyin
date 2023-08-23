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
	c.JSON(http.StatusOK, Response{StatusCode: int32(statusCode), StatusMsg: errMsg})
}

var (
	commentService *service.CommentService
)

func init() {
	commentService = service.GetCommentServiceInstance()
	userService = service.NewUserServiceInstance()
}

type CommentListResponse struct {
	Response
	CommentList []CommentResponse `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment CommentResponse `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	userId := c.GetInt64("token_user_id")
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
		var commentResponse = ConvertDBCommentToResponse(commentRes, userId)
		//评论操作时
		if err != nil {
			c.JSON(http.StatusInternalServerError, CommentActionResponse{
				Response: Response{StatusCode: -1, StatusMsg: "comment failed"},
			})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0,
				StatusMsg: "comment success"},
			Comment: commentResponse,
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
		err = commentService.DeleteComment(videoId, commentId)
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
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	token, _ := strconv.ParseInt(c.Query("token_user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CommentListResponse{
			Response: Response{StatusCode: -1, StatusMsg: "comment videoId json invalid"},
		})
	}
	commentList, err := commentService.GetCommentList(videoId)
	var commentResponseList []CommentResponse
	for i, comment := range commentList {
		/*
			userId := comment.UserId
			UIdU := redis.Clients.UserId_UserR
			key := config.UserId_User_KEY_PREFIX + strconv.FormatInt(userId, 10)
			userdao, err := redis.GetKeysAndUpdateExpiration(UIdU, key)
			var UserDao dao.UserDao
			if userdao == nil || err != nil {
				//从数据库中获得用户信息
				userDao := dao.UserDao{}
				UserDao, _ = userDao.GetUserById(userId)
			} else {
				UserDao, _ = userdao.(dao.UserDao)
			}
			//todo 获得评论者的信息，进行转化 User := dao.GetUserById(userId)
			//todo 获得FavoriteCount int64, FollowCount int64, FollowerCount int64, IsFollow bool, TotalFavorited string, WorkCount int64
			UserResponse := util.ConvertDBUserToResponse(UserDao)
		*/
		commentResponseList[i] = ConvertDBCommentToResponse(comment, token)
		commentResponseList = append(commentResponseList, commentResponseList[i])
	}

	if err != nil {
		respondWithError(c, -1, err.Error())
		return
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentResponseList,
	})
	return
}
