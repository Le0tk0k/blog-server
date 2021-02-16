package main

import (
	"github.com/Le0tk0k/blog-server/config"
	"github.com/Le0tk0k/blog-server/log"
	"github.com/Le0tk0k/blog-server/repository"
	"github.com/Le0tk0k/blog-server/service"
	"github.com/Le0tk0k/blog-server/web"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	logger := log.New()

	db, err := repository.NewDB()
	if err != nil {
		logger.Fatal(err)
	}

	m, err := migrate.New("file://db/migrations", "mysql://"+config.DSN())
	if err != nil {
		logger.Fatal(err)
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal(err)
	}

	postRepository := repository.NewPostRepository(db)
	postService := service.NewPostService(postRepository)
	tagRepository := repository.NewTagRepository(db)
	tagService := service.NewTagService(tagRepository)
	postTagRepository := repository.NewPostTagRepository(db)
	postTagService := service.NewPostTagService(postTagRepository)
	e := web.NewServer(postService, tagService, postTagService)

	e.Logger.Fatal(e.Start(":1323"))
}
