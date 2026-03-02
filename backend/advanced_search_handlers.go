package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// AdvancedSearchRequest represents an advanced search with filters
type AdvancedSearchRequest struct {
	Query         string     `json:"query"`
	ChatID        *uint      `json:"chat_id,omitempty"`
	UserID        *uint      `json:"user_id,omitempty"`
	FromUser      *uint      `json:"from_user,omitempty"`
	StartDate     *time.Time `json:"start_date,omitempty"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	MessageType   *string    `json:"message_type,omitempty"` // text, image, file
	HasAttachment *bool      `json:"has_attachment,omitempty"`
	IsEdited      *bool      `json:"is_edited,omitempty"`
	Limit         int        `json:"limit"`
	Offset        int        `json:"offset"`
	SortBy        string     `json:"sort_by"` // created_at, relevance
}

// POST /api/search/advanced - Advanced message search with filters
func (s *Server) advancedSearchHandler(w http.ResponseWriter, r *http.Request) {
	var req AdvancedSearchRequest
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
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	
	var messages []Message
	query := s.db
	
	// Build advanced search filters
	if req.Query != "" {
		query = query.Where("content ILIKE ?", "%"+req.Query+"%")
	}
	
	if req.ChatID != nil {
		query = query.Where("chat_id = ?", *req.ChatID)
	}
	
	if req.FromUser != nil {
		query = query.Where("user_id = ?", *req.FromUser)
	}
	
	if req.MessageType != nil {
		query = query.Where("type = ?", *req.MessageType)
	}
	
	if req.IsEdited != nil {
		query = query.Where("is_edited = ?", *req.IsEdited)
	}
	
	if req.StartDate != nil {
		query = query.Where("created_at >= ?", *req.StartDate)
	}
	
	if req.EndDate != nil {
		query = query.Where("created_at <= ?", *req.EndDate)
	}
	
	// Sort
	if req.SortBy == "relevance" {
		query = query.Order("CASE WHEN content ILIKE ? THEN 1 ELSE 2 END, created_at DESC", req.Query+"%")
	} else {
		query = query.Order("created_at DESC")
	}
	
	result := query.
		Preload("User").
		Preload("Chat").
		Limit(req.Limit).
		Offset(req.Offset).
		Find(&messages)
	
	if result.Error != nil {
		http.Error(w, "Failed to search messages", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"results": messages,
		"count":   len(messages),
		"limit":   req.Limit,
		"offset":  req.Offset,
	})
}

// GET /api/search/trending - Get trending messages/topics
func (s *Server) getTrendingHandler(w http.ResponseWriter, r *http.Request) {
	// Get messages from last 24 hours, grouped by word frequency
	var messages []Message
	
	result := s.db.
		Where("created_at >= NOW() - INTERVAL '24 hours'").
		Preload("User").
		Preload("Chat").
		Order("created_at DESC").
		Limit(50).
		Find(&messages)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch trending", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"trending_messages": messages,
		"count":             len(messages),
	})
}

// GET /api/search/history/{userID} - Get user's search history
func (s *Server) getSearchHistoryHandler(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, you would have a search_history table
	// For now, we'll return recent messages by the user
	userID := r.URL.Query().Get("user_id")
	limit := 20
	
	var messages []Message
	result := s.db.
		Where("user_id = ?", userID).
		Preload("User").
		Preload("Chat").
		Order("created_at DESC").
		Limit(limit).
		Find(&messages)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch history", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}
