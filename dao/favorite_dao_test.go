package dao

import (
	"testing"
)

func TestInsertFavoriteInfo(t *testing.T) {
	favorite := FavoriteDao{
		UserId:    1,
		VideoId:   1,
		Favorited: ISFAVORITE,
	}

	err := InsertFavoriteInfo(favorite)
	if err != nil {
		t.Errorf("Failed to insert favorite info: %v", err)
	}

	// Cleanup
	err = DeleteFavoriteInfo(favorite.UserId, favorite.VideoId)
	if err != nil {
		t.Errorf("Failed to delete favorite info after test: %v", err)
	}
}

func TestUpdateFavoriteInfo(t *testing.T) {
	favorite := FavoriteDao{
		UserId:    1,
		VideoId:   1,
		Favorited: ISFAVORITE,
	}

	// Insert a dummy record to update
	err := InsertFavoriteInfo(favorite)
	if err != nil {
		t.Errorf("Failed to insert favorite info for update test: %v", err)
	}

	err = UpdateFavoriteInfo(favorite.UserId, favorite.VideoId, ISNOTFAVORITE)
	if err != nil {
		t.Errorf("Failed to update favorite info: %v", err)
	}

	// Cleanup
	err = DeleteFavoriteInfo(favorite.UserId, favorite.VideoId)
	if err != nil {
		t.Errorf("Failed to delete favorite info after test: %v", err)
	}
}

func TestDeleteFavoriteInfo(t *testing.T) {
	favorite := FavoriteDao{
		UserId:    1,
		VideoId:   1,
		Favorited: ISFAVORITE,
	}

	// Insert a dummy record to delete
	err := InsertFavoriteInfo(favorite)
	if err != nil {
		t.Errorf("Failed to insert favorite info for delete test: %v", err)
	}

	err = DeleteFavoriteInfo(favorite.UserId, favorite.VideoId)
	if err != nil {
		t.Errorf("Failed to delete favorite info: %v", err)
	}
}

func TestGetFavoriteIdListByUserId(t *testing.T) {
	favorite := FavoriteDao{
		UserId:    1,
		VideoId:   1,
		Favorited: ISFAVORITE,
	}

	// Insert a dummy record
	err := InsertFavoriteInfo(favorite)
	if err != nil {
		t.Errorf("Failed to insert favorite info for fetch test: %v", err)
	}

	favoriteIdList, _, err := GetFavoriteIdListByUserId(favorite.UserId)
	if err != nil || len(favoriteIdList) == 0 || favoriteIdList[0] != favorite.VideoId {
		t.Errorf("Failed to fetch favorite info by user id or data mismatched: %v", err)
	}

	// Cleanup
	err = DeleteFavoriteInfo(favorite.UserId, favorite.VideoId)
	if err != nil {
		t.Errorf("Failed to delete favorite info after test: %v", err)
	}
}

func TestIsVideoFavoritedByUser(t *testing.T) {
	// 创建一个测试的Favorite数据
	testFavorite := FavoriteDao{
		UserId:    1,
		VideoId:   1,
		Favorited: ISFAVORITE,
	}

	err := InsertFavoriteInfo(testFavorite)
	if err != nil {
		t.Errorf("Failed to insert test favorite: %v", err)
	}

	isFavorited, err := IsVideoFavoritedByUser(1, 1)
	if err != nil {
		t.Errorf("Failed to check if video is favorited: %v", err)
	}

	if !isFavorited {
		t.Errorf("Expected video to be favorited by user but it was not")
	}

	// 清理测试数据
	err = DeleteFavoriteInfo(1, 1)
	if err != nil {
		t.Errorf("Failed to delete test favorite: %v", err)
	}
}

func TestVideoFavoritedCount(t *testing.T) {
	// 创建一个测试的Favorite数据
	testFavorite := FavoriteDao{
		UserId:    1,
		VideoId:   1,
		Favorited: ISFAVORITE,
	}

	err := InsertFavoriteInfo(testFavorite)
	if err != nil {
		t.Errorf("Failed to insert test favorite: %v", err)
	}

	favoriteCnt, err := VideoFavoritedCount(1)
	if err != nil {
		t.Errorf("Failed to get favorite count: %v", err)
	}

	if favoriteCnt != 1 {
		t.Errorf("Expected favorite count to be 1 but got %d", favoriteCnt)
	}

	// 清理测试数据
	err = DeleteFavoriteInfo(1, 1)
	if err != nil {
		t.Errorf("Failed to delete test favorite: %v", err)
	}
}
