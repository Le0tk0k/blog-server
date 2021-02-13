package repository

import (
	"fmt"

	"github.com/Le0tk0k/blog-server/model"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type TagRepository interface {
	StoreTag(tag *model.Tag) error
}

type tagRepository struct {
	db *sqlx.DB
}

type tagDTO struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func NewTagRepository(db *sqlx.DB) TagRepository {
	return &tagRepository{db: db}
}

// StoreTag はタグを新規保存する
func (t *tagRepository) StoreTag(tag *model.Tag) error {
	dto := tagToDTO(tag)
	_, err := t.db.Exec("INSERT INTO tags (id, name) VALUES (?, ?)", dto.ID, dto.Name)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return fmt.Errorf("StorePost: cannot store post: %w", model.ErrPostAlreadyExisted)
		}
		return fmt.Errorf("StorePost: cannot store post: %w", err)
	}
	return nil
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
