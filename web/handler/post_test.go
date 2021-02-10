package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Le0tk0k/blog-server/model"
	"github.com/Le0tk0k/blog-server/service/mock_service"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func TestPostHandler_CreatePost(t *testing.T) {
	tests := []struct {
		name                     string
		prepareMockPostServiceFn func(mock *mock_service.MockPostService)
		body                     string
		wantErr                  bool
		wantCode                 int
	}{
		{
			name: "正常に投稿を保存できたときは201を返す",
			prepareMockPostServiceFn: func(mock *mock_service.MockPostService) {
				mock.EXPECT().CreatePost().Return(model.NewPost(), nil)
			},
			body:     ``,
			wantErr:  false,
			wantCode: http.StatusCreated,
		},
		{
			name: "投稿の保存に失敗したときは500を返す",
			prepareMockPostServiceFn: func(mock *mock_service.MockPostService) {
				mock.EXPECT().CreatePost().Return(nil, errors.New("error"))
			},
			body:     ``,
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ms := mock_service.NewMockPostService(ctrl)
			tt.prepareMockPostServiceFn(ms)
			ph := &PostHandler{postService: ms}

			e := echo.New()
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			w := httptest.NewRecorder()
			c := e.NewContext(r, w)
			c.SetPath("/posts")

			postErr := ph.CreatePost(c)
			if (postErr != nil) != tt.wantErr {
				t.Errorf("CreatePost() error = %v, wantErr %v", postErr, tt.wantErr)
			}
			if er, ok := postErr.(*echo.HTTPError); (ok && er.Code != tt.wantCode) || (!ok && w.Code != tt.wantCode) {
				t.Errorf("CreatePost() code = %d, want = %d", w.Code, tt.wantCode)
			}
		})
	}
}
