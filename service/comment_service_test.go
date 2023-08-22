package service

import (
	"reflect"
	"testing"
)

func TestGetCommentServiceInstance(t *testing.T) {
	instance1 := GetCommentServiceInstance()
	instance2 := GetCommentServiceInstance()

	if instance1 == nil || instance2 == nil {
		t.Fatal("Expected non-nil instance, got nil")
	}

	if instance1 != instance2 {
		t.Fatal("Expected singleton instances to be the same")
	}
}

func TestComment(t *testing.T) {
	cs := GetCommentServiceInstance()

	_, err := cs.Comment(1, 1, "Test comment content")
	if err != nil {
		t.Fatalf("Failed to create comment: %v", err)
	}
	// You can further verify if the comment is correctly added in your mock database
}

func TestDeleteComment(t *testing.T) {
	cs := GetCommentServiceInstance()

	err := cs.DeleteComment(1, 1) // Assuming a comment with videoId 1 and commentId 1 exists
	if err != nil {
		t.Fatalf("Failed to delete comment: %v", err)
	}
	// You can further verify if the comment is correctly deleted from your mock database
}

func TestGetCommentList(t *testing.T) {
	cs := GetCommentServiceInstance()

	comments, err := cs.GetCommentList(1)
	if err != nil {
		t.Fatalf("Failed to get comments list: %v", err)
	}

	if len(comments) == 0 {
		t.Fatalf("Expected non-empty comments list")
	}
}

func TestGetCommentListFromRedis(t *testing.T) {
	cs := GetCommentServiceInstance()

	commentsFromRedis, err := cs.getCommentListFromRedis(1)
	if err != nil {
		t.Fatalf("Failed to get comments from Redis: %v", err)
	}

	commentsFromDB, err := cs.getCommentListFromDB(1)
	if err != nil {
		t.Fatalf("Failed to get comments from DB: %v", err)
	}

	if !reflect.DeepEqual(commentsFromRedis, commentsFromDB) {
		t.Fatalf("Comments from Redis and DB do not match")
	}
}

func TestGetCommentListFromDB(t *testing.T) {
	cs := GetCommentServiceInstance()

	comments, err := cs.getCommentListFromDB(1)
	if err != nil {
		t.Fatalf("Failed to get comments from DB: %v", err)
	}

	if len(comments) == 0 {
		t.Fatalf("Expected non-empty comments list from DB")
	}
}

// Additional tests can be added for other helper functions like updateCommentRedis and deleteCommentRedis.
