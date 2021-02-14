package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Le0tk0k/blog-server/model"

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

func TestTagService_GetTag(t *testing.T) {
	existsTag := &model.Tag{
		ID:   "tag_id",
		Name: "tag",
	}

	tests := []struct {
		name                 string
		id                   string
		prepareMockTagRepoFn func(mock *mock_repository.MockTagRepository)
		want                 *model.Tag
		wantErr              bool
	}{
		{
			name: "タグを返す",
			id:   "tag_id",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTagRepository) {
				mock.EXPECT().FindTagByID("tag_id").Return(existsTag, nil)
			},
			want: &model.Tag{
				ID:   "tag_id",
				Name: "tag",
			},
			wantErr: false,
		},
		{
			name: "記事の取得に失敗したときエラーを返す",
			id:   "not_found",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTagRepository) {
				mock.EXPECT().FindTagByID("not_found").Return(nil, errors.New("error"))
			},
			want:    nil,
			wantErr: true,
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

			got, err := ts.GetTag(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTag() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTag() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestTagService_GetTags(t *testing.T) {
	existsTags := []*model.Tag{{
		ID:   "tag_id_1",
		Name: "tag1",
	}, {
		ID:   "tag_id_2",
		Name: "tag2",
	}}

	tests := []struct {
		name                 string
		prepareMockTagRepoFn func(mock *mock_repository.MockTagRepository)
		want                 []*model.Tag
		wantErr              bool
	}{
		{
			name: "全タグを返す",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTagRepository) {
				mock.EXPECT().FindAllTags().Return(existsTags, nil)
			},
			want: []*model.Tag{
				{
					ID:   "tag_id_1",
					Name: "tag1",
				},
				{
					ID:   "tag_id_2",
					Name: "tag2",
				},
			},
			wantErr: false,
		},
		{
			name: "タグの取得に失敗したときエラーを返す",
			prepareMockTagRepoFn: func(mock *mock_repository.MockTagRepository) {
				mock.EXPECT().FindAllTags().Return(nil, errors.New("error"))
			},
			want:    nil,
			wantErr: true,
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

			got, err := ts.GetTags()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTags() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTags() got = %v, want = %v", got, tt.want)
			}
		})
	}
}
