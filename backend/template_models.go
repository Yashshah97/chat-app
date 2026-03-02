package main

import "gorm.io/gorm"

// MessageTemplate represents a reusable message template
type MessageTemplate struct {
	gorm.Model
	Name        string `gorm:"not null;index" json:"name"`
	Description string `json:"description"`
	Content     string `gorm:"type:text" json:"content"`
	Category    string `json:"category"` // greeting, support, announcement, etc
	CreatedByID uint   `json:"created_by_id"`
	CreatedBy   *User  `gorm:"foreignKey:CreatedByID" json:"-"`
	ChatID      *uint  `gorm:"index" json:"chat_id"`
	Chat        *Chat  `gorm:"foreignKey:ChatID" json:"-"`
	IsPublic    bool   `gorm:"default:false" json:"is_public"`
	Tags        string `json:"tags"` // comma-separated
	Variables   string `json:"variables"` // JSON array of variable names
	UsageCount  int64  `gorm:"default:0" json:"usage_count"`
	IsFavorite  bool   `gorm:"default:false" json:"is_favorite"`
}

// QuickReply represents quick reply buttons/suggestions
type QuickReply struct {
	gorm.Model
	Title   string `gorm:"not null" json:"title"`
	Payload string `json:"payload"`
	Icon    string `json:"icon"`
	Type    string `json:"type"` // text, url, postback
}

// ChatBot represents an automated chat bot
type ChatBot struct {
	gorm.Model
	Name        string `gorm:"not null;index" json:"name"`
	Description string `json:"description"`
	ChatID      uint   `gorm:"not null;index" json:"chat_id"`
	Chat        *Chat  `gorm:"foreignKey:ChatID" json:"-"`
	CreatedByID uint   `json:"created_by_id"`
	CreatedBy   *User  `gorm:"foreignKey:CreatedByID" json:"-"`
	Token       string `gorm:"not null;uniqueIndex" json:"token"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	Webhook     string `json:"webhook"`
	IntentModels string `json:"intent_models"` // JSON array of models
	ResponseLatency int `json:"response_latency"` // ms
	MessageCount int64 `gorm:"default:0" json:"message_count"`
	SuccessRate float64 `json:"success_rate"`
}

// BotIntent represents intent patterns for bots
type BotIntent struct {
	gorm.Model
	BotID       uint   `gorm:"not null;index" json:"bot_id"`
	Bot         *ChatBot `gorm:"foreignKey:BotID" json:"-"`
	Pattern     string `gorm:"not null" json:"pattern"` // regex or keyword pattern
	Response    string `gorm:"type:text" json:"response"`
	Priority    int    `json:"priority"` // higher = checked first
	Enabled     bool   `gorm:"default:true" json:"enabled"`
	MatchCount  int64  `gorm:"default:0" json:"match_count"`
}

// ChatStatistics represents time-series statistics
type ChatStatistics struct {
	gorm.Model
	ChatID      uint   `gorm:"not null;index" json:"chat_id"`
	Chat        *Chat  `gorm:"foreignKey:ChatID" json:"-"`
	Date        string `gorm:"index" json:"date"`
	MessageCount int   `json:"message_count"`
	ActiveUsers int   `json:"active_users"`
	NewUsers    int   `json:"new_users"`
	AverageResponseTime int `json:"average_response_time"`
	MostActiveHour int `json:"most_active_hour"`
}
