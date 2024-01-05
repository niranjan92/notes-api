package entity

import "time"

// User represents a user.
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u User) TableName() string {
	return "users"
}

// GetID returns the user ID.
func (u User) GetID() string {
	return u.ID
}

// GetName returns the user name.
func (u User) GetName() string {
	return u.Name
}
