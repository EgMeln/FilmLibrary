package repository

import (
	"database/sql"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
)

// UserManager represents an interface for managing actors in the system.
type UserManager interface {
	Create(user *model.User) error
	IfExist(username string) (bool, error)
	GetByUsername(username string) (*model.User, error)
}

// NewUserManager returns a new instance of the user repository.
func NewUserManager(db *sql.DB) UserManager {
	return &userManager{
		db: db,
	}
}

// userManager implements CRUD methods for users.
type userManager struct {
	db *sql.DB
}

// Create inserts a new user record into the database.
func (um *userManager) Create(user *model.User) error {
	query := "INSERT INTO users (id, username, password) VALUES ($1, $2, $3)"

	_, err := um.db.Exec(query, user.ID, user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

// IfExist checks the existence of a user with the given username in the database.
func (um *userManager) IfExist(username string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"

	var exists bool
	err := um.db.QueryRow(query, username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// GetByUsername retrieves user information from the database based on the provided username.
func (um *userManager) GetByUsername(username string) (*model.User, error) {
	query := "SELECT id, username, password, role FROM users WHERE username = $1"

	var user model.User

	err := um.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
