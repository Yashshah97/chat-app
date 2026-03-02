package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// POST /api/users/{id}/block/{targetID} - Block a user
func (s *Server) blockUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	targetID := chi.URLParam(r, "targetID")
	
	uID, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	tID, err := strconv.ParseUint(targetID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid target ID", http.StatusBadRequest)
		return
	}
	
	var reqBody struct {
		Reason string `json:"reason"`
	}
	
	json.NewDecoder(r.Body).Decode(&reqBody)
	
	blocked := BlockedUser{
		BlockerID: uint(uID),
		BlockedID: uint(tID),
		Reason:    reqBody.Reason,
	}
	
	result := s.db.Create(&blocked)
	if result.Error != nil {
		http.Error(w, "Failed to block user", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(blocked)
}

// DELETE /api/users/{id}/block/{targetID} - Unblock a user
func (s *Server) unblockUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	targetID := chi.URLParam(r, "targetID")
	
	uID, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	tID, err := strconv.ParseUint(targetID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid target ID", http.StatusBadRequest)
		return
	}
	
	result := s.db.Where("blocker_id = ? AND blocked_id = ?", uID, tID).Delete(&BlockedUser{})
	
	if result.Error != nil {
		http.Error(w, "Failed to unblock user", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User unblocked successfully"})
}

// GET /api/users/{id}/blocked - Get list of blocked users
func (s *Server) getBlockedUsersHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	var blocked []BlockedUser
	result := s.db.
		Where("blocker_id = ?", id).
		Preload("Blocked").
		Find(&blocked)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch blocked users", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(blocked)
}

// POST /api/users/{id}/mute/{targetID} - Mute a user
func (s *Server) muteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	targetID := chi.URLParam(r, "targetID")
	
	uID, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	tID, err := strconv.ParseUint(targetID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid target ID", http.StatusBadRequest)
		return
	}
	
	var reqBody struct {
		DurationMinutes int    `json:"duration_minutes"` // 0 = permanent
		Reason          string `json:"reason"`
	}
	
	json.NewDecoder(r.Body).Decode(&reqBody)
	
	var muteUntil *time.Time
	if reqBody.DurationMinutes > 0 {
		t := time.Now().Add(time.Duration(reqBody.DurationMinutes) * time.Minute)
		muteUntil = &t
	}
	
	muted := MutedUser{
		MuterID:   uint(uID),
		MutedID:   uint(tID),
		MuteUntil: muteUntil,
		Reason:    reqBody.Reason,
	}
	
	result := s.db.Create(&muted)
	if result.Error != nil {
		http.Error(w, "Failed to mute user", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(muted)
}

// DELETE /api/users/{id}/mute/{targetID} - Unmute a user
func (s *Server) unmuteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	targetID := chi.URLParam(r, "targetID")
	
	uID, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	tID, err := strconv.ParseUint(targetID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid target ID", http.StatusBadRequest)
		return
	}
	
	result := s.db.Where("muter_id = ? AND muted_id = ?", uID, tID).Delete(&MutedUser{})
	
	if result.Error != nil {
		http.Error(w, "Failed to unmute user", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User unmuted successfully"})
}

// GET /api/users/{id}/muted - Get list of muted users
func (s *Server) getMutedUsersHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	var muted []MutedUser
	result := s.db.
		Where("muter_id = ? AND (mute_until IS NULL OR mute_until > NOW())", id).
		Preload("Muted").
		Find(&muted)
	
	if result.Error != nil {
		http.Error(w, "Failed to fetch muted users", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(muted)
}

// GET /api/users/{id}/is-blocked-by/{targetID} - Check if user is blocked by another user
func (s *Server) checkBlockStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	targetID := chi.URLParam(r, "targetID")
	
	uID, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	tID, err := strconv.ParseUint(targetID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid target ID", http.StatusBadRequest)
		return
	}
	
	var blocked BlockedUser
	result := s.db.Where("blocker_id = ? AND blocked_id = ?", uID, tID).First(&blocked)
	
	isBlocked := result.Error == nil
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"is_blocked": isBlocked,
	})
}
