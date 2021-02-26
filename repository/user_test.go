package repository

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Le0tk0k/blog-server/model"
)

func TestUserRepository_FindByEmail(t *testing.T) {
	existUser := &model.User{
		Email:    "example@example.com",
		Password: "password",
	}
	_, err := db.Exec("INSERT INTO users VALUES (?, ?, ?)", "1", existUser.Email, existUser.Password)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		email   string
		want    *model.User
		wantErr error
	}{
		{
			name:  "存在するuserを正常に取得できる",
			email: "example@example.com",
			want: &model.User{
				Email:    "example@example.com",
				Password: "password",
			},
			wantErr: nil,
		},
		{
			name:    "存在しないIDの場合ErrTagNotFoundを返す",
			email:   "not_found",
			want:    nil,
			wantErr: model.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepository{db: db}
			got, err := r.FindByEmail(tt.email)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("FindByEmail()  error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByEmail() got = %v, want = %v", got, tt.want)
			}
		})
	}

	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		t.Fatal(err)
	}
}
