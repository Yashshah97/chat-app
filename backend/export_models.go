package main

import "gorm.io/gorm"

// DataExport represents user data export requests
type DataExport struct {
	gorm.Model
	UserID      uint   `gorm:"not null;index" json:"user_id"`
	User        *User  `gorm:"foreignKey:UserID" json:"-"`
	Format      string `json:"format"` // json, csv, pdf, zip
	IncludeData string `json:"include_data"` // messages, chats, profile, files
	Status      string `gorm:"default:'pending'" json:"status"` // pending, processing, ready, failed, expired
	FileURL     string `json:"file_url"`
	ExpiresAt   string `json:"expires_at"`
	DownloadCount int  `json:"download_count"`
	DataSize    int64  `json:"data_size"`
}

// FileManagement represents file organization and management
type FileManagement struct {
	gorm.Model
	ChatID      uint   `gorm:"not null;index" json:"chat_id"`
	Chat        *Chat  `gorm:"foreignKey:ChatID" json:"-"`
	FileName    string `json:"file_name"`
	FileSize    int64  `json:"file_size"`
	FileType    string `json:"file_type"`
	FileHash    string `json:"file_hash"`
	StorageKey  string `json:"storage_key"` // S3/cloud storage key
	UploadedByID uint  `json:"uploaded_by_id"`
	UploadedBy  *User  `gorm:"foreignKey:UploadedByID" json:"-"`
	IsShared    bool   `gorm:"default:false" json:"is_shared"`
	Tags        string `json:"tags"`
	Description string `json:"description"`
	Downloads   int64  `gorm:"default:0" json:"downloads"`
}

// UserDataRequest represents GDPR data access requests
type UserDataRequest struct {
	gorm.Model
	UserID      uint   `gorm:"not null;index" json:"user_id"`
	User        *User  `gorm:"foreignKey:UserID" json:"-"`
	RequestType string `json:"request_type"` // access, delete, rectify, port
	Status      string `gorm:"default:'pending'" json:"status"` // pending, approved, denied, completed
	Reason      string `json:"reason"`
	ApprovedAt  *string `json:"approved_at"`
	CompletedAt *string `json:"completed_at"`
	DataURL     string `json:"data_url"`
	RequestedAt string `json:"requested_at"`
	DeadlineAt  string `json:"deadline_at"`
}

// AuditLogEntry represents detailed audit logs
type AuditLogEntry struct {
	gorm.Model
	UserID      uint   `gorm:"not null;index" json:"user_id"`
	User        *User  `gorm:"foreignKey:UserID" json:"-"`
	Action      string `gorm:"index" json:"action"`
	ResourceType string `json:"resource_type"`
	ResourceID  *uint  `json:"resource_id"`
	OldValue    string `json:"old_value"`
	NewValue    string `json:"new_value"`
	IPAddress   string `json:"ip_address"`
	UserAgent   string `json:"user_agent"`
	Status      string `json:"status"`
	Details     string `json:"details"`
}
