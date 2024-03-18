package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

func TestAuthAdminMiddleware(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		role           string
		token          string
		expectedStatus int
	}{
		{
			name:           "ValidAdminToken",
			role:           "admin",
			token:          "valid_admin_token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "InvalidUserToken",
			role:           "user",
			token:          "valid_user_token",
			expectedStatus: http.StatusForbidden,
		},

		{
			name:           "InvalidRole",
			role:           "invalid",
			token:          "invalid_role_token",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "EmptyHeader",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := jwt.MapClaims{
				"role": tt.role,
				"sub":  tt.name,
				"exp":  time.Now().Add(time.Hour * 72).Unix(),
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

			tk, _ := token.SignedString([]byte("your-secret-key"))

			req := httptest.NewRequest(http.MethodGet, "/", nil)

			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tk)
			}

			recorder := httptest.NewRecorder()

			AuthAdminMiddleware(handler).ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, recorder.Code)
			}
		})
	}
}

func TestAuthUserMiddleware(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		role           string
		token          string
		expectedStatus int
	}{
		{
			name:           "ValidAdminToken",
			role:           "admin",
			token:          "valid_admin_token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ValidUserToken",
			role:           "user",
			token:          "valid_user_token",
			expectedStatus: http.StatusOK,
		},

		{
			name:           "InvalidRole",
			role:           "invalid",
			token:          "invalid_role_token",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "EmptyHeader",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := jwt.MapClaims{
				"role": tt.role,
				"sub":  tt.name,
				"exp":  time.Now().Add(time.Hour * 72).Unix(),
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

			tk, _ := token.SignedString([]byte("your-secret-key"))

			req := httptest.NewRequest(http.MethodGet, "/", nil)

			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tk)
			}

			recorder := httptest.NewRecorder()

			AuthUserMiddleware(handler).ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, recorder.Code)
			}
		})
	}
}
