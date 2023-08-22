package service

import (
	"bytes"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"simple_douyin/config"
	"simple_douyin/dao"
	"sync"
	"time"
)

type VideoServiceImpl struct{}

var (
	videoServiceImp  *VideoServiceImpl
	videoServiceOnce sync.Once
)

func init() {
	videoServiceImp = GetVideoServiceInstance()
}

func GetVideoServiceInstance() *VideoServiceImpl {
	videoServiceOnce.Do(func() {
		videoServiceImp = &VideoServiceImpl{}
	})
	return videoServiceImp
}

// Publish 将视频上传到oss并保存到数据库中
func (videoService *VideoServiceImpl) Publish(data *multipart.FileHeader, title string, userId int64) error {
	// 保证唯一的 videoName
	videoName := uuid.New().String()

	err := UploadVideoToOSS(data, videoName)
	if err != nil {
		return err
	}
	log.Println("视频存入OSS成功！")

	err = dao.UploadVideo(videoName, userId, title)
	if err != nil {
		log.Println("视频存入数据库失败！")
		return err
	}
	log.Println("视频存入数据库成功！")
	return nil
}

func UploadVideoToOSS(file *multipart.FileHeader, title string) error {
	client, err := oss.New(config.OssEndpoint, config.OssAccessKeyId, config.OssAccessKeySecret)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(config.OssBucketName)
	if err != nil {
		return err
	}

	videoContent, err := os.Open(file.Filename)
	if err != nil {
		return err
	}
	defer videoContent.Close()

	videoData, err := ioutil.ReadAll(videoContent)
	if err != nil {
		return err
	}

	videoReader := bytes.NewReader(videoData)

	err = bucket.PutObject(config.OssVideoDir+title+".mp4", videoReader)
	if err != nil {
		return err
	}
	return nil
}

//===================已测试完成接口===================

// PublishList 获取某位用户的视频信息list
func (videoService *VideoServiceImpl) PublishList(userId int64) ([]Video, error) {
	videos := make([]Video, 0, config.VideoInitNum)
	plainVideos, err := dao.GetVideosByUserId(userId)
	if err != nil {
		log.Println("GetVideosByUserId:", err)
		return nil, err
	}
	err = videoService.getRespVideos(&videos, &plainVideos, userId)
	if err != nil {
		log.Println("getRespVideos:", err)
		return nil, err
	}
	return videos, nil
}

// Feed 按投稿时间倒序的视频list
func (videoService *VideoServiceImpl) Feed(latestTime time.Time, userId int64) ([]Video, time.Time, error) {
	videos := make([]Video, 0, config.VideoInitNumPerRefresh)
	plainVideos, err := dao.GetVideosByLatestTime(latestTime)
	if err != nil {
		log.Println("GetVideosByLatestTime:", err)
		return nil, time.Time{}, err
	}
	if plainVideos == nil || len(plainVideos) == 0 {
		return videos, time.Time{}, nil
	}
	err = videoService.getRespVideos(&videos, &plainVideos, userId)
	if err != nil {
		log.Println("getRespVideos:", err)
		return nil, time.Time{}, err
	}
	return videos, plainVideos[len(plainVideos)-1].CreatedAt, nil
}

// 获取视频信息list
func (videoService *VideoServiceImpl) getRespVideos(videos *[]Video, plainVideos *[]dao.Video, userId int64) error {
	for _, tmpVideo := range *plainVideos {
		var video Video
		//ConvertDBVideoToResponse(&video, &tmpVideo, userId)
		err := videoService.combineVideo(&video, &tmpVideo, userId)
		if err != nil {
			return err
		}
		*videos = append(*videos, video)
	}
	return nil
}

// TODO TEST dao.video (--> service.video )--> Responce
func (videoService *VideoServiceImpl) combineVideo(video *Video, plainVideo *dao.Video, tokenId int64) error {
	var wg sync.WaitGroup
	wg.Add(4)
	video.Video = *plainVideo
	go func(v *Video) {
		user := dao.User{
			UserId:   1,
			Username: "test",
		}
		v.Author = user
		wg.Done()
	}(video)

	// 视频点赞数量
	go func(v *Video) {
		favoriteCount := int64(100)
		v.FavoriteCount = favoriteCount
		wg.Done()
	}(video)

	// 视频评论数量
	go func(v *Video) {
		count := int64(10)
		v.CommentCount = count
		wg.Done()
	}(video)

	// 当前登录用户/游客是否对该视频点过赞
	go func(v *Video) {
		v.IsFavorite = 0
		wg.Done()
	}(video)

	wg.Wait()
	return nil
}

// GetVideoListById 根据videoIdList查询视频信息
func (videoService *VideoServiceImpl) GetVideoListById(videoIdList []int64, userId int64) ([]Video, error) {
	videoList := make([]Video, 0, config.VideoInitNum)
	plainVideoList, _ := dao.GetVideoListById(videoIdList)
	err := videoService.getRespVideos(&videoList, &plainVideoList, userId)
	if err != nil {
		log.Println("getRespVideos:", err)
		return nil, err
	}
	return videoList, nil
}

// GetVideoCnt 根据userId获取作品数量
func (videoService *VideoServiceImpl) GetVideoCnt(userId int64) (int64, error) {
	return dao.GetVideoCnt(userId)
}
