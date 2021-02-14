package model

import "github.com/google/uuid"

type Tag struct {
	ID   string
	Name string
}

// NewTag はTagのポインタを返す
func NewTag(name string) *Tag {
	return &Tag{
		ID:   uuid.New().String(),
		Name: name,
	}
}
