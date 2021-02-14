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

func TestTagHandler_CreateTag(t *testing.T) {
	tests := []struct {
		name                    string
		prepareMockTagServiceFn func(mock *mock_service.MockTagService)
		body                    string
		wantErr                 bool
		wantCode                int
	}{
		{
			name: "正常にタグを保存できたときは201を返す",
			prepareMockTagServiceFn: func(mock *mock_service.MockTagService) {
				mock.EXPECT().CreateTag(gomock.Any()).Return(nil)
			},
			body: `{
					"name": "new_tag"
					}`,
			wantErr:  false,
			wantCode: http.StatusCreated,
		},
		{
			name: "タグの保存に失敗したときは500を返す",
			prepareMockTagServiceFn: func(mock *mock_service.MockTagService) {
				mock.EXPECT().CreateTag(gomock.Any()).Return(errors.New("error"))
			},
			body: `{
					"name": "new_tag"
					}`,
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ms := mock_service.NewMockTagService(ctrl)
			tt.prepareMockTagServiceFn(ms)
			th := &TagHandler{tagService: ms}

			e := echo.New()
			r := httptest.NewRequest(http.MethodPost, "/v1/tags", strings.NewReader(tt.body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			w := httptest.NewRecorder()
			c := e.NewContext(r, w)

			err := th.CreateTag(c)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTag() error = %v, wantErr %v", err, tt.wantErr)
			}
			if er, ok := err.(*echo.HTTPError); (ok && er.Code != tt.wantCode) || (!ok && w.Code != tt.wantCode) {
				t.Errorf("CreateTag() code = %d, want = %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestTagHandler_GetTag(t *testing.T) {
	existsTag := &model.Tag{
		ID:   "tag_id",
		Name: "tag",
	}

	tests := []struct {
		name                    string
		id                      string
		prepareMockTagServiceFn func(mock *mock_service.MockTagService)
		wantErr                 bool
		wantCode                int
	}{
		{
			name: "正常にタグを取得できたときは200を返す",
			id:   "tag_id",
			prepareMockTagServiceFn: func(mock *mock_service.MockTagService) {
				mock.EXPECT().GetTag(gomock.Any()).Return(existsTag, nil)
			},
			wantErr:  false,
			wantCode: http.StatusOK,
		},
		{
			name: "タグの取得に失敗した場合は500を返す",
			id:   "not_found",
			prepareMockTagServiceFn: func(mock *mock_service.MockTagService) {
				mock.EXPECT().GetTag(gomock.Any()).Return(nil, errors.New("error"))
			},
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ms := mock_service.NewMockTagService(ctrl)
			tt.prepareMockTagServiceFn(ms)
			th := &TagHandler{tagService: ms}

			e := echo.New()
			r := httptest.NewRequest(http.MethodGet, "/v1/tags/"+tt.id, nil)
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			w := httptest.NewRecorder()
			c := e.NewContext(r, w)

			err := th.GetTag(c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTag() error = %v, wantErr %v", err, tt.wantErr)
			}
			if er, ok := err.(*echo.HTTPError); (ok && er.Code != tt.wantCode) || (!ok && w.Code != tt.wantCode) {
				t.Errorf("GetTag() code = %d, want = %d", w.Code, tt.wantCode)
			}
		})
	}
}
