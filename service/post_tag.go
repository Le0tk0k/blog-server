package service

import (
	"fmt"

	"github.com/Le0tk0k/blog-server/model"
	"github.com/Le0tk0k/blog-server/repository"
)

type PostTagService interface {
	LinkPostTag(post *model.Post) error
}

type postTagService struct {
	postTagRepository repository.PostTagRepository
}

// NewPostService はPostTagServiceを返す
func NewPostTagService(postTagRepository repository.PostTagRepository) PostTagService {
	return &postTagService{
		postTagRepository: postTagRepository,
	}
}

// LinkPostTag は記事とタグを関連付ける
func (s *postTagService) LinkPostTag(post *model.Post) error {
	err := s.postTagRepository.DeleteByPostID(post.ID)
	if err != nil {
		return fmt.Errorf("LinkPostTag: cannot link postTag: %w", err)
	}

	for _, tag := range post.Tags {
		postTag := model.NewPostTag(post.ID, tag.ID)
		err = s.postTagRepository.StorePostTag(postTag)
		if err != nil {
			return fmt.Errorf("LinkPostTag: cannot link postTag: %w", err)
		}
	}

	return nil
}
