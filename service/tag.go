package service

import (
	"fmt"

	"github.com/Le0tk0k/blog-server/model"
	"github.com/Le0tk0k/blog-server/repository"
)

type TagService interface {
	CreateTag(name string) error
	GetTag(id string) (*model.Tag, error)
}

type tagService struct {
	tagRepository repository.TagRepository
}

// NewTagService はTagServiceを返す
func NewTagService(tagRepository repository.TagRepository) TagService {
	return &tagService{tagRepository: tagRepository}
}

// CreateTag は新しい記事を生成、保存する
func (t *tagService) CreateTag(name string) error {
	tag := model.NewTag(name)
	if err := t.tagRepository.StoreTag(tag); err != nil {
		return fmt.Errorf("CreateTag: cannot create tag: %w", err)
	}
	return nil
}

// GetTag はidを持つタグを取得する
func (t *tagService) GetTag(id string) (*model.Tag, error) {
	tag, err := t.tagRepository.FindTagByID(id)
	if err != nil {
		return nil, fmt.Errorf("GetTag: cannot get tag: %w", err)
	}
	return tag, nil
}
