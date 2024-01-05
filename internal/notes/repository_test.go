package notes

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/internal/test"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var errCRUD = errors.New("error crud")

// import (
// 	"context"
// 	"database/sql"
// 	"github.com/qiangxue/go-rest-api/internal/entity"
// 	"github.com/qiangxue/go-rest-api/internal/test"
// 	"github.com/qiangxue/go-rest-api/pkg/log"
// 	"github.com/stretchr/testify/assert"
// 	"testing"
// 	"time"
// )

func TestRepository(t *testing.T) {
	logger, _ := log.NewForTest()
	db := test.DB(t)
	gormDSN := "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable TimeZone=Asia/Shanghai"
	gormDB, err := gorm.Open(postgres.Open(gormDSN), &gorm.Config{})
	if err != nil {
		logger.Error(err)
		t.FailNow()
	}
	test.ResetTables(t, db, "notes")
	repo := NewRepository(gormDB, db, logger)

	ctx := context.Background()

	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	// create
	err = repo.Create(ctx, entity.Note{
		ID:             "test1",
		Title:          "title1",
		Text:           "text1",
		TextSearchable: "text1",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, 1, count2-count)

	// get
	note, err := repo.Get(ctx, "test1")
	assert.Nil(t, err)
	assert.Equal(t, "title1", note.Title)
	_, err = repo.Get(ctx, "test0")
	assert.Equal(t, sql.ErrNoRows, err)

	// update
	err = repo.Update(ctx, entity.Note{
		ID:        "test1",
		Title:     "title1 updated",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	assert.Nil(t, err)
	note, _ = repo.Get(ctx, "test1")
	assert.Equal(t, "title1 updated", note.Title)

	// query
	notes, err := repo.Query(ctx, 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, len(notes))

	// delete
	err = repo.Delete(ctx, "test1")
	assert.Nil(t, err)
	_, err = repo.Get(ctx, "test1")
	assert.Equal(t, sql.ErrNoRows, err)
	err = repo.Delete(ctx, "test1")
	assert.Equal(t, sql.ErrNoRows, err)
}

type mockNoteRepo struct {
	items []entity.Note
}

func (m *mockNoteRepo) Get(ctx context.Context, id string) (entity.Note, error) {
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return entity.Note{}, sql.ErrNoRows
}

func (m *mockNoteRepo) Count(ctx context.Context) (int, error) {
	return len(m.items), nil
}

func (m *mockNoteRepo) Query(ctx context.Context, offset, limit int) ([]entity.Note, error) {
	return m.items, nil
}

func (m *mockNoteRepo) QueryByUserID(ctx context.Context, userID string) ([]entity.Note, error) {
	return m.items, nil
}

func (m *mockNoteRepo) Create(ctx context.Context, note entity.Note) error {
	if note.Title == "error" {
		return errCRUD
	}
	m.items = append(m.items, note)
	return nil
}

func (m *mockNoteRepo) Update(ctx context.Context, note entity.Note) error {
	if note.Title == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.ID == note.ID {
			m.items[i] = note
			break
		}
	}
	return nil
}

func (m *mockNoteRepo) Delete(ctx context.Context, id string) error {
	for i, item := range m.items {
		if item.ID == id {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}

func (m *mockNoteRepo) SharedNoteCreate(ctx context.Context, note *entity.SharedNote) error {
	return nil
}

func (m *mockNoteRepo) GetSharedNoteByID(ctx context.Context, id string) (entity.SharedNote, error) {
	return entity.SharedNote{}, nil
}

func (m *mockNoteRepo) QuerySharedNotes(ctx context.Context, userID string) ([]entity.Note, error) {
	return []entity.Note{}, nil
}

func (m *mockNoteRepo) SearchNotes(ctx context.Context, userID string, query string) ([]entity.Note, error) {
	return []entity.Note{}, nil
}
