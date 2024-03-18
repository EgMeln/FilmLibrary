package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/EgMeln/filmLibraryPrivate/internal/model"
	"github.com/EgMeln/filmLibraryPrivate/internal/service"
)

// jwtSecretKey is the secret key used for JWT token signing.
var jwtSecretKey = []byte("very-secret-key")

// UserHandler handles HTTP requests related to users.
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register handles the HTTP request to register a new user.
// @Summary Register a new user
// @Description Register a new user with a username and password
// @Tags users
// @Accept json
// @Produce json
// @Param username formData string true "Username"
// @Param password formData string true "Password"
// @Success 201 {string} string "User created successfully"
// @Failure 400 {string} string "Failed to parse form or username and password are required"
// @Failure 500 {string} string "Failed to create user"
// @Router /register [post]
func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		log.Printf("Failed to parse form: %v", err)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		log.Printf("Username and password are required")
		return
	}
	user := &model.User{
		Username: username,
		Password: password,
	}
	err = uh.userService.Register(user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Printf("Failed to create user: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login handles the HTTP request for user login.
// @Summary Login
// @Description Log in an existing user with a username and password
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {string} string "Login successful"
// @Failure 400 {string} string "Unable to decode request body"
// @Failure 401 {string} string "Unauthorized"
// @Router /login [post]
func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Unable to decode request body", http.StatusBadRequest)
		return
	}

	err = uh.userService.Login(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	payload := jwt.MapClaims{
		"role": user.Role,
		"sub":  user.Username,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		log.Printf("JWT token signing")
		http.Error(w, "JWT token signing", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]string{"token": t})
}
