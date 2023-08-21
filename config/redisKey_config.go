package config

const (
	//命名规则：redis 数据库的单个key的命名规则为：去除了 R 的结尾的数据库名_KEY = "前半部分的key"+:（英文冒号)

	//  VideoId_CommentIdR *redis.Client
	VideoId_CommentI_KEY = "VideoId:"

	//	CommentId_CommentR *redis.Client
	CommentId_Comment_KEY = "CommentId:"

	//	//F : favorite
	//	UserId_FVideoIdR *redis.Client
	UserId_FVideoId_KEY = "UserId:"

	//	VideoId_VideoR   *redis.Client
	VideoId_Video_KEY = "VideoId:"
)
