package repository

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Le0tk0k/blog-server/model"
)

func TestTagRepository_StoreTag(t *testing.T) {
	tests := []struct {
		name    string
		tag     *model.Tag
		wantErr error
	}{
		{
			name: "新規タグを正しく保存できる",
			tag: &model.Tag{
				ID:   "new_id",
				Name: "new_tag_1",
			},
			wantErr: nil,
		},
		{
			name: "既に存在するIDの場合ErrTagAlreadyExistedエラーを返す",
			tag: &model.Tag{
				ID:   "new_id",
				Name: "new_tag_2",
			},
			wantErr: model.ErrTagAlreadyExisted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &tagRepository{db: db}
			if err := r.StoreTag(tt.tag); !errors.Is(err, tt.wantErr) {
				t.Errorf("StoreTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	_, err := db.Exec("DELETE FROM tags")
	if err != nil {
		t.Fatal(err)
	}
}

func TestTagRepository_FindTagByID(t *testing.T) {
	existTag := &tagDTO{
		ID:   "tag_id",
		Name: "tag1",
	}
	_, err := db.Exec("INSERT INTO tags VALUES (?, ?)", existTag.ID, existTag.Name)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		id      string
		want    *model.Tag
		wantErr error
	}{
		{
			name: "存在するタグを正常に取得できる",
			id:   "tag_id",
			want: &model.Tag{
				ID:   "tag_id",
				Name: "tag1",
			},
			wantErr: nil,
		},
		{
			name:    "存在しないIDの場合ErrTagNotFoundを返す",
			id:      "not_found",
			want:    nil,
			wantErr: model.ErrTagNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &tagRepository{db: db}
			got, err := r.FindTagByID(tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("FindTagByID()  error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindTagByID() got = %v, want = %v", got, tt.want)
			}
		})
	}

	_, err = db.Exec("DELETE FROM tags")
	if err != nil {
		t.Fatal(err)
	}
}

func TestTagRepository_FindAllTags(t *testing.T) {
	existTags := []*tagDTO{{
		ID:   "tag_id_1",
		Name: "tag1",
	}, {
		ID:   "tag_id_2",
		Name: "tag2",
	}}

	for _, tag := range existTags {
		_, err := db.Exec("INSERT INTO tags VALUES (?, ?)", tag.ID, tag.Name)
		if err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name    string
		wantErr error
	}{
		{
			name:    "タグを全件取得できる",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &tagRepository{db: db}
			got, err := r.FindAllTags()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("FindAllTags()  error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(got) != len(existTags) {
				t.Errorf("FindAllTags() does not fetch all tags got = %v, want = %v", got, existTags)
			}
		})
	}

	_, err := db.Exec("DELETE FROM posts")
	if err != nil {
		t.Fatal(err)
	}
}
