package repository

import (
	"strings"
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

type postDTOWithTags struct {
	ID          string     `db:"id"`
	Title       string     `db:"title"`
	Content     string     `db:"content"`
	Slug        string     `db:"slug"`
	Draft       bool       `db:"draft"`
	PublishedAt *time.Time `db:"published_at"`
	TagIDs      string     `db:"tag_id"`
	TagNames    string     `db:"tags"`
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

func postDTOWithTagsTOPostAndTags(dto *postDTOWithTags) (*model.Post, []*model.Tag) {
	post := &model.Post{
		ID:          dto.ID,
		Title:       dto.Title,
		Content:     dto.Content,
		Slug:        dto.Slug,
		Draft:       dto.Draft,
		PublishedAt: dto.PublishedAt,
	}
	tagIDs := strings.Split(dto.TagIDs, ",")
	tagNames := strings.Split(dto.TagNames, ",")
	tags := make([]*model.Tag, len(tagIDs))
	for i := 0; i < len(tagIDs); i++ {
		tags[i] = &model.Tag{
			ID:   tagIDs[i],
			Name: tagNames[i],
		}
	}
	return post, tags
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
