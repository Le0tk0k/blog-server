package handler

import (
	"time"

	"github.com/Le0tk0k/blog-server/model"
)

type postJSON struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Slug        string     `json:"slug"`
	Draft       bool       `json:"draft"`
	PublishedAt *time.Time `json:"published_at"`
}

type tagJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type postWithTagsJSON struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Slug        string     `json:"slug"`
	Draft       bool       `json:"draft"`
	PublishedAt *time.Time `json:"published_at"`
	Tags        []*tagJSON `json:"tags"`
}

func postToJSON(post *model.Post) *postJSON {
	return &postJSON{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		Slug:        post.Slug,
		Draft:       post.Draft,
		PublishedAt: post.PublishedAt,
	}
}

func jsonToPOST(json *postJSON) *model.Post {
	return &model.Post{
		ID:          json.ID,
		Title:       json.Title,
		Content:     json.Content,
		Slug:        json.Slug,
		Draft:       json.Draft,
		PublishedAt: json.PublishedAt,
	}
}

func postToPostWithTagsJSON(post *model.Post) *postWithTagsJSON {
	tagsJSON := make([]*tagJSON, len(post.Tags))
	for i, tag := range post.Tags {
		tagsJSON[i] = tagToJSON(tag)
	}
	return &postWithTagsJSON{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		Slug:        post.Slug,
		Draft:       post.Draft,
		PublishedAt: post.PublishedAt,
		Tags:        tagsJSON,
	}
}

func tagToJSON(tag *model.Tag) *tagJSON {
	return &tagJSON{
		ID:   tag.ID,
		Name: tag.Name,
	}
}

func jsonToTag(json *tagJSON) *model.Tag {
	return &model.Tag{
		ID:   json.ID,
		Name: json.Name,
	}
}
