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
	"simple_douyin/controller"
	"simple_douyin/dao"
	"simple_douyin/util"
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

// Publish ==================== publish接口 ====================
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

// ======================= feed接口 =======================
// getRespVideos dao.video --> FeedResponse
func (videoService *VideoServiceImpl) getRespVideos(plainVideos []dao.Video, userId int64) ([]controller.VideoResponse, error) {
	var douyinVideos []controller.VideoResponse
	for _, video := range plainVideos {
		responce, err := util.ConvertDBVideoToResponse(video, userId)
		if err != nil {
			log.Println("getRespVideos:", err)
			return []controller.VideoResponse{}, nil
		}
		douyinVideos = append(douyinVideos, responce)
	}
	return douyinVideos, nil
}

// PublishList 获取某位用户的视频信息list
func (videoService *VideoServiceImpl) PublishList(userId int64) ([]controller.VideoResponse, error) {
	videos := make([]controller.VideoResponse, 0, config.VideoInitNum)
	plainVideos, err := dao.GetVideosByUserId(userId)
	if err != nil {
		log.Println("GetVideosByUserId:", err)
		return nil, err
	}
	videos, err = videoService.getRespVideos(plainVideos, userId)
	if err != nil {
		log.Println("getRespVideos:", err)
		return nil, err
	}
	return videos, nil
}

// Feed 按投稿时间倒序的视频list
func (videoService *VideoServiceImpl) Feed(latestTime time.Time, userId int64) ([]controller.VideoResponse, time.Time, error) {
	douyinVideos := make([]controller.VideoResponse, 0, config.VideoInitNumPerRefresh)
	plainVideos, err := dao.GetVideosByLatestTime(latestTime)
	if err != nil {
		log.Println("GetVideosByLatestTime:", err)
		return nil, time.Time{}, err
	}
	if plainVideos == nil || len(plainVideos) == 0 {
		return []controller.VideoResponse{}, time.Time{}, nil
	}

	douyinVideos, err = videoServiceImp.getRespVideos(plainVideos, userId)
	return douyinVideos, plainVideos[len(plainVideos)-1].CreatedAt, nil
}

// GetVideoListById 根据videoIdList查询视频信息
func (videoService *VideoServiceImpl) GetVideoListById(videoIdList []int64, userId int64) ([]controller.VideoResponse, error) {
	videoList := make([]controller.VideoResponse, 0, config.VideoInitNum)
	plainVideoList, _ := dao.GetVideoListById(videoIdList)
	videoList, err := videoService.getRespVideos(plainVideoList, userId)
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
