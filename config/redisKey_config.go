package config

const (
	//命名规则：redis 数据库的单个key的命名规则为：去除了 R 的结尾的数据库名_KEY_PREFIX = "前半部分的key"+:（英文冒号)

	//  VideoId_CommentIdR *redis.Client
	VideoId_CommentId_KEY_PREFIX = "VideoId:"

	//	CommentId_CommentR *redis.Client
	CommentId_Comment_KEY_PREFIX = "CommentId:"

	//	//F : favorite
	//	UserId_FVideoIdR *redis.Client
	UserId_FVideoId_KEY_PREFIX = "UserId:"

	//	VideoId_VideoR   *redis.Client
	VideoId_Video_KEY_PREFIX = "VideoId:"

	//	UserId_UserR     *redis.Client
	UserId_User_KEY_PREFIX = "UserId:"
)
