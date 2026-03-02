package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// PUT /api/messages/{id} - Edit a message (creates history entry)
func (s *Server) editMessageHandler(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "id")
	mID, err := strconv.ParseUint(messageID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	
	var reqBody struct {
		NewContent string `json:"content"`
		EditReason string `json:"edit_reason"`
		UserID     uint   `json:"user_id"`
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
	
	// Create edit history entry
	edit := MessageEdit{
		MessageID:  uint(mID),
		UserID:     reqBody.UserID,
		OldContent: message.Content,
		NewContent: reqBody.NewContent,
		EditReason: reqBody.EditReason,
	}
	
	if result := s.db.Create(&edit); result.Error != nil {
		http.Error(w, "Failed to record edit history", http.StatusInternalServerError)
		return
	}
	
	// Update message
	message.Content = reqBody.NewContent
	message.IsEdited = true
	
	if result := s.db.Save(&message); result.Error != nil {
		http.Error(w, "Failed to update message", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": message,
		"edit":    edit,
	})
}

// GET /api/messages/{id}/history - Get message edit history
func (s *Server) getMessageHistoryHandler(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "id")
	mID, err := strconv.ParseUint(messageID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	
	var edits []MessageEdit
	result := s.db.
		Where("message_id = ?", mID).
		Preload("User").
		Order("created_at DESC").
		Find(&edits)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch edit history", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(edits)
}

// GET /api/messages/{id}/edit-count - Get number of edits
func (s *Server) getEditCountHandler(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "id")
	mID, err := strconv.ParseUint(messageID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	
	var count int64
	s.db.Model(&MessageEdit{}).Where("message_id = ?", mID).Count(&count)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message_id": mID,
		"edit_count": count,
	})
}

// DELETE /api/message-edits/{id} - Delete an edit history entry (admin only)
func (s *Server) deleteEditHistoryHandler(w http.ResponseWriter, r *http.Request) {
	editID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(editID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid edit ID", http.StatusBadRequest)
		return
	}
	
	result := s.db.Delete(&MessageEdit{}, id)
	
	if result.Error != nil {
		http.Error(w, "Failed to delete edit history", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Edit history deleted successfully"})
}

// GET /api/chats/{id}/edited-messages - Get all edited messages in a chat
func (s *Server) getEditedMessagesHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cID, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	
	var messages []Message
	result := s.db.
		Where("chat_id = ? AND is_edited = ?", cID, true).
		Preload("User").
		Order("updated_at DESC").
		Find(&messages)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch edited messages", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}
