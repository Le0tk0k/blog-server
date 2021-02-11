package handler

import (
	"net/http"

	"github.com/Le0tk0k/blog-server/log"
	"github.com/Le0tk0k/blog-server/service"

	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	postService service.PostService
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
	return c.JSON(http.StatusCreated, post)
}
