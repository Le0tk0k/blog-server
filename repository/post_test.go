package repository

import (
	"errors"
	"os"
	"testing"

	"github.com/Le0tk0k/blog-server/config"
	"github.com/golang-migrate/migrate/v4"

	"github.com/Le0tk0k/blog-server/model"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	mig, err := migrate.New("file://../db/migrations", "mysql://"+config.DSN())
	if err != nil {
		panic(err)
	}
	if err = mig.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	db = NewTestDB()
	code := m.Run()
	os.Exit(code)
}

func TestPostRepository_StorePost(t *testing.T) {
	tests := []struct {
		name    string
		post    *model.Post
		wantErr error
	}{
		{
			name: "新規postを正しく保存できる",
			post: &model.Post{
				ID:      "new_id",
				Title:   "new_post_1",
				Content: "new_content_1",
				Slug:    "new_slug_1",
				Draft:   true,
			},
			wantErr: nil,
		},
		{
			name: "既に存在するIDの場合ErrPostAlreadyExistedエラーを返す",
			post: &model.Post{
				ID:      "new_id",
				Title:   "new_post_2",
				Content: "new_content_2",
				Slug:    "new_slug_2",
				Draft:   true,
			},
			wantErr: model.ErrPostAlreadyExisted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &postRepository{db: db}
			if err := r.StorePost(tt.post); !errors.Is(err, tt.wantErr) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	db.Exec("DELETE FROM posts")
}
