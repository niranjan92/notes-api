package notes

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/pkg/dbcontext"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"gorm.io/gorm"
)

// Repository encapsulates the logic to access notes from the data source.
type Repository interface {
	// Get returns the note with the specified note ID.
	Get(ctx context.Context, id string) (entity.Note, error)
	// Count returns the number of notes.
	Count(ctx context.Context) (int, error)
	// Query returns the list of notes with the given offset and limit.
	Query(ctx context.Context, offset, limit int) ([]entity.Note, error)
	QueryByUserID(ctx context.Context, userID string) ([]entity.Note, error)
	// Create saves a new note in the storage.
	Create(ctx context.Context, note entity.Note) error
	// Update updates the note with given ID in the storage.
	Update(ctx context.Context, note entity.Note) error
	// Delete removes the note with given ID from the storage.
	Delete(ctx context.Context, id string) error

	SharedNoteCreate(ctx context.Context, note *entity.SharedNote) error
	GetSharedNoteByID(ctx context.Context, id string) (entity.SharedNote, error)

	QuerySharedNotes(ctx context.Context, userID string) ([]entity.Note, error) // returns notes that are shared with the user
	SearchNotes(ctx context.Context, userID string, query string) ([]entity.Note, error)
}

// repository persists notes in database
type repository struct {
	gormDB *gorm.DB
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new note repository
func NewRepository(gormDB *gorm.DB, db *dbcontext.DB, logger log.Logger) Repository {
	return repository{gormDB, db, logger}
}

// Get reads the note with the specified ID from the database.
func (r repository) Get(ctx context.Context, id string) (entity.Note, error) {
	var note entity.Note
	err := r.db.With(ctx).Select().Model(id, &note)
	return note, err
}

// Create saves a new note record in the database.
// It returns the ID of the newly inserted note record.
func (r repository) Create(ctx context.Context, note entity.Note) error {
	return r.db.With(ctx).Model(&note).Insert()
}

// Update saves the changes to an note in the database.
func (r repository) Update(ctx context.Context, note entity.Note) error {
	return r.db.With(ctx).Model(&note).Update()
}

// Delete deletes an note with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id string) error {
	note, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&note).Delete()
}

// Count returns the number of the note records in the database.
func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("notes").Row(&count)
	return count, err
}

// Query retrieves the note records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int) ([]entity.Note, error) {
	var notes []entity.Note
	err := r.db.With(ctx).
		Select().
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&notes)
	return notes, err
}

func (r repository) QueryByUserID(ctx context.Context, userID string) ([]entity.Note, error) {
	var notes []entity.Note
	err := r.db.With(ctx).
		Select().
		Where(dbx.HashExp{"user_id": userID}).
		OrderBy("id").
		All(&notes)
	return notes, err
}

func (r repository) SharedNoteCreate(ctx context.Context, note *entity.SharedNote) error {
	return r.db.With(ctx).Model(note).Insert()
}

func (r repository) GetSharedNoteByID(ctx context.Context, id string) (entity.SharedNote, error) {
	var note entity.SharedNote
	err := r.db.With(ctx).Select().Model(id, &note)
	return note, err
}

func (r repository) QuerySharedNotes(ctx context.Context, userID string) ([]entity.Note, error) {
	var notes []entity.Note

	tx := r.gormDB.Raw("SELECT notes.* FROM notes LEFT JOIN shared_notes ON shared_notes.note_id = notes.id WHERE shared_notes.shared_user_id = ?", userID).Scan(&notes)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return notes, nil
}

func (r repository) SearchNotes(ctx context.Context, userID string, query string) ([]entity.Note, error) {
	var notes []entity.Note

	tx := r.gormDB.Raw("SELECT notes.* FROM notes LEFT JOIN shared_notes ON shared_notes.note_id = notes.id WHERE (shared_notes.shared_user_id = ? OR notes.user_id = ?) "+
		" AND notes.text @@ to_tsquery('english', ?)", userID, userID, query).Scan(&notes)

	if tx.Error != nil {
		return nil, tx.Error
	}
	return notes, nil

}
