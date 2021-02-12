package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Le0tk0k/blog-server/model"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type PostRepository interface {
	StorePost(post *model.Post) error
	FindPostByID(id string) (*model.Post, error)
	FindAllPosts() ([]*model.Post, error)
	DeletePostByID(id string) error
}

type postRepository struct {
	db *sqlx.DB
}

type postDTO struct {
	ID          string     `db:"id"`
	Title       string     `db:"title"`
	Content     string     `db:"content"`
	Slug        string     `db:"slug"`
	Draft       bool       `db:"draft"`
	PublishedAt *time.Time `db:"published_at"`
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
	var dto postDTO
	if err := p.db.Get(&dto, "SELECT * FROM posts WHERE id = ?", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("FindPostByID: cannot find post: %w", model.ErrPostNotFound)
		}
		return nil, fmt.Errorf("FindPostByID: cannot find post: %w", err)
	}
	return dtoToPost(&dto), nil
}

// FindAllPosts は全記事を取得する
func (p *postRepository) FindAllPosts() ([]*model.Post, error) {
	var dtos []*postDTO
	if err := p.db.Select(&dtos, "SELECT * FROM posts"); err != nil {
		return nil, fmt.Errorf("FindAllPosts: cannot find post: %w", err)
	}

	posts := make([]*model.Post, len(dtos))
	for i, dto := range dtos {
		posts[i] = dtoToPost(dto)
	}

	return posts, nil
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
