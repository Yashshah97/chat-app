package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
	router *chi.Mux
}

func main() {
	// Initialize database
	db, err := initDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run migrations
	err = db.AutoMigrate(&User{}, &Message{}, &Chat{})
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize server
	server := &Server{
		db: db,
		router: chi.NewRouter(),
	}

	// Setup routes
	server.setupRoutes()

	// Start server
	port := ":8080"
	log.Printf("Starting server on %s", port)
	if err := http.ListenAndServe(port, server.router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func initDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func (s *Server) setupRoutes() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(corsMiddleware)

	// Auth routes
	s.router.Post("/api/auth/register", s.registerHandler)
	s.router.Post("/api/auth/login", s.loginHandler)

	// Chat routes
	s.router.With(authMiddleware).Get("/api/chats", s.getChatsHandler)
	s.router.With(authMiddleware).Post("/api/chats", s.createChatHandler)
	s.router.With(authMiddleware).Post("/api/chats/group", s.createGroupChatHandler)
	s.router.With(authMiddleware).Get("/api/chats/{id}/messages", s.getChatMessagesHandler)
	
	// WebSocket
	s.router.With(authMiddleware).Get("/ws/chat/{id}", s.websocketHandler)

	// Health check
	s.router.Get("/health", s.healthHandler)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}
		
		// Verify token (simplified for now)
		if len(token) < 7 || token[:7] != "Bearer " {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// Handler stubs
func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"token":"jwt-token","username":"user"}`))
}

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"token":"jwt-token","username":"user"}`))
}

func (s *Server) getChatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`[]`))
}

func (s *Server) createChatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"id":1,"name":"Chat","type":"private"}`))
}

func (s *Server) createGroupChatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"id":1,"name":"Group","type":"group"}`))
}

func (s *Server) getChatMessagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`[]`))
}

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"websocket connection established"}`))
}
