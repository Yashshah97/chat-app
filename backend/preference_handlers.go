package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GET /api/users/{id}/preferences - Get user preferences
func (s *Server) getUserPreferencesHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var prefs UserPreference
	result := s.db.Where("user_id = ?", uid).First(&prefs)
	if result.Error != nil {
		http.Error(w, "Preferences not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prefs)
}

// PUT /api/users/{id}/preferences - Update user preferences
func (s *Server) updateUserPreferencesHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Theme           string `json:"theme"`
		Language        string `json:"language"`
		Privacy         string `json:"privacy"`
		Notifications   bool   `json:"notifications"`
		EmailDigest     bool   `json:"email_digest"`
		ShowActivity    bool   `json:"show_activity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result := s.db.Model(&UserPreference{}).Where("user_id = ?", uid).Updates(reqBody)
	if result.Error != nil {
		http.Error(w, "Failed to update preferences", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// POST /api/surveys - Create survey
func (s *Server) createSurveyHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	survey := Survey{
		Title:       reqBody.Title,
		Description: reqBody.Description,
		CreatedByID: uint(userID),
		StartDate:   reqBody.StartDate,
		EndDate:     reqBody.EndDate,
	}

	result := s.db.Create(&survey)
	if result.Error != nil {
		http.Error(w, "Failed to create survey", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(survey)
}

// GET /api/surveys - List surveys
func (s *Server) listSurveysHandler(w http.ResponseWriter, r *http.Request) {
	var surveys []Survey
	result := s.db.Order("created_at DESC").Find(&surveys)
	if result.Error != nil {
		http.Error(w, "Failed to fetch surveys", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(surveys)
}

// POST /api/surveys/{id}/responses - Submit survey response
func (s *Server) submitSurveyResponseHandler(w http.ResponseWriter, r *http.Request) {
	surveyID := chi.URLParam(r, "id")
	sid, err := strconv.ParseUint(surveyID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid survey ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Answers string `json:"answers"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	response := SurveyResponse{
		SurveyID: uint(sid),
		UserID:   (*uint)(&[]uint{uint(userID)}[0]),
		Answers:  reqBody.Answers,
	}

	result := s.db.Create(&response)
	if result.Error != nil {
		http.Error(w, "Failed to submit response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GET /api/users/{id}/recommendations - Get recommendations for user
func (s *Server) getUserRecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	var recommendations []Recommendation
	result := s.db.
		Where("user_id = ?", uid).
		Order("score DESC").
		Limit(limit).
		Find(&recommendations)

	if result.Error != nil {
		http.Error(w, "Failed to fetch recommendations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(recommendations)
}

// POST /api/recommendations/{id}/click - Track recommendation click
func (s *Server) trackRecommendationClickHandler(w http.ResponseWriter, r *http.Request) {
	recID := chi.URLParam(r, "id")
	rid, err := strconv.ParseUint(recID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid recommendation ID", http.StatusBadRequest)
		return
	}

	result := s.db.Model(&Recommendation{}).Where("id = ?", rid).Update("clicked", true)
	if result.Error != nil {
		http.Error(w, "Failed to track click", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "tracked"})
}

// GET /api/surveys/{id}/responses - Get survey responses
func (s *Server) getSurveyResponsesHandler(w http.ResponseWriter, r *http.Request) {
	surveyID := chi.URLParam(r, "id")
	sid, err := strconv.ParseUint(surveyID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid survey ID", http.StatusBadRequest)
		return
	}

	var responses []SurveyResponse
	result := s.db.
		Where("survey_id = ?", sid).
		Order("created_at DESC").
		Find(&responses)

	if result.Error != nil {
		http.Error(w, "Failed to fetch responses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses)
}
