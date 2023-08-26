package config

import "time"

const GO_STARTER_TIME = "2006-01-02 15:04:05"

// RedisAddr redis地址
const (
	RedisAddr  = "localhost:6379"
	RedisPwd   = ""
	THREASHOLD = 10
	ExpireTime = 24 * time.Hour
)

// Mq 消息队列
const (
	MqUrl = "amqp://guest:guest@localhost:5672/"
)

// 视频模块相关配置
const (
	// VideoInitNum 每位作者初始展示的视频数量
	VideoInitNum = 10
	// VideoInitNumPerRefresh 表示每次刷新时展示的视频数量，值为 6。
	VideoInitNumPerRefresh = 6

	OssVideoDir = "OSS_VIDEO/"
	// PlayUrlPrefix 播放视频的URL前缀
	PlayUrlPrefix = OssVideoDir

	// CoverUrlSuffix 频封面图的URL后缀，通过该后缀可以对视频进行截图并获取封面图
	CoverUrlSuffix = "?x-oss-process=video/snapshot,t_2000,m_fast"

	OssAccessKeyId     = "LTAI5tNUGBD7bCfAoUddQVVw"
	OssAccessKeySecret = "1DM9nfo7MlGuN00sWFU0zrtkUELDtG"
	OssBucketName      = "simple-tiktok"
	OssEndpoint        = "https://oss-cn-beijing.aliyuncs.com"
	//	https://simple-tiktok.oss-cn-beijing.aliyuncs.com
)

// 消息模块相关配置
const (
	// MessageInitNum 每次刷新时展示的消息数量
	MessageInitNum = 10
)

// JWT配置
const (
	// TokenExpireDuration token过期时间
	// 测试阶段为了方便，设置为1天，正式上线后应该设置为2小时
	TokenExpireDuration = time.Hour * 24
	// JWTSECRET jwt加密串
	JWTSECRET = "hello"
)
