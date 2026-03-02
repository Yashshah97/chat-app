package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// PinnedMessage represents a pinned message in a chat
type PinnedMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	MessageID uint      `json:"message_id"`
	Message   Message   `gorm:"foreignKey:MessageID" json:"message,omitempty"`
	ChatID    uint      `json:"chat_id"`
	Chat      Chat      `gorm:"foreignKey:ChatID" json:"chat,omitempty"`
	PinnedBy  uint      `json:"pinned_by"`
	PinnedUser User     `gorm:"foreignKey:PinnedBy" json:"pinned_user,omitempty"`
	Reason    string    `json:"reason"`
}

// TableName defines custom table name
func (PinnedMessage) TableName() string {
	return "pinned_messages"
}

// POST /api/messages/{id}/pin - Pin a message
func (s *Server) pinMessageHandler(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "id")
	mID, err := strconv.ParseUint(messageID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	
	var reqBody struct {
		ChatID  uint   `json:"chat_id"`
		PinnedBy uint  `json:"pinned_by"`
		Reason  string `json:"reason"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Get the message first to ensure it exists and belongs to the chat
	var message Message
	if result := s.db.First(&message, mID); result.Error != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}
	
	pinnedMsg := PinnedMessage{
		MessageID: uint(mID),
		ChatID:    reqBody.ChatID,
		PinnedBy:  reqBody.PinnedBy,
		Reason:    reqBody.Reason,
	}
	
	result := s.db.Create(&pinnedMsg)
	if result.Error != nil {
		http.Error(w, "Failed to pin message", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pinnedMsg)
}

// DELETE /api/messages/{id}/pin - Unpin a message
func (s *Server) unpinMessageHandler(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "id")
	mID, err := strconv.ParseUint(messageID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	
	result := s.db.Where("message_id = ?", mID).Delete(&PinnedMessage{})
	
	if result.Error != nil {
		http.Error(w, "Failed to unpin message", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Message unpinned successfully"})
}

// GET /api/chats/{id}/pinned - Get pinned messages in a chat
func (s *Server) getPinnedMessagesHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cID, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	
	var pinnedMsgs []PinnedMessage
	result := s.db.
		Where("chat_id = ?", cID).
		Preload("Message").
		Preload("Message.User").
		Preload("PinnedUser").
		Order("created_at DESC").
		Find(&pinnedMsgs)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch pinned messages", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pinnedMsgs)
}

// GET /api/messages/{id}/pin-status - Check if message is pinned
func (s *Server) checkPinStatusHandler(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "id")
	mID, err := strconv.ParseUint(messageID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	
	var pinnedMsg PinnedMessage
	result := s.db.Where("message_id = ?", mID).First(&pinnedMsg)
	
	isPinned := result.Error == nil
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"is_pinned": isPinned,
		"pinned_message": pinnedMsg,
	})
}
