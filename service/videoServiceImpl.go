package service

import (
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	videoServiceImpl *VideoServiceImpl
	videoServiceOnce sync.Once
)

func init() {
	videoServiceImpl = GetVideoServiceInstance()
}

func GetVideoServiceInstance() *VideoServiceImpl {
	videoServiceOnce.Do(func() {
		videoServiceImpl = &VideoServiceImpl{}
	})
	return videoServiceImpl
}

// Publish ==================== publish接口 ====================
func (videoService *VideoServiceImpl) Publish(c *gin.Context, data *multipart.FileHeader, videoTitle string, authorId int64) error {
	// 保证唯一的 videoName
	videoName := uuid.New().String()
	log.Println("videoName:", videoName)
	err := UploadVideoToOSS(c, data, videoName)
	if err != nil {
		return err
	}
	log.Println("视频存入OSS成功！")

	err = dao.UploadVideo(videoName, authorId, videoTitle)
	if err != nil {
		log.Println("视频存入数据库失败！")
		return err
	}
	log.Println("视频存入数据库成功！")
	return nil
}

func UploadVideoToOSS(c *gin.Context, file *multipart.FileHeader, title string) error {
	// 1. 保存上传的文件到临时目录
	tempFilePath := "/tmp/uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		return err
	}

	// 2. 打开临时文件
	videoContent, err := os.Open(tempFilePath)
	if err != nil {
		log.Print("os.Open:", tempFilePath)
		return err
	}

	// 使用defer确保在函数结束后关闭文件
	defer func() {
		_ = videoContent.Close()
		_ = os.Remove(tempFilePath) // 删除临时文件
	}()

	// 3. 上传到OSS
	client, err := oss.New(config.OssEndpoint, config.OssAccessKeyId, config.OssAccessKeySecret)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(config.OssBucketName)
	if err != nil {
		return err
	}

	err = bucket.PutObject(config.OssVideoDir+title+".mp4", videoContent)
	if err != nil {
		return err
	}
	return nil
}

// ======================= feed接口 =======================

// PublishList 获取某位用户的视频信息list
func (videoService *VideoServiceImpl) PublishList(userId int64) ([]dao.Video, error) {
	plainVideos, err := dao.GetVideosByUserId(userId)
	if err != nil {
		log.Println("GetVideosByUserId:", err)
		return nil, err
	}
	return plainVideos, nil
}

// Feed 按投稿时间倒序的视频list
func (videoService *VideoServiceImpl) Feed(latestTime time.Time) ([]dao.Video, time.Time, error) {
	plainVideos, err := dao.GetVideosByLatestTime(latestTime)
	if err != nil {
		log.Println("GetVideosByLatestTime:", err)
		return nil, time.Time{}, err
	}
	log.Println("plainVideos:", plainVideos)
	if len(plainVideos) == 0 {
		return []dao.Video{}, time.Time{}, errors.New("plainVideos is empty")
	}
	return plainVideos, plainVideos[len(plainVideos)-1].CreatedAt, nil

}

// GetVideoCnt 根据userId获取作品数量
func (videoService *VideoServiceImpl) GetVideoCnt(userId int64) (int64, error) {
	return dao.GetVideoCnt(userId)
}
