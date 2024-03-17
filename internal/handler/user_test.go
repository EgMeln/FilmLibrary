package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
)

type mockUserService struct {
	RegisterFunc func(user *model.User) error
	LoginFunc    func(user *model.User) error
}

func (m *mockUserService) Register(user *model.User) error {
	return m.RegisterFunc(user)
}

func (m *mockUserService) Login(user *model.User) error {
	return m.LoginFunc(user)
}

func TestUserHandler_Register(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		formData           map[string]string
		registerFunc       func(user *model.User) error
		expectedStatusCode int
	}{
		{
			name: "Success",
			formData: map[string]string{
				"username": "testuser",
				"password": "testpassword",
			},
			registerFunc: func(user *model.User) error {
				return nil
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "InvalidFormData",
			formData: map[string]string{
				"password": "testpassword",
			},
			registerFunc:       nil,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "ServiceError",
			formData: map[string]string{
				"username": "testuser",
				"password": "testpassword",
			},
			registerFunc: func(user *model.User) error {
				return errors.New("service error")
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userService := &mockUserService{
				RegisterFunc: tc.registerFunc,
			}
			userHandler := NewUserHandler(userService)

			form := url.Values{}
			for key, value := range tc.formData {
				form.Add(key, value)
			}
			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(form.Encode()))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			recorder := httptest.NewRecorder()
			userHandler.Register(recorder, req)

			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
			}
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		user               model.User
		loginFunc          func(user *model.User) error
		expectedStatusCode int
	}{
		{
			name: "Success",
			user: model.User{
				Username: "testuser",
				Password: "testpassword",
			},
			loginFunc: func(user *model.User) error {
				return nil
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Unauthorized",
			user: model.User{
				Username: "testuser",
				Password: "testpassword",
			},
			loginFunc: func(user *model.User) error {
				return errors.New("unauthorized")
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userService := &mockUserService{
				LoginFunc: tc.loginFunc,
			}
			userHandler := NewUserHandler(userService)

			requestBody, _ := json.Marshal(tc.user)
			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			userHandler.Login(recorder, req)

			if recorder.Code != tc.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatusCode, recorder.Code)
			}
		})
	}
}
