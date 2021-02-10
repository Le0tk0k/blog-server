package service

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/Le0tk0k/blog-server/repository/mock_repository"
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
