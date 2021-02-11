package model

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          string
	Title       string
	Content     string
	Slug        string
	Draft       bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt time.Time
}

// NewPost はPostのポインタを返す
func NewPost() *Post {
	return &Post{
		ID:      uuid.New().String(),
		Title:   "",
		Content: "",
		Slug:    "",
		Draft:   true,
	}
}
