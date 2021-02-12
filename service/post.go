package service

import (
	"fmt"

	"github.com/Le0tk0k/blog-server/model"
	"github.com/Le0tk0k/blog-server/repository"
)

type PostService interface {
	CreatePost() (*model.Post, error)
	GetPost(id string) (*model.Post, error)
	GetPosts() ([]*model.Post, error)
	DeletePost(id string) error
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

// GetPost はidを持つ記事を取得する
func (p *postService) GetPost(id string) (*model.Post, error) {
	post, err := p.postRepository.FindPostByID(id)
	if err != nil {
		return nil, fmt.Errorf("GetPost: cannot get post: %w", err)
	}
	return post, nil
}

// GetPosts は全記事を取得する
func (p *postService) GetPosts() ([]*model.Post, error) {
	posts, err := p.postRepository.FindAllPosts()
	if err != nil {
		return nil, fmt.Errorf("GetPosts: cannot get posts: %w", err)
	}
	return posts, nil
}

// DeletePost はidを持つ記事を削除する
func (p *postService) DeletePost(id string) error {
	err := p.postRepository.DeletePostByID(id)
	if err != nil {
		return fmt.Errorf("DeletePost: cannot delete post: %w", err)
	}
	return nil
}
