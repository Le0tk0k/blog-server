package service

import (
	"testing"
	"time"

	"github.com/Le0tk0k/blog-server/model"
	"github.com/Le0tk0k/blog-server/repository/mock_repository"

	"github.com/golang/mock/gomock"
)

func TestPostTagService_LinkPostTag(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name                     string
		post                     *model.Post
		prepareMockPostTagRepoFn func(mock *mock_repository.MockPostTagRepository)
		wantErr                  error
	}{
		{
			name: "postとtagを関連付ける",
			post: &model.Post{
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
			prepareMockPostTagRepoFn: func(mock *mock_repository.MockPostTagRepository) {
				mock.EXPECT().DeleteByPostID(gomock.Any()).Return(nil)
				mock.EXPECT().StorePostTag(gomock.Any()).Return(nil).AnyTimes()
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockPostTagRepository(ctrl)
			tt.prepareMockPostTagRepoFn(mr)
			pts := &postTagService{
				postTagRepository: mr,
			}

			err := pts.LinkPostTag(tt.post)
			if err != tt.wantErr {
				t.Errorf("LinkPostTag() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}

}
