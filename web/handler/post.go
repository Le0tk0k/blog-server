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

// GetPost は GET /post/:id に対するhandler
func (p *PostHandler) GetPost(c echo.Context) error {
	logger := log.New()

	id := c.Param("id")
	post, err := p.postService.GetPost(id)
	if err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, postToResponse(post))
}

// GetPosts は GET /posts に対するhandler
func (p *PostHandler) GetPosts(c echo.Context) error {
	logger := log.New()

	posts, err := p.postService.GetPosts()
	if err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	postsRes := make([]*postResponse, len(posts))
	for i, post := range posts {
		postsRes[i] = postToResponse(post)
	}

	return c.JSON(http.StatusOK, postsRes)
}

// DeletePost は DELETE /post/:id に対するhandler
func (p *PostHandler) DeletePost(c echo.Context) error {
	logger := log.New()

	id := c.Param("id")
	err := p.postService.DeletePost(id)
	if err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, "successfully deleted")
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
