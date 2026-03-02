package main

import (
	"time"
)

// MessageEdit represents an edit history entry for a message
type MessageEdit struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	MessageID uint      `json:"message_id"`
	Message   Message   `gorm:"foreignKey:MessageID" json:"message,omitempty"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	OldContent string    `json:"old_content"`
	NewContent string    `json:"new_content"`
	EditReason string    `json:"edit_reason"`
}

// TableName defines custom table name
func (MessageEdit) TableName() string {
	return "message_edits"
}
