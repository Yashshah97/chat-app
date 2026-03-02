package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// POST /api/chats/{id}/calls - Initiate a call
func (s *Server) initiateCallHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		CallType string `json:"call_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	call := VoiceCall{
		ChatID:      uint(cid),
		InitiatorID: uint(userID),
		Status:      "ringing",
		CallType:    reqBody.CallType,
	}

	result := s.db.Create(&call)
	if result.Error != nil {
		http.Error(w, "Failed to initiate call", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(call)
}

// GET /api/calls/{id} - Get call details
func (s *Server) getCallDetailsHandler(w http.ResponseWriter, r *http.Request) {
	callID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(callID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid call ID", http.StatusBadRequest)
		return
	}

	var call VoiceCall
	result := s.db.
		Preload("Chat").
		Preload("Initiator").
		First(&call, cid)

	if result.Error != nil {
		http.Error(w, "Call not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(call)
}

// POST /api/calls/{id}/join - Join a call
func (s *Server) joinCallHandler(w http.ResponseWriter, r *http.Request) {
	callID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(callID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid call ID", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	participant := CallParticipant{
		CallID:   uint(cid),
		UserID:   uint(userID),
		JoinedAt: "",
	}

	result := s.db.Create(&participant)
	if result.Error != nil {
		http.Error(w, "Failed to join call", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(participant)
}

// POST /api/calls/{id}/end - End a call
func (s *Server) endCallHandler(w http.ResponseWriter, r *http.Request) {
	callID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(callID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid call ID", http.StatusBadRequest)
		return
	}

	result := s.db.Model(&VoiceCall{}).Where("id = ?", cid).Update("status", "ended")
	if result.Error != nil {
		http.Error(w, "Failed to end call", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ended"})
}

// GET /api/calls/{id}/participants - Get call participants
func (s *Server) getCallParticipantsHandler(w http.ResponseWriter, r *http.Request) {
	callID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(callID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid call ID", http.StatusBadRequest)
		return
	}

	var participants []CallParticipant
	result := s.db.
		Where("call_id = ?", cid).
		Preload("User").
		Find(&participants)

	if result.Error != nil {
		http.Error(w, "Failed to fetch participants", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(participants)
}

// GET /api/chats/{id}/call-history - Get call history for chat
func (s *Server) getCallHistoryHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var calls []VoiceCall
	result := s.db.
		Where("chat_id = ?", cid).
		Order("created_at DESC").
		Limit(50).
		Find(&calls)

	if result.Error != nil {
		http.Error(w, "Failed to fetch call history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(calls)
}

// POST /api/call-logs - Create call log entry
func (s *Server) createCallLogHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		CallID    uint   `json:"call_id"`
		EventType string `json:"event_type"`
		UserID    *uint  `json:"user_id"`
		Metadata  string `json:"metadata"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log := CallLog{
		CallID:    reqBody.CallID,
		EventType: reqBody.EventType,
		UserID:    reqBody.UserID,
		Metadata:  reqBody.Metadata,
	}

	result := s.db.Create(&log)
	if result.Error != nil {
		http.Error(w, "Failed to create log", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(log)
}

// GET /api/calls/{id}/logs - Get call logs
func (s *Server) getCallLogsHandler(w http.ResponseWriter, r *http.Request) {
	callID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(callID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid call ID", http.StatusBadRequest)
		return
	}

	var logs []CallLog
	result := s.db.
		Where("call_id = ?", cid).
		Order("created_at ASC").
		Find(&logs)

	if result.Error != nil {
		http.Error(w, "Failed to fetch logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logs)
}
