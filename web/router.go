package web

import (
	"os"

	"github.com/Le0tk0k/blog-server/config"
	"github.com/Le0tk0k/blog-server/service"
	"github.com/Le0tk0k/blog-server/web/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// NewServer はhandlerが登録されたechoの構造体を返す
func NewServer(postService service.PostService, tagService service.TagService, postTagService service.PostTagService, userService service.UserService) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{config.CORSAllowOrigin()},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	cnf := middleware.JWTConfig{
		SigningKey:  []byte(os.Getenv("AUTH_KEY")),
		TokenLookup: "cookie:jwt",
	}

	postHandler := handler.NewPostHandler(postService, postTagService)
	tagHandler := handler.NewTagHandler(tagService)
	authMiddleware := NewAuthMiddleware(userService)

	v1 := e.Group("/v1")

	v1.POST("/login", authMiddleware.Login)

	posts := v1.Group("/posts")
	posts.GET("", postHandler.GetPosts)
	posts.GET("/:id", postHandler.GetPost)
	posts.PUT("/:id", postHandler.UpdatePost, middleware.JWTWithConfig(cnf))
	posts.DELETE("/:id", postHandler.DeletePost, middleware.JWTWithConfig(cnf))
	posts.POST("", postHandler.CreatePost, middleware.JWTWithConfig(cnf))

	tags := v1.Group("/tags")
	tags.POST("", tagHandler.CreateTag, middleware.JWTWithConfig(cnf))
	tags.GET("/:id", tagHandler.GetTag)
	tags.GET("", tagHandler.GetTags)
	tags.DELETE("/:id", tagHandler.DeleteTag, middleware.JWTWithConfig(cnf))

	return e
}
