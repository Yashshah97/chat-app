package main

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// TwoFactorAuth represents two-factor authentication configuration
type TwoFactorAuth struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Method       string    `json:"method"` // "email", "sms", "totp", "backup_codes"
	Enabled      bool      `json:"enabled"`
	Secret       string    `json:"-"` // TOTP secret, encrypted
	PhoneNumber  string    `json:"phone_number"`
	BackupCodes  []string  `gorm:"type:text[]" json:"backup_codes,omitempty"`
	LastUsedAt   time.Time `json:"last_used_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// SecurityKey represents hardware security keys
type SecurityKey struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	KeyName   string         `json:"key_name"`
	PublicKey string         `json:"-"`
	Counter   uint32         `json:"counter"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// LoginAttempt tracks login attempts for security
type LoginAttempt struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Success   bool      `json:"success"`
	Reason    string    `json:"reason"` // "invalid_password", "account_locked", etc.
	CreatedAt time.Time `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// SessionToken represents active user sessions
type SessionToken struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Token       string    `gorm:"index" json:"-"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	LastUsedAt  time.Time `json:"last_used_at"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// PasswordHistory tracks previous passwords to prevent reuse
type PasswordHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Hash      string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

// SecurityPolicy defines security policies for the system
type SecurityPolicy struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	MinPasswordLength     int            `json:"min_password_length"`
	MaxLoginAttempts      int            `json:"max_login_attempts"`
	LockoutDuration       int            `json:"lockout_duration"` // in minutes
	SessionTimeout        int            `json:"session_timeout"`  // in minutes
	Require2FA            bool           `json:"require_2fa"`
	EnforceSSL            bool           `json:"enforce_ssl"`
	AllowWeakPasswords    bool           `json:"allow_weak_passwords"`
	PasswordExpiryDays    int            `json:"password_expiry_days"`
	IPWhitelist           datatypes.JSON `gorm:"type:jsonb" json:"ip_whitelist"`
	AllowedCountries      []string       `gorm:"type:text[]" json:"allowed_countries"`
	DisableLoginMethod    string         `json:"disable_login_method"` // temporary disable
	DisableReason         string         `json:"disable_reason"`
	DisabledUntil         time.Time      `json:"disabled_until"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}

// FailedPasswordAttempt tracks password change failures
type FailedPasswordAttempt struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

// TrustedDevice allows users to mark devices as trusted
type TrustedDevice struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Name      string    `json:"name"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	Verified  bool      `json:"verified"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// SecurityIncident tracks security-related incidents
type SecurityIncident struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `json:"user_id"`
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	IncidentType string         `json:"incident_type"` // "suspicious_login", "brute_force", etc.
	Severity     string         `json:"severity"`      // "low", "medium", "high", "critical"
	Description  string         `json:"description"`
	Details      datatypes.JSON `gorm:"type:jsonb" json:"details"`
	IPAddress    string         `json:"ip_address"`
	Resolved     bool           `json:"resolved"`
	ResolvedAt   time.Time      `json:"resolved_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// OAuthToken stores OAuth tokens for third-party integrations
type OAuthToken struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Provider     string    `json:"provider"`
	AccessToken  string    `json:"-"`
	RefreshToken string    `json:"-"`
	ExpiresAt    time.Time `json:"expires_at"`
	Scope        string    `json:"scope"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
