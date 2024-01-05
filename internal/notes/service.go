package notes

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/pkg/log"
)

// Service encapsulates usecase logic for notes.
type Service interface {
	Get(ctx context.Context, id string) (Note, error)
	Query(ctx context.Context, offset, limit int) ([]Note, error)
	QueryByUser(ctx context.Context, userId string) ([]Note, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input CreateNoteRequest) (Note, error)
	Update(ctx context.Context, id string, input UpdateNoteRequest) (Note, error)
	Delete(ctx context.Context, id string) (Note, error)
	ShareNote(ctx context.Context, noteID string, input ShareNoteRequest) (SharedNote, error)
	QuerySharedNotes(ctx context.Context, userID string) ([]Note, error)
	SearchNotes(ctx context.Context, userID string, query string) ([]Note, error)
}

// Note represents the data about an note.
//
//	type Note struct {
//		entity.Note
//	}
type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SharedNote struct {
	entity.SharedNote
}

// CreateNoteRequest represents an note creation request.
type CreateNoteRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID string `json:"user_id"`
}

// ShareNoteRequest represents an note sharing request.
type ShareNoteRequest struct {
	ID           string `json:"id"`
	NoteID       string `json:"note_id"`
	SharedUserID string `json:"shared_user_id"`
}

func (s ShareNoteRequest) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.SharedUserID, validation.Required),
		validation.Field(&s.NoteID, validation.Required),
	)
}

func (s service) ShareNote(ctx context.Context, noteID string, req ShareNoteRequest) (SharedNote, error) {
	if err := req.Validate(); err != nil {
		return SharedNote{}, err
	}
	id := entity.GenerateID()
	sharedNote := entity.SharedNote{
		ID:           id,
		NoteID:       noteID,
		SharedUserID: req.SharedUserID,
	}
	err := s.repo.SharedNoteCreate(ctx, &sharedNote)
	if err != nil {
		return SharedNote{}, err
	}
	return s.GetSharedNoteByID(ctx, id)
}

func (s service) SearchNotes(ctx context.Context, userID string, query string) ([]Note, error) {
	notes, err := s.repo.SearchNotes(ctx, userID, query)
	if err != nil {
		return nil, err
	}
	result := []Note{}
	for _, note := range notes {
		result = append(result, Note{
			ID:        note.ID,
			Title:     note.Title,
			Text:      note.Text,
			UserID:    note.UserID,
			CreatedAt: note.CreatedAt,
			UpdatedAt: note.UpdatedAt,
		})
	}
	return result, nil
}

func (s service) GetSharedNoteByID(ctx context.Context, id string) (SharedNote, error) {
	sharedNote, err := s.repo.GetSharedNoteByID(ctx, id)
	if err != nil {
		return SharedNote{}, err
	}
	return SharedNote{sharedNote}, nil
}

// Validate validates the CreateNoteRequest fields.
func (m CreateNoteRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Title, validation.Required, validation.Length(0, 128)),
		validation.Field(&m.Text, validation.Required, validation.Length(0, 1024)),
	)
}

// UpdateNoteRequest represents an note update request.
type UpdateNoteRequest struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// Validate validates the CreateNoteRequest fields.
func (m UpdateNoteRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Title, validation.Length(1, 128)),
		validation.Field(&m.Text, validation.Length(1, 1024)),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new note service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the note with the specified the note ID.
func (s service) Get(ctx context.Context, id string) (Note, error) {
	note, err := s.repo.Get(ctx, id)
	if err != nil {
		return Note{}, err
	}
	return Note{
		ID:        note.ID,
		Title:     note.Title,
		Text:      note.Text,
		UserID:    note.UserID,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}, nil
}

// Create creates a new note.
func (s service) Create(ctx context.Context, req CreateNoteRequest) (Note, error) {
	if err := req.Validate(); err != nil {
		return Note{}, err
	}
	id := entity.GenerateID()
	now := time.Now()
	note := entity.Note{
		ID:             id,
		Title:          req.Title,
		Text:           req.Text,
		TextSearchable: req.Text,
		UserID:         req.UserID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	err := s.repo.Create(ctx, note)
	if err != nil {
		return Note{}, err
	}
	return s.Get(ctx, id)
}

// Update updates the note with the specified ID.
func (s service) Update(ctx context.Context, id string, req UpdateNoteRequest) (Note, error) {
	if err := req.Validate(); err != nil {
		return Note{}, err
	}

	note, err := s.Get(ctx, id)
	if err != nil {
		return note, err
	}
	note.Title = req.Title
	note.Text = req.Text
	note.UpdatedAt = time.Now()

	noteE := entity.Note{
		ID:             note.ID,
		Title:          note.Title,
		Text:           note.Text,
		TextSearchable: note.Text,
		UserID:         note.UserID,
		CreatedAt:      note.CreatedAt,
		UpdatedAt:      note.UpdatedAt,
	}
	if err := s.repo.Update(ctx, noteE); err != nil {
		return note, err
	}
	return note, nil
}

// Delete deletes the note with the specified ID.
func (s service) Delete(ctx context.Context, id string) (Note, error) {
	note, err := s.Get(ctx, id)
	if err != nil {
		return Note{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return Note{}, err
	}
	return note, nil
}

// Count returns the number of notes.
func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

// Query returns the notes with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int) ([]Note, error) {
	notes, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	result := []Note{}
	for _, note := range notes {
		result = append(result, Note{
			ID:        note.ID,
			Title:     note.Title,
			Text:      note.Text,
			UserID:    note.UserID,
			CreatedAt: note.CreatedAt,
			UpdatedAt: note.UpdatedAt,
		})
	}
	return result, nil
}

func (s service) QueryByUser(ctx context.Context, userID string) ([]Note, error) {
	items, err := s.repo.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := []Note{}
	for _, item := range items {
		result = append(result, Note{
			ID:        item.ID,
			Title:     item.Title,
			Text:      item.Text,
			UserID:    item.UserID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}
	return result, nil
}

func (s service) QuerySharedNotes(ctx context.Context, userID string) ([]Note, error) {
	items, err := s.repo.QuerySharedNotes(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := []Note{}
	for _, item := range items {
		result = append(result, Note{
			ID:        item.ID,
			Title:     item.Title,
			Text:      item.Text,
			UserID:    item.UserID,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}
	return result, nil
}
