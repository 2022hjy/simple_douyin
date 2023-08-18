package config

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

// JWT配置
const (
	// TokenExpireDuration token过期时间
	TokenExpireDuration = time.Hour * 2
	// JWTSECRET jwt加密串
	JWTSECRET = "hello"
)
