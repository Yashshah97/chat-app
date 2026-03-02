package main

import (
	"time"
)

// ChatSettings represents configuration options for a chat
type ChatSettings struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	ChatID            uint      `json:"chat_id"`
	Chat              Chat      `gorm:"foreignKey:ChatID" json:"chat,omitempty"`
	AllowNotifications bool      `json:"allow_notifications"`
	AllowMentions     bool      `json:"allow_mentions"`
	MuteNotifications bool      `json:"mute_notifications"`
	AutoArchive       bool      `json:"auto_archive"`
	ArchiveAfterDays  int       `json:"archive_after_days"` // 0 = disabled
	AllowFileSharing  bool      `json:"allow_file_sharing"`
	AllowVoice        bool      `json:"allow_voice"`
	AllowVideo        bool      `json:"allow_video"`
	DisappearingMsgs  int       `json:"disappearing_msgs"` // 0 = disabled, seconds
	Encryption        bool      `json:"encryption"`
	ReadReceiptEnabled bool     `json:"read_receipt_enabled"`
	TypingIndicator   bool      `json:"typing_indicator"`
}

// UserChatPreference represents user-specific settings for a chat
type UserChatPreference struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	UserID          uint      `json:"user_id"`
	User            User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ChatID          uint      `json:"chat_id"`
	Chat            Chat      `gorm:"foreignKey:ChatID" json:"chat,omitempty"`
	IsMuted         bool      `json:"is_muted"`
	MuteUntil       *time.Time `json:"mute_until,omitempty"`
	IsPinned        bool      `json:"is_pinned"`
	IsArchived      bool      `json:"is_archived"`
	NotificationType string    `json:"notification_type"` // all, mentions_only, none
	SoundEnabled    bool      `json:"sound_enabled"`
	VibrationEnabled bool     `json:"vibration_enabled"`
	CustomColor     *string   `json:"custom_color,omitempty"`
}

// NotificationPreference represents user notification settings
type NotificationPreference struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	UserID          uint      `json:"user_id"`
	User            User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AllowPush       bool      `json:"allow_push"`
	AllowEmail      bool      `json:"allow_email"`
	AllowSound      bool      `json:"allow_sound"`
	AllowVibration  bool      `json:"allow_vibration"`
	QuietHoursStart string    `json:"quiet_hours_start"` // HH:MM format
	QuietHoursEnd   string    `json:"quiet_hours_end"`   // HH:MM format
	MuteAllChats    bool      `json:"mute_all_chats"`
}

// TableName defines custom table names
func (ChatSettings) TableName() string {
	return "chat_settings"
}

func (UserChatPreference) TableName() string {
	return "user_chat_preferences"
}

func (NotificationPreference) TableName() string {
	return "notification_preferences"
}
