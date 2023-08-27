package controller

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func Error(code int32, msg string) Response {
	return Response{
		StatusCode: code,
		StatusMsg:  msg,
	}
}

func Success() Response {
	return Response{
		StatusCode: 0,
		StatusMsg:  "success",
	}
}

// VideoResponse data 内部返回给前端的结构体
type VideoResponse struct {
	Id            int64        `json:"id"`
	Author        UserResponse `json:"author"`
	PlayUrl       string       `json:"play_url"`
	CoverUrl      string       `json:"cover_url"`
	FavoriteCount int64        `json:"favorite_count"`
	CommentCount  int64        `json:"comment_count"`
	IsFavorite    bool         `json:"is_favorite"`
	Title         string       `json:"title"`
}

// UserResponse  data 返回给前端的结构体
type UserResponse struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	Signature       string `json:"signature"`        // 个人简介
	FollowCount     int64  `json:"follow_count"`
	FollowerCount   int64  `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	FavoriteCount   int64  `json:"favorite_count"`  // 喜欢数
	TotalFavorited  int64  `json:"total_favorited"` // 获赞数量
	WorkCount       int64  `json:"work_count"`      // 作品数
}

type FriendUser struct {
	User
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
	Avatar        string `json:"avatar"`   //头像
	Message       string `json:"message"`  //聊天记录
	MsgType       int64  `json:"msg_type"` //消息类型
}

type CommentResponse struct {
	Id         int64        `json:"id"`
	User       UserResponse `json:"user"`
	Content    string       `json:"content"`
	CreateDate string       `json:"create_date"`
}

type User struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	Signature       string `json:"signature"`        // 个人简介
	Password        string `json:"-"`                // 不返回给前端
}

type MessageSave struct {
	Id         int64  `json:"id"`
	FromUserId int64  `json:"from_user_id"`
	ToUserId   int64  `json:"to_user_id"`
	Content    string `json:"content"`
	CreateTime string `json:"create_time"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id"`
	ToUserId   int64  `json:"to_user_id"`
	MsgContent string `json:"msg_content"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id"`
	MsgContent string `json:"msg_content"`
}
