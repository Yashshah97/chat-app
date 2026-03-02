package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GET /api/users/{id}/profile - Get user profile
func (s *Server) getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var profile UserProfile
	result := s.db.Where("user_id = ?", uid).First(&profile)
	if result.Error != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

// PUT /api/users/{id}/profile - Update user profile
func (s *Server) updateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Bio        string  `json:"bio"`
		Avatar     string  `json:"avatar"`
		CoverPhoto string  `json:"cover_photo"`
		Location   string  `json:"location"`
		Website    string  `json:"website"`
		BirthDate  *string `json:"birth_date"`
		TimeZone   string  `json:"time_zone"`
		Language   string  `json:"language"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result := s.db.Model(&UserProfile{}).Where("user_id = ?", uid).Updates(reqBody)
	if result.Error != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// GET /api/users/{id}/status - Get user status
func (s *Server) getUserStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var status UserStatus
	result := s.db.Where("user_id = ?", uid).First(&status)
	if result.Error != nil {
		http.Error(w, "Status not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// PUT /api/users/{id}/status - Update user status
func (s *Server) updateUserStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Status         string `json:"status"`
		StatusMessage  string `json:"status_message"`
		Device         string `json:"device"`
		Browser        string `json:"browser"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result := s.db.Model(&UserStatus{}).Where("user_id = ?", uid).Updates(reqBody)
	if result.Error != nil {
		http.Error(w, "Failed to update status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// POST /api/saved-searches - Save a search
func (s *Server) saveSearchHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Name    string `json:"name"`
		Query   string `json:"query"`
		Filters string `json:"filters"`
		IsPublic bool  `json:"is_public"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	search := SavedSearch{
		UserID:   uint(userID),
		Name:     reqBody.Name,
		Query:    reqBody.Query,
		Filters:  reqBody.Filters,
		IsPublic: reqBody.IsPublic,
	}

	result := s.db.Create(&search)
	if result.Error != nil {
		http.Error(w, "Failed to save search", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(search)
}

// GET /api/users/{id}/saved-searches - Get saved searches
func (s *Server) getSavedSearchesHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var searches []SavedSearch
	result := s.db.
		Where("user_id = ?", uid).
		Order("created_at DESC").
		Find(&searches)

	if result.Error != nil {
		http.Error(w, "Failed to fetch searches", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(searches)
}

// GET /api/search/advanced - Advanced search with filters
func (s *Server) advancedSearchWithIndexHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	docType := r.URL.Query().Get("type")

	searchQuery := s.db
	if query != "" {
		searchQuery = searchQuery.Where("content ILIKE ?", "%"+query+"%")
	}
	if docType != "" {
		searchQuery = searchQuery.Where("document_type = ?", docType)
	}

	var results []SearchIndex
	result := searchQuery.
		Order("created_at DESC").
		Limit(50).
		Find(&results)

	if result.Error != nil {
		http.Error(w, "Search failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

// POST /api/content-policies - Create content policy
func (s *Server) createContentPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Rules       string `json:"rules"`
		Category    string `json:"category"`
		Severity    int    `json:"severity"`
		Actions     string `json:"actions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	policy := ContentPolicy{
		Name:        reqBody.Name,
		Description: reqBody.Description,
		Rules:       reqBody.Rules,
		Category:    reqBody.Category,
		Severity:    reqBody.Severity,
		Actions:     reqBody.Actions,
	}

	result := s.db.Create(&policy)
	if result.Error != nil {
		http.Error(w, "Failed to create policy", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(policy)
}

// GET /api/content-policies - List content policies
func (s *Server) listContentPoliciesHandler(w http.ResponseWriter, r *http.Request) {
	var policies []ContentPolicy
	result := s.db.
		Where("is_active = ?", true).
		Find(&policies)

	if result.Error != nil {
		http.Error(w, "Failed to fetch policies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(policies)
}
