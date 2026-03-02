package main

import "gorm.io/gorm"

// RateLimitConfig represents rate limiting configuration
type RateLimitConfig struct {
	gorm.Model
	UserID    *uint  `gorm:"uniqueIndex" json:"user_id"`
	User      *User  `gorm:"foreignKey:UserID" json:"-"`
	ChatID    *uint  `gorm:"uniqueIndex" json:"chat_id"`
	Chat      *Chat  `gorm:"foreignKey:ChatID" json:"-"`
	Endpoint  string `gorm:"not null;index" json:"endpoint"`
	RequestsPerMinute int `json:"requests_per_minute"`
	RequestsPerHour int `json:"requests_per_hour"`
	RequestsPerDay int `json:"requests_per_day"`
	IsActive  bool   `gorm:"default:true" json:"is_active"`
	Priority  int    `json:"priority"` // Higher priority = less restrictive
	WhitelistIPs string `json:"whitelist_ips"` // comma-separated
}

// RateLimitUsage tracks rate limit usage
type RateLimitUsage struct {
	gorm.Model
	UserID    *uint  `gorm:"index" json:"user_id"`
	User      *User  `gorm:"foreignKey:UserID" json:"-"`
	Endpoint  string `gorm:"not null;index" json:"endpoint"`
	Period    string `json:"period"` // minute, hour, day
	RequestCount int `json:"request_count"`
	ResetAt   string `json:"reset_at"`
	IsBlocked bool   `gorm:"default:false" json:"is_blocked"`
	IPAddress string `json:"ip_address"`
}

// RateLimitViolation logs rate limit violations
type RateLimitViolation struct {
	gorm.Model
	UserID    *uint  `gorm:"index" json:"user_id"`
	User      *User  `gorm:"foreignKey:UserID" json:"-"`
	Endpoint  string `json:"endpoint"`
	Attempts  int    `json:"attempts"`
	Limit     int    `json:"limit"`
	Period    string `json:"period"` // minute, hour, day
	Action    string `json:"action"` // blocked, throttled, warned
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Reason    string `json:"reason"`
}

// APIVersion represents API version information
type APIVersion struct {
	gorm.Model
	Version         string `gorm:"not null;uniqueIndex" json:"version"`
	Status          string `json:"status"` // active, deprecated, sunset
	ReleasedAt      string `json:"released_at"`
	SunsetDate      *string `json:"sunset_date"`
	DeprecationDate *string `json:"deprecation_date"`
	Description     string `json:"description"`
	Changes         string `json:"changes"` // JSON array of changes
	MigrationGuide  string `json:"migration_guide"`
	IsDefault       bool   `gorm:"default:false" json:"is_default"`
}

// APIEndpoint represents API endpoint metadata
type APIEndpoint struct {
	gorm.Model
	Path          string `gorm:"not null;index" json:"path"`
	Method        string `gorm:"not null;index" json:"method"` // GET, POST, PUT, DELETE
	Version       string `gorm:"not null;index" json:"version"`
	Description   string `json:"description"`
	Authentication bool  `gorm:"default:true" json:"authentication"`
	RateLimit     int    `json:"rate_limit"` // requests per minute
	IsDeprecated  bool   `gorm:"default:false" json:"is_deprecated"`
	DeprecatedAt  *string `json:"deprecated_at"`
	ReplacementPath *string `json:"replacement_path"`
	Examples      string `json:"examples"` // JSON array
	Parameters    string `json:"parameters"` // JSON schema
}
