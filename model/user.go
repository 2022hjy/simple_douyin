package model

type UserResponse struct {
	Response
	UserInfo `json:"user"`
}
