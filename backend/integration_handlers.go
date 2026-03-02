package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// POST /api/integrations - Create integration
func (s *Server) createIntegrationHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Name       string `json:"name"`
		Type       string `json:"type"`
		APIKey     string `json:"api_key"`
		APISecret  string `json:"api_secret"`
		WebhookURL string `json:"webhook_url"`
		Config     string `json:"config"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	integration := Integration{
		Name:        reqBody.Name,
		Type:        reqBody.Type,
		APIKey:      reqBody.APIKey,
		APISecret:   reqBody.APISecret,
		WebhookURL:  reqBody.WebhookURL,
		Config:      reqBody.Config,
		CreatedByID: uint(userID),
	}

	result := s.db.Create(&integration)
	if result.Error != nil {
		http.Error(w, "Failed to create integration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(integration)
}

// GET /api/integrations - List integrations
func (s *Server) listIntegrationsHandler(w http.ResponseWriter, r *http.Request) {
	intType := r.URL.Query().Get("type")

	query := s.db
	if intType != "" {
		query = query.Where("type = ?", intType)
	}

	var integrations []Integration
	result := query.Find(&integrations)
	if result.Error != nil {
		http.Error(w, "Failed to fetch integrations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(integrations)
}

// GET /api/integrations/{id} - Get integration details
func (s *Server) getIntegrationHandler(w http.ResponseWriter, r *http.Request) {
	integrationID := chi.URLParam(r, "id")
	iid, err := strconv.ParseUint(integrationID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid integration ID", http.StatusBadRequest)
		return
	}

	var integration Integration
	result := s.db.First(&integration, iid)
	if result.Error != nil {
		http.Error(w, "Integration not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(integration)
}

// PUT /api/integrations/{id} - Update integration
func (s *Server) updateIntegrationHandler(w http.ResponseWriter, r *http.Request) {
	integrationID := chi.URLParam(r, "id")
	iid, err := strconv.ParseUint(integrationID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid integration ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Name      string `json:"name"`
		IsActive  bool   `json:"is_active"`
		Config    string `json:"config"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result := s.db.Model(&Integration{}).Where("id = ?", iid).Updates(reqBody)
	if result.Error != nil {
		http.Error(w, "Failed to update integration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// DELETE /api/integrations/{id} - Delete integration
func (s *Server) deleteIntegrationHandler(w http.ResponseWriter, r *http.Request) {
	integrationID := chi.URLParam(r, "id")
	iid, err := strconv.ParseUint(integrationID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid integration ID", http.StatusBadRequest)
		return
	}

	result := s.db.Delete(&Integration{}, iid)
	if result.Error != nil {
		http.Error(w, "Failed to delete integration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// POST /api/integrations/{id}/mappings - Create integration mapping
func (s *Server) createIntegrationMappingHandler(w http.ResponseWriter, r *http.Request) {
	integrationID := chi.URLParam(r, "id")
	iid, err := strconv.ParseUint(integrationID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid integration ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		ChatID       uint   `json:"chat_id"`
		ExternalID   string `json:"external_id"`
		ExternalName string `json:"external_name"`
		SyncMessages bool   `json:"sync_messages"`
		SyncUsers    bool   `json:"sync_users"`
		Direction    string `json:"direction"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	mapping := IntegrationMapping{
		IntegrationID: uint(iid),
		ChatID:        reqBody.ChatID,
		ExternalID:    reqBody.ExternalID,
		ExternalName:  reqBody.ExternalName,
		SyncMessages:  reqBody.SyncMessages,
		SyncUsers:     reqBody.SyncUsers,
		Direction:     reqBody.Direction,
	}

	result := s.db.Create(&mapping)
	if result.Error != nil {
		http.Error(w, "Failed to create mapping", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(mapping)
}

// GET /api/integrations/{id}/mappings - Get integration mappings
func (s *Server) getIntegrationMappingsHandler(w http.ResponseWriter, r *http.Request) {
	integrationID := chi.URLParam(r, "id")
	iid, err := strconv.ParseUint(integrationID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid integration ID", http.StatusBadRequest)
		return
	}

	var mappings []IntegrationMapping
	result := s.db.
		Where("integration_id = ?", iid).
		Preload("Chat").
		Find(&mappings)

	if result.Error != nil {
		http.Error(w, "Failed to fetch mappings", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mappings)
}

// GET /api/integrations/{id}/events - Get integration events
func (s *Server) getIntegrationEventsHandler(w http.ResponseWriter, r *http.Request) {
	integrationID := chi.URLParam(r, "id")
	iid, err := strconv.ParseUint(integrationID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid integration ID", http.StatusBadRequest)
		return
	}

	var events []IntegrationEvent
	result := s.db.
		Where("integration_id = ?", iid).
		Order("created_at DESC").
		Limit(100).
		Find(&events)

	if result.Error != nil {
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

// POST /api/integrations/{id}/sync - Trigger integration sync
func (s *Server) triggerIntegrationSyncHandler(w http.ResponseWriter, r *http.Request) {
	integrationID := chi.URLParam(r, "id")
	iid, err := strconv.ParseUint(integrationID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid integration ID", http.StatusBadRequest)
		return
	}

	result := s.db.Model(&Integration{}).Where("id = ?", iid).Update("last_sync_at", "now()")
	if result.Error != nil {
		http.Error(w, "Failed to trigger sync", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "sync_triggered"})
}

// GET /api/integrations/types/available - Get available integration types
func (s *Server) getAvailableIntegrationTypesHandler(w http.ResponseWriter, r *http.Request) {
	types := []map[string]string{
		{"name": "slack", "description": "Slack integration"},
		{"name": "discord", "description": "Discord integration"},
		{"name": "github", "description": "GitHub integration"},
		{"name": "jira", "description": "Jira integration"},
		{"name": "trello", "description": "Trello integration"},
		{"name": "microsoft_teams", "description": "Microsoft Teams integration"},
		{"name": "telegram", "description": "Telegram integration"},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types)
}
