package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// AddReactionRequest represents a request to add a reaction
type AddReactionRequest struct {
	MessageID uint   `json:"message_id"`
	Emoji     string `json:"emoji"`
}

// TypingStatusRequest represents a typing status update
type TypingStatusRequest struct {
	IsTyping bool `json:"is_typing"`
}

// reactionHandlers creates reaction-related HTTP handlers
func (s *Server) reactionHandlers() {
	reactionService := NewMessageReactionService(s.db)

	// Add reaction
	s.router.With(authMiddleware).Post("/api/messages/{messageID}/reactions", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Add reaction request from %s", r.RemoteAddr)

		messageID, err := strconv.ParseUint(chi.URLParam(r, "messageID"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid message ID", http.StatusBadRequest)
			return
		}

		var req AddReactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get user ID from token (simplified)
		userID := uint(1) // This should be extracted from JWT

		reaction, err := reactionService.AddReaction(uint(messageID), userID, req.Emoji)
		if err != nil {
			http.Error(w, "Error adding reaction", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(reaction)
	})

	// Get reactions
	s.router.Get("/api/messages/{messageID}/reactions", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Get reactions request")

		messageID, err := strconv.ParseUint(chi.URLParam(r, "messageID"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid message ID", http.StatusBadRequest)
			return
		}

		reactions, err := reactionService.GetMessageReactions(uint(messageID))
		if err != nil {
			http.Error(w, "Error fetching reactions", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(reactions)
	})

	// Remove reaction
	s.router.With(authMiddleware).Delete("/api/reactions/{reactionID}", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Remove reaction request")

		reactionID, err := strconv.ParseUint(chi.URLParam(r, "reactionID"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid reaction ID", http.StatusBadRequest)
			return
		}

		err = reactionService.RemoveReaction(uint(reactionID))
		if err != nil {
			http.Error(w, "Error removing reaction", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

// typingHandlers creates typing indicator HTTP handlers
func (s *Server) typingHandlers() {
	typingService := NewTypingIndicatorService(s.db)

	// Update typing status
	s.router.With(authMiddleware).Post("/api/chats/{chatID}/typing", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Update typing status request")

		chatID, err := strconv.ParseUint(chi.URLParam(r, "chatID"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}

		var req TypingStatusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get user ID from token (simplified)
		userID := uint(1) // This should be extracted from JWT

		indicator, err := typingService.UpdateTypingStatus(uint(chatID), userID, req.IsTyping)
		if err != nil {
			http.Error(w, "Error updating typing status", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(indicator)
	})

	// Get typing users
	s.router.Get("/api/chats/{chatID}/typing", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Get typing users request")

		chatID, err := strconv.ParseUint(chi.URLParam(r, "chatID"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}

		indicators, err := typingService.GetTypingUsers(uint(chatID))
		if err != nil {
			http.Error(w, "Error fetching typing users", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(indicators)
	})
}
