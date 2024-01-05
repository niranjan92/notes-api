package notes

import (
	"net/http"
	"testing"
	"time"

	"github.com/qiangxue/go-rest-api/internal/auth"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/internal/test"
	"github.com/qiangxue/go-rest-api/pkg/log"
)

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	router := test.MockRouter(logger)

	now := time.Now()
	repo := &mockNoteRepo{items: []entity.Note{
		{"123", "note123", "text123", "text_searchable123", "testuser", now, now},
	}}

	// ignore rate limiter and use mock auth handler itself for now
	RegisterHandlers(router.Group(""), NewService(repo, logger), auth.MockAuthHandler, auth.MockAuthHandler, logger)
	header := auth.MockAuthHeader()

	tests := []test.APITestCase{
		{"get 123", "GET", "/notes/123", "", header, http.StatusOK, `*text123*`},
		{"get all", "GET", "/notes", "", header, http.StatusOK, `*text123*`},
		{"get unknown", "GET", "/albums/1234", "", header, http.StatusNotFound, ""},
		{"create ok", "POST", "/notes", `{"title":"test", "text": "text1"}`, header, http.StatusCreated, "*test*"},
		{"create ok count", "GET", "/notes", "", header, http.StatusOK, `*"total_count":2*`},
		{"create auth error", "POST", "/notes", `{"title":"test2", "text": "text2"}`, nil, http.StatusUnauthorized, ""},
		{"create input error", "POST", "/notes", `{"title":"test2"}`, header, http.StatusBadRequest, ""},
		{"update ok", "PUT", "/notes/123", `{"title":"test_changed"}`, header, http.StatusOK, "*test_changed*"},
		{"update verify", "GET", "/notes/123", "", header, http.StatusOK, `*test_changed*`},
		{"update auth error", "PUT", "/notes/123", `{"title":"notesxyz"}`, nil, http.StatusUnauthorized, ""},
		{"update input error", "PUT", "/notes/123", `"name":"notesxyz"}`, header, http.StatusBadRequest, ""},
		{"delete ok", "DELETE", "/notes/123", ``, header, http.StatusOK, "*test_changed*"},
		{"delete verify", "DELETE", "/notes/123", ``, header, http.StatusNotFound, ""},
		{"delete auth error", "DELETE", "/notes/123", ``, nil, http.StatusUnauthorized, ""},
	}
	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
