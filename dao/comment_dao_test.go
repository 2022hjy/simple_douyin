package dao

import (
	"gorm.io/gorm"
	"simple_douyin/middleware/database"
	"testing"
)

var testDb *gorm.DB

func init() {
	database.Init()
}

func TestAddComment(t *testing.T) {
	comment := CommentDao{
		UserId:  1,
		VideoId: 1,
		Content: "Test Comment",
	}

	addedComment, err := AddComment(comment)
	if err != nil {
		t.Errorf("Failed to add comment: %v", err)
	}

	if addedComment.Content != comment.Content {
		t.Errorf("Expected content %s but got %s", comment.Content, addedComment.Content)
	}

	// Cleanup
	err = DeleteComment(addedComment.Id)
	if err != nil {
		t.Errorf("Failed to delete comment after test: %v", err)
	}
}

func TestDeleteComment(t *testing.T) {
	comment := CommentDao{
		UserId:  2,
		VideoId: 2,
		Content: "Another Test Comment",
	}

	addedComment, _ := AddComment(comment)
	err := DeleteComment(addedComment.Id)
	if err != nil {
		t.Errorf("Failed to delete comment: %v", err)
	}
}

func TestGetCommentList(t *testing.T) {
	videoID := int64(3)

	_, err := GetCommentList(videoID)
	if err != nil {
		t.Errorf("Failed to fetch comment list: %v", err)
	}
}

func TestGetCommentCnt(t *testing.T) {
	videoID := int64(4)

	_, err := GetCommentCnt(videoID)
	if err != nil {
		t.Errorf("Failed to get comment count: %v", err)
	}
}

func TestGetUserFromCommentId(t *testing.T) {
	comment := CommentDao{
		UserId:  5,
		VideoId: 5,
		Content: "Yet Another Test Comment",
	}

	addedComment, _ := AddComment(comment)
	_, err := GetUserFromCommentId(addedComment.Id)
	if err != nil {
		t.Errorf("Failed to get user from comment ID: %v", err)
	}

	// Cleanup
	DeleteComment(addedComment.Id)
}
