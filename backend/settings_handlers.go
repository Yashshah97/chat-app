package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// GET /api/chats/{id}/settings - Get chat settings
func (s *Server) getChatSettingsHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	
	var settings ChatSettings
	result := s.db.Where("chat_id = ?", id).First(&settings)
	
	// If settings don't exist, create defaults
	if result.Error != nil {
		settings = ChatSettings{
			ChatID:             uint(id),
			AllowNotifications: true,
			AllowMentions:      true,
			AllowFileSharing:   true,
			ReadReceiptEnabled: true,
			TypingIndicator:    true,
		}
		s.db.Create(&settings)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(settings)
}

// PUT /api/chats/{id}/settings - Update chat settings
func (s *Server) updateChatSettingsHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	
	var req ChatSettings
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	result := s.db.Model(&ChatSettings{}, id).Updates(req)
	
	if result.Error != nil {
		http.Error(w, "Failed to update settings", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Settings updated successfully"})
}

// GET /api/users/{id}/preferences/{chatID} - Get user preferences for a chat
func (s *Server) getUserChatPreferencesHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	chatID := chi.URLParam(r, "chatID")
	
	uID, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	cID, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	
	var prefs UserChatPreference
	result := s.db.Where("user_id = ? AND chat_id = ?", uID, cID).First(&prefs)
	
	if result.Error != nil {
		http.Error(w, "Preferences not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prefs)
}

// PUT /api/users/{id}/preferences/{chatID} - Update user preferences for a chat
func (s *Server) updateUserChatPreferencesHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	chatID := chi.URLParam(r, "chatID")
	
	uID, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	cID, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	
	var req UserChatPreference
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	result := s.db.Where("user_id = ? AND chat_id = ?", uID, cID).Updates(req)
	
	if result.Error != nil {
		http.Error(w, "Failed to update preferences", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Preferences updated successfully"})
}

// GET /api/users/{id}/notifications - Get user notification preferences
func (s *Server) getNotificationPreferencesHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	var prefs NotificationPreference
	result := s.db.Where("user_id = ?", id).First(&prefs)
	
	// Create defaults if doesn't exist
	if result.Error != nil {
		prefs = NotificationPreference{
			UserID:          uint(id),
			AllowPush:       true,
			AllowEmail:      true,
			AllowSound:      true,
			AllowVibration:  true,
			QuietHoursStart: "22:00",
			QuietHoursEnd:   "08:00",
		}
		s.db.Create(&prefs)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prefs)
}

// PUT /api/users/{id}/notifications - Update notification preferences
func (s *Server) updateNotificationPreferencesHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	var req NotificationPreference
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	result := s.db.Where("user_id = ?", id).Updates(req)
	
	if result.Error != nil {
		http.Error(w, "Failed to update preferences", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Notification preferences updated"})
}

// POST /api/chats/{id}/mute/{userID} - Mute chat for user
func (s *Server) muteChatForUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	chatID := chi.URLParam(r, "id")
	
	uID, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	cID, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	
	var duration int // in minutes, optional
	json.NewDecoder(r.Body).Decode(&map[string]interface{}{"duration": &duration})
	
	var muteUntil *time.Time
	if duration > 0 {
		t := time.Now().Add(time.Duration(duration) * time.Minute)
		muteUntil = &t
	}
	
	prefs := UserChatPreference{
		UserID:    uint(uID),
		ChatID:    uint(cID),
		IsMuted:   true,
		MuteUntil: muteUntil,
	}
	
	result := s.db.Where("user_id = ? AND chat_id = ?", uID, cID).Updates(prefs)
	
	if result.Error != nil {
		http.Error(w, "Failed to mute chat", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Chat muted successfully",
		"muted_until": muteUntil,
	})
}
