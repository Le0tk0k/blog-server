package handler

import (
	"net/http"

	"github.com/Le0tk0k/blog-server/log"
	"github.com/Le0tk0k/blog-server/model"
	"github.com/Le0tk0k/blog-server/service"

	"github.com/labstack/echo/v4"
)

type TagHandler struct {
	tagService service.TagService
}

type tagJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NewTagHandler はTagHandlerを返す
func NewTagHandler(tagService service.TagService) TagHandler {
	return TagHandler{tagService: tagService}
}

// CreateTag は POST /tags に対するhandler
func (t *TagHandler) CreateTag(c echo.Context) error {
	logger := log.New()

	req := new(tagJSON)
	if err := c.Bind(req); err != nil {
		logger.Errorj(map[string]interface{}{"message": "failed to bind", "error": err.Error()})
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	tag := jsonToTag(req)
	err := t.tagService.CreateTag(tag.Name)
	if err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, "successfully created")
}

func tagToJSON(tag *model.Tag) *tagJSON {
	return &tagJSON{
		ID:   tag.ID,
		Name: tag.Name,
	}
}

func jsonToTag(json *tagJSON) *model.Tag {
	return &model.Tag{
		ID:   json.ID,
		Name: json.Name,
	}
}
