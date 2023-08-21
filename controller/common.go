package controller

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// VideoResponse data 内部返回给前端的结构体
type VideoResponse struct {
	Id            int64        `json:"id,omitempty"`
	Author        UserResponse `json:"author"`
	PlayUrl       string       `json:"play_url,omitempty"`
	CoverUrl      string       `json:"cover_url,omitempty"`
	FavoriteCount int64        `json:"favorite_count,omitempty"`
	CommentCount  int64        `json:"comment_count,omitempty"`
	IsFavorite    bool         `json:"is_favorite,omitempty"`
	Title         string       `json:"title,omitempty"`
}

// UserResponse  data 返回给前端的结构体
type UserResponse struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	Signature       string `json:"signature"`        // 个人简介
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow,omitempty"`
	FavoriteCount   int64  `json:"favorite_count"`  // 喜欢数
	TotalFavorited  string `json:"total_favorited"` // 获赞数量
	WorkCount       int64  `json:"work_count"`      // 作品数
}

type FriendUser struct {
	User
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
	Avatar        string `json:"avatar"`            //头像
	Message       string `json:"message,omitempty"` //聊天记录
	MsgType       int64  `json:"msg_type"`          //消息类型
}

type CommentResponse struct {
	Id         int64        `json:"id,omitempty"`
	User       UserResponse `json:"user"`
	Content    string       `json:"content,omitempty"`
	CreateDate string       `json:"create_date,omitempty"`
}

type User struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	Signature       string `json:"signature"`        // 个人简介
	Password        string `json:"-"`                // 不返回给前端
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}
