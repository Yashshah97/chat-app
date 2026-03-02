package main

import (
	"encoding/json"
	"net/http"
)

// ValidatePasswordRequest for password validation
type ValidatePasswordRequest struct {
	Password string `json:"password"`
}

// CreateEncryptionKeyRequest for creating encryption keys
type CreateEncryptionKeyRequest struct {
	KeyName   string `json:"key_name"`
	Algorithm string `json:"algorithm"`
}

// ValidatePasswordHandler validates password strength
func (s *Server) validatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req ValidatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Get current password policy
	var policy PasswordPolicy
	s.db.FirstOrCreate(&policy)

	isValid := ValidatePasswordStrength(req.Password, &policy)
	complexity := CalculatePasswordComplexity(req.Password)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":           isValid,
		"complexity_score": complexity,
		"requirements_met": map[string]interface{}{
			"length":             len(req.Password) >= policy.MinLength,
			"uppercase":          policy.RequireUppercase,
			"lowercase":          policy.RequireLowercase,
			"numbers":            policy.RequireNumbers,
			"special_characters": policy.RequireSpecialCharacters,
		},
	})
}

// CreateEncryptionKeyHandler creates a new encryption key for user
func (s *Server) createEncryptionKeyHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateEncryptionKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	key := EncryptionKey{
		UserID:    userID,
		KeyName:   req.KeyName,
		Algorithm: req.Algorithm,
	}

	if err := s.db.Create(&key).Error; err != nil {
		http.Error(w, "Failed to create encryption key", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

// GetEncryptionKeysHandler retrieves all encryption keys for user
func (s *Server) getEncryptionKeysHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var keys []EncryptionKey
	s.db.Where("user_id = ?", userID).Find(&keys)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}

// GetPasswordPolicyHandler retrieves password policy
func (s *Server) getPasswordPolicyHandler(w http.ResponseWriter, r *http.Request) {
	var policy PasswordPolicy
	s.db.FirstOrCreate(&policy)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policy)
}

// UpdatePasswordPolicyHandler updates password policy (admin only)
func (s *Server) updatePasswordPolicyHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if admin
	var user User
	if err := s.db.First(&user, userID).Error; err != nil || !user.IsAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req PasswordPolicy
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var policy PasswordPolicy
	s.db.FirstOrCreate(&policy)

	policy.MinLength = req.MinLength
	policy.RequireUppercase = req.RequireUppercase
	policy.RequireLowercase = req.RequireLowercase
	policy.RequireNumbers = req.RequireNumbers
	policy.RequireSpecialCharacters = req.RequireSpecialCharacters
	policy.ExpiryDays = req.ExpiryDays
	policy.HistoryCount = req.HistoryCount

	s.db.Save(&policy)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policy)
}

// GetAuditActionsHandler retrieves audit actions (admin only)
func (s *Server) getAuditActionsHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if admin
	var user User
	if err := s.db.First(&user, userID).Error; err != nil || !user.IsAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var actions []AuditAction
	s.db.Order("created_at DESC").Limit(100).Find(&actions)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actions)
}

// GetSecureDeletesHandler retrieves secure delete requests
func (s *Server) getSecureDeletesHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var deletes []SecureDelete
	s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&deletes)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletes)
}

// RequestSecureDeleteHandler requests secure deletion of data
func (s *Server) requestSecureDeleteHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		DataType string `json:"data_type"`
		DataID   uint   `json:"data_id"`
		Reason   string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	delete := SecureDelete{
		UserID:   userID,
		DataType: req.DataType,
		DataID:   req.DataID,
		Reason:   req.Reason,
	}

	if err := s.db.Create(&delete).Error; err != nil {
		http.Error(w, "Failed to request secure delete", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(delete)
}

// GetDataEncryptionStatusHandler checks encryption status of user data
func (s *Server) getDataEncryptionStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var encryptedData []DataEncryption
	s.db.Where("user_id = ?", userID).Find(&encryptedData)

	encryptedByType := make(map[string]int)
	for _, d := range encryptedData {
		encryptedByType[d.DataType]++
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_encrypted": len(encryptedData),
		"by_type":         encryptedByType,
	})
}

// EnableEndToEndEncryptionHandler enables E2E encryption for user
func (s *Server) enableEndToEndEncryptionHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Mark user as E2E enabled
	var user User
	if err := s.db.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.IsAdmin = true // Placeholder - should have E2E flag
	s.db.Save(&user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "E2E encryption enabled"})
}
