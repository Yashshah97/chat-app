package main

import (
	"time"
)

// UserPresence tracks user online/offline status and activity
type UserPresence struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	UserID          uint      `json:"user_id"`
	User            User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status          string    `json:"status"` // online, offline, away, dnd (do not disturb)
	LastSeen        time.Time `json:"last_seen"`
	CurrentChatID   *uint     `json:"current_chat_id"` // nullable - which chat is user currently in
	OnlineAtTime    time.Time `json:"online_at_time"`
	SessionID       string    `json:"session_id"` // for tracking multiple sessions
	DeviceType      string    `json:"device_type"` // web, mobile, desktop
	IPAddress       string    `json:"ip_address"`
}

// PresenceHistory tracks historical presence data for analytics
type PresenceHistory struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UserID        uint      `json:"user_id"`
	User          User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Status        string    `json:"status"`
	DurationSecs  int64     `json:"duration_secs"` // how long user was in this status
	ChatID        *uint     `json:"chat_id"`
	EventType     string    `json:"event_type"` // online, offline, away, active_chat
}

// TableName defines custom table names
func (UserPresence) TableName() string {
	return "user_presence"
}

func (PresenceHistory) TableName() string {
	return "presence_history"
}
