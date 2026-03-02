package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// POST /api/badges - Create badge
func (s *Server) createBadgeHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		IconURL     string `json:"icon_url"`
		Category    string `json:"category"`
		Condition   string `json:"condition"`
		Rarity      string `json:"rarity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	badge := Badge{
		Name:        reqBody.Name,
		Description: reqBody.Description,
		IconURL:     reqBody.IconURL,
		Category:    reqBody.Category,
		Condition:   reqBody.Condition,
		Rarity:      reqBody.Rarity,
	}

	result := s.db.Create(&badge)
	if result.Error != nil {
		http.Error(w, "Failed to create badge", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(badge)
}

// GET /api/users/{id}/badges - Get user badges
func (s *Server) getUserBadgesHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var badges []UserBadge
	result := s.db.
		Where("user_id = ?", uid).
		Preload("Badge").
		Order("earned_at DESC").
		Find(&badges)

	if result.Error != nil {
		http.Error(w, "Failed to fetch badges", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(badges)
}

// POST /api/chats/{id}/invitations - Send invitation
func (s *Server) sendInvitationHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Email  string `json:"email"`
		UserID *uint  `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	invitation := Invitation{
		ChatID:      uint(cid),
		InvitedByID: uint(userID),
		InvitedUserID: reqBody.UserID,
		Email:       reqBody.Email,
		InviteToken: generateToken(),
	}

	result := s.db.Create(&invitation)
	if result.Error != nil {
		http.Error(w, "Failed to send invitation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invitation)
}

// GET /api/chats/{id}/invitations - Get chat invitations
func (s *Server) getChatInvitationsHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var invitations []Invitation
	result := s.db.
		Where("chat_id = ?", cid).
		Order("created_at DESC").
		Find(&invitations)

	if result.Error != nil {
		http.Error(w, "Failed to fetch invitations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invitations)
}

// POST /api/invitations/{token}/accept - Accept invitation
func (s *Server) acceptInvitationHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	result := s.db.
		Model(&Invitation{}).
		Where("invite_token = ?", token).
		Update("status", "accepted")

	if result.Error != nil {
		http.Error(w, "Failed to accept invitation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
}

// GET /api/trending - Get trending topics
func (s *Server) getTrendingTopicsHandler(w http.ResponseWriter, r *http.Request) {
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	var topics []TrendingTopic
	result := s.db.
		Order("rank ASC").
		Limit(limit).
		Find(&topics)

	if result.Error != nil {
		http.Error(w, "Failed to fetch trending topics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(topics)
}

// GET /api/statistics/summary - Get summary statistics
func (s *Server) getSummaryStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "day"
	}

	var stats SummaryStatistics
	result := s.db.
		Where("period = ?", period).
		Order("created_at DESC").
		First(&stats)

	if result.Error != nil {
		http.Error(w, "Statistics not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

// POST /api/statistics/summary - Create summary statistics
func (s *Server) createSummaryStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Period                string  `json:"period"`
		StartDate             string  `json:"start_date"`
		EndDate               string  `json:"end_date"`
		TotalUsers            int64   `json:"total_users"`
		ActiveUsers           int64   `json:"active_users"`
		TotalChats            int64   `json:"total_chats"`
		TotalMessages         int64   `json:"total_messages"`
		NewUsers              int64   `json:"new_users"`
		ChatEngagement        float64 `json:"chat_engagement"`
		MessageGrowth         float64 `json:"message_growth"`
		RetentionRate         float64 `json:"retention_rate"`
		AverageSessionLength  int     `json:"average_session_length"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stats := SummaryStatistics{
		Period:              reqBody.Period,
		StartDate:           reqBody.StartDate,
		EndDate:             reqBody.EndDate,
		TotalUsers:          reqBody.TotalUsers,
		ActiveUsers:         reqBody.ActiveUsers,
		TotalChats:          reqBody.TotalChats,
		TotalMessages:       reqBody.TotalMessages,
		NewUsers:            reqBody.NewUsers,
		ChatEngagement:      reqBody.ChatEngagement,
		MessageGrowth:       reqBody.MessageGrowth,
		RetentionRate:       reqBody.RetentionRate,
		AverageSessionLength: reqBody.AverageSessionLength,
	}

	result := s.db.Create(&stats)
	if result.Error != nil {
		http.Error(w, "Failed to create statistics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(stats)
}

// GET /api/badges - Get all badges
func (s *Server) getAllBadgesHandler(w http.ResponseWriter, r *http.Request) {
	var badges []Badge
	result := s.db.Find(&badges)
	if result.Error != nil {
		http.Error(w, "Failed to fetch badges", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(badges)
}
