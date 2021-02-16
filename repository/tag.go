package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Le0tk0k/blog-server/model"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type TagRepository interface {
	StoreTag(tag *model.Tag) error
	FindTagByID(id string) (*model.Tag, error)
	FindAllTags() ([]*model.Tag, error)
	DeleteTagByID(id string) error
}

type tagRepository struct {
	db *sqlx.DB
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
			return fmt.Errorf("StoreTag: cannot store tag: %w", model.ErrTagAlreadyExisted)
		}
		return fmt.Errorf("StoreTag: cannot store tag: %w", err)
	}
	return nil
}

// FindTagByID はidを持つ記事を取得する
func (t *tagRepository) FindTagByID(id string) (*model.Tag, error) {
	var dto tagDTO
	if err := t.db.Get(&dto, "SELECT * FROM tags WHERE id = ?", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("FindTagByID: cannot find tag: %w", model.ErrTagNotFound)
		}
		return nil, fmt.Errorf("FindTagByID: cannot find tag: %w", err)
	}
	return dtoToTag(&dto), nil
}

// FindAllTags は全記事を取得する
func (t *tagRepository) FindAllTags() ([]*model.Tag, error) {
	var dtos []*tagDTO
	if err := t.db.Select(&dtos, "SELECT * FROM tags"); err != nil {
		return nil, fmt.Errorf("FindAllTags: cannot find tag: %w", err)
	}

	tags := make([]*model.Tag, len(dtos))
	for i, dto := range dtos {
		tags[i] = dtoToTag(dto)
	}

	return tags, nil
}

// DeleteTagByID はidを持つ記事を削除する
func (t *tagRepository) DeleteTagByID(id string) error {
	result, err := t.db.Exec("DELETE FROM tags WHERE id = ?", id)
	if rowsAffected, err := result.RowsAffected(); rowsAffected == 0 && err == nil {
		return fmt.Errorf("DeleteTabByID: cannot find tag: %w", model.ErrTagNotFound)
	}
	if err != nil {
		return fmt.Errorf("DeleteTabByID: cannot delete tag: %w", err)
	}
	return nil
}
