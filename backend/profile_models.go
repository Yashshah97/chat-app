package main

import "gorm.io/gorm"

// UserProfile extends user with additional fields
type UserProfile struct {
	gorm.Model
	UserID      uint   `gorm:"not null;uniqueIndex" json:"user_id"`
	User        *User  `gorm:"foreignKey:UserID" json:"-"`
	Bio         string `json:"bio"`
	Avatar      string `json:"avatar"`
	CoverPhoto  string `json:"cover_photo"`
	Location    string `json:"location"`
	Website     string `json:"website"`
	BirthDate   *string `json:"birth_date"`
	LastSeen    string `json:"last_seen"`
	TimeZone    string `json:"time_zone"`
	Language    string `json:"language"`
	ProfileViews int64 `gorm:"default:0" json:"profile_views"`
	IsVerified  bool   `gorm:"default:false" json:"is_verified"`
	IsBanned    bool   `gorm:"default:false" json:"is_banned"`
}

// UserStatus represents user online/away status
type UserStatus struct {
	gorm.Model
	UserID      uint   `gorm:"not null;uniqueIndex" json:"user_id"`
	User        *User  `gorm:"foreignKey:UserID" json:"-"`
	Status      string `json:"status"` // online, away, dnd, offline, invisible
	StatusMessage string `json:"status_message"`
	LastActive  string `json:"last_active"`
	Device      string `json:"device"` // web, mobile, desktop
	Browser     string `json:"browser"`
}

// SearchIndex represents indexed documents for search
type SearchIndex struct {
	gorm.Model
	DocumentType string `gorm:"index" json:"document_type"` // message, chat, user, file
	DocumentID  uint   `gorm:"index" json:"document_id"`
	Title       string `json:"title"`
	Content     string `gorm:"type:text" json:"content"`
	SearchVector string `json:"search_vector"`
	Metadata    string `json:"metadata"`
	IndexedAt   string `json:"indexed_at"`
}

// SavedSearch represents user's saved searches
type SavedSearch struct {
	gorm.Model
	UserID    uint   `gorm:"not null;index" json:"user_id"`
	User      *User  `gorm:"foreignKey:UserID" json:"-"`
	Name      string `json:"name"`
	Query     string `json:"query"`
	Filters   string `json:"filters"`
	IsPublic  bool   `gorm:"default:false" json:"is_public"`
	LastRunAt string `json:"last_run_at"`
}

// ContentPolicy represents platform content policies
type ContentPolicy struct {
	gorm.Model
	Name        string `gorm:"not null;uniqueIndex" json:"name"`
	Description string `json:"description"`
	Rules       string `gorm:"type:jsonb" json:"rules"`
	Category    string `json:"category"` // spam, abuse, adult, violence, etc
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	Severity    int    `json:"severity"` // 1-5
	Actions     string `json:"actions"` // warn, mute, ban
}
