package main

import (
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// PerformanceAlert triggers when metrics exceed thresholds
type PerformanceAlert struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Name         string         `json:"name"`
	MetricName   string         `json:"metric_name"`
	Threshold    float64        `json:"threshold"`
	Operator     string         `json:"operator"` // "greater_than", "less_than", "equals"
	AlertLevel   string         `json:"alert_level"` // "warning", "critical"
	IsActive     bool           `json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// CPUMetric tracks CPU usage
type CPUMetric struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Usage     float64   `json:"usage"` // percentage
	Cores     int       `json:"cores"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

// MemoryMetric tracks memory usage
type MemoryMetric struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Used      int64     `json:"used"` // bytes
	Total     int64     `json:"total"`
	Available int64     `json:"available"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

// DatabaseMetric tracks database performance
type DatabaseMetric struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	QueryCount      int       `json:"query_count"`
	AvgQueryTime    int64     `json:"avg_query_time"` // milliseconds
	SlowQueryCount  int       `json:"slow_query_count"`
	ConnectionPool  int       `json:"connection_pool"`
	Timestamp       time.Time `json:"timestamp"`
	CreatedAt       time.Time `json:"created_at"`
}

// ServiceMetric tracks individual service performance
type ServiceMetric struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	ServiceName      string         `json:"service_name"`
	RequestCount     int            `json:"request_count"`
	ErrorCount       int            `json:"error_count"`
	AvgResponseTime  int64          `json:"avg_response_time"`
	P95ResponseTime  int64          `json:"p95_response_time"`
	P99ResponseTime  int64          `json:"p99_response_time"`
	MaxResponseTime  int64          `json:"max_response_time"`
	Uptime           int            `json:"uptime"` // seconds
	HealthStatus     string         `json:"health_status"`
	DependencyStatus datatypes.JSON `gorm:"type:jsonb" json:"dependency_status"`
	Timestamp        time.Time      `json:"timestamp"`
	CreatedAt        time.Time      `json:"created_at"`
}

// GetCPUMetricsHandler retrieves CPU metrics
func (s *Server) getCPUMetricsHandler(w http.ResponseWriter, r *http.Request) {
	hours := 24
	startTime := time.Now().Add(-time.Duration(hours) * time.Hour)

	var metrics []CPUMetric
	s.db.Where("timestamp >= ?", startTime).Order("timestamp DESC").Find(&metrics)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// GetMemoryMetricsHandler retrieves memory metrics
func (s *Server) getMemoryMetricsHandler(w http.ResponseWriter, r *http.Request) {
	hours := 24
	startTime := time.Now().Add(-time.Duration(hours) * time.Hour)

	var metrics []MemoryMetric
	s.db.Where("timestamp >= ?", startTime).Order("timestamp DESC").Find(&metrics)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// GetDatabaseMetricsHandler retrieves database performance metrics
func (s *Server) getDatabaseMetricsHandler(w http.ResponseWriter, r *http.Request) {
	var metrics []DatabaseMetric
	s.db.Order("timestamp DESC").Limit(100).Find(&metrics)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// GetServiceMetricsHandler retrieves service metrics
func (s *Server) getServiceMetricsHandler(w http.ResponseWriter, r *http.Request) {
	var metrics []ServiceMetric
	s.db.Order("timestamp DESC").Find(&metrics)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// CreatePerformanceAlertHandler creates a new performance alert
func (s *Server) createPerformanceAlertHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req PerformanceAlert
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := s.db.Create(&req).Error; err != nil {
		http.Error(w, "Failed to create alert", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

// GetPerformanceAlertsHandler retrieves active alerts
func (s *Server) getPerformanceAlertsHandler(w http.ResponseWriter, r *http.Request) {
	var alerts []PerformanceAlert
	s.db.Where("is_active = ?", true).Find(&alerts)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alerts)
}

// GetSystemPerformanceHandler returns overall system performance
func (s *Server) getSystemPerformanceHandler(w http.ResponseWriter, r *http.Request) {
	// Get latest CPU metric
	var cpuMetric CPUMetric
	s.db.Order("timestamp DESC").First(&cpuMetric)

	// Get latest memory metric
	var memMetric MemoryMetric
	s.db.Order("timestamp DESC").First(&memMetric)

	// Get latest database metric
	var dbMetric DatabaseMetric
	s.db.Order("timestamp DESC").First(&dbMetric)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"cpu":      cpuMetric,
		"memory":   memMetric,
		"database": dbMetric,
		"timestamp": time.Now(),
	})
}
