package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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
			name: "正常に記事を保存できたときは201を返す",
			prepareMockPostServiceFn: func(mock *mock_service.MockPostService) {
				mock.EXPECT().CreatePost().Return(model.NewPost(), nil)
			},
			body:     ``,
			wantErr:  false,
			wantCode: http.StatusCreated,
		},
		{
			name: "記事の保存に失敗したときは500を返す",
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
			r := httptest.NewRequest(http.MethodPost, "/v1/posts", strings.NewReader(tt.body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			w := httptest.NewRecorder()
			c := e.NewContext(r, w)

			err := ph.CreatePost(c)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePost() error = %v, wantErr %v", err, tt.wantErr)
			}
			if er, ok := err.(*echo.HTTPError); (ok && er.Code != tt.wantCode) || (!ok && w.Code != tt.wantCode) {
				t.Errorf("CreatePost() code = %d, want = %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestPostHandler_GetPost(t *testing.T) {
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
		name                     string
		id                       string
		prepareMockPostServiceFn func(mock *mock_service.MockPostService)
		wantErr                  bool
		wantCode                 int
	}{
		{
			name: "正常に記事を取得できたときは200を返す",
			id:   "post_id_1",
			prepareMockPostServiceFn: func(mock *mock_service.MockPostService) {
				mock.EXPECT().GetPost(gomock.Any()).Return(existsPost, nil)
			},
			wantErr:  false,
			wantCode: http.StatusOK,
		},
		{
			name: "記事の取得に失敗した場合は500を返す",
			id:   "not_found",
			prepareMockPostServiceFn: func(mock *mock_service.MockPostService) {
				mock.EXPECT().GetPost(gomock.Any()).Return(nil, errors.New("error"))
			},
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
			r := httptest.NewRequest(http.MethodGet, "/v1/posts/"+tt.id, nil)
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			w := httptest.NewRecorder()
			c := e.NewContext(r, w)

			err := ph.GetPost(c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPosts() error = %v, wantErr %v", err, tt.wantErr)
			}
			if er, ok := err.(*echo.HTTPError); (ok && er.Code != tt.wantCode) || (!ok && w.Code != tt.wantCode) {
				t.Errorf("GetPosts() code = %d, want = %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestPostHandler_GetPosts(t *testing.T) {
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
		name                     string
		prepareMockPostServiceFn func(mock *mock_service.MockPostService)
		wantErr                  bool
		wantCode                 int
	}{
		{
			name: "正常に記事を取得できたときは200を返す",
			prepareMockPostServiceFn: func(mock *mock_service.MockPostService) {
				mock.EXPECT().GetPosts(gomock.Any()).Return(existsPosts, nil)
			},
			wantErr:  false,
			wantCode: http.StatusOK,
		},
		{
			name: "記事が0件でもエラーにならずに200を返す",
			prepareMockPostServiceFn: func(mock *mock_service.MockPostService) {
				mock.EXPECT().GetPosts(gomock.Any()).Return(nil, nil)
			},
			wantErr:  false,
			wantCode: http.StatusOK,
		},
		{
			name: "記事の取得に失敗した場合は500を返す",
			prepareMockPostServiceFn: func(mock *mock_service.MockPostService) {
				mock.EXPECT().GetPosts(gomock.Any()).Return(nil, errors.New("error"))
			},
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
			r := httptest.NewRequest(http.MethodGet, "/v1/posts", nil)
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			w := httptest.NewRecorder()
			c := e.NewContext(r, w)

			err := ph.GetPosts(c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPosts() error = %v, wantErr %v", err, tt.wantErr)
			}
			if er, ok := err.(*echo.HTTPError); (ok && er.Code != tt.wantCode) || (!ok && w.Code != tt.wantCode) {
				t.Errorf("GetPosts() code = %d, want = %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestPostHandler_UpdatePost(t *testing.T) {
	tests := []struct {
		name                        string
		id                          string
		prepareMockPostServiceFn    func(mock *mock_service.MockPostService)
		prepareMockPostTagServiceFn func(mock *mock_service.MockPostTagService)
		body                        string
		wantErr                     bool
		wantCode                    int
	}{
		{
			name: "正常に記事を更新できたときは200を返す",
			id:   "new_post_id",
			prepareMockPostServiceFn: func(mock *mock_service.MockPostService) {
				mock.EXPECT().UpdatePost(gomock.Any()).Return(nil)
			},
			prepareMockPostTagServiceFn: func(mock *mock_service.MockPostTagService) {
				mock.EXPECT().LinkPostTag(gomock.Any()).Return(nil)
			},
			body: `{
					"id": "new_post_id",
					"title": "new_post_title",
					"content": "new_post_content",
					"slug": "new-post-slug",
					"draft": false,
					"published_at": "2021-02-13T14:03:55+09:00",
					"tags": [
        				{
            				"id": "2",
            				"name": "tag2"
        				},
        				{
            				"id": "4",
            				"name": "tag4"
        				}
    				]
					}`,
			wantErr:  false,
			wantCode: http.StatusOK,
		},
		{
			name: "記事の更新に失敗したときは500を返す-1",
			id:   "new_post_id",
			prepareMockPostServiceFn: func(mock *mock_service.MockPostService) {
				mock.EXPECT().UpdatePost(gomock.Any()).Return(errors.New("error"))
			},
			prepareMockPostTagServiceFn: func(mock *mock_service.MockPostTagService) {
			},
			body: `{
					"id": "new_post_id",
					"title": "new_post_title",
					"content": "new_post_content",
					"slug": "new-post-slug",
					"draft": false,
					"published_at": "2021-02-13T14:03:55+09:00",
					"tags": [
        				{
            				"id": "2",
            				"name": "tag2"
        				},
        				{
            				"id": "4",
            				"name": "tag4"
        				}
    				]
					}`,
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "記事の更新に失敗したときは500を返す-2",
			id:   "new_post_id",
			prepareMockPostServiceFn: func(mock *mock_service.MockPostService) {
				mock.EXPECT().UpdatePost(gomock.Any()).Return(nil)
			},
			prepareMockPostTagServiceFn: func(mock *mock_service.MockPostTagService) {
				mock.EXPECT().LinkPostTag(gomock.Any()).Return(errors.New("error"))
			},
			body: `{
					"id": "new_post_id",
					"title": "new_post_title",
					"content": "new_post_content",
					"slug": "new-post-slug",
					"draft": false,
					"published_at": "2021-02-13T14:03:55+09:00",
					"tags": [
        				{
            				"id": "2",
            				"name": "tag2"
        				},
        				{
            				"id": "4",
            				"name": "tag4"
        				}
    				]
					}`,
			wantErr:  true,
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mps := mock_service.NewMockPostService(ctrl)
			mpts := mock_service.NewMockPostTagService(ctrl)
			tt.prepareMockPostServiceFn(mps)
			tt.prepareMockPostTagServiceFn(mpts)
			ph := &PostHandler{
				postService:    mps,
				postTagService: mpts,
			}

			e := echo.New()
			r := httptest.NewRequest(http.MethodPut, "/v1/posts/"+tt.id, strings.NewReader(tt.body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			w := httptest.NewRecorder()
			c := e.NewContext(r, w)

			err := ph.UpdatePost(c)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePost() error = %v, wantErr %v", err, tt.wantErr)
			}
			if er, ok := err.(*echo.HTTPError); (ok && er.Code != tt.wantCode) || (!ok && w.Code != tt.wantCode) {
				t.Errorf("UpdatePost() code = %d, want = %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestPostHandler_DeletePost(t *testing.T) {
	tests := []struct {
		name                     string
		id                       string
		prepareMockPostServiceFn func(mock *mock_service.MockPostService)
		wantErr                  bool
		wantCode                 int
	}{
		{
			name: "正常に記事を削除できたときは200を返す",
			id:   "post_id_1",
			prepareMockPostServiceFn: func(mock *mock_service.MockPostService) {
				mock.EXPECT().DeletePost(gomock.Any()).Return(nil)
			},
			wantErr:  false,
			wantCode: http.StatusOK,
		},
		{
			name: "記事の削除に失敗した場合は500を返す",
			id:   "not_found",
			prepareMockPostServiceFn: func(mock *mock_service.MockPostService) {
				mock.EXPECT().DeletePost(gomock.Any()).Return(errors.New("error"))
			},
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
			r := httptest.NewRequest(http.MethodDelete, "/v1/posts/"+tt.id, nil)
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			w := httptest.NewRecorder()
			c := e.NewContext(r, w)

			err := ph.DeletePost(c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPosts() error = %v, wantErr %v", err, tt.wantErr)
			}
			if er, ok := err.(*echo.HTTPError); (ok && er.Code != tt.wantCode) || (!ok && w.Code != tt.wantCode) {
				t.Errorf("GetPosts() code = %d, want = %d", w.Code, tt.wantCode)
			}
		})
	}
}
