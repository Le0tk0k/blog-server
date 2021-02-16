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
		Tags: []*model.Tag{
			{
				ID:   "tag_id_1",
				Name: "tag1",
			},
			{
				ID:   "tag_id_2",
				Name: "tag2",
			},
		},
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
				Tags: []*model.Tag{
					{
						ID:   "tag_id_1",
						Name: "tag1",
					},
					{
						ID:   "tag_id_2",
						Name: "tag2",
					},
				},
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
		Tags: []*model.Tag{
			{
				ID:   "tag_id_1",
				Name: "tag1",
			},
			{
				ID:   "tag_id_2",
				Name: "tag2",
			},
		},
	}, {
		ID:          "post_id_2",
		Title:       "post_title_2",
		Content:     "post_content_2",
		Slug:        "post-slug-2",
		Draft:       false,
		PublishedAt: &now,
		Tags: []*model.Tag{
			{
				ID:   "tag_id_1",
				Name: "tag1",
			},
		},
	}}

	tests := []struct {
		name                  string
		prepareMockPostRepoFn func(mock *mock_repository.MockPostRepository)
		conditions            []string
		want                  []*model.Post
		wantErr               bool
	}{
		{
			name: "全記事を返す",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPostRepository) {
				mock.EXPECT().FindAllPosts(gomock.Any()).Return(existsPosts, nil)
			},
			conditions: []string{},
			want: []*model.Post{
				{
					ID:          "post_id_1",
					Title:       "post_title_1",
					Content:     "pot_content_1",
					Slug:        "post-slug-1",
					Draft:       true,
					PublishedAt: &now,
					Tags: []*model.Tag{
						{
							ID:   "tag_id_1",
							Name: "tag1",
						},
						{
							ID:   "tag_id_2",
							Name: "tag2",
						},
					},
				},
				{

					ID:          "post_id_2",
					Title:       "post_title_2",
					Content:     "post_content_2",
					Slug:        "post-slug-2",
					Draft:       false,
					PublishedAt: &now,
					Tags: []*model.Tag{
						{
							ID:   "tag_id_1",
							Name: "tag1",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "公開記事を返す",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPostRepository) {
				mock.EXPECT().FindAllPosts(gomock.Any()).Return([]*model.Post{
					{
						ID:          "post_id_2",
						Title:       "post_title_2",
						Content:     "post_content_2",
						Slug:        "post-slug-2",
						Draft:       false,
						PublishedAt: &now,
						Tags: []*model.Tag{
							{
								ID:   "tag_id_1",
								Name: "tag1",
							},
						},
					},
				}, nil)
			},
			conditions: []string{"draft = false"},
			want: []*model.Post{
				{

					ID:          "post_id_2",
					Title:       "post_title_2",
					Content:     "post_content_2",
					Slug:        "post-slug-2",
					Draft:       false,
					PublishedAt: &now,
					Tags: []*model.Tag{
						{
							ID:   "tag_id_1",
							Name: "tag1",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "タグによって取得",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPostRepository) {
				mock.EXPECT().FindAllPosts(gomock.Any()).Return([]*model.Post{
					{
						ID:          "post_id_1",
						Title:       "post_title_1",
						Content:     "pot_content_1",
						Slug:        "post-slug-1",
						Draft:       true,
						PublishedAt: &now,
						Tags: []*model.Tag{
							{
								ID:   "tag_id_1",
								Name: "tag1",
							},
						},
					},
					{
						ID:          "post_id_2",
						Title:       "post_title_2",
						Content:     "post_content_2",
						Slug:        "post-slug-2",
						Draft:       false,
						PublishedAt: &now,
						Tags: []*model.Tag{
							{
								ID:   "tag_id_1",
								Name: "tag1",
							},
						},
					},
				}, nil)
			},
			conditions: []string{"tags.name = 'tag_id_1'"},
			want: []*model.Post{
				{
					ID:          "post_id_1",
					Title:       "post_title_1",
					Content:     "pot_content_1",
					Slug:        "post-slug-1",
					Draft:       true,
					PublishedAt: &now,
					Tags: []*model.Tag{
						{
							ID:   "tag_id_1",
							Name: "tag1",
						},
					},
				},
				{

					ID:          "post_id_2",
					Title:       "post_title_2",
					Content:     "post_content_2",
					Slug:        "post-slug-2",
					Draft:       false,
					PublishedAt: &now,
					Tags: []*model.Tag{
						{
							ID:   "tag_id_1",
							Name: "tag1",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "記事の取得に失敗したときエラーを返す",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPostRepository) {
				mock.EXPECT().FindAllPosts(gomock.Any()).Return(nil, errors.New("error"))
			},
			conditions: []string{},
			want:       nil,
			wantErr:    true,
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

			got, err := ps.GetPosts(tt.conditions)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPosts() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPosts() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestPostService_UpdatePost(t *testing.T) {

	existPost := &model.Post{
		ID:          "post_id_1",
		Title:       "post_title_1",
		Content:     "pot_content_1",
		Slug:        "post-slug-1",
		Draft:       true,
		PublishedAt: func() *time.Time { t := time.Now(); return &t }(),
	}

	tests := []struct {
		name                  string
		post                  *model.Post
		prepareMockPostRepoFn func(mock *mock_repository.MockPostRepository)
		wantErr               bool
	}{
		{
			name: "記事を更新できたときはエラーを返さない",
			post: &model.Post{
				ID:          "post_id_1_update",
				Title:       "post_title_1_update",
				Content:     "pot_content_1_update",
				Slug:        "post-slug-1-update",
				Draft:       false,
				PublishedAt: func() *time.Time { t := time.Now(); return &t }(),
			},
			prepareMockPostRepoFn: func(mock *mock_repository.MockPostRepository) {
				mock.EXPECT().FindPostByID(gomock.Any()).Return(existPost, nil)
				mock.EXPECT().UpdatePost(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "更新対象の記事がないときはエラーを返す",
			post: &model.Post{
				ID:          "post_id_1_update",
				Title:       "post_title_1_update",
				Content:     "pot_content_1_update",
				Slug:        "post-slug-1-update",
				Draft:       false,
				PublishedAt: func() *time.Time { t := time.Now(); return &t }(),
			},
			prepareMockPostRepoFn: func(mock *mock_repository.MockPostRepository) {
				mock.EXPECT().FindPostByID(gomock.Any()).Return(nil, errors.New("error"))
			},
			wantErr: true,
		},
		{
			name: "記事の更新に失敗したときエラーを返す",
			post: &model.Post{
				ID:          "post_id_1_update",
				Title:       "post_title_1_update",
				Content:     "pot_content_1_update",
				Slug:        "post-slug-1-update",
				Draft:       false,
				PublishedAt: func() *time.Time { t := time.Now(); return &t }(),
			},
			prepareMockPostRepoFn: func(mock *mock_repository.MockPostRepository) {
				mock.EXPECT().FindPostByID(gomock.Any()).Return(existPost, nil)
				mock.EXPECT().UpdatePost(gomock.Any()).Return(errors.New("error"))
			},
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

			err := ps.UpdatePost(tt.post)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPosts() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostService_DeletePost(t *testing.T) {
	tests := []struct {
		name                  string
		id                    string
		prepareMockPostRepoFn func(mock *mock_repository.MockPostRepository)
		wantErr               bool
	}{
		{
			name: "記事を削除できたときはエラーを返さない",
			id:   "post_id_1",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPostRepository) {
				mock.EXPECT().DeletePostByID(gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "記事の削除に失敗したときエラーを返す",
			id:   "not_found",
			prepareMockPostRepoFn: func(mock *mock_repository.MockPostRepository) {
				mock.EXPECT().DeletePostByID("not_found").Return(model.ErrPostNotFound)
			},
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

			err := ps.DeletePost(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeletePost() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
