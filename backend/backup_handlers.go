package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// POST /api/chats/{id}/backup - Create a chat backup
func (s *Server) createChatBackupHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Format       string `json:"format"`
		IncludeMedia bool   `json:"include_media"`
		StartDate    *string `json:"start_date"`
		EndDate      *string `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	backup := ChatBackup{
		ChatID:       uint(cid),
		CreatedByID:  uint(userID),
		Format:       reqBody.Format,
		IncludeMedia: reqBody.IncludeMedia,
		StartDate:    reqBody.StartDate,
		EndDate:      reqBody.EndDate,
		Status:       "pending",
		Progress:     0,
	}

	result := s.db.Create(&backup)
	if result.Error != nil {
		http.Error(w, "Failed to create backup", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(backup)
}

// GET /api/chats/{id}/backups - List chat backups
func (s *Server) listChatBackupsHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var backups []ChatBackup
	result := s.db.
		Where("chat_id = ?", cid).
		Order("created_at DESC").
		Find(&backups)

	if result.Error != nil {
		http.Error(w, "Failed to fetch backups", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(backups)
}

// GET /api/backups/{id} - Get backup details
func (s *Server) getBackupDetailsHandler(w http.ResponseWriter, r *http.Request) {
	backupID := chi.URLParam(r, "id")
	bid, err := strconv.ParseUint(backupID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid backup ID", http.StatusBadRequest)
		return
	}

	var backup ChatBackup
	result := s.db.First(&backup, bid)
	if result.Error != nil {
		http.Error(w, "Backup not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(backup)
}

// POST /api/chats/{id}/export - Create a chat export
func (s *Server) createChatExportHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Format      string `json:"format"`
		FilterType  string `json:"filter_type"`
		FilterValue string `json:"filter_value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	export := ChatExport{
		ChatID:      uint(cid),
		UserID:      uint(userID),
		Format:      reqBody.Format,
		Status:      "pending",
		FilterType:  reqBody.FilterType,
		FilterValue: reqBody.FilterValue,
	}

	result := s.db.Create(&export)
	if result.Error != nil {
		http.Error(w, "Failed to create export", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(export)
}

// GET /api/exports/{id} - Get export details
func (s *Server) getExportDetailsHandler(w http.ResponseWriter, r *http.Request) {
	exportID := chi.URLParam(r, "id")
	eid, err := strconv.ParseUint(exportID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid export ID", http.StatusBadRequest)
		return
	}

	var export ChatExport
	result := s.db.First(&export, eid)
	if result.Error != nil {
		http.Error(w, "Export not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(export)
}

// POST /api/chats/{id}/backup-schedule - Create backup schedule
func (s *Server) createBackupScheduleHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Frequency       string `json:"frequency"`
		Format          string `json:"format"`
		MaxBackups      int    `json:"max_backups"`
		IncludeMedia    bool   `json:"include_media"`
		StorageLocation string `json:"storage_location"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	schedule := BackupSchedule{
		ChatID:          uint(cid),
		CreatedByID:     uint(userID),
		Frequency:       reqBody.Frequency,
		Format:          reqBody.Format,
		MaxBackups:      reqBody.MaxBackups,
		IncludeMedia:    reqBody.IncludeMedia,
		StorageLocation: reqBody.StorageLocation,
		IsActive:        true,
	}

	result := s.db.Create(&schedule)
	if result.Error != nil {
		http.Error(w, "Failed to create schedule", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(schedule)
}

// GET /api/chats/{id}/backup-schedule - Get backup schedule
func (s *Server) getBackupScheduleHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var schedule BackupSchedule
	result := s.db.Where("chat_id = ?", cid).First(&schedule)
	if result.Error != nil {
		http.Error(w, "Schedule not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schedule)
}

// DELETE /api/backups/{id} - Delete backup
func (s *Server) deleteBackupHandler(w http.ResponseWriter, r *http.Request) {
	backupID := chi.URLParam(r, "id")
	bid, err := strconv.ParseUint(backupID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid backup ID", http.StatusBadRequest)
		return
	}

	result := s.db.Delete(&ChatBackup{}, bid)
	if result.Error != nil {
		http.Error(w, "Failed to delete backup", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// POST /api/archives/message - Archive a message
func (s *Server) archiveMessageHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		MessageID uint   `json:"message_id"`
		ChatID    uint   `json:"chat_id"`
		Reason    string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	archive := ArchiveMessage{
		MessageID:    reqBody.MessageID,
		ChatID:       reqBody.ChatID,
		ArchivedByID: uint(userID),
		Reason:       reqBody.Reason,
	}

	result := s.db.Create(&archive)
	if result.Error != nil {
		http.Error(w, "Failed to archive message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(archive)
}

// GET /api/chats/{id}/archived-messages - Get archived messages
func (s *Server) getArchivedMessagesHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var archived []ArchiveMessage
	result := s.db.
		Where("chat_id = ?", cid).
		Order("created_at DESC").
		Find(&archived)

	if result.Error != nil {
		http.Error(w, "Failed to fetch archived messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(archived)
}

// POST /api/backups/{id}/download - Download backup
func (s *Server) downloadBackupHandler(w http.ResponseWriter, r *http.Request) {
	backupID := chi.URLParam(r, "id")
	bid, err := strconv.ParseUint(backupID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid backup ID", http.StatusBadRequest)
		return
	}

	var backup ChatBackup
	result := s.db.First(&backup, bid)
	if result.Error != nil {
		http.Error(w, "Backup not found", http.StatusNotFound)
		return
	}

	if backup.Status != "completed" {
		http.Error(w, "Backup not ready for download", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"download_url": backup.FileURL})
}
