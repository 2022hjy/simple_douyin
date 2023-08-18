package model

type User struct {
	UserId          int64  `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Username        string `json:"name" gorm:"unique;not null"`
	Password        string `json:"-"` // 不返回给前端
	Avatar          string `json:"avatar,omitempty" gorm:"default:''"`
	BackgroundImage string `json:"background_image,omitempty" gorm:"default:''"`
	Signature       string `json:"signature,omitempty" gorm:"default:''"`
}

func (User) TableName() string {
	return "user"
}

type LoginRequest struct {
	UserName string
	Password string
}

type RegisterRequest struct {
	UserName string
	Password string
}

type LoginResponse struct {
	Response

	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type RegisterResponse struct {
	Response

	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response

	User

	FollowCount   int64 `json:"follow_count,omitempty"`
	FollowerCount int64 `json:"follower_count,omitempty"`
	IsFollow      bool  `json:"is_follow,omitempty"`

	TotalFavorited int64 `json:"total_favorited,omitempty"`
	WorkCount      int64 `json:"work_count,omitempty"`
	FavoriteCount  int64 `json:"favorite_count,omitempty"`
}
