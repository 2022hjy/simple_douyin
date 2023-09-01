package dao

import (
	"log"
	"simple_douyin/config"
	"simple_douyin/middleware/database"
	"time"
)

type Video struct {
	Id            int64     `json:"id" gorm:"primaryKey;autoIncrement"` // 视频 Id
	UserInfoId    int64     `json:"user_info_id"`                       // 视频作者 Id
	Title         string    `json:"title"`                              // 视频标题
	PlayUrl       string    `json:"play_url"`                           // 视频播放地址
	CoverUrl      string    `json:"cover_url"`                          // 视频封面地址
	IsFavorite    int       `json:"is_favorite"`                        // 是否被like
	FavoriteCount int64     `json:"favorite_count"`                     // like数
	CommentCount  int64     `json:"comment_count"`                      // 评论数
	CreatedAt     time.Time // 视频创建时间
	UpdatedAt     time.Time // 视频更新时间
}

// TableName 数据库表名映射到 Video 结构体
func (Video) TableName() string {
	return "video"
}

// SaveVideo 保存视频记录到数据库中
func SaveVideo(video Video) error {
	log.Printf("要保存的 video：%+v\n", video) // 添加日志输出
	result := database.Db.Create(&video)
	if result.Error != nil {
		log.Println("数据库保存视频失败！", result.Error)
		return result.Error
	}
	return nil
}

// GetVideosByUserId 根据用户 Id 获取该用户已发布的所有视频
func GetVideosByUserId(userId int64) ([]Video, error) {
	var videos []Video
	result := database.Db.Where("user_info_id = ?", userId).Limit(config.VideoInitNum).Find(&videos)
	if result.Error != nil {
		log.Println("获取用户已发布视频失败！", result.Error)
		return nil, result.Error
	}
	return videos, nil
}

// GetVideosByLatestTime 按投稿时间倒序的视频列表
func GetVideosByLatestTime(latestTime time.Time) ([]Video, error) {
	videos := make([]Video, config.VideoInitNumPerRefresh)
	result := database.Db.Where("created_at < ?", latestTime).
		Order("created_at desc").
		Limit(config.VideoInitNumPerRefresh).
		Find(&videos)
	if result.RowsAffected == 0 {
		log.Println("没有更多视频了！")
		return videos, nil
	}
	if result.Error != nil {
		log.Println("获取视频 Feed 失败！")
		return nil, result.Error
	}
	return videos, nil
}

// GetVideoByVideoId 根据视频 Id 获取视频信息
func GetVideoByVideoId(videoId int64) (Video, error) {
	var video Video
	//Take:没有找到记录，返回的结构体会保持零值（即默认值），不会报错    --区别First
	result := database.Db.Where("id = ?", videoId).Take(&video)
	if result.Error != nil {
		log.Println("根据视频 Id 获取视频失败！")
		return video, result.Error
	}
	return video, nil
}

// GetVideoListById 根据videoIdList查询视频信息
func GetVideoListById(videoIdList []int64) ([]Video, error) {
	var videoList []Video
	result := database.Db.Model(Video{}).
		Where("id in (?)", videoIdList).
		Find(&videoList)
	if result.Error != nil {
		return videoList, result.Error
	}
	return videoList, nil
}

// GetVideoCnt 根据userId获取作品数量
func GetVideoCnt(userId int64) (int64, error) {
	var count int64
	result := database.Db.Model(Video{}).
		Where("user_info_id = ?", userId).
		Count(&count)
	if result.Error != nil {
		log.Println("根据userId获取作品数量失败！")
		return 0, result.Error
	}
	return count, nil
}

// UploadVideo 上传视频
func UploadVideo(videoName string, authorId int64, videoTitle string) error {
	video := Video{
		UserInfoId: authorId,
		Title:      videoTitle,
		PlayUrl:    config.PlayUrlPrefix + videoName + ".mp4",
		CoverUrl:   config.PlayUrlPrefix + videoName + ".mp4" + config.CoverUrlSuffix,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	return SaveVideo(video)
}

//// GetVideoIdListByUserId 根据用户 Id 获取该用户已发布的所有视频
//func GetVideoIdListByUserId(userId int64) ([]int64, error) {
//	var videoIdList []int64
//	result := database.Db.Model(Video{}).Where("user_info_id = ?", userId).Pluck("id", &videoIdList)
//	if result.Error != nil {
//		log.Println("获取用户已发布视频失败！", result.Error)
//		return nil, result.Error
//	}
//	return videoIdList, nil
//}
