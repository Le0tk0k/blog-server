package repository

import (
	"fmt"

	"github.com/Le0tk0k/blog-server/model"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type PostRepository interface {
	StorePost(post *model.Post) error
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
	_, err := p.db.Exec("INSERT INTO posts (id, title, content, slug, draft) VALUES (?, ?, ?, ?, ?)", post.ID, post.Title, post.Content, post.Slug, post.Draft)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return fmt.Errorf("StorePost: cannot store post: %w", model.ErrPostAlreadyExisted)
		}
		return fmt.Errorf("StorePost: cannot store post: %w", err)
	}
	return nil
}
