package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// NotificationChannel represents different notification channels
type NotificationChannel struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"` // "email", "sms", "push", "webhook"
	IsActive    bool      `json:"is_active"`
	Config      datatypes.JSON `gorm:"type:jsonb" json:"config"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NotificationSchedule schedules notifications for later delivery
type NotificationSchedule struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	NotificationID uint      `json:"notification_id"`
	Notification   Notification `gorm:"foreignKey:NotificationID" json:"notification,omitempty"`
	ScheduledFor   time.Time `json:"scheduled_for"`
	Sent           bool      `json:"sent"`
	SentAt         time.Time `json:"sent_at"`
	CreatedAt      time.Time `json:"created_at"`
}

// NotificationBatch allows bulk sending of notifications
type NotificationBatch struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `json:"name"`
	TotalCount    int       `json:"total_count"`
	SentCount     int       `json:"sent_count"`
	FailedCount   int       `json:"failed_count"`
	Status        string    `json:"status"` // "pending", "in_progress", "completed"
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// NotificationLog tracks notification history
type NotificationLog struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	NotificationID uint      `json:"notification_id"`
	UserID         uint      `json:"user_id"`
	Channel        string    `json:"channel"`
	Status         string    `json:"status"` // "sent", "failed", "pending"
	ErrorMessage   string    `json:"error_message"`
	SentAt         time.Time `json:"sent_at"`
	ReadAt         time.Time `json:"read_at"`
	CreatedAt      time.Time `json:"created_at"`
}

// GetNotificationChannelsHandler retrieves available notification channels
func (s *Server) getNotificationChannelsHandler(w http.ResponseWriter, r *http.Request) {
	var channels []NotificationChannel
	s.db.Where("is_active = ?", true).Find(&channels)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(channels)
}

// CreateNotificationChannelHandler creates a new notification channel
func (s *Server) createNotificationChannelHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req NotificationChannel
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := s.db.Create(&req).Error; err != nil {
		http.Error(w, "Failed to create notification channel", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

// ScheduleNotificationHandler schedules a notification for later delivery
func (s *Server) scheduleNotificationHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		NotificationID uint      `json:"notification_id"`
		ScheduledFor   time.Time `json:"scheduled_for"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	schedule := NotificationSchedule{
		NotificationID: req.NotificationID,
		ScheduledFor:   req.ScheduledFor,
	}

	if err := s.db.Create(&schedule).Error; err != nil {
		http.Error(w, "Failed to schedule notification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedule)
}

// GetScheduledNotificationsHandler retrieves scheduled notifications
func (s *Server) getScheduledNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var schedules []NotificationSchedule
	s.db.Where("sent = ?", false).Find(&schedules)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schedules)
}

// CreateNotificationBatchHandler creates a batch notification job
func (s *Server) createNotificationBatchHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name       string `json:"name"`
		TotalCount int    `json:"total_count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	batch := NotificationBatch{
		Name:       req.Name,
		TotalCount: req.TotalCount,
		Status:     "pending",
	}

	if err := s.db.Create(&batch).Error; err != nil {
		http.Error(w, "Failed to create notification batch", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(batch)
}

// GetNotificationBatchesHandler retrieves notification batches
func (s *Server) getNotificationBatchesHandler(w http.ResponseWriter, r *http.Request) {
	var batches []NotificationBatch
	s.db.Order("created_at DESC").Limit(50).Find(&batches)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(batches)
}

// GetNotificationLogsHandler retrieves notification logs
func (s *Server) getNotificationLogsHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var logs []NotificationLog
	s.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(100).Find(&logs)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// GetNotificationStatsHandler retrieves notification statistics
func (s *Server) getNotificationStatsHandler(w http.ResponseWriter, r *http.Request) {
	var totalSent int64
	var totalFailed int64
	var totalPending int64

	s.db.Model(&NotificationLog{}).Where("status = ?", "sent").Count(&totalSent)
	s.db.Model(&NotificationLog{}).Where("status = ?", "failed").Count(&totalFailed)
	s.db.Model(&NotificationLog{}).Where("status = ?", "pending").Count(&totalPending)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_sent":    totalSent,
		"total_failed":  totalFailed,
		"total_pending": totalPending,
		"success_rate":  float64(totalSent) / float64(totalSent+totalFailed),
	})
}

// CancelScheduledNotificationHandler cancels a scheduled notification
func (s *Server) cancelScheduledNotificationHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	scheduleID := chi.URLParam(r, "scheduleID")
	var schedule NotificationSchedule
	if err := s.db.Where("id = ?", scheduleID).First(&schedule).Error; err != nil {
		http.Error(w, "Schedule not found", http.StatusNotFound)
		return
	}

	s.db.Delete(&schedule)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Scheduled notification cancelled"})
}
