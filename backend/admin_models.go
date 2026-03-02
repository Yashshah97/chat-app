package main

import (
	"time"
)

// AdminUser represents an admin with special privileges
type AdminUser struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role      string    `json:"role"` // super_admin, moderator, support
	Reason    string    `json:"reason"`
}

// AdminAction tracks admin actions for audit purposes
type AdminAction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	AdminID   uint      `json:"admin_id"`
	Admin     AdminUser `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	Action    string    `json:"action"` // suspend_user, delete_chat, warn_user, etc
	TargetID  uint      `json:"target_id"`
	TargetType string   `json:"target_type"` // user, chat, message
	Reason    string    `json:"reason"`
}

// UserReport represents user reports for moderation
type UserReport struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ReporterID uint     `json:"reporter_id"`
	Reporter   User     `gorm:"foreignKey:ReporterID" json:"reporter,omitempty"`
	TargetID  uint      `json:"target_id"`
	Target    User     `gorm:"foreignKey:TargetID" json:"target,omitempty"`
	Reason    string    `json:"reason"`
	Status    string    `json:"status"` // pending, investigating, resolved, dismissed
	Resolution string   `json:"resolution"`
}

// TableName defines custom table names for admin models
func (AdminUser) TableName() string {
	return "admin_users"
}

func (AdminAction) TableName() string {
	return "admin_actions"
}

func (UserReport) TableName() string {
	return "user_reports"
}
