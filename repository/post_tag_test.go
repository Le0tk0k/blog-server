package repository

import (
	"errors"
	"testing"

	"github.com/Le0tk0k/blog-server/model"
)

func TestPostTagRepository_StorePostTag(t *testing.T) {
	existPost := &postDTO{
		ID:      "post_id_post_tag_test",
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
	existTag := &tagDTO{
		ID:   "tag_id_post_tag_test",
		Name: "tag1",
	}
	_, err = db.Exec("INSERT INTO tags VALUES (?, ?)", existTag.ID, existTag.Name)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		id      string
		postTag *model.PostTag
		wantErr error
	}{
		{
			name: "postとtagを関連づけられる",
			id:   "post_id",
			postTag: &model.PostTag{
				ID:     "post_tag_id",
				PostID: "post_id_post_tag_test",
				TagID:  "tag_id_post_tag_test",
			},
			wantErr: nil,
		},
		{
			name: "既に存在するIDの場合ErrPostTagAlreadyExistedエラーを返す",
			id:   "not_found",
			postTag: &model.PostTag{
				ID:     "post_tag_id",
				PostID: "post_id_post_tag_test",
				TagID:  "tag_id_post_tag_test",
			},
			wantErr: model.ErrPostTagAlreadyExisted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &postTagRepository{db: db}
			if err := r.StorePostTag(tt.postTag); !errors.Is(err, tt.wantErr) {
				t.Errorf("StorePostTag() error = %v, wantErr %v", err, tt.wantErr)
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

func TestPostTagRepository_DeleteByPostID(t *testing.T) {
	existPost := &postDTO{
		ID:      "post_id_post_tag_test",
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
	existTag := &tagDTO{
		ID:   "tag_id_post_tag_test",
		Name: "tag1",
	}
	_, err = db.Exec("INSERT INTO tags VALUES (?, ?)", existTag.ID, existTag.Name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec("INSERT INTO posts_tags VALUES (?, ?, ?)", "1", existPost.ID, existTag.ID); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name:    "存在するpost_tagを正常に削除できる",
			id:      "post_id_post_tag_test",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &postTagRepository{db: db}
			err := r.DeleteByPostID(tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DeleteByPostID()  error = %v, wantErr %v", err, tt.wantErr)
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
