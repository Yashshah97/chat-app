package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// POST /api/notifications - Create a notification
func (s *Server) createNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		UserID   uint   `json:"user_id"`
		Type     string `json:"type"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		RefID    *uint  `json:"ref_id,omitempty"`
		RefType  string `json:"ref_type"`
		Priority string `json:"priority"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	notification := Notification{
		UserID:   reqBody.UserID,
		Type:     reqBody.Type,
		Title:    reqBody.Title,
		Content:  reqBody.Content,
		RefID:    reqBody.RefID,
		RefType:  reqBody.RefType,
		Priority: reqBody.Priority,
		IsRead:   false,
	}
	
	result := s.db.Create(&notification)
	if result.Error != nil {
		http.Error(w, "Failed to create notification", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(notification)
}

// GET /api/users/{id}/notifications - Get user's notifications
func (s *Server) getUserNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	unreadOnly := r.URL.Query().Get("unread_only") == "true"
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}
	
	var notifications []Notification
	query := s.db.Where("user_id = ?", id)
	
	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}
	
	result := query.
		Order("created_at DESC").
		Limit(limit).
		Find(&notifications)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch notifications", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notifications)
}

// PUT /api/notifications/{id}/read - Mark notification as read
func (s *Server) markNotificationReadHandler(w http.ResponseWriter, r *http.Request) {
	notifID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(notifID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}
	
	now := time.Now()
	result := s.db.Model(&Notification{}, id).Updates(map[string]interface{}{
		"is_read": true,
		"read_at": now,
	})
	
	if result.Error != nil {
		http.Error(w, "Failed to update notification", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Notification marked as read"})
}

// POST /api/users/{id}/notifications/mark-all-read - Mark all as read
func (s *Server) markAllNotificationsReadHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	now := time.Now()
	result := s.db.
		Where("user_id = ? AND is_read = ?", id, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
		})
	
	if result.Error != nil {
		http.Error(w, "Failed to update notifications", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":         "All notifications marked as read",
		"updated_count":   result.RowsAffected,
	})
}

// DELETE /api/notifications/{id} - Delete a notification
func (s *Server) deleteNotificationHandler(w http.ResponseWriter, r *http.Request) {
	notifID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(notifID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}
	
	result := s.db.Delete(&Notification{}, id)
	
	if result.Error != nil {
		http.Error(w, "Failed to delete notification", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Notification deleted successfully"})
}

// GET /api/users/{id}/notifications/unread-count - Get unread notification count
func (s *Server) getUnreadNotificationCountHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	var count int64
	s.db.Model(&Notification{}).Where("user_id = ? AND is_read = ?", id, false).Count(&count)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int64{"unread_count": count})
}
