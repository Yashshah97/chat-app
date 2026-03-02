package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GET /api/analytics/system - Get system-wide analytics
func (s *Server) getSystemAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	var analytics SystemAnalytics
	
	// Fetch latest system analytics
	result := s.db.Order("created_at DESC").First(&analytics)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch analytics", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(analytics)
}

// GET /api/analytics/chat/{id} - Get analytics for a specific chat
func (s *Server) getChatAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	
	var analytics ChatAnalytics
	result := s.db.Where("chat_id = ?", id).Order("created_at DESC").First(&analytics)
	
	if result.Error != nil {
		http.Error(w, "Analytics not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(analytics)
}

// GET /api/analytics/user/{id} - Get analytics for a specific user
func (s *Server) getUserAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	var analytics UserAnalytics
	result := s.db.Where("user_id = ?", id).Order("created_at DESC").First(&analytics)
	
	if result.Error != nil {
		http.Error(w, "Analytics not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(analytics)
}

// POST /api/analytics/compute - Compute/update analytics
func (s *Server) computeAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	// Count total users
	var totalUsers int64
	s.db.Model(&User{}).Count(&totalUsers)
	
	// Count total chats
	var totalChats int64
	s.db.Model(&Chat{}).Count(&totalChats)
	
	// Count total messages
	var totalMessages int64
	s.db.Model(&Message{}).Count(&totalMessages)
	
	// Count active users online
	var activeUsersOnline int64
	s.db.Model(&User{}).Where("status = ?", "online").Count(&activeUsersOnline)
	
	systemAnalytics := SystemAnalytics{
		TotalUsers:        totalUsers,
		ActiveUsersOnline: activeUsersOnline,
		TotalChats:        totalChats,
		TotalMessages:     totalMessages,
		MessagesPerMinute: float64(totalMessages) / 1440, // Assuming since creation
	}
	
	result := s.db.Create(&systemAnalytics)
	if result.Error != nil {
		http.Error(w, "Failed to compute analytics", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(systemAnalytics)
}

// GET /api/analytics/dashboard - Get dashboard summary
func (s *Server) getAnalyticsDashboardHandler(w http.ResponseWriter, r *http.Request) {
	var systemAnalytics SystemAnalytics
	s.db.Order("created_at DESC").First(&systemAnalytics)
	
	var allChatAnalytics []ChatAnalytics
	s.db.Order("created_at DESC").Limit(10).Find(&allChatAnalytics)
	
	var allUserAnalytics []UserAnalytics
	s.db.Order("created_at DESC").Limit(10).Find(&allUserAnalytics)
	
	dashboard := map[string]interface{}{
		"system_analytics": systemAnalytics,
		"top_chats":        allChatAnalytics,
		"top_users":        allUserAnalytics,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dashboard)
}
