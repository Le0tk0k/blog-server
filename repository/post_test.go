package repository

import (
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/Le0tk0k/blog-server/config"
	"github.com/Le0tk0k/blog-server/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
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
			name: "新規記事を正しく保存できる",
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
	_, err := db.Exec("DELETE FROM posts")
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostRepository_FindPostByID(t *testing.T) {

	existPost := &postDTO{
		ID:      "post_id_1",
		Title:   "post_title_1",
		Content: "pot_content_1",
		Slug:    "post-slug-1",
		Draft:   true,
		// TODO mysqlとgoで時間の表示形式が違う
		PublishedAt: nil,
	}
	_, err := db.Exec("INSERT INTO posts VALUES (?, ?, ?, ?, ?, ?)", existPost.ID, existPost.Title, existPost.Content, existPost.Slug, existPost.Draft, existPost.PublishedAt)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		id      string
		want    *model.Post
		wantErr error
	}{
		{
			name: "存在する記事を正常に取得できる",
			id:   "post_id_1",
			want: &model.Post{
				ID:          "post_id_1",
				Title:       "post_title_1",
				Content:     "pot_content_1",
				Slug:        "post-slug-1",
				Draft:       true,
				PublishedAt: nil,
			},
			wantErr: nil,
		},
		{
			name:    "存在しないIDの場合ErrPostNotFoundを返す",
			id:      "not_found",
			want:    nil,
			wantErr: model.ErrPostNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &postRepository{db: db}
			got, err := r.FindPostByID(tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("FindPostByID()  error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindPostByID() got = %v, want = %v", got, tt.want)
			}
		})
	}

	_, err = db.Exec("DELETE FROM posts")
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostRepository_FindAllPosts(t *testing.T) {
	now := time.Now()

	existsPosts := []*postDTO{{
		ID:          "post_id_1",
		Title:       "post_title_1",
		Content:     "pot_content_1",
		Slug:        "post-slug-1",
		Draft:       true,
		PublishedAt: &now,
	}, {
		ID:          "post_id_2",
		Title:       "post_title_2",
		Content:     "post_content_2",
		Slug:        "post-slug-2",
		Draft:       false,
		PublishedAt: &now,
	}}

	for _, post := range existsPosts {
		_, err := db.Exec("INSERT INTO posts VALUES (?, ?, ?, ?, ?, ?)", post.ID, post.Title, post.Content, post.Slug, post.Draft, post.PublishedAt)
		if err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "記事を全件取得できる",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &postRepository{db: db}
			got, err := r.FindAllPosts()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("FindAll()  error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(got) != len(existsPosts) {
				t.Errorf("FindAll() does not fetch all posts got = %v, want = %v", got, existsPosts)
			}
		})
	}

	_, err := db.Exec("DELETE FROM posts")
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostRepository_UpdatePost(t *testing.T) {
	now := time.Now()
	existPost := &postDTO{
		ID:          "post_id_1",
		Title:       "post_title_1",
		Content:     "pot_content_1",
		Slug:        "post-slug-1",
		Draft:       true,
		PublishedAt: &now,
	}
	_, err := db.Exec("INSERT INTO posts VALUES (?, ?, ?, ?, ?, ?)", existPost.ID, existPost.Title, existPost.Content, existPost.Slug, existPost.Draft, existPost.PublishedAt)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		post    *model.Post
		wantErr error
	}{
		{
			name: "存在する記事を正常に更新できる",
			post: &model.Post{
				ID:          "post_id_1_update",
				Title:       "post_title_1_update",
				Content:     "pot_content_1_update",
				Slug:        "post-slug-1-update",
				Draft:       false,
				PublishedAt: &now,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &postRepository{db: db}
			err := r.UpdatePost(tt.post)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UpdatePost()  error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	_, err = db.Exec("DELETE FROM posts")
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostRepository_DeletePostByID(t *testing.T) {
	now := time.Now()
	existPost := &postDTO{
		ID:          "post_id_1",
		Title:       "post_title_1",
		Content:     "pot_content_1",
		Slug:        "post-slug-1",
		Draft:       true,
		PublishedAt: &now,
	}
	_, err := db.Exec("INSERT INTO posts VALUES (?, ?, ?, ?, ?, ?)", existPost.ID, existPost.Title, existPost.Content, existPost.Slug, existPost.Draft, existPost.PublishedAt)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name:    "存在する記事を正常に削除できる",
			id:      "post_id_1",
			wantErr: nil,
		},
		{
			name:    "存在しないIDの場合ErrPostNotFoundを返す",
			id:      "not_found",
			wantErr: model.ErrPostNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &postRepository{db: db}
			err := r.DeletePostByID(tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DeletePostByID()  error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
