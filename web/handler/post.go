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

// CreatePost は POST /posts に対するhandler
func (p *PostHandler) CreatePost(c echo.Context) error {
	logger := log.New()

	post, err := p.postService.CreatePost()
	if err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, postToJSON(post))
}

// GetPost は GET /posts/:id に対するhandler
func (p *PostHandler) GetPost(c echo.Context) error {
	logger := log.New()

	id := c.Param("id")
	post, tags, err := p.postService.GetPost(id)
	if err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	tagsJSON := make([]*tagJSON, len(tags))
	for i, tag := range tags {
		tagsJSON[i] = tagToJSON(tag)
	}
	res := postWithTagsToPostJSONWithTags(post, tags)
	return c.JSON(http.StatusOK, res)
}

// GetPosts は GET /posts に対するhandler
func (p *PostHandler) GetPosts(c echo.Context) error {
	logger := log.New()

	posts, err := p.postService.GetPosts()
	if err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	postsJSON := make([]*postJSON, len(posts))
	for i, post := range posts {
		postsJSON[i] = postToJSON(post)
	}

	return c.JSON(http.StatusOK, postsJSON)
}

// UpdatePost は PUT /posts/:id に対するhandler
func (p *PostHandler) UpdatePost(c echo.Context) error {
	logger := log.New()

	req := new(postJSON)
	if err := c.Bind(req); err != nil {
		logger.Errorj(map[string]interface{}{"message": "failed to bind", "error": err.Error()})
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	post := jsonToPOST(req)
	err := p.postService.UpdatePost(post)
	if err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, "successfully updated")
}

// DeletePost は DELETE /posts/:id に対するhandler
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
