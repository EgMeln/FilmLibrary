package service

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
)

type mockUserManager struct {
	IfExistFunc       func(username string) (bool, error)
	CreateFunc        func(user *model.User) error
	GetByUsernameFunc func(username string) (*model.User, error)
}

func (m *mockUserManager) IfExist(username string) (bool, error) {
	return m.IfExistFunc(username)
}

func (m *mockUserManager) Create(user *model.User) error {
	return m.CreateFunc(user)
}

func (m *mockUserManager) GetByUsername(username string) (*model.User, error) {
	return m.GetByUsernameFunc(username)
}

func TestUserService_Register(t *testing.T) {
	t.Parallel()

	var (
		users = make(map[string]*model.User)
	)

	mockUserManager := &mockUserManager{
		IfExistFunc: func(username string) (bool, error) {
			_, ok := users[username]
			return ok, nil
		},
		CreateFunc: func(user *model.User) error {
			users[user.Username] = user
			return nil
		},
	}

	tests := []struct {
		name           string
		user           *model.User
		expectedResult error
	}{
		{
			name:           "Success",
			user:           &model.User{Username: "KenRyanGosling", Password: "MargoRobbieTheBest"},
			expectedResult: nil,
		},
		{
			name:           "UserExists",
			user:           &model.User{Username: "KenRyanGosling", Password: "MargoRobbieTheBest"},
			expectedResult: errors.New("the user already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &userService{
				userManager: mockUserManager,
			}

			err := us.Register(tt.user)

			if err != nil && err.Error() != tt.expectedResult.Error() {
				t.Errorf("Expected error: %v, got: %v", tt.expectedResult, err)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	t.Parallel()
	pass, _ := bcrypt.GenerateFromPassword([]byte("MargoRobbieTheBest"), 10)

	mockUserManager := &mockUserManager{
		IfExistFunc: func(username string) (bool, error) {
			existingUsernames := map[string]bool{
				"KenRyanGosling": true,
			}
			_, ok := existingUsernames[username]
			return ok, nil
		},

		GetByUsernameFunc: func(username string) (*model.User, error) {
			users := map[string]*model.User{
				"KenRyanGosling": {
					Username: "KenRyanGosling",
					Password: string(pass),
				},
			}
			user, ok := users[username]
			if !ok {
				return nil, errors.New("user not found")
			}
			return user, nil
		},
	}

	tests := []struct {
		name           string
		user           *model.User
		expectedResult error
	}{
		{
			name:           "Success",
			user:           &model.User{Username: "KenRyanGosling", Password: "MargoRobbieTheBest"},
			expectedResult: nil,
		},
		{
			name:           "UserNotExists",
			user:           &model.User{Username: "BarbiRyanGosling", Password: "MargoRobbieTheBest"},
			expectedResult: errors.New("the user not exists"),
		},
		{
			name:           "IncorrectPassword",
			user:           &model.User{Username: "KenRyanGosling", Password: "MargoRobbieTheWorst"},
			expectedResult: errors.New("password is incorrect"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &userService{
				userManager: mockUserManager,
			}

			err := us.Login(tt.user)

			if err != nil && err.Error() != tt.expectedResult.Error() {
				t.Errorf("Expected error: %v, got: %v", tt.expectedResult, err)
			}
		})
	}
}
