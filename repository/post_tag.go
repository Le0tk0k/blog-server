package repository

import (
	"fmt"

	"github.com/Le0tk0k/blog-server/model"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type PostTagRepository interface {
	StorePostTag(postTag *model.PostTag) error
	DeleteByPostID(postID string) error
}

type postTagRepository struct {
	db *sqlx.DB
}

func NewPostTagRepository(db *sqlx.DB) PostTagRepository {
	return &postTagRepository{db: db}
}

func (r *postTagRepository) StorePostTag(postTag *model.PostTag) error {
	dto := postTagToDTO(postTag)
	_, err := r.db.Exec("INSERT INTO posts_tags VALUES (?, ?, ?)", dto.ID, dto.PostID, dto.TagID)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return fmt.Errorf("StorePostTag: cannot store postTag: %w", model.ErrPostTagAlreadyExisted)
		}
		return fmt.Errorf("StorePostTag: cannot store postTag: %w", err)
	}
	return nil
}

func (r *postTagRepository) DeleteByPostID(postID string) error {
	_, err := r.db.Exec("DELETE FROM posts_tags WHERE post_id = ?", postID)
	if err != nil {
		return fmt.Errorf("DeleteByPostID: cannot delete postTag: %w", err)
	}
	return nil
}
