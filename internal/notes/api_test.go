package notes

import (
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
		{"123", "note123", "text123", "text_searchable123", "user123", now, now},
	}}

	// ignore rate limiter and use mock auth handler itself for now
	RegisterHandlers(router.Group(""), NewService(repo, logger), auth.MockAuthHandler, auth.MockAuthHandler, logger)
	// header := auth.MockAuthHeader()

	tests := []test.APITestCase{
		// TODO:
		// {"get 123", "GET", "/notes/123", "", nil, http.StatusOK, `*text123*`},
	}
	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
