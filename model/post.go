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
	PublishedAt *time.Time
	Tags        []*Tag
}

// NewPost はPostのポインタを返す
// PublishedAt は初公開時に設定する。db側のdefaultはnull
func NewPost() *Post {
	return &Post{
		ID:      uuid.New().String(),
		Title:   "",
		Content: "",
		Slug:    "",
		Draft:   true,
	}
}
