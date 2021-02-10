package service

import (
	"fmt"

	"github.com/Le0tk0k/blog-server/model"
	"github.com/Le0tk0k/blog-server/repository"
)

type PostService interface {
	CreatePost() (*model.Post, error)
}

type postService struct {
	postRepository repository.PostRepository
}

// NrePostService はPostServiceを返す
func NewPostService(postRepository repository.PostRepository) PostService {
	return &postService{postRepository: postRepository}
}

// CreatePost は新しい記事を生成、保存する
func (p *postService) CreatePost() (*model.Post, error) {
	post := model.NewPost()
	if err := p.postRepository.StorePost(post); err != nil {
		return nil, fmt.Errorf("CreatePost: cannot create post: %w", err)
	}
	return post, nil
}
