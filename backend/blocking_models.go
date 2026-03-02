package main

import (
	"time"
)

// BlockedUser represents a user blocked by another user
type BlockedUser struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	BlockerID uint      `json:"blocker_id"`
	Blocker   User      `gorm:"foreignKey:BlockerID" json:"blocker,omitempty"`
	BlockedID uint      `json:"blocked_id"`
	Blocked   User      `gorm:"foreignKey:BlockedID" json:"blocked,omitempty"`
	Reason    string    `json:"reason"`
}

// MutedUser represents a user muted by another user
type MutedUser struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	MuterID   uint      `json:"muter_id"`
	Muter     User      `gorm:"foreignKey:MuterID" json:"muter,omitempty"`
	MutedID   uint      `json:"muted_id"`
	Muted     User      `gorm:"foreignKey:MutedID" json:"muted,omitempty"`
	MuteUntil *time.Time `json:"mute_until,omitempty"`
	Reason    string    `json:"reason"`
}

// TableName defines custom table names
func (BlockedUser) TableName() string {
	return "blocked_users"
}

func (MutedUser) TableName() string {
	return "muted_users"
}
