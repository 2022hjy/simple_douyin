package model

type FriendUser struct {
	User
	Avatar  string `json:"avatar"`            //头像
	Message string `json:"message,omitempty"` //聊天记录
	MsgType int64  `json:"msg_type"`          //消息类型
}
