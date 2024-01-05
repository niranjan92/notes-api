package entity

import (
	"time"
)

// Note represents an note record.
type Note struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Text           string    `json:"text"`
	TextSearchable string    `json:"text_searchable"`
	UserID         string    `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (u Note) TableName() string {
	return "notes"
}

// SharedNote used to share note with other users
type SharedNote struct {
	ID           string `json:"id"`
	NoteID       string `json:"note_id"`
	SharedUserID string `json:"shared_user_id"`
}

func (u SharedNote) TableName() string {
	return "shared_notes"
}
