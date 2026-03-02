package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// POST /api/users/{id}/data-export - Request data export
func (s *Server) requestDataExportHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Format      string `json:"format"`
		IncludeData string `json:"include_data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	export := DataExport{
		UserID:      uint(uid),
		Format:      reqBody.Format,
		IncludeData: reqBody.IncludeData,
		Status:      "pending",
	}

	result := s.db.Create(&export)
	if result.Error != nil {
		http.Error(w, "Failed to request export", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(export)
}

// GET /api/users/{id}/data-exports - List user's data exports
func (s *Server) listDataExportsHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var exports []DataExport
	result := s.db.
		Where("user_id = ?", uid).
		Order("created_at DESC").
		Find(&exports)

	if result.Error != nil {
		http.Error(w, "Failed to fetch exports", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exports)
}

// POST /api/chats/{id}/files - Upload file to chat
func (s *Server) uploadFileToChatHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		FileName    string `json:"file_name"`
		FileSize    int64  `json:"file_size"`
		FileType    string `json:"file_type"`
		FileHash    string `json:"file_hash"`
		StorageKey  string `json:"storage_key"`
		Tags        string `json:"tags"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.ParseUint(r.Header.Get("User-ID"), 10, 32)

	file := FileManagement{
		ChatID:      uint(cid),
		FileName:    reqBody.FileName,
		FileSize:    reqBody.FileSize,
		FileType:    reqBody.FileType,
		FileHash:    reqBody.FileHash,
		StorageKey:  reqBody.StorageKey,
		UploadedByID: uint(userID),
		Tags:        reqBody.Tags,
		Description: reqBody.Description,
	}

	result := s.db.Create(&file)
	if result.Error != nil {
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(file)
}

// GET /api/chats/{id}/files - List files in chat
func (s *Server) listChatFilesHandler(w http.ResponseWriter, r *http.Request) {
	chatID := chi.URLParam(r, "id")
	cid, err := strconv.ParseUint(chatID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid chat ID", http.StatusBadRequest)
		return
	}

	var files []FileManagement
	result := s.db.
		Where("chat_id = ?", cid).
		Order("created_at DESC").
		Find(&files)

	if result.Error != nil {
		http.Error(w, "Failed to fetch files", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(files)
}

// POST /api/users/{id}/data-request - Submit GDPR data request
func (s *Server) submitUserDataRequestHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		RequestType string `json:"request_type"`
		Reason      string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	request := UserDataRequest{
		UserID:      uint(uid),
		RequestType: reqBody.RequestType,
		Reason:      reqBody.Reason,
		Status:      "pending",
	}

	result := s.db.Create(&request)
	if result.Error != nil {
		http.Error(w, "Failed to submit request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(request)
}

// GET /api/users/{id}/data-requests - Get user data requests
func (s *Server) getUserDataRequestsHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var requests []UserDataRequest
	result := s.db.
		Where("user_id = ?", uid).
		Order("created_at DESC").
		Find(&requests)

	if result.Error != nil {
		http.Error(w, "Failed to fetch requests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(requests)
}

// POST /api/audit-logs - Create audit log entry
func (s *Server) createAuditLogEntryHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		UserID       uint   `json:"user_id"`
		Action       string `json:"action"`
		ResourceType string `json:"resource_type"`
		ResourceID   *uint  `json:"resource_id"`
		OldValue     string `json:"old_value"`
		NewValue     string `json:"new_value"`
		IPAddress    string `json:"ip_address"`
		UserAgent    string `json:"user_agent"`
		Status       string `json:"status"`
		Details      string `json:"details"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	entry := AuditLogEntry{
		UserID:       reqBody.UserID,
		Action:       reqBody.Action,
		ResourceType: reqBody.ResourceType,
		ResourceID:   reqBody.ResourceID,
		OldValue:     reqBody.OldValue,
		NewValue:     reqBody.NewValue,
		IPAddress:    reqBody.IPAddress,
		UserAgent:    reqBody.UserAgent,
		Status:       reqBody.Status,
		Details:      reqBody.Details,
	}

	result := s.db.Create(&entry)
	if result.Error != nil {
		http.Error(w, "Failed to create audit log", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(entry)
}

// GET /api/audit-logs - Get audit logs
func (s *Server) getAuditLogsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	action := r.URL.Query().Get("action")

	query := s.db
	if userID != "" {
		if uid, err := strconv.ParseUint(userID, 10, 32); err == nil {
			query = query.Where("user_id = ?", uid)
		}
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}

	var logs []AuditLogEntry
	result := query.
		Order("created_at DESC").
		Limit(100).
		Find(&logs)

	if result.Error != nil {
		http.Error(w, "Failed to fetch logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logs)
}
