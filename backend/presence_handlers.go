package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// POST /api/presence/update - Update user presence status
func (s *Server) updatePresenceHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		UserID     uint   `json:"user_id"`
		Status     string `json:"status"` // online, offline, away, dnd
		ChatID     *uint  `json:"chat_id,omitempty"`
		DeviceType string `json:"device_type"`
		IPAddress  string `json:"ip_address"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Update or create presence record
	presence := UserPresence{
		UserID:      reqBody.UserID,
		Status:      reqBody.Status,
		LastSeen:    time.Now(),
		CurrentChatID: reqBody.ChatID,
		OnlineAtTime: time.Now(),
		DeviceType:  reqBody.DeviceType,
		IPAddress:   reqBody.IPAddress,
	}
	
	result := s.db.Where("user_id = ?", reqBody.UserID).Save(&presence)
	
	if result.Error != nil {
		http.Error(w, "Failed to update presence", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(presence)
}

// GET /api/presence/user/{id} - Get user presence status
func (s *Server) getUserPresenceHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	var presence UserPresence
	result := s.db.Where("user_id = ?", id).First(&presence)
	
	if result.Error != nil {
		http.Error(w, "Presence not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(presence)
}

// GET /api/presence/chat/{id}/members - Get online members in a chat
func (s *Server) getChatMembersPresenceHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	
	var presences []UserPresence
	result := s.db.Where("current_chat_id = ? AND status = ?", id, "online").Find(&presences)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch presence data", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(presences)
}

// GET /api/presence/online - Get all online users count
func (s *Server) getOnlineUsersCountHandler(w http.ResponseWriter, r *http.Request) {
	var count int64
	s.db.Model(&UserPresence{}).Where("status = ?", "online").Count(&count)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int64{"online_users": count})
}

// POST /api/presence/set-away - Set user as away
func (s *Server) setUserAwayHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		UserID uint `json:"user_id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	result := s.db.Model(&UserPresence{}, reqBody.UserID).Updates(map[string]interface{}{
		"status":    "away",
		"last_seen": time.Now(),
	})
	
	if result.Error != nil {
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User set to away"})
}

// POST /api/presence/history - Log presence history
func (s *Server) logPresenceHistoryHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		UserID       uint   `json:"user_id"`
		Status       string `json:"status"`
		DurationSecs int64  `json:"duration_secs"`
		ChatID       *uint  `json:"chat_id,omitempty"`
		EventType    string `json:"event_type"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	history := PresenceHistory{
		UserID:       reqBody.UserID,
		Status:       reqBody.Status,
		DurationSecs: reqBody.DurationSecs,
		ChatID:       reqBody.ChatID,
		EventType:    reqBody.EventType,
	}
	
	result := s.db.Create(&history)
	if result.Error != nil {
		http.Error(w, "Failed to log history", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(history)
}
