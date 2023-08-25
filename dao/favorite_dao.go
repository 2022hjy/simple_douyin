package dao

import (
	"errors"
	"log"
	"simple_douyin/config"
	"simple_douyin/middleware/redis"
	"strconv"
	"time"
)

/**
{
//
    "status_code": "200",
    "status_msg": "Success",
    "video_list": [
        {
            "id": 123,
            "author": {
                "id": 456,
                "name": "JohnDoe",
                "follow_count": 120,
                "follower_count": 5000,
                "is_follow": false,
                "avatar": "profile.jpg",
                "background_image": "background.jpg",
                "signature": "Passionate about travel and adventure",
                "total_favorited": "10K",
                "work_count": 25,
                "favorite_count": 300
            },
            "play_url": "video123.mp4",
            "cover_url": "cover123.jpg",
            "favorite_count": 50,
            "comment_count": 15,
            "is_favorite": true,
            "title": "Exploring the Amazon Rainforest"
        },

        {
            "id": 124,
            "author": {
                "id": 457,
                "name": "JaneSmith",
                "follow_count": 85,
                "follower_count": 2500,
                "is_follow": true,
                "avatar": "profile_jane.jpg",
                "background_image": "background_jane.jpg",
                "signature": "Nature lover and photographer",
                "total_favorited": "5.2K",
                "work_count": 50,
                "favorite_count": 120
            },
            "play_url": "video124.mp4",
            "cover_url": "cover124.jpg",
            "favorite_count": 30,
            "comment_count": 8,
            "is_favorite": false,
            "title": "Capturing the Serenity of Mountains"
        }
    ]
}

{
    "status_code": "404",
    "status_msg": "Videos not found",
    "video_list": []
}
*/

type FavoriteDao struct {
	Id        int64
	UserId    int64
	VideoId   int64
	Favorited int // 0: not favorite, 1: favorite
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (FavoriteDao) TableName() string {
	return "favorite"
}

const (
	ISFAVORITE    = 1
	ISNOTFAVORITE = 0
)

func GetFavoriteIdListByUserId(userId int64) ([]int64, int64, error) {
	var FavoriteIdList []int64
	result := Db.Model(&FavoriteDao{}).
		Where("user_id=? and is_favorite = ?", userId, ISFAVORITE).
		Order("created_at desc").
		Pluck("video_id", &FavoriteIdList)
	favoriteCnt := result.RowsAffected
	if result.Error != nil {
		log.Println("FavoritedVideoIdList:", result.Error.Error())
		return nil, -1, result.Error
	}
	return FavoriteIdList, favoriteCnt, nil
}

func VideoFavoritedCount(videoId int64) (int64, error) {
	var count int64
	err := Db.Model(FavoriteDao{}).Where("video_id =? and is_favorite = ?", videoId, ISFAVORITE).Count(&count).Error
	if err != nil {
		log.Println("FavoriteDao-Count: return count failed")
		return -1, errors.New("find favorites count failed")
	}
	log.Println("FavoriteDao-Count: return count success")
	return count, nil
}

func UsersOfFavoriteVideo(videoId int64) ([]int64, int64, error) {
	var userIdList []int64
	result := Db.Model(&FavoriteDao{}).
		Where("video_id=? and is_favorite=?", videoId, 1).
		Pluck("user_id", &userIdList)
	favoriteCnt := result.RowsAffected
	if favoriteCnt == 0 {
		return nil, 0, result.Error
	}
	if result.Error != nil {
		log.Println("UsersOfFavoriteVideo:", result.Error.Error())
		return nil, 0, result.Error
	}
	return userIdList, favoriteCnt, nil
}

// UpdateFavoriteInfo update favorite info
func UpdateFavoriteInfo(userId int64, videoId int64, favorited int8) error {
	result := Db.Model(FavoriteDao{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).Update("is_favorite", favorited)
	if result.RowsAffected == 0 {
		return errors.New("update favorite failed, record not exists")
	}
	log.Println("FavoriteDao-UpdateFavoriteInfo: return success")
	return nil
}

func InsertFavoriteInfo(favorite FavoriteDao) error {
	err := Db.Model(FavoriteDao{}).Create(&favorite).Error
	if err != nil {
		log.Println(err.Error())
		return errors.New("insert favorites failed")
	}
	return nil
}

func DeleteFavoriteInfo(userId int64, videoId int64) error {
	err := Db.Model(FavoriteDao{}).Where("user_id = ? and video_id = ?", userId, videoId).Delete(&FavoriteDao{}).Error
	if err != nil {
		log.Println(err.Error())
		return errors.New("delete favorite failed")
	}
	return nil
}

func IsVideoFavoritedByUser(userId int64, videoId int64) (bool, error) {
	var isFavorited int8
	result := Db.Model(FavoriteDao{}).Select("is_favorite").Where("user_id= ? and video_id= ?", userId, videoId).First(&isFavorited)
	c := result.RowsAffected
	if result == nil {
		return false, errors.New("current video is not exist")
	}
	if c == 0 {
		return false, errors.New("current user haven not favorited current video")
	}
	if result.Error != nil {
		log.Println(result.Error)
	}
	return true, nil
}

func GetFavoriteCountByUser(userId int64) (int64, error) {
	var count int64
	err := Db.Model(FavoriteDao{}).Where(map[string]interface{}{"user_id": userId, "is_favorite": 1}).Count(&count).Error
	if err != nil {
		log.Println("FavoriteDao-Count: return count failed")
		return -1, errors.New("find favorites count failed")
	}
	log.Println("FavoriteDao-Count: return count success")
	return count, nil
}

func IsFavoritedByUser(userId int64, videoId int64) (bool, error) {
	var favorite FavoriteDao
	result := Db.Model(FavoriteDao{}).
		Where("user_id = ? and video_id = ? and is_favorite = ?", userId, videoId, 1).
		First(&favorite)
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func GettotalFavorited(userId int64) (int, error) {
	UIdFVIdR := redis.Clients.UserId_FavoriteVideoIdR
	key := config.UserId_FVideoId_KEY_PREFIX + strconv.FormatInt(userId, 10)

	// 尝试从 Redis 中获取数据
	count, err := redis.CountElements(UIdFVIdR, key)
	if err != nil || count == 0 { // 如果 Redis 返回错误或计数为0
		// 从数据库中获取点赞的视频数量
		_, countFromDb, errDb := GetFavoriteIdListByUserId(userId)
		if errDb != nil {
			log.Println("从数据库获取点赞的视频数量失败：", errDb)
			return 0, errDb
		}
		// 将数据重新设置到 Redis
		redis.SetValueWithRandomExp(UIdFVIdR, key, countFromDb)
		return int(countFromDb), nil
	}
	return count, nil
}

func GetfavoriteCount(videoId int64) (int, error) {
	UIdFVIdR := redis.Clients.UserId_FavoriteVideoIdR
	key := config.VideoId_FavoritebUserId_KEY_PREFIX + strconv.FormatInt(videoId, 10)

	// 尝试从 Redis 中获取数据
	count, err := redis.CountElements(UIdFVIdR, key)
	if err != nil || count == 0 { // 如果 Redis 返回错误或计数为0
		// 从数据库中获取点赞的用户数量
		countFromDb, errDb := VideoFavoritedCount(videoId)
		if errDb != nil {
			log.Println("从数据库获取点赞的用户数量失败：", errDb)
			return 0, errDb
		}
		// 将数据重新设置到 Redis
		err := redis.SetValueWithRandomExp(UIdFVIdR, key, countFromDb)
		if err != nil {
			return 0, err
		}
		return int(countFromDb), nil
	}
	return count, nil
}
