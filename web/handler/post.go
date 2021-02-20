package handler

import (
	"net/http"

	"github.com/Le0tk0k/blog-server/log"
	"github.com/Le0tk0k/blog-server/service"

	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	postService    service.PostService
	postTagService service.PostTagService
}

// NewPostHandler はPostHandlerを返す
func NewPostHandler(postService service.PostService, postTagService service.PostTagService) PostHandler {
	return PostHandler{
		postService:    postService,
		postTagService: postTagService,
	}
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

// GetPost は GET /posts/:slug に対するhandler
func (p *PostHandler) GetPost(c echo.Context) error {
	logger := log.New()

	slug := c.Param("slug")
	post, err := p.postService.GetPost(slug)
	if err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, postToPostWithTagsJSON(post))
}

// GetPosts は GET /posts に対するhandler
func (p *PostHandler) GetPosts(c echo.Context) error {
	logger := log.New()

	conditions := make([]string, 0)

	if draft := c.QueryParam("draft"); draft != "" {
		conditions = append(conditions, "draft = "+draft)
	}

	if tag := c.QueryParam("tag"); tag != "" {
		conditions = append(conditions, "tags.name = '"+tag+"'")
	}

	posts, err := p.postService.GetPosts(conditions)
	if err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	postsJSON := make([]*postWithTagsJSON, len(posts))
	for i, post := range posts {
		postsJSON[i] = postToPostWithTagsJSON(post)
	}

	return c.JSON(http.StatusOK, postsJSON)
}

// UpdatePost は PUT /posts/:id に対するhandler
func (p *PostHandler) UpdatePost(c echo.Context) error {
	logger := log.New()

	req := new(postWithTagsJSON)
	if err := c.Bind(req); err != nil {
		logger.Errorj(map[string]interface{}{"message": "failed to bind", "error": err.Error()})
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	post := jsonToPostWithTags(req)

	err := p.postService.UpdatePost(post)
	if err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = p.postTagService.LinkPostTag(post)
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
