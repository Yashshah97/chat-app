package main

import (
	"time"
)

// UserActivityLog tracks user actions for audit purposes
type UserActivityLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action    string    `json:"action"` // login, logout, create_message, edit_message, etc
	ResourceType string `json:"resource_type"` // message, chat, user, file
	ResourceID   *uint  `json:"resource_id,omitempty"`
	Metadata   string    `json:"metadata"` // JSON string with additional data
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	Status     string    `json:"status"` // success, failure
}

// TableName defines custom table name
func (UserActivityLog) TableName() string {
	return "user_activity_logs"
}
