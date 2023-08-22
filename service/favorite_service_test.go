package service

import (
	"testing"
)

func TestFavoriteAction(t *testing.T) {
	fs := GetFavoriteServiceInstance()

	err := fs.FavoriteAction(1, 1)
	if err != nil {
		t.Errorf("Failed to perform favorite action: %v", err)
	}
}

func TestGetFavoriteList(t *testing.T) {
	fs := GetFavoriteServiceInstance()

	videos, err := fs.GetFavoriteList(1)
	if err != nil {
		t.Errorf("Failed to get favorite list: %v", err)
	}

	if len(videos) == 0 {
		t.Errorf("Expected videos, got none")
	}

	for _, video := range videos {
		if video.Id == 0 {
			t.Errorf("Expected valid video, got %v", video)
		}
	}
}

func TestGettotalFavorited(t *testing.T) {
	count, err := GettotalFavorited(1)
	if err != nil {
		t.Errorf("Failed to get total favorited: %v", err)
	}

	if count < 0 {
		t.Errorf("Expected count >= 0, got %d", count)
	}
}

func TestGetfavoriteCount(t *testing.T) {
	count, err := GetfavoriteCount(1)
	if err != nil {
		t.Errorf("Failed to get favorite count: %v", err)
	}

	if count < 0 {
		t.Errorf("Expected count >= 0, got %d", count)
	}
}
