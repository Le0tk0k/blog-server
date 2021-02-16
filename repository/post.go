package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Le0tk0k/blog-server/model"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type PostRepository interface {
	StorePost(post *model.Post) error
	FindPostByID(id string) (*model.Post, error)
	FindAllPosts() ([]*model.Post, error)
	UpdatePost(post *model.Post) error
	DeletePostByID(id string) error
}

type postRepository struct {
	db *sqlx.DB
}

// NewPostRepository はPostRepositoryを返す
func NewPostRepository(db *sqlx.DB) PostRepository {
	return &postRepository{db: db}
}

// StorePost は記事を新規保存する
func (p *postRepository) StorePost(post *model.Post) error {
	dto := postToDTO(post)
	_, err := p.db.Exec("INSERT INTO posts (id, title, content, slug, draft) VALUES (?, ?, ?, ?, ?)", dto.ID, dto.Title, dto.Content, dto.Slug, dto.Draft)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return fmt.Errorf("StorePost: cannot store post: %w", model.ErrPostAlreadyExisted)
		}
		return fmt.Errorf("StorePost: cannot store post: %w", err)
	}
	return nil
}

// FindPostByID はidを持つ記事を取得する
func (p *postRepository) FindPostByID(id string) (*model.Post, error) {
	var dto postWithTagsDTO
	query := "SELECT posts.id, posts.title, posts.content, posts.slug, posts.draft, posts.published_at, GROUP_CONCAT(tags.id) AS tag_id, GROUP_CONCAT(tags.name) AS tags FROM posts LEFT JOIN posts_tags on posts.id = posts_tags.post_id LEFT JOIN tags on posts_tags.tag_id = tags.id WHERE posts.id = ? GROUP BY posts.id"
	if err := p.db.Get(&dto, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("FindPostByID: cannot find post: %w", model.ErrPostNotFound)
		}
		return nil, fmt.Errorf("FindPostByID: cannot find post: %w", err)
	}
	post := postWithTagsDTOTOPost(&dto)
	return post, nil
}

// FindAllPosts は全記事を取得する
func (p *postRepository) FindAllPosts() ([]*model.Post, error) {
	var dtos []*postWithTagsDTO
	query := "SELECT posts.id, posts.title, posts.content, posts.slug, posts.draft, posts.published_at, GROUP_CONCAT(tags.id) AS tag_id, GROUP_CONCAT(tags.name) AS tags FROM posts LEFT JOIN posts_tags on posts.id = posts_tags.post_id LEFT JOIN tags on posts_tags.tag_id = tags.id GROUP BY posts.id"
	if err := p.db.Select(&dtos, query); err != nil {
		return nil, fmt.Errorf("FindAllPosts: cannot find post: %w", err)
	}

	posts := make([]*model.Post, len(dtos))
	for i, dto := range dtos {
		posts[i] = postWithTagsDTOTOPost(dto)
	}

	return posts, nil
}

// UpdatePost はidを持つ記事を更新する
func (p *postRepository) UpdatePost(post *model.Post) error {
	dto := postToDTO(post)
	_, err := p.db.Exec("UPDATE posts SET title = ?, content = ?, slug = ?, draft = ?, published_at = ? WHERE id = ?", dto.Title, dto.Content, dto.Slug, dto.Draft, dto.PublishedAt, dto.ID)
	if err != nil {
		return fmt.Errorf("UpdatePost: cannot update post: %w", err)
	}
	return nil
}

// DeletePostByID はidを持つ記事を削除する
func (p *postRepository) DeletePostByID(id string) error {
	result, err := p.db.Exec("DELETE FROM posts WHERE id = ?", id)
	if rowsAffected, err := result.RowsAffected(); rowsAffected == 0 && err == nil {
		return fmt.Errorf("DeletePostByID: cannot find post: %w", model.ErrPostNotFound)
	}
	if err != nil {
		return fmt.Errorf("DeletePostByID: cannot delete post: %w", err)
	}
	return nil
}
