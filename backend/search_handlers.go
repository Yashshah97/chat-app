package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// SearchRequest represents a message search query
type SearchRequest struct {
	Query  string `json:"query"`
	ChatID *uint  `json:"chat_id,omitempty"`
	UserID *uint  `json:"user_id,omitempty"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// SearchResult wraps message search results with metadata
type SearchResult struct {
	Message   *Message `json:"message"`
	ChatName  string   `json:"chat_name"`
	Username  string   `json:"username"`
	Relevance float64  `json:"relevance"`
}

// POST /api/search/messages - Search for messages
func (s *Server) searchMessagesHandler(w http.ResponseWriter, r *http.Request) {
	var req SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Set defaults
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	
	var messages []Message
	query := s.db
	
	// Build search query
	if req.Query != "" {
		query = query.Where("content ILIKE ?", "%"+req.Query+"%")
	}
	
	if req.ChatID != nil {
		query = query.Where("chat_id = ?", *req.ChatID)
	}
	
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}
	
	result := query.
		Preload("User").
		Preload("Chat").
		Order("created_at DESC").
		Limit(req.Limit).
		Offset(req.Offset).
		Find(&messages)
	
	if result.Error != nil {
		http.Error(w, "Failed to search messages", http.StatusInternalServerError)
		return
	}
	
	// Build results with metadata
	var searchResults []SearchResult
	for _, msg := range messages {
		result := SearchResult{
			Message:  &msg,
			ChatName: msg.Chat.Name,
			Username: msg.User.Username,
		}
		searchResults = append(searchResults, result)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(searchResults)
}

// GET /api/search/messages/chat/{id} - Search messages in a specific chat
func (s *Server) searchChatMessagesHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	
	query := r.URL.Query().Get("q")
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}
	
	var messages []Message
	result := s.db.
		Where("chat_id = ? AND content ILIKE ?", id, "%"+query+"%").
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Find(&messages)
	
	if result.Error != nil {
		http.Error(w, "Failed to search messages", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}

// GET /api/search/users - Search for users
func (s *Server) searchUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}
	
	var users []User
	result := s.db.
		Where("username ILIKE ? OR email ILIKE ?", "%"+query+"%", "%"+query+"%").
		Limit(limit).
		Find(&users)
	
	if result.Error != nil {
		http.Error(w, "Failed to search users", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GET /api/search/chats - Search for chats
func (s *Server) searchChatsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}
	
	var chats []Chat
	result := s.db.
		Where("name ILIKE ?", "%"+query+"%").
		Preload("Members").
		Limit(limit).
		Find(&chats)
	
	if result.Error != nil {
		http.Error(w, "Failed to search chats", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chats)
}
