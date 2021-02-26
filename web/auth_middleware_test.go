package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Le0tk0k/blog-server/model"
	"github.com/Le0tk0k/blog-server/service/mock_service"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func TestAuthMiddleware_Login(t *testing.T) {
	existsUser := &model.User{
		Email:    "example@example.com",
		Password: "$2a$12$xbkwp1EqD1z5EZnsbN/pA.ECLHpMfxgW.lAfbp2/NphThAer4yw4C",
	}

	tests := []struct {
		name                     string
		prepareMockUserServiceFn func(mock *mock_service.MockUserService)
		body                     string
		wantErr                  bool
		wantCode                 int
	}{
		{
			name: "正常にログインできたときは200を返す",
			prepareMockUserServiceFn: func(mock *mock_service.MockUserService) {
				mock.EXPECT().GetUser(gomock.Any()).Return(existsUser, nil)
			},
			body: `{
					"email": "example@example.com",
					"password": "password"
					}`,
			wantErr:  false,
			wantCode: http.StatusOK,
		},
		{
			name: "ログインに失敗した場合は500を返す",
			prepareMockUserServiceFn: func(mock *mock_service.MockUserService) {
				mock.EXPECT().GetUser(gomock.Any()).Return(nil, errors.New("error"))
			},
			body: `{
					"email": "example@example.com",
					"password": "password"
					}`,
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ms := mock_service.NewMockUserService(ctrl)
			tt.prepareMockUserServiceFn(ms)
			m := &AuthMiddleware{userService: ms}

			e := echo.New()
			r := httptest.NewRequest(http.MethodGet, "/v1/login", strings.NewReader(tt.body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			w := httptest.NewRecorder()
			c := e.NewContext(r, w)

			err := m.Login(c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
			if er, ok := err.(*echo.HTTPError); (ok && er.Code != tt.wantCode) || (!ok && w.Code != tt.wantCode) {
				t.Errorf("Login() code = %d, want = %d", w.Code, tt.wantCode)
			}
		})
	}
}
