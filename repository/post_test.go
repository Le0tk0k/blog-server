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
				t.Errorf("StorePost() error = %v, wantErr %v", err, tt.wantErr)
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
	existTags := []*tagDTO{
		{
			ID:   "tag_id_1_post_test",
			Name: "tag1",
		},
		{
			ID:   "tag_id_2_post_test",
			Name: "tag2",
		},
	}
	tmp := []string{"1", "2"}
	for i, tag := range existTags {
		_, err = db.Exec("INSERT INTO tags VALUES (?, ?)", tag.ID, tag.Name)
		if err != nil {
			t.Fatal(err)
		}
		// TODO posts_tags 実装後リファクタリング
		_, err = db.Exec("INSERT INTO posts_tags VALUES (?, ?, ?)", tmp[i], existPost.ID, tag.ID)
		if err != nil {
			t.Fatal(err)
		}
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
				Tags: []*model.Tag{
					{
						ID:   "tag_id_1_post_test",
						Name: "tag1",
					},
					{
						ID:   "tag_id_2_post_test",
						Name: "tag2",
					},
				},
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
	_, err = db.Exec("DELETE FROM tags")
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostRepository_FindPostBySlug(t *testing.T) {

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
	existTags := []*tagDTO{
		{
			ID:   "tag_id_1_post_test",
			Name: "tag1",
		},
		{
			ID:   "tag_id_2_post_test",
			Name: "tag2",
		},
	}
	tmp := []string{"1", "2"}
	for i, tag := range existTags {
		_, err = db.Exec("INSERT INTO tags VALUES (?, ?)", tag.ID, tag.Name)
		if err != nil {
			t.Fatal(err)
		}
		// TODO posts_tags 実装後リファクタリング
		_, err = db.Exec("INSERT INTO posts_tags VALUES (?, ?, ?)", tmp[i], existPost.ID, tag.ID)
		if err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name    string
		slug    string
		want    *model.Post
		wantErr error
	}{
		{
			name: "存在する記事を正常に取得できる",
			slug: "post-slug-1",
			want: &model.Post{
				ID:          "post_id_1",
				Title:       "post_title_1",
				Content:     "pot_content_1",
				Slug:        "post-slug-1",
				Draft:       true,
				PublishedAt: nil,
				Tags: []*model.Tag{
					{
						ID:   "tag_id_1_post_test",
						Name: "tag1",
					},
					{
						ID:   "tag_id_2_post_test",
						Name: "tag2",
					},
				},
			},
			wantErr: nil,
		},
		{
			name:    "存在しないIDの場合ErrPostNotFoundを返す",
			slug:    "not_found",
			want:    nil,
			wantErr: model.ErrPostNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &postRepository{db: db}
			got, err := r.FindPostBySlug(tt.slug)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("FindPostBySlug()  error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindPostBySlug() got = %v, want = %v", got, tt.want)
			}
		})
	}

	_, err = db.Exec("DELETE FROM posts")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM tags")
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostRepository_FindAllPosts(t *testing.T) {
	now := time.Now()

	existsPosts := []*postDTO{{
		ID:          "post_id_1_post_test",
		Title:       "post_title_1",
		Content:     "pot_content_1",
		Slug:        "post-slug-1",
		Draft:       true,
		PublishedAt: &now,
	}, {
		ID:          "post_id_2_post_test",
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
	existTags := []*tagDTO{
		{
			ID:   "tag_id_1",
			Name: "tag1",
		},
		{
			ID:   "tag_id_2",
			Name: "tag2",
		},
	}
	for _, tag := range existTags {
		_, err := db.Exec("INSERT INTO tags VALUES (?, ?)", tag.ID, tag.Name)
		if err != nil {
			t.Fatal(err)
		}
	}
	// TODO posts_tags 実装後リファクタリング
	_, err := db.Exec("INSERT INTO posts_tags VALUES (?, ?, ?)", "1", existsPosts[0].ID, existTags[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO posts_tags VALUES (?, ?, ?)", "2", existsPosts[1].ID, existTags[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO posts_tags VALUES (?, ?, ?)", "3", existsPosts[1].ID, existTags[1].ID)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		postCount  int
		conditions []string
		wantErr    error
	}{
		{
			name:       "記事を全件取得できる",
			postCount:  2,
			conditions: []string{},
			wantErr:    nil,
		},
		{
			name:       "クエリパラメータにdraftを設定",
			postCount:  1,
			conditions: []string{"draft = true"},
			wantErr:    nil,
		},
		{
			name:       "クエリパラメータにtagを設定",
			postCount:  1,
			conditions: []string{"tags.name = 'tag2'"},
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &postRepository{db: db}
			got, err := r.FindAllPosts(tt.conditions)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("FindAll()  error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(got) != tt.postCount {
				t.Errorf("FindAll() does not fetch all posts got = %v, want = %v", got, existsPosts)
			}
		})
	}

	_, err = db.Exec("DELETE FROM posts")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("DELETE FROM tags")
	if err != nil {
		t.Fatal(err)
	}
}

func TestPostRepository_UpdatePost(t *testing.T) {
	now := time.Now()
	existPost := &postDTO{
		ID:          "post_id_1_post_test",
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
		ID:          "post_id_post_test",
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
			id:      "post_id_post_test",
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
