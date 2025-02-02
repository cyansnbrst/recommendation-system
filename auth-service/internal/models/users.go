package models

import "time"

// User struct
type User struct {
	ID           string
	Email        string
	Name         string // not stored in db
	PasswordHash string
	IsAdmin      bool
	CreatedAt    time.Time
}
