package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/Le0tk0k/blog-server/model"
	"github.com/Le0tk0k/blog-server/repository/mock_repository"
)

func TestUserService_GetUser(t *testing.T) {
	existsUser := &model.User{
		Email:    "example@example.com",
		Password: "password",
	}

	tests := []struct {
		name                  string
		email                 string
		prepareMockUserRepoFn func(mock *mock_repository.MockUserRepository)
		want                  *model.User
		wantErr               bool
	}{
		{
			name:  "Userを返す",
			email: "example@example.com",
			prepareMockUserRepoFn: func(mock *mock_repository.MockUserRepository) {
				mock.EXPECT().FindByEmail("example@example.com").Return(existsUser, nil)
			},
			want: &model.User{
				Email:    "example@example.com",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name:  "Userの取得に失敗したときエラーを返す",
			email: "not_found",
			prepareMockUserRepoFn: func(mock *mock_repository.MockUserRepository) {
				mock.EXPECT().FindByEmail("not_found").Return(nil, errors.New("error"))
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mr := mock_repository.NewMockUserRepository(ctrl)
			tt.prepareMockUserRepoFn(mr)
			us := &userService{
				userRepository: mr,
			}

			got, err := us.GetUser(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want = %v", got, tt.want)
			}
		})
	}
}
