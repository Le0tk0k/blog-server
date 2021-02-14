package web

import (
	"github.com/Le0tk0k/blog-server/service"
	"github.com/Le0tk0k/blog-server/web/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewServer はhandlerが登録されたechoの構造体を返す
func NewServer(postService service.PostService, tagService service.TagService) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	postHandler := handler.NewPostHandler(postService)
	tagHandler := handler.NewTagHandler(tagService)

	v1 := e.Group("/v1")

	posts := v1.Group("/posts")
	posts.GET("", postHandler.GetPosts)
	posts.GET("/:id", postHandler.GetPost)
	posts.POST("", postHandler.CreatePost)
	posts.PUT("/:id", postHandler.UpdatePost)
	posts.DELETE("/:id", postHandler.DeletePost)

	tags := v1.Group("/tags")
	tags.POST("", tagHandler.CreateTag)
	tags.GET("/:id", tagHandler.GetTag)

	return e
}
