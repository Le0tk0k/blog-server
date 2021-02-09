package main

import (
	"log"
	"net/http"

	"github.com/Le0tk0k/blog-server/config"

	"github.com/Le0tk0k/blog-server/repository"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
)

func main() {
	_, err := repository.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.New("file://db/migrations", "mysql://"+config.DSN())
	if err != nil {
		log.Fatal(err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
