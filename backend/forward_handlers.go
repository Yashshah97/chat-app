package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// ForwardedMessage represents a forwarded message
type ForwardedMessage struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	OriginalID      uint      `json:"original_id"`
	OriginalMessage Message   `gorm:"foreignKey:OriginalID" json:"original_message,omitempty"`
	ForwardedChatID uint      `json:"forwarded_chat_id"`
	ForwardedChat   Chat      `gorm:"foreignKey:ForwardedChatID" json:"forwarded_chat,omitempty"`
	ForwardedBy     uint      `json:"forwarded_by"`
	ForwardedUser   User      `gorm:"foreignKey:ForwardedBy" json:"forwarded_user,omitempty"`
	CustomNote      string    `json:"custom_note"` // Optional note added while forwarding
}

// TableName defines custom table name
func (ForwardedMessage) TableName() string {
	return "forwarded_messages"
}

// POST /api/messages/{id}/forward - Forward a message to another chat
func (s *Server) forwardMessageHandler(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "id")
	mID, err := strconv.ParseUint(messageID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	
	var reqBody struct {
		TargetChatID uint   `json:"target_chat_id"`
		ForwardedBy  uint   `json:"forwarded_by"`
		CustomNote   string `json:"custom_note"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Get original message
	var message Message
	if result := s.db.First(&message, mID); result.Error != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}
	
	forwarded := ForwardedMessage{
		OriginalID:      uint(mID),
		ForwardedChatID: reqBody.TargetChatID,
		ForwardedBy:     reqBody.ForwardedBy,
		CustomNote:      reqBody.CustomNote,
	}
	
	result := s.db.Create(&forwarded)
	if result.Error != nil {
		http.Error(w, "Failed to forward message", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(forwarded)
}

// GET /api/chats/{id}/forwarded - Get all forwarded messages in a chat
func (s *Server) getForwardedMessagesHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cID, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	
	var forwarded []ForwardedMessage
	result := s.db.
		Where("forwarded_chat_id = ?", cID).
		Preload("OriginalMessage").
		Preload("OriginalMessage.User").
		Preload("OriginalMessage.Chat").
		Preload("ForwardedUser").
		Order("created_at DESC").
		Find(&forwarded)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch forwarded messages", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(forwarded)
}

// POST /api/messages/{id}/forward-to-multiple - Forward to multiple chats
func (s *Server) forwardToMultipleHandler(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "id")
	mID, err := strconv.ParseUint(messageID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	
	var reqBody struct {
		TargetChatIDs []uint `json:"target_chat_ids"`
		ForwardedBy   uint   `json:"forwarded_by"`
		CustomNote    string `json:"custom_note"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	var message Message
	if result := s.db.First(&message, mID); result.Error != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}
	
	// Create forwarded messages for each target chat
	var forwardedMessages []ForwardedMessage
	for _, chatID := range reqBody.TargetChatIDs {
		forwarded := ForwardedMessage{
			OriginalID:      uint(mID),
			ForwardedChatID: chatID,
			ForwardedBy:     reqBody.ForwardedBy,
			CustomNote:      reqBody.CustomNote,
		}
		forwardedMessages = append(forwardedMessages, forwarded)
	}
	
	result := s.db.CreateInBatches(forwardedMessages, len(forwardedMessages))
	if result.Error != nil {
		http.Error(w, "Failed to forward message", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Message forwarded successfully",
		"count":    len(forwardedMessages),
		"forwarded": forwardedMessages,
	})
}
