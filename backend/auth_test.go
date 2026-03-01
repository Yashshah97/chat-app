package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	// Setup
	db, _ := initTestDatabase()
	server := &Server{
		db:     db,
		router: setupTestRouter(db),
	}

	payload := RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.registerHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp AuthResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Token == "" {
		t.Error("Expected token to be set")
	}

	if resp.User.Username != "testuser" {
		t.Errorf("Expected username testuser, got %s", resp.User.Username)
	}
}

func TestLogin(t *testing.T) {
	// Setup
	db, _ := initTestDatabase()
	server := &Server{
		db:     db,
		router: setupTestRouter(db),
	}

	// Create user first
	hashedPassword, _ := hashPassword("password123")
	user := User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: hashedPassword,
	}
	db.Create(&user)

	payload := LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.loginHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp AuthResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Token == "" {
		t.Error("Expected token to be set")
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	db, _ := initTestDatabase()
	server := &Server{
		db:     db,
		router: setupTestRouter(db),
	}

	payload := LoginRequest{
		Username: "nonexistent",
		Password: "password",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.loginHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestGenerateToken(t *testing.T) {
	token, err := generateToken(1, "testuser")

	if err != nil {
		t.Errorf("Failed to generate token: %v", err)
	}

	if token == "" {
		t.Error("Token should not be empty")
	}
}

func TestVerifyToken(t *testing.T) {
	// Generate token
	token, _ := generateToken(123, "testuser")

	// Verify token
	userID, username, err := verifyToken(token)

	if err != nil {
		t.Errorf("Failed to verify token: %v", err)
	}

	if userID != 123 {
		t.Errorf("Expected userID 123, got %d", userID)
	}

	if username != "testuser" {
		t.Errorf("Expected username testuser, got %s", username)
	}
}

func TestVerifyTokenInvalid(t *testing.T) {
	_, _, err := verifyToken("invalid-token")

	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

// Helper functions
func initTestDatabase() (*gorm.DB, error) {
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=chat_app_test sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Use in-memory SQLite for testing
		db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		return db, nil
	}
	return db, nil
}

func setupTestRouter(db *gorm.DB) *chi.Mux {
	router := chi.NewRouter()
	server := &Server{db: db, router: router}
	server.setupRoutes()
	return router
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}
