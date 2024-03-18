package model

import "github.com/google/uuid"

// User represents information about an user.
type User struct {
	ID       uuid.UUID
	Username string // Name of the user
	Password string // Gender of the user
	Role     string // Role of the user
}
