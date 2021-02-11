package handler

import (
	"net/http"
	"time"

	"github.com/Le0tk0k/blog-server/model"

	"github.com/Le0tk0k/blog-server/log"
	"github.com/Le0tk0k/blog-server/service"

	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	postService service.PostService
}

type postResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Slug        string     `json:"slug"`
	Draft       bool       `json:"draft"`
	PublishedAt *time.Time `json:"published_at"`
}

// NewPostHandler はPostHandlerを返す
func NewPostHandler(postService service.PostService) PostHandler {
	return PostHandler{postService: postService}
}

// CreatePost は POST /post に対するhandler
func (p *PostHandler) CreatePost(c echo.Context) error {
	logger := log.New()

	post, err := p.postService.CreatePost()
	if err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, postToResponse(post))
}

// GetPosts は GET /posts に対するhandler
func (p *PostHandler) GetPosts(c echo.Context) error {
	logger := log.New()

	posts, err := p.postService.GetPosts()
	if err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	var postsRes []*postResponse
	for _, post := range posts {
		postRes := postToResponse(post)
		postsRes = append(postsRes, postRes)
	}

	return c.JSON(http.StatusOK, postsRes)
}

func postToResponse(post *model.Post) *postResponse {
	return &postResponse{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		Slug:        post.Slug,
		Draft:       post.Draft,
		PublishedAt: post.PublishedAt,
	}
}
