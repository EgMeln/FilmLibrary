package service

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
	"github.com/EgMeln/filmLibraryPrivate/internal/repository"
)

// UserService represents a service for managing user accounts.
type UserService interface {
	Register(user *model.User) error
	Login(user *model.User) error
}

type userService struct {
	userManager repository.UserManager
}

// NewUserService creates a new instance of the UserService with the provided UserManager.
func NewUserService(userManager repository.UserManager) UserService {
	return &userService{
		userManager: userManager,
	}
}

// Register creates a new user account if the provided username is unique,
func (us *userService) Register(user *model.User) error {
	ifExist, err := us.userManager.IfExist(user.Username)
	if err != nil {
		return err
	}
	if ifExist {
		return errors.New("the user already exists")
	}
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return err
	}

	user.ID = uuid.New()
	err = us.userManager.Create(user)
	if err != nil {
		return err
	}
	return err
}

// Login authenticates the user by checking if the username exists,
func (us *userService) Login(user *model.User) error {
	ifExist, err := us.userManager.IfExist(user.Username)
	if err != nil {
		return err
	}
	if !ifExist {
		return errors.New("the user not exists")
	}

	getUser, err := us.userManager.GetByUsername(user.Username)
	if err != nil {
		return err
	}
	if !comparePassword(getUser.Password, user.Password) {
		return errors.New("password is incorrect")

	}
	return err

}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func comparePassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
