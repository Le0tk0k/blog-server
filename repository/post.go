package repository

import (
	"fmt"
	"time"

	"github.com/Le0tk0k/blog-server/model"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type PostRepository interface {
	StorePost(post *model.Post) error
	FindAllPosts() ([]*model.Post, error)
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
