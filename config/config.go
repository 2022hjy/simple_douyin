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

	CustomDomain = "CUSTOM_DOMAIN"
	OssVideoDir  = "OSS_VIDEO_DIR"
	// PlayUrlPrefix 播放视频的URL前缀
	PlayUrlPrefix = CustomDomain + OssVideoDir
	// CoverUrlSuffix 频封面图的URL后缀，通过该后缀可以对视频进行截图并获取封面图
	CoverUrlSuffix     = "?x-oss-process=video/snapshot,t_2000,m_fast"
	OssAccessKeyId     = "OSS_ACCESS_KEY_ID"
	OssAccessKeySecret = "OSS_ACCESS_KEY_SECRET"
	OssBucketName      = "OSS_BUCKET_NAME"
	OssEndpoint        = "OSS_ENDPOINT"
)

// 消息模块相关配置
const (
	// MessageInitNum 每次刷新时展示的消息数量
	MessageInitNum = 10
)

// JWT配置
const (
	// TokenExpireDuration token过期时间
	TokenExpireDuration = time.Hour * 2
	// JWTSECRET jwt加密串
	JWTSECRET = "hello"
)
