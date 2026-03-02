package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// EncryptionKey manages encryption keys for sensitive data
type EncryptionKey struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	KeyName   string         `json:"key_name"`
	Algorithm string         `json:"algorithm"` // "AES-256", "RSA", etc.
	PublicKey string         `json:"-"`
	Metadata  datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// PasswordPolicy enforces password requirements
type PasswordPolicy struct {
	ID                        uint           `gorm:"primaryKey" json:"id"`
	MinLength                 int            `json:"min_length"`
	RequireUppercase          bool           `json:"require_uppercase"`
	RequireLowercase          bool           `json:"require_lowercase"`
	RequireNumbers            bool           `json:"require_numbers"`
	RequireSpecialCharacters  bool           `json:"require_special_characters"`
	ExpiryDays                int            `json:"expiry_days"`
	HistoryCount              int            `json:"history_count"` // Prevent reuse
	MaxConsecutiveCharacters  int            `json:"max_consecutive_characters"`
	MinUniqueCharacters       int            `json:"min_unique_characters"`
	ForbiddenPatterns         []string       `gorm:"type:text[]" json:"forbidden_patterns"`
	ComplexityScore           int            `json:"complexity_score"`
	CreatedAt                 time.Time      `json:"created_at"`
	UpdatedAt                 time.Time      `json:"updated_at"`
	DeletedAt                 gorm.DeletedAt `gorm:"index" json:"-"`
}

// DataEncryption stores encrypted data
type DataEncryption struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	DataType     string    `json:"data_type"` // "message", "file", "backup"
	EncryptedKey string    `json:"-"`
	IV           string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// SecureDelete marks data for secure deletion
type SecureDelete struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	DataType  string    `json:"data_type"`
	DataID    uint      `json:"data_id"`
	Reason    string    `json:"reason"`
	DeletedAt time.Time `json:"deleted_at"`
	CreatedAt time.Time `json:"created_at"`
}

// AuditAction tracks every action in the system
type AuditAction struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action    string         `json:"action"` // "login", "delete_message", etc.
	Resource  string         `json:"resource"`
	ResourceID uint          `json:"resource_id"`
	Changes   datatypes.JSON `gorm:"type:jsonb" json:"changes"`
	IPAddress string         `json:"ip_address"`
	Status    string         `json:"status"` // "success", "failure"
	ErrorMsg  string         `json:"error_msg,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

// Encryption utilities
func EncryptAES(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptAES(ciphertext string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext2 := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, []byte(ciphertext2), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// ValidatePasswordStrength checks password against policy
func ValidatePasswordStrength(password string, policy *PasswordPolicy) bool {
	if len(password) < policy.MinLength {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	if policy.RequireUppercase && !hasUpper {
		return false
	}
	if policy.RequireLowercase && !hasLower {
		return false
	}
	if policy.RequireNumbers && !hasDigit {
		return false
	}
	if policy.RequireSpecialCharacters && !hasSpecial {
		return false
	}

	return true
}

// CalculatePasswordComplexity returns a score 0-100
func CalculatePasswordComplexity(password string) int {
	score := 0

	if len(password) >= 8 {
		score += 20
	}
	if len(password) >= 12 {
		score += 10
	}
	if len(password) >= 16 {
		score += 10
	}

	hasLower, hasUpper, hasDigit, hasSpecial := false, false, false, false
	for _, char := range password {
		switch {
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= '0' && char <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	if hasLower {
		score += 15
	}
	if hasUpper {
		score += 15
	}
	if hasDigit {
		score += 15
	}
	if hasSpecial {
		score += 15
	}

	return score
}
