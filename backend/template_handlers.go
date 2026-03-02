package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// POST /api/message-templates - Create message template
func (s *Server) createMessageTemplateHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Name       string `json:"name"`
		Description string `json:"description"`
		Content    string `json:"content"`
		Category   string `json:"category"`
		ChatID     *uint  `json:"chat_id"`
		IsPublic   bool   `json:"is_public"`
		Tags       string `json:"tags"`
		Variables  string `json:"variables"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	template := MessageTemplate{
		Name:        reqBody.Name,
		Description: reqBody.Description,
		Content:     reqBody.Content,
		Category:    reqBody.Category,
		ChatID:      reqBody.ChatID,
		CreatedByID: uint(userID),
		IsPublic:    reqBody.IsPublic,
		Tags:        reqBody.Tags,
		Variables:   reqBody.Variables,
	}

	result := s.db.Create(&template)
	if result.Error != nil {
		http.Error(w, "Failed to create template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(template)
}

// GET /api/message-templates - List message templates
func (s *Server) listMessageTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	query := s.db
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var templates []MessageTemplate
	result := query.Find(&templates)
	if result.Error != nil {
		http.Error(w, "Failed to fetch templates", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(templates)
}

// GET /api/message-templates/{id} - Get template
func (s *Server) getMessageTemplateHandler(w http.ResponseWriter, r *http.Request) {
	templateID := chi.URLParam(r, "id")
	tid, err := strconv.ParseUint(templateID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid template ID", http.StatusBadRequest)
		return
	}

	var template MessageTemplate
	result := s.db.First(&template, tid)
	if result.Error != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(template)
}

// POST /api/chatbots - Create chatbot
func (s *Server) createChatBotHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ChatID      uint   `json:"chat_id"`
		Webhook     string `json:"webhook"`
		IntentModels string `json:"intent_models"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	bot := ChatBot{
		Name:        reqBody.Name,
		Description: reqBody.Description,
		ChatID:      reqBody.ChatID,
		CreatedByID: uint(userID),
		Webhook:     reqBody.Webhook,
		IntentModels: reqBody.IntentModels,
		Token:       generateToken(),
	}

	result := s.db.Create(&bot)
	if result.Error != nil {
		http.Error(w, "Failed to create chatbot", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bot)
}

// GET /api/chats/{id}/chatbots - List chatbots for chat
func (s *Server) listChatBotsHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var bots []ChatBot
	result := s.db.Where("chat_id = ?", cid).Find(&bots)
	if result.Error != nil {
		http.Error(w, "Failed to fetch chatbots", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bots)
}

// POST /api/chatbots/{id}/intents - Create bot intent
func (s *Server) createBotIntentHandler(w http.ResponseWriter, r *http.Request) {
	botID := chi.URLParam(r, "id")
	bid, err := strconv.ParseUint(botID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid bot ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Pattern  string `json:"pattern"`
		Response string `json:"response"`
		Priority int    `json:"priority"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	intent := BotIntent{
		BotID:    uint(bid),
		Pattern:  reqBody.Pattern,
		Response: reqBody.Response,
		Priority: reqBody.Priority,
	}

	result := s.db.Create(&intent)
	if result.Error != nil {
		http.Error(w, "Failed to create intent", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(intent)
}

// GET /api/chatbots/{id}/intents - Get bot intents
func (s *Server) getBotIntentsHandler(w http.ResponseWriter, r *http.Request) {
	botID := chi.URLParam(r, "id")
	bid, err := strconv.ParseUint(botID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid bot ID", http.StatusBadRequest)
		return
	}

	var intents []BotIntent
	result := s.db.
		Where("bot_id = ?", bid).
		Order("priority DESC").
		Find(&intents)

	if result.Error != nil {
		http.Error(w, "Failed to fetch intents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(intents)
}

// GET /api/chats/{id}/statistics - Get chat statistics
func (s *Server) getChatStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var stats []ChatStatistics
	result := s.db.
		Where("chat_id = ?", cid).
		Order("date DESC").
		Limit(30).
		Find(&stats)

	if result.Error != nil {
		http.Error(w, "Failed to fetch statistics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

// POST /api/chats/{id}/statistics - Record chat statistics
func (s *Server) recordChatStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Date              string `json:"date"`
		MessageCount      int    `json:"message_count"`
		ActiveUsers       int    `json:"active_users"`
		NewUsers          int    `json:"new_users"`
		AverageResponseTime int  `json:"average_response_time"`
		MostActiveHour    int    `json:"most_active_hour"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stat := ChatStatistics{
		ChatID:          uint(cid),
		Date:            reqBody.Date,
		MessageCount:    reqBody.MessageCount,
		ActiveUsers:     reqBody.ActiveUsers,
		NewUsers:        reqBody.NewUsers,
		AverageResponseTime: reqBody.AverageResponseTime,
		MostActiveHour:  reqBody.MostActiveHour,
	}

	result := s.db.Create(&stat)
	if result.Error != nil {
		http.Error(w, "Failed to record statistics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(stat)
}
