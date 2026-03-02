package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Enable2FARequest represents a request to enable 2FA
type Enable2FARequest struct {
	Method      string `json:"method"` // "email", "sms", "totp"
	PhoneNumber string `json:"phone_number,omitempty"`
}

// Verify2FARequest represents a request to verify 2FA code
type Verify2FARequest struct {
	Code string `json:"code"`
}

// SecurityPolicyRequest represents a security policy update
type SecurityPolicyRequest struct {
	MinPasswordLength  int         `json:"min_password_length"`
	MaxLoginAttempts   int         `json:"max_login_attempts"`
	LockoutDuration    int         `json:"lockout_duration"`
	SessionTimeout     int         `json:"session_timeout"`
	Require2FA         bool        `json:"require_2fa"`
	PasswordExpiryDays int         `json:"password_expiry_days"`
	AllowedCountries   []string    `json:"allowed_countries"`
}

// Enable2FAHandler enables 2FA for a user
func (s *Server) enable2FAHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req Enable2FARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Create 2FA record
	twoFA := TwoFactorAuth{
		UserID:  userID,
		Method:  req.Method,
		Enabled: false, // Not enabled until verified
	}

	if req.Method == "sms" {
		twoFA.PhoneNumber = req.PhoneNumber
	}

	if err := s.db.Create(&twoFA).Error; err != nil {
		http.Error(w, "Failed to enable 2FA", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(twoFA)
}

// Verify2FAHandler verifies a 2FA code
func (s *Server) verify2FAHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req Verify2FARequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	twoFAID := chi.URLParam(r, "id")
	var twoFA TwoFactorAuth
	if err := s.db.Where("id = ? AND user_id = ?", twoFAID, userID).First(&twoFA).Error; err != nil {
		http.Error(w, "2FA not found", http.StatusNotFound)
		return
	}

	// In a real app, validate the code against TOTP/email/SMS
	twoFA.Enabled = true
	s.db.Save(&twoFA)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "2FA enabled successfully", "enabled": true})
}

// GetSecurityStatusHandler returns user's security status
func (s *Server) getSecurityStatusHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var twoFAs []TwoFactorAuth
	s.db.Where("user_id = ?", userID).Find(&twoFAs)

	var securityKeys []SecurityKey
	s.db.Where("user_id = ?", userID).Find(&securityKeys)

	var sessions []SessionToken
	s.db.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(&sessions)

	var incidents []SecurityIncident
	s.db.Where("user_id = ? AND resolved = ?", userID, false).Find(&incidents)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"two_factor_methods": len(twoFAs),
		"security_keys":      len(securityKeys),
		"active_sessions":    len(sessions),
		"security_incidents": len(incidents),
	})
}

// GetLoginAttemptsHandler retrieves login attempts for a user
func (s *Server) getLoginAttemptsHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var attempts []LoginAttempt
	s.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(50).Find(&attempts)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attempts)
}

// GetSessionsHandler retrieves active sessions
func (s *Server) getSessionsHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var sessions []SessionToken
	s.db.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(&sessions)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

// RevokeSessionHandler revokes a specific session
func (s *Server) revokeSessionHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionID := chi.URLParam(r, "sessionID")
	var session SessionToken
	if err := s.db.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error; err != nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	s.db.Delete(&session)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Session revoked"})
}

// GetSecurityIncidentsHandler retrieves security incidents
func (s *Server) getSecurityIncidentsHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var incidents []SecurityIncident
	s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&incidents)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incidents)
}

// UpdateSecurityPolicyHandler updates security policies (admin only)
func (s *Server) updateSecurityPolicyHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if user is admin
	var user User
	if err := s.db.First(&user, userID).Error; err != nil || !user.IsAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req SecurityPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var policy SecurityPolicy
	s.db.FirstOrCreate(&policy)

	policy.MinPasswordLength = req.MinPasswordLength
	policy.MaxLoginAttempts = req.MaxLoginAttempts
	policy.LockoutDuration = req.LockoutDuration
	policy.SessionTimeout = req.SessionTimeout
	policy.Require2FA = req.Require2FA
	policy.PasswordExpiryDays = req.PasswordExpiryDays
	policy.AllowedCountries = req.AllowedCountries

	s.db.Save(&policy)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policy)
}

// RegisterTrustedDeviceHandler registers a device as trusted
func (s *Server) registerTrustedDeviceHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	device := TrustedDevice{
		UserID:    userID,
		Name:      req.Name,
		Verified:  true,
		ExpiresAt: time.Now().AddDate(1, 0, 0), // 1 year
	}

	if err := s.db.Create(&device).Error; err != nil {
		http.Error(w, "Failed to register device", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}

// GetTrustedDevicesHandler retrieves trusted devices
func (s *Server) getTrustedDevicesHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var devices []TrustedDevice
	s.db.Where("user_id = ? AND expires_at > ?", userID, time.Now()).Find(&devices)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}
