package repository

import (
	"time"

	"github.com/Le0tk0k/blog-server/model"
)

type postDTO struct {
	ID          string     `db:"id"`
	Title       string     `db:"title"`
	Content     string     `db:"content"`
	Slug        string     `db:"slug"`
	Draft       bool       `db:"draft"`
	PublishedAt *time.Time `db:"published_at"`
}

type tagDTO struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func postToDTO(post *model.Post) *postDTO {
	return &postDTO{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		Slug:        post.Slug,
		Draft:       post.Draft,
		PublishedAt: post.PublishedAt,
	}
}

func dtoToPost(dto *postDTO) *model.Post {
	return &model.Post{
		ID:          dto.ID,
		Title:       dto.Title,
		Content:     dto.Content,
		Slug:        dto.Slug,
		Draft:       dto.Draft,
		PublishedAt: dto.PublishedAt,
	}
}

func tagToDTO(tag *model.Tag) *tagDTO {
	return &tagDTO{
		ID:   tag.ID,
		Name: tag.Name,
	}
}

func dtoToTag(dto *tagDTO) *model.Tag {
	return &model.Tag{
		ID:   dto.ID,
		Name: dto.Name,
	}
}
