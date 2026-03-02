package main

import (
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// FilterRule defines reusable filtering rules
type FilterRule struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	UserID    uint           `json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Conditions datatypes.JSON `gorm:"type:jsonb" json:"conditions"`
	Actions   datatypes.JSON `gorm:"type:jsonb" json:"actions"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// MessageFilter provides advanced message filtering
type MessageFilter struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Keyword      string    `json:"keyword"`
	From         string    `json:"from"`
	To           string    `json:"to"`
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	HasAttachment bool     `json:"has_attachment"`
	MinLength    int       `json:"min_length"`
	MaxLength    int       `json:"max_length"`
	CreatedAt    time.Time `json:"created_at"`
}

// CreateFilterRuleHandler creates a new filter rule
func (s *Server) createFilterRuleHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name       string      `json:"name"`
		Conditions interface{} `json:"conditions"`
		Actions    interface{} `json:"actions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	condJSON, _ := json.Marshal(req.Conditions)
	actionJSON, _ := json.Marshal(req.Actions)

	rule := FilterRule{
		Name:       req.Name,
		UserID:     userID,
		Conditions: condJSON,
		Actions:    actionJSON,
		IsActive:   true,
	}

	if err := s.db.Create(&rule).Error; err != nil {
		http.Error(w, "Failed to create filter rule", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rule)
}

// GetFilterRulesHandler retrieves all filter rules for user
func (s *Server) getFilterRulesHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var rules []FilterRule
	s.db.Where("user_id = ?", userID).Find(&rules)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rules)
}

// CreateMessageFilterHandler creates an advanced message filter
func (s *Server) createMessageFilterHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req MessageFilter
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	req.UserID = userID
	if err := s.db.Create(&req).Error; err != nil {
		http.Error(w, "Failed to create message filter", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

// ApplyMessageFilterHandler applies message filter and returns results
func (s *Server) applyMessageFilterHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req MessageFilter
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var messages []Message
	query := s.db.Where("1=1")

	if req.Keyword != "" {
		query = query.Where("body ILIKE ?", "%"+req.Keyword+"%")
	}

	if req.StartDate.After(time.Time{}) {
		query = query.Where("created_at >= ?", req.StartDate)
	}

	if req.EndDate.After(time.Time{}) {
		query = query.Where("created_at <= ?", req.EndDate)
	}

	if req.MinLength > 0 {
		query = query.Where("LENGTH(body) >= ?", req.MinLength)
	}

	if req.MaxLength > 0 {
		query = query.Where("LENGTH(body) <= ?", req.MaxLength)
	}

	query.Order("created_at DESC").Limit(100).Find(&messages)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// GetMessageFiltersHandler retrieves all message filters for user
func (s *Server) getMessageFiltersHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var filters []MessageFilter
	s.db.Where("user_id = ?", userID).Find(&filters)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filters)
}
