package service

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/Le0tk0k/blog-server/model"
	"github.com/Le0tk0k/blog-server/repository/mock_repository"

	"github.com/golang/mock/gomock"
)

func TestService_CreatePost(t *testing.T) {
	tests := []struct {
		name                  string
		prepareMockPostRepoFn func(mock *mock_repository.MockPostRepository)
		wantErr               error
	}{
		{
			name: "新規記事を生成、保存し、その記事を返す",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPostRepository) {
				mock.EXPECT().StorePost(gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockPostRepository(ctrl)
			tt.prepareMockPostRepoFn(mr)
			ps := &postService{
				postRepository: mr,
			}

			_, err := ps.CreatePost()
			if err != tt.wantErr {
				t.Errorf("CreatePost() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostService_GetPost(t *testing.T) {
	now := time.Now()
	existsPost := &model.Post{
		ID:          "post_id_1",
		Title:       "post_title_1",
		Content:     "pot_content_1",
		Slug:        "post-slug-1",
		Draft:       true,
		PublishedAt: &now,
	}

	tests := []struct {
		name                  string
		id                    string
		prepareMockPostRepoFn func(mock *mock_repository.MockPostRepository)
		want                  *model.Post
		wantErr               bool
	}{
		{
			name: "記事を返す",
			id:   "post_id_1",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPostRepository) {
				mock.EXPECT().FindPostByID("post_id_1").Return(existsPost, nil)
			},
			want: &model.Post{
				ID:          "post_id_1",
				Title:       "post_title_1",
				Content:     "pot_content_1",
				Slug:        "post-slug-1",
				Draft:       true,
				PublishedAt: &now,
			},
			wantErr: false,
		},
		{
			name: "記事の取得に失敗したときエラーを返す",
			id:   "not_found",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPostRepository) {
				mock.EXPECT().FindPostByID("not_found").Return(nil, errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockPostRepository(ctrl)
			tt.prepareMockPostRepoFn(mr)
			ps := &postService{
				postRepository: mr,
			}

			got, err := ps.GetPost(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPost() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPost() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestPostService_GetPosts(t *testing.T) {
	now := time.Now()
	existsPosts := []*model.Post{{
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

	tests := []struct {
		name                  string
		prepareMockPostRepoFn func(mock *mock_repository.MockPostRepository)
		want                  []*model.Post
		wantErr               bool
	}{
		{
			name: "全記事を返す",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPostRepository) {
				mock.EXPECT().FindAllPosts().Return(existsPosts, nil)
			},
			want: []*model.Post{
				{
					ID:          "post_id_1",
					Title:       "post_title_1",
					Content:     "pot_content_1",
					Slug:        "post-slug-1",
					Draft:       true,
					PublishedAt: &now,
				},
				{

					ID:          "post_id_2",
					Title:       "post_title_2",
					Content:     "post_content_2",
					Slug:        "post-slug-2",
					Draft:       false,
					PublishedAt: &now,
				},
			},
			wantErr: false,
		},
		{
			name: "記事の取得に失敗したときエラーを返す",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPostRepository) {
				mock.EXPECT().FindAllPosts().Return(nil, errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockPostRepository(ctrl)
			tt.prepareMockPostRepoFn(mr)
			ps := &postService{
				postRepository: mr,
			}

			got, err := ps.GetPosts()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPosts() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPosts() got = %v, want = %v", got, tt.want)
			}
		})
	}
}
