package notes

import (
	"context"
	"testing"

	"github.com/qiangxue/go-rest-api/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestCreateNoteRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     CreateNoteRequest
		wantError bool
	}{
		{"success", CreateNoteRequest{Title: "test", Text: "text213"}, false},
		{"required", CreateNoteRequest{Title: ""}, true},
		{"too long", CreateNoteRequest{Title: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func Test_service_CRUD(t *testing.T) {
	logger, _ := log.NewForTest()
	s := NewService(&mockNoteRepo{}, logger)

	ctx := context.Background()

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, 0, count)

	// successful creation
	note, err := s.Create(ctx, CreateNoteRequest{Title: "test", Text: "text1"})
	assert.Nil(t, err)
	assert.NotEmpty(t, note.ID)
	id := note.ID
	assert.Equal(t, "test", note.Title)
	assert.Equal(t, "text1", note.Text)
	assert.NotEmpty(t, note.CreatedAt)
	assert.NotEmpty(t, note.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// validation error in creation
	_, err = s.Create(ctx, CreateNoteRequest{Title: "", Text: "text1"})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// unexpected error in creation
	_, err = s.Create(ctx, CreateNoteRequest{Title: "error", Text: "text1"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	_, _ = s.Create(ctx, CreateNoteRequest{Title: "test2", Text: "text2"})

	// update
	note, err = s.Update(ctx, id, UpdateNoteRequest{Title: "test updated"})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", note.Title)
	_, err = s.Update(ctx, "none", UpdateNoteRequest{Title: "test updated"})
	assert.NotNil(t, err)

	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// unexpected error in update
	_, err = s.Update(ctx, id, UpdateNoteRequest{Title: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	note, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", note.Title)
	assert.Equal(t, id, note.ID)

	// query
	notes, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, len(notes))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	note, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, note.ID)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)
}
