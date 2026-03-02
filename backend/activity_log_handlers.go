package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// POST /api/audit/log - Log a user activity
func (s *Server) logUserActivityHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		UserID       uint   `json:"user_id"`
		Action       string `json:"action"`
		ResourceType string `json:"resource_type"`
		ResourceID   *uint  `json:"resource_id,omitempty"`
		Metadata     string `json:"metadata"`
		IPAddress    string `json:"ip_address"`
		UserAgent    string `json:"user_agent"`
		Status       string `json:"status"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	activity := UserActivityLog{
		UserID:       reqBody.UserID,
		Action:       reqBody.Action,
		ResourceType: reqBody.ResourceType,
		ResourceID:   reqBody.ResourceID,
		Metadata:     reqBody.Metadata,
		IPAddress:    reqBody.IPAddress,
		UserAgent:    reqBody.UserAgent,
		Status:       reqBody.Status,
	}
	
	result := s.db.Create(&activity)
	if result.Error != nil {
		http.Error(w, "Failed to log activity", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(activity)
}

// GET /api/users/{id}/activity - Get user activity logs
func (s *Server) getUserActivityLogsHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 200 {
			limit = parsed
		}
	}
	
	var activities []UserActivityLog
	result := s.db.
		Where("user_id = ?", id).
		Order("created_at DESC").
		Limit(limit).
		Find(&activities)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch activity logs", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activities)
}

// GET /api/audit/logs - Get all activity logs (admin only)
func (s *Server) getAllActivityLogsHandler(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	resourceType := r.URL.Query().Get("resource_type")
	limit := 100
	
	query := s.db
	
	if action != "" {
		query = query.Where("action = ?", action)
	}
	
	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}
	
	var activities []UserActivityLog
	result := query.
		Order("created_at DESC").
		Limit(limit).
		Find(&activities)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch activity logs", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activities)
}

// GET /api/audit/stats - Get activity statistics
func (s *Server) getActivityStatsHandler(w http.ResponseWriter, r *http.Request) {
	var stats struct {
		TotalLogs      int64
		UniqueUsers    int64
		FailedAttempts int64
		LoginCount     int64
		MessageCount   int64
	}
	
	s.db.Model(&UserActivityLog{}).Count(&stats.TotalLogs)
	s.db.Model(&UserActivityLog{}).Distinct("user_id").Count(&stats.UniqueUsers)
	s.db.Model(&UserActivityLog{}).Where("status = ?", "failure").Count(&stats.FailedAttempts)
	s.db.Model(&UserActivityLog{}).Where("action = ?", "login").Count(&stats.LoginCount)
	s.db.Model(&UserActivityLog{}).Where("action = ?", "create_message").Count(&stats.MessageCount)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

// DELETE /api/audit/logs/{id} - Delete activity log entry
func (s *Server) deleteActivityLogHandler(w http.ResponseWriter, r *http.Request) {
	logID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(logID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid log ID", http.StatusBadRequest)
		return
	}
	
	result := s.db.Delete(&UserActivityLog{}, id)
	
	if result.Error != nil {
		http.Error(w, "Failed to delete log", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Activity log deleted successfully"})
}

// GET /api/audit/suspicious - Get suspicious activities
func (s *Server) getSuspiciousActivityHandler(w http.ResponseWriter, r *http.Request) {
	var activities []UserActivityLog
	
	// Find failed login attempts and deletions
	result := s.db.
		Where("status = ? OR action = ?", "failure", "delete_message").
		Order("created_at DESC").
		Limit(50).
		Find(&activities)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch suspicious activities", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activities)
}
