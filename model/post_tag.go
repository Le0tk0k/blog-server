package model

import "github.com/google/uuid"

type PostTag struct {
	ID     string
	PostID string
	TagID  string
}

// NewPostTag はPostTagのポインタを返す
func NewPostTag(postID, tagID string) *PostTag {
	return &PostTag{
		ID:     uuid.New().String(),
		PostID: postID,
		TagID:  tagID,
	}
}
