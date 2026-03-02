package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// MarkAsReadRequest represents a request to mark message as read
type MarkAsReadRequest struct {
	MessageID uint `json:"message_id"`
}

// readReceiptHandlers creates read receipt HTTP handlers
func (s *Server) readReceiptHandlers() {
	readReceiptService := NewReadReceiptService(s.db)

	// Mark message as read
	s.router.With(authMiddleware).Post("/api/messages/{messageID}/read", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Mark message as read request")

		messageID, err := strconv.ParseUint(chi.URLParam(r, "messageID"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid message ID", http.StatusBadRequest)
			return
		}

		// Get user ID from token (simplified)
		userID := uint(1) // This should be extracted from JWT

		receipt, err := readReceiptService.MarkAsRead(uint(messageID), userID)
		if err != nil {
			http.Error(w, "Error marking message as read", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(receipt)
	})

	// Get read receipts for a message
	s.router.Get("/api/messages/{messageID}/read-receipts", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Get read receipts request")

		messageID, err := strconv.ParseUint(chi.URLParam(r, "messageID"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid message ID", http.StatusBadRequest)
			return
		}

		receipts, err := readReceiptService.GetReadReceipts(uint(messageID))
		if err != nil {
			http.Error(w, "Error fetching read receipts", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(receipts)
	})

	// Get unread message count
	s.router.With(authMiddleware).Get("/api/messages/unread/count", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Get unread message count request")

		// Get user ID from token (simplified)
		userID := uint(1) // This should be extracted from JWT

		count, err := readReceiptService.GetUnreadMessages(userID)
		if err != nil {
			http.Error(w, "Error fetching unread count", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"unread_count": count,
		})
	})

	// Mark entire chat as read
	s.router.With(authMiddleware).Post("/api/chats/{chatID}/mark-read", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Mark chat as read request")

		chatID, err := strconv.ParseUint(chi.URLParam(r, "chatID"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid chat ID", http.StatusBadRequest)
			return
		}

		// Get user ID from token (simplified)
		userID := uint(1) // This should be extracted from JWT

		err = readReceiptService.MarkChatAsRead(uint(chatID), userID)
		if err != nil {
			http.Error(w, "Error marking chat as read", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"message": "Chat marked as read",
		})
	})
}
