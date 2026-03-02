package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// POST /api/rate-limit-config - Create rate limit config
func (s *Server) createRateLimitConfigHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		UserID              *uint  `json:"user_id"`
		ChatID              *uint  `json:"chat_id"`
		Endpoint            string `json:"endpoint"`
		RequestsPerMinute   int    `json:"requests_per_minute"`
		RequestsPerHour     int    `json:"requests_per_hour"`
		RequestsPerDay      int    `json:"requests_per_day"`
		Priority            int    `json:"priority"`
		WhitelistIPs        string `json:"whitelist_ips"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	config := RateLimitConfig{
		UserID:              reqBody.UserID,
		ChatID:              reqBody.ChatID,
		Endpoint:            reqBody.Endpoint,
		RequestsPerMinute:   reqBody.RequestsPerMinute,
		RequestsPerHour:     reqBody.RequestsPerHour,
		RequestsPerDay:      reqBody.RequestsPerDay,
		Priority:            reqBody.Priority,
		WhitelistIPs:        reqBody.WhitelistIPs,
	}

	result := s.db.Create(&config)
	if result.Error != nil {
		http.Error(w, "Failed to create config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
}

// GET /api/rate-limit-config - List rate limit configs
func (s *Server) listRateLimitConfigsHandler(w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Query().Get("endpoint")

	query := s.db
	if endpoint != "" {
		query = query.Where("endpoint = ?", endpoint)
	}

	var configs []RateLimitConfig
	result := query.Find(&configs)
	if result.Error != nil {
		http.Error(w, "Failed to fetch configs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(configs)
}

// GET /api/rate-limit-usage - Get rate limit usage
func (s *Server) getRateLimitUsageHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	endpoint := r.URL.Query().Get("endpoint")

	query := s.db
	if userID != "" {
		if uid, err := strconv.ParseUint(userID, 10, 32); err == nil {
			query = query.Where("user_id = ?", uid)
		}
	}
	if endpoint != "" {
		query = query.Where("endpoint = ?", endpoint)
	}

	var usage []RateLimitUsage
	result := query.Find(&usage)
	if result.Error != nil {
		http.Error(w, "Failed to fetch usage", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usage)
}

// GET /api/rate-limit-violations - Get violations
func (s *Server) getRateLimitViolationsHandler(w http.ResponseWriter, r *http.Request) {
	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	var violations []RateLimitViolation
	result := s.db.
		Order("created_at DESC").
		Limit(limit).
		Find(&violations)

	if result.Error != nil {
		http.Error(w, "Failed to fetch violations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(violations)
}

// POST /api/api-versions - Create API version
func (s *Server) createAPIVersionHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Version        string  `json:"version"`
		Status         string  `json:"status"`
		ReleasedAt     string  `json:"released_at"`
		SunsetDate     *string `json:"sunset_date"`
		Description    string  `json:"description"`
		Changes        string  `json:"changes"`
		MigrationGuide string  `json:"migration_guide"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	version := APIVersion{
		Version:        reqBody.Version,
		Status:         reqBody.Status,
		ReleasedAt:     reqBody.ReleasedAt,
		SunsetDate:     reqBody.SunsetDate,
		Description:    reqBody.Description,
		Changes:        reqBody.Changes,
		MigrationGuide: reqBody.MigrationGuide,
	}

	result := s.db.Create(&version)
	if result.Error != nil {
		http.Error(w, "Failed to create version", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(version)
}

// GET /api/api-versions - List API versions
func (s *Server) listAPIVersionsHandler(w http.ResponseWriter, r *http.Request) {
	var versions []APIVersion
	result := s.db.Order("created_at DESC").Find(&versions)
	if result.Error != nil {
		http.Error(w, "Failed to fetch versions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(versions)
}

// GET /api/api-versions/{version} - Get version details
func (s *Server) getAPIVersionHandler(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")

	var apiVersion APIVersion
	result := s.db.Where("version = ?", version).First(&apiVersion)
	if result.Error != nil {
		http.Error(w, "Version not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiVersion)
}

// POST /api/api-endpoints - Register endpoint
func (s *Server) registerAPIEndpointHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Path            string  `json:"path"`
		Method          string  `json:"method"`
		Version         string  `json:"version"`
		Description     string  `json:"description"`
		Authentication  bool    `json:"authentication"`
		RateLimit       int     `json:"rate_limit"`
		IsDeprecated    bool    `json:"is_deprecated"`
		ReplacementPath *string `json:"replacement_path"`
		Examples        string  `json:"examples"`
		Parameters      string  `json:"parameters"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	endpoint := APIEndpoint{
		Path:            reqBody.Path,
		Method:          reqBody.Method,
		Version:         reqBody.Version,
		Description:     reqBody.Description,
		Authentication:  reqBody.Authentication,
		RateLimit:       reqBody.RateLimit,
		IsDeprecated:    reqBody.IsDeprecated,
		ReplacementPath: reqBody.ReplacementPath,
		Examples:        reqBody.Examples,
		Parameters:      reqBody.Parameters,
	}

	result := s.db.Create(&endpoint)
	if result.Error != nil {
		http.Error(w, "Failed to register endpoint", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(endpoint)
}

// GET /api/api-endpoints - List endpoints
func (s *Server) listAPIEndpointsHandler(w http.ResponseWriter, r *http.Request) {
	version := r.URL.Query().Get("version")

	query := s.db
	if version != "" {
		query = query.Where("version = ?", version)
	}

	var endpoints []APIEndpoint
	result := query.Find(&endpoints)
	if result.Error != nil {
		http.Error(w, "Failed to fetch endpoints", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(endpoints)
}

// GET /api/api-docs - Get API documentation
func (s *Server) getAPIDocsHandler(w http.ResponseWriter, r *http.Request) {
	version := r.URL.Query().Get("version")

	var endpoints []APIEndpoint
	query := s.db
	if version != "" {
		query = query.Where("version = ?", version)
	}
	query.Find(&endpoints)

	docs := map[string]interface{}{
		"version":   version,
		"endpoints": endpoints,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(docs)
}
