package service

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	//提供了生成和处理 UUID（通用唯一标识符）的功能
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"simple_douyin/config"
	"simple_douyin/dao"
	"sync"
	"time"
)

// VideoServiceImpl TODO  会调用到以下接口的函数
type VideoServiceImpl struct {
	//CommentService
	//LikeService
	//UserService
}

var (
	videoServiceImp  *VideoServiceImpl
	videoServiceOnce sync.Once
)

// GetVideoServiceInstance Go 单例模式：https://www.liwenzhou.com/posts/Go/singleton/
func GetVideoServiceInstance() *VideoServiceImpl {
	videoServiceOnce.Do(func() {
		videoServiceImp = &VideoServiceImpl{
			//UserService:    &UserServiceImpl{},
			//CommentService: &CommentServiceImpl{},
			//LikeService:    &LikeServiceImpl{},
		}
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
	err = dao.UploadVideo(videoName, userId, title)
	if err != nil {
		log.Println("视频存入数据库失败！")
		return err
	}
	return nil
}

// UploadVideoToOSS 将上传的视频文件保存到阿里云 OSS
func UploadVideoToOSS(file *multipart.FileHeader, videoName string) error {
	// ⭐ 创建OSSClient实例，使用了配置文件中的 OSS 相关信息
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	client, err := oss.New(config.OSS_ENDPOINT, config.OSS_ACCESS_KEY_ID, config.OSS_ACCESS_KEY_SECRET)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// 获取存储空间(Bucket)
	bucket, err := client.Bucket(config.OSS_BUCKET_NAME)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	//打开文件
	fd, err := file.Open()
	if err != nil {
		log.Println("file open failed!")
		return err
	}
	//defer 关键字确保文件在函数结束时被关闭
	defer func(fd multipart.File) {
		err := fd.Close()
		if err != nil {

		}
	}(fd)

	//⭐ 将打开的文件对象上传到指定路径 OSS Bucket 中
	err = bucket.PutObject(config.OSS_VIDEO_DIR+videoName+".mp4", fd)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}

//===================已测试完成接口===================

// PublishList 获取某位用户的视频信息list
func (videoService *VideoServiceImpl) PublishList(userId int64) ([]Video, error) {
	videos := make([]Video, 0, config.VIDEO_INIT_NUM)
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
	videos := make([]Video, 0, config.VIDEO_NUM_PER_REFRESH)
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
		err := videoService.combineVideo(&video, &tmpVideo, userId)
		if err != nil {
			return err
		}
		*videos = append(*videos, video)
	}
	return nil
}

// TODO TEST 组装 controller 层所需的 Video 结构   获取视频信息
func (videoService *VideoServiceImpl) combineVideo(video *Video, plainVideo *dao.Video, userId int64) error {
	//videoServiceNew := GetVideoServiceInstance()
	// 解决循环依赖

	//建立协程组，确保所有协程的任务完成后才退出本方法
	var wg sync.WaitGroup
	wg.Add(4)
	video.Video = *plainVideo
	//视频作者信息 —— 登录情况下返回用户信息, 第一个id是视频作者的id，第二个id是我们用户的id
	go func(v *Video) {
		//TODO: UserServiceImpl 待补全
		//user, err := videoServiceNew.GetUserLoginInfoByIdWithCurId(v.AuthorId, userId)
		//if err != nil {
		//	return
		//}
		user := User{
			Id:   1,
			Name: "test",
		}
		v.Author = user
		wg.Done()
	}(video)

	// 视频点赞数量
	go func(v *Video) {
		//TODO LikeServiceImpl 待补全
		//favoriteCount, err := videoServiceNew.GetVideoLikedCount(v.Id)
		//if err != nil {
		//	return
		//}
		favoriteCount := int64(100)
		v.FavoriteCount = favoriteCount
		wg.Done()
	}(video)

	// 视频评论数量
	go func(v *Video) {
		//TODO CommentServiceImpl 待补全
		//count, err := videoServiceNew.GetCommentCnt(v.Id)
		//if err != nil {
		//	return
		//}
		count := int64(10)
		v.CommentCount = count
		wg.Done()
	}(video)

	// 当前登录用户/游客是否对该视频点过赞
	go func(v *Video) {
		//TODO LikeServiceImpl 待补全
		//isFavorite, err := videoServiceNew.IsLikedByUser(userId, v.Id)
		//if err != nil {
		//	return
		//}
		isFavorite := false
		// 等待点赞服务，获取是否点赞
		v.IsFavorite = isFavorite
		wg.Done()
	}(video)

	wg.Wait()
	return nil
}

// GetVideoListById 根据videoIdList查询视频信息
func (videoService *VideoServiceImpl) GetVideoListById(videoIdList []int64, userId int64) ([]Video, error) {
	videoList := make([]Video, 0, config.VIDEO_INIT_NUM)
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
