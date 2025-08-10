// backend/internal/models/user.go
package models

import "time"

// User corresponds to the "users" table in the database.
type User struct {
	ID           string    `json:"id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`               // The dash ensures this is never sent in JSON responses
	Phone        *string   `json:"phone,omitempty"` // Pointer for nullable fields
	RoleID       *string   `json:"role_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Role corresponds to the "roles" table.
type Role struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
