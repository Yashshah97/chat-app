package main

import (
	"time"
)

// ChatAnalytics tracks metrics for a chat
type ChatAnalytics struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	ChatID          uint      `json:"chat_id"`
	Chat            Chat      `gorm:"foreignKey:ChatID" json:"chat,omitempty"`
	TotalMessages   int64     `json:"total_messages"`
	ActiveMembers   int64     `json:"active_members"`
	MessageCount24h int64     `json:"message_count_24h"`
	AverageResTime  float64   `json:"average_response_time"` // in milliseconds
}

// UserAnalytics tracks metrics for a user
type UserAnalytics struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	UserID           uint      `json:"user_id"`
	User             User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TotalMessages    int64     `json:"total_messages"`
	TotalChats       int64     `json:"total_chats"`
	LastActive       time.Time `json:"last_active"`
	MessageCount24h  int64     `json:"message_count_24h"`
	AverageResTime   float64   `json:"average_response_time"` // in milliseconds
	OnlineTime24h    int64     `json:"online_time_24h"` // in seconds
}

// SystemAnalytics tracks overall system metrics
type SystemAnalytics struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	TotalUsers            int64     `json:"total_users"`
	ActiveUsers24h        int64     `json:"active_users_24h"`
	ActiveUsersOnline     int64     `json:"active_users_online"`
	TotalChats            int64     `json:"total_chats"`
	TotalMessages         int64     `json:"total_messages"`
	MessagesPerMinute     float64   `json:"messages_per_minute"`
	AverageResponseTime   float64   `json:"average_response_time"` // in milliseconds
	NewUsersToday         int64     `json:"new_users_today"`
	AverageChatSize       float64   `json:"average_chat_size"`
}

// TableName defines custom table names for analytics models
func (ChatAnalytics) TableName() string {
	return "chat_analytics"
}

func (UserAnalytics) TableName() string {
	return "user_analytics"
}

func (SystemAnalytics) TableName() string {
	return "system_analytics"
}
