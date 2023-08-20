package config

import "time"

const GO_STARTER_TIME = "2006-01-02 15:04:05"

// RedisAddr redis地址
const (
	RedisAddr  = "localhost:6379"
	RedisPwd   = ""
	THREASHOLD = 10
)

// Mq 消息队列
const (
	MqUrl = "amqp://guest:guest@localhost:5672/"
)

// 视频模块相关配置
const (
	// VIDEO_INIT_NUM 每位作者初始展示的视频数量
	VIDEO_INIT_NUM = 10
	// VIDEO_NUM_PER_REFRESH 表示每次刷新时展示的视频数量，值为 6。
	VIDEO_NUM_PER_REFRESH = 6

	CUSTOM_DOMAIN = "CUSTOM_DOMAIN"
	OSS_VIDEO_DIR = "OSS_VIDEO_DIR"
	// PLAY_URL_PREFIX 播放视频的URL前缀
	PLAY_URL_PREFIX = CUSTOM_DOMAIN + OSS_VIDEO_DIR
	// COVER_URL_SUFFIX 频封面图的URL后缀，通过该后缀可以对视频进行截图并获取封面图
	COVER_URL_SUFFIX      = "?x-oss-process=video/snapshot,t_2000,m_fast"
	OSS_ACCESS_KEY_ID     = "OSS_ACCESS_KEY_ID"
	OSS_ACCESS_KEY_SECRET = "OSS_ACCESS_KEY_SECRET"
	OSS_BUCKET_NAME       = "OSS_BUCKET_NAME"
	OSS_ENDPOINT          = "OSS_ENDPOINT"
)

// JWT配置
const (
	// TokenExpireDuration token过期时间
	TokenExpireDuration = time.Hour * 2
	// JWTSECRET jwt加密串
	JWTSECRET = "hello"
)
