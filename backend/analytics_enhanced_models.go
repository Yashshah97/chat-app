package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AnalyticsMetric tracks system metrics
type AnalyticsMetric struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	MetricKey string         `gorm:"index" json:"metric_key"`
	Value     float64        `json:"value"`
	Labels    datatypes.JSON `gorm:"type:jsonb" json:"labels"`
	Timestamp time.Time      `json:"timestamp"`
	CreatedAt time.Time      `json:"created_at"`
}

// PerformanceMetric tracks performance data
type PerformanceMetric struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Endpoint      string    `json:"endpoint"`
	Method        string    `json:"method"`
	ResponseTime  int64     `json:"response_time"` // milliseconds
	StatusCode    int       `json:"status_code"`
	UserID        uint      `json:"user_id"`
	RequestSize   int64     `json:"request_size"`
	ResponseSize  int64     `json:"response_size"`
	Timestamp     time.Time `json:"timestamp"`
	CreatedAt     time.Time `json:"created_at"`
}

// UsageStatistics tracks usage patterns
type UsageStatistics struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	UserID              uint      `json:"user_id"`
	User                User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ActiveMinutesToday  int       `json:"active_minutes_today"`
	MessagesCreated     int       `json:"messages_created"`
	ChatsAccessed       int       `json:"chats_accessed"`
	FilesUploaded       int       `json:"files_uploaded"`
	StorageUsed         int64     `json:"storage_used"`
	LastActivityAt      time.Time `json:"last_activity_at"`
	Date                time.Time `json:"date"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// ErrorMetric tracks errors and exceptions
type ErrorMetric struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ErrorType string         `json:"error_type"`
	Message   string         `json:"message"`
	Stack     string         `json:"stack"`
	Context   datatypes.JSON `gorm:"type:jsonb" json:"context"`
	Resolved  bool           `json:"resolved"`
	Count     int            `json:"count"`
	FirstSeen time.Time      `json:"first_seen"`
	LastSeen  time.Time      `json:"last_seen"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// RecordAnalyticsMetricHandler records a metric
func (s *Server) recordAnalyticsMetricHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MetricKey string                 `json:"metric_key"`
		Value     float64                `json:"value"`
		Labels    map[string]interface{} `json:"labels"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	labelsJSON, _ := json.Marshal(req.Labels)
	metric := AnalyticsMetric{
		MetricKey: req.MetricKey,
		Value:     req.Value,
		Labels:    labelsJSON,
		Timestamp: time.Now(),
	}

	if err := s.db.Create(&metric).Error; err != nil {
		http.Error(w, "Failed to record metric", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metric)
}

// GetAnalyticsMetricsHandler retrieves metrics for a given key
func (s *Server) getAnalyticsMetricsHandler(w http.ResponseWriter, r *http.Request) {
	metricKey := chi.URLParam(r, "metricKey")
	hours := 24 // default

	var metrics []AnalyticsMetric
	startTime := time.Now().Add(-time.Duration(hours) * time.Hour)
	s.db.Where("metric_key = ? AND timestamp >= ?", metricKey, startTime).
		Order("timestamp DESC").Find(&metrics)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// RecordPerformanceMetricHandler records performance data
func (s *Server) recordPerformanceMetricHandler(w http.ResponseWriter, r *http.Request) {
	var req PerformanceMetric
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	req.Timestamp = time.Now()
	if err := s.db.Create(&req).Error; err != nil {
		http.Error(w, "Failed to record performance metric", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

// GetPerformanceStatsHandler retrieves performance statistics
func (s *Server) getPerformanceStatsHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now().Add(-24 * time.Hour)

	var metrics []PerformanceMetric
	s.db.Where("timestamp >= ?", startTime).Find(&metrics)

	// Calculate stats
	stats := map[string]interface{}{
		"total_requests": len(metrics),
		"avg_response_time": 0,
		"p99_response_time": 0,
		"error_rate": 0,
	}

	if len(metrics) > 0 {
		totalTime := int64(0)
		errors := 0
		for _, m := range metrics {
			totalTime += m.ResponseTime
			if m.StatusCode >= 400 {
				errors++
			}
		}
		stats["avg_response_time"] = totalTime / int64(len(metrics))
		stats["error_rate"] = float64(errors) / float64(len(metrics))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetUserUsageStatsHandler retrieves usage statistics for a user
func (s *Server) getUserUsageStatsHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var stats []UsageStatistics
	s.db.Where("user_id = ?", userID).Order("date DESC").Limit(30).Find(&stats)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetErrorMetricsHandler retrieves error metrics
func (s *Server) getErrorMetricsHandler(w http.ResponseWriter, r *http.Request) {
	var errors []ErrorMetric
	s.db.Where("resolved = ?", false).Order("count DESC").Limit(50).Find(&errors)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(errors)
}

// ReportErrorHandler reports an error
func (s *Server) reportErrorHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ErrorType string `json:"error_type"`
		Message   string `json:"message"`
		Stack     string `json:"stack"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var existing ErrorMetric
	result := s.db.Where("error_type = ? AND message = ?", req.ErrorType, req.Message).First(&existing)

	if result.Error == nil {
		// Update existing error
		existing.Count++
		existing.LastSeen = time.Now()
		s.db.Save(&existing)
	} else if result.Error == gorm.ErrRecordNotFound {
		// Create new error
		newError := ErrorMetric{
			ErrorType: req.ErrorType,
			Message:   req.Message,
			Stack:     req.Stack,
			Count:     1,
			FirstSeen: time.Now(),
			LastSeen:  time.Now(),
		}
		s.db.Create(&newError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "error reported"})
}

// GetSystemHealthHandler returns overall system health
func (s *Server) getSystemHealthHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now().Add(-1 * time.Hour)

	// Get error count
	var errorCount int64
	s.db.Model(&ErrorMetric{}).Where("created_at >= ?", startTime).Count(&errorCount)

	// Get active users
	var activeUserCount int64
	s.db.Model(&UsageStatistics{}).Where("last_activity_at >= ?", time.Now().Add(-15*time.Minute)).Distinct("user_id").Count(&activeUserCount)

	// Get recent errors
	var recentErrors []ErrorMetric
	s.db.Where("resolved = ?", false).Order("last_seen DESC").Limit(10).Find(&recentErrors)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":             "healthy",
		"errors_last_hour":   errorCount,
		"active_users":       activeUserCount,
		"recent_errors":      recentErrors,
		"timestamp":          time.Now(),
	})
}
