package util

import (
	"testing"

	"simple_douyin/model"
)

func TestMapToStruct(t *testing.T) {
	res := model.UserResponse{
		Response: model.Response{},
		User: model.User{
			UserId:          0,
			Username:        "123",
			Avatar:          "123",
			Signature:       "123",
			BackgroundImage: "123",
		},
		FollowCount:    0,
		FollowerCount:  0,
		IsFollow:       false,
		TotalFavorited: 0,
		WorkCount:      0,
		FavoriteCount:  0,
	}
	toMap := StructToMap(res)
	t.Log(toMap)
	toStruct := model.UserResponse{}
	t.Log(toStruct)
}
