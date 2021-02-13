package service

import (
	"testing"

	"github.com/Le0tk0k/blog-server/repository/mock_repository"
	"github.com/golang/mock/gomock"
)

func TestTagService_CreateTag(t *testing.T) {
	tests := []struct {
		name                 string
		tagName              string
		prepareMockTagRepoFn func(mock *mock_repository.MockTagRepository)
		wantErr              error
	}{
		{
			name:    "新規タグを生成し保存する",
			tagName: "tag_name",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTagRepository) {
				mock.EXPECT().StoreTag(gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockTagRepository(ctrl)
			tt.prepareMockTagRepoFn(mr)
			ts := &tagService{
				tagRepository: mr,
			}

			err := ts.CreateTag(tt.tagName)
			if err != tt.wantErr {
				t.Errorf("CreateTag() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}
