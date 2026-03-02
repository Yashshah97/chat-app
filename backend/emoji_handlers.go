package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// POST /api/emoji-packs - Create new emoji pack
func (s *Server) createEmojiPackHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Version     string `json:"version"`
		IconURL     string `json:"icon_url"`
		IsPublic    bool   `json:"is_public"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)
	pack := EmojiPack{
		Name:        reqBody.Name,
		Description: reqBody.Description,
		CreatedByID: uint(userID),
		Version:     reqBody.Version,
		IconURL:     reqBody.IconURL,
		IsPublic:    reqBody.IsPublic,
	}

	result := s.db.Create(&pack)
	if result.Error != nil {
		http.Error(w, "Failed to create emoji pack", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pack)
}

// GET /api/emoji-packs - List all public emoji packs
func (s *Server) listEmojiPacksHandler(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 200 {
			limit = parsed
		}
	}

	var packs []EmojiPack
	result := s.db.
		Where("is_public = ?", true).
		Order("downloads DESC").
		Limit(limit).
		Find(&packs)

	if result.Error != nil {
		http.Error(w, "Failed to fetch emoji packs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(packs)
}

// GET /api/emoji-packs/{id} - Get emoji pack details
func (s *Server) getEmojiPackHandler(w http.ResponseWriter, r *http.Request) {
	packID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(packID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid pack ID", http.StatusBadRequest)
		return
	}

	var pack EmojiPack
	result := s.db.First(&pack, id)
	if result.Error != nil {
		http.Error(w, "Pack not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pack)
}

// POST /api/emoji-packs/{id}/emojis - Add emoji to pack
func (s *Server) addEmojiHandler(w http.ResponseWriter, r *http.Request) {
	packID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(packID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid pack ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Code      string `json:"code"`
		Alias     string `json:"alias"`
		ImageURL  string `json:"image_url"`
		Category  string `json:"category"`
		Tags      string `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	emoji := Emoji{
		PackID:   uint(id),
		Code:     reqBody.Code,
		Alias:    reqBody.Alias,
		ImageURL: reqBody.ImageURL,
		Category: reqBody.Category,
		Tags:     reqBody.Tags,
	}

	result := s.db.Create(&emoji)
	if result.Error != nil {
		http.Error(w, "Failed to add emoji", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emoji)
}

// GET /api/emoji-packs/{id}/emojis - Get emojis in pack
func (s *Server) getPackEmojisHandler(w http.ResponseWriter, r *http.Request) {
	packID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(packID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid pack ID", http.StatusBadRequest)
		return
	}

	var emojis []Emoji
	result := s.db.
		Where("pack_id = ?", id).
		Order("created_at ASC").
		Find(&emojis)

	if result.Error != nil {
		http.Error(w, "Failed to fetch emojis", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emojis)
}

// POST /api/users/{id}/emoji-packs/{packID}/subscribe - Subscribe to emoji pack
func (s *Server) subscribeEmojiPackHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	packID := chi.URLParam(r, "packID")

	uid, _ := strconv.ParseUint(userID, 10, 32)
	pid, err := strconv.ParseUint(packID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid pack ID", http.StatusBadRequest)
		return
	}

	subscription := UserEmojiPack{
		UserID:    uint(uid),
		EmojiPack: uint(pid),
	}

	result := s.db.Create(&subscription)
	if result.Error != nil {
		http.Error(w, "Failed to subscribe to pack", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subscription)
}

// GET /api/users/{id}/emoji-packs - Get user's subscribed emoji packs
func (s *Server) getUserEmojiPacksHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var packs []UserEmojiPack
	result := s.db.
		Where("user_id = ?", uid).
		Preload("Pack").
		Find(&packs)

	if result.Error != nil {
		http.Error(w, "Failed to fetch emoji packs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(packs)
}

// POST /api/messages/{id}/emoji-react - Add emoji reaction to message
func (s *Server) addMessageEmojiHandler(w http.ResponseWriter, r *http.Request) {
	msgID := chi.URLParam(r, "id")
	mid, err := strconv.ParseUint(msgID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		EmojiCode string `json:"emoji_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	var existing MessageEmoji
	s.db.Where("message_id = ? AND user_id = ? AND emoji_code = ?", mid, userID, reqBody.EmojiCode).First(&existing)

	if existing.ID == 0 {
		emoji := MessageEmoji{
			MessageID: uint(mid),
			UserID:    uint(userID),
			EmojiCode: reqBody.EmojiCode,
			Count:     1,
		}
		s.db.Create(&emoji)
	} else {
		s.db.Model(&existing).Update("count", existing.Count+1)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// GET /api/emoji-packs/{id}/reviews - Get emoji pack reviews
func (s *Server) getPackReviewsHandler(w http.ResponseWriter, r *http.Request) {
	packID := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(packID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid pack ID", http.StatusBadRequest)
		return
	}

	var reviews []EmojiPackReview
	result := s.db.
		Where("pack_id = ?", id).
		Order("created_at DESC").
		Find(&reviews)

	if result.Error != nil {
		http.Error(w, "Failed to fetch reviews", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reviews)
}

// POST /api/emoji-packs/{id}/reviews - Create pack review
func (s *Server) createPackReviewHandler(w http.ResponseWriter, r *http.Request) {
	packID := chi.URLParam(r, "id")
	pid, err := strconv.ParseUint(packID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid pack ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Rating int    `json:"rating"`
		Comment string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	review := EmojiPackReview{
		PackID:  uint(pid),
		UserID:  uint(userID),
		Rating:  reqBody.Rating,
		Comment: reqBody.Comment,
	}

	result := s.db.Create(&review)
	if result.Error != nil {
		http.Error(w, "Failed to create review", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}
