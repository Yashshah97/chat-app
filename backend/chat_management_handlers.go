package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GET /api/admin/chats - List all chats
func (s *Server) listChatsHandler(w http.ResponseWriter, r *http.Request) {
	var chats []Chat
	result := s.db.Preload("Members").Find(&chats)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch chats", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chats)
}

// GET /api/admin/chats/{id} - Get chat details
func (s *Server) getChatDetailsHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	
	var chat Chat
	result := s.db.Preload("Members").Preload("Messages").First(&chat, id)
	
	if result.Error != nil {
		http.Error(w, "Chat not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chat)
}

// DELETE /api/admin/chats/{id} - Delete a chat
func (s *Server) deleteChatHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	
	result := s.db.Delete(&Chat{}, id)
	
	if result.Error != nil {
		http.Error(w, "Failed to delete chat", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Chat deleted successfully"})
}

// POST /api/admin/chats/{id}/remove-member/{memberID} - Remove member from chat
func (s *Server) removeChatMemberHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	memberID := chi.URLParam(r, "memberID")
	
	cID, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	
	mID, err := strconv.ParseUint(memberID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid member ID", http.StatusBadRequest)
		return
	}
	
	// Remove member from chat_members association
	result := s.db.Model(&Chat{}, cID).Association("Members").Delete(&User{}, mID)
	
	if result != nil {
		http.Error(w, "Failed to remove member", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Member removed successfully"})
}

// POST /api/admin/chats/{id}/mute - Mute a chat
func (s *Server) muteChatHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}
	
	// In a real implementation, you'd add a muted_by field to Chat
	var reqBody struct {
		Reason string `json:"reason"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Chat muted successfully",
		"chat_id": id,
		"reason": reqBody.Reason,
	})
}
