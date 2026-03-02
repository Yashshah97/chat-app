package main

import "gorm.io/gorm"

// ChatBackup represents a backup of a chat
type ChatBackup struct {
	gorm.Model
	ChatID      uint   `gorm:"not null;index" json:"chat_id"`
	Chat        *Chat  `gorm:"foreignKey:ChatID" json:"-"`
	CreatedByID uint   `json:"created_by_id"`
	CreatedBy   *User  `gorm:"foreignKey:CreatedByID" json:"-"`
	Format      string `gorm:"not null" json:"format"` // json, csv, pdf, html
	FileURL     string `json:"file_url"`
	FileName    string `json:"file_name"`
	BackupSize  int64  `json:"backup_size"` // bytes
	MessageCount int64 `json:"message_count"`
	IncludeMedia bool  `gorm:"default:false" json:"include_media"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
	Status      string `gorm:"default:'pending'" json:"status"` // pending, processing, completed, failed
	Progress    int    `gorm:"default:0" json:"progress"` // 0-100
}

// ChatExport represents an on-demand export of chat data
type ChatExport struct {
	gorm.Model
	ChatID      uint   `gorm:"not null;index" json:"chat_id"`
	Chat        *Chat  `gorm:"foreignKey:ChatID" json:"-"`
	UserID      uint   `gorm:"not null;index" json:"user_id"`
	User        *User  `gorm:"foreignKey:UserID" json:"-"`
	Format      string `gorm:"not null" json:"format"` // json, csv, xlsx, pdf
	Status      string `gorm:"default:'pending'" json:"status"` // pending, ready, expired
	DownloadURL string `json:"download_url"`
	ExpiresAt   *string `json:"expires_at"`
	FilterType  string `json:"filter_type"` // all, date_range, members
	FilterValue string `json:"filter_value"` // json encoded
	DownloadCount int `gorm:"default:0" json:"download_count"`
}

// BackupSchedule represents scheduled backups for a chat
type BackupSchedule struct {
	gorm.Model
	ChatID      uint   `gorm:"not null;uniqueIndex" json:"chat_id"`
	Chat        *Chat  `gorm:"foreignKey:ChatID" json:"-"`
	CreatedByID uint   `json:"created_by_id"`
	CreatedBy   *User  `gorm:"foreignKey:CreatedByID" json:"-"`
	Frequency   string `gorm:"not null" json:"frequency"` // daily, weekly, monthly
	Format      string `gorm:"not null" json:"format"` // json, csv, pdf
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	LastRunAt   *string `json:"last_run_at"`
	NextRunAt   string `json:"next_run_at"`
	MaxBackups  int    `gorm:"default:10" json:"max_backups"` // Keep last N backups
	IncludeMedia bool  `gorm:"default:false" json:"include_media"`
	StorageLocation string `json:"storage_location"` // local, s3, gdrive, dropbox
}

// ArchiveMessage represents messages that have been archived/deleted from active view
type ArchiveMessage struct {
	gorm.Model
	MessageID   uint   `gorm:"not null;uniqueIndex" json:"message_id"`
	Message     *Message `gorm:"foreignKey:MessageID" json:"-"`
	ChatID      uint   `gorm:"not null;index" json:"chat_id"`
	Chat        *Chat  `gorm:"foreignKey:ChatID" json:"-"`
	ArchivedByID uint  `json:"archived_by_id"`
	ArchivedBy  *User  `gorm:"foreignKey:ArchivedByID" json:"-"`
	OriginalData string `json:"original_data"` // JSON blob of original message
	Reason      string `json:"reason"` // deleted, archived, purged
}
