package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// POST /api/webhooks - Create outgoing webhook
func (s *Server) createWebhookHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Name       string `json:"name"`
		URL        string `json:"url"`
		ChatID     *uint  `json:"chat_id"`
		Events     string `json:"events"`
		Secret     string `json:"secret"`
		Headers    string `json:"headers"`
		RateLimit  int    `json:"rate_limit"`
		Retries    int    `json:"retries"`
		RetryDelay int    `json:"retry_delay"`
		Timeout    int    `json:"timeout"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	webhook := Webhook{
		Name:        reqBody.Name,
		URL:         reqBody.URL,
		ChatID:      reqBody.ChatID,
		CreatedByID: uint(userID),
		Events:      reqBody.Events,
		Secret:      reqBody.Secret,
		Headers:     reqBody.Headers,
		RateLimit:   reqBody.RateLimit,
		Retries:     reqBody.Retries,
		RetryDelay:  reqBody.RetryDelay,
		Timeout:     reqBody.Timeout,
	}

	result := s.db.Create(&webhook)
	if result.Error != nil {
		http.Error(w, "Failed to create webhook", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(webhook)
}

// GET /api/webhooks - List webhooks
func (s *Server) listWebhooksHandler(w http.ResponseWriter, r *http.Request) {
	chatID := r.URL.Query().Get("chat_id")

	query := s.db
	if chatID != "" {
		if cid, err := strconv.ParseUint(chatID, 10, 32); err == nil {
			query = query.Where("chat_id = ?", cid)
		}
	}

	var webhooks []Webhook
	result := query.Find(&webhooks)
	if result.Error != nil {
		http.Error(w, "Failed to fetch webhooks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(webhooks)
}

// GET /api/webhooks/{id} - Get webhook details
func (s *Server) getWebhookHandler(w http.ResponseWriter, r *http.Request) {
	webhookID := chi.URLParam(r, "id")
	wid, err := strconv.ParseUint(webhookID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid webhook ID", http.StatusBadRequest)
		return
	}

	var webhook Webhook
	result := s.db.First(&webhook, wid)
	if result.Error != nil {
		http.Error(w, "Webhook not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(webhook)
}

// PUT /api/webhooks/{id} - Update webhook
func (s *Server) updateWebhookHandler(w http.ResponseWriter, r *http.Request) {
	webhookID := chi.URLParam(r, "id")
	wid, err := strconv.ParseUint(webhookID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid webhook ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Name       string `json:"name"`
		URL        string `json:"url"`
		Events     string `json:"events"`
		IsActive   bool   `json:"is_active"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result := s.db.Model(&Webhook{}).Where("id = ?", wid).Updates(reqBody)
	if result.Error != nil {
		http.Error(w, "Failed to update webhook", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// DELETE /api/webhooks/{id} - Delete webhook
func (s *Server) deleteWebhookHandler(w http.ResponseWriter, r *http.Request) {
	webhookID := chi.URLParam(r, "id")
	wid, err := strconv.ParseUint(webhookID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid webhook ID", http.StatusBadRequest)
		return
	}

	result := s.db.Delete(&Webhook{}, wid)
	if result.Error != nil {
		http.Error(w, "Failed to delete webhook", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// GET /api/webhooks/{id}/events - Get webhook events
func (s *Server) getWebhookEventsHandler(w http.ResponseWriter, r *http.Request) {
	webhookID := chi.URLParam(r, "id")
	wid, err := strconv.ParseUint(webhookID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid webhook ID", http.StatusBadRequest)
		return
	}

	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	var events []WebhookEvent
	result := s.db.
		Where("webhook_id = ?", wid).
		Order("created_at DESC").
		Limit(limit).
		Find(&events)

	if result.Error != nil {
		http.Error(w, "Failed to fetch events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

// GET /api/webhooks/{id}/logs - Get webhook logs
func (s *Server) getWebhookLogsHandler(w http.ResponseWriter, r *http.Request) {
	webhookID := chi.URLParam(r, "id")
	wid, err := strconv.ParseUint(webhookID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid webhook ID", http.StatusBadRequest)
		return
	}

	var logs []WebhookLog
	result := s.db.
		Where("webhook_id = ?", wid).
		Order("created_at DESC").
		Limit(100).
		Find(&logs)

	if result.Error != nil {
		http.Error(w, "Failed to fetch logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logs)
}

// POST /api/chats/{id}/incoming-webhooks - Create incoming webhook
func (s *Server) createIncomingWebhookHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Name       string `json:"name"`
		AllowedIP  string `json:"allowed_ip"`
		Format     string `json:"format"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	incoming := IncomingWebhook{
		ChatID:      uint(cid),
		CreatedByID: uint(userID),
		Name:        reqBody.Name,
		AllowedIP:   reqBody.AllowedIP,
		Format:      reqBody.Format,
		Token:       generateToken(), // Implement token generation
	}

	result := s.db.Create(&incoming)
	if result.Error != nil {
		http.Error(w, "Failed to create incoming webhook", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(incoming)
}

// GET /api/chats/{id}/incoming-webhooks - List incoming webhooks
func (s *Server) listIncomingWebhooksHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var incoming []IncomingWebhook
	result := s.db.
		Where("chat_id = ?", cid).
		Find(&incoming)

	if result.Error != nil {
		http.Error(w, "Failed to fetch incoming webhooks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(incoming)
}

// Helper function to generate webhook tokens
func generateToken() string {
	return "whk_" + randomString(32)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[0]
	}
	return string(result)
}
