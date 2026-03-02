package main

import "gorm.io/gorm"

// Integration represents a third-party integration
type Integration struct {
	gorm.Model
	Name        string `gorm:"not null;uniqueIndex" json:"name"`
	Type        string `json:"type"` // slack, discord, github, jira, trello, etc
	APIKey      string `json:"api_key"`
	APISecret   string `json:"api_secret"`
	WebhookURL  string `json:"webhook_url"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	CreatedByID uint   `json:"created_by_id"`
	CreatedBy   *User  `gorm:"foreignKey:CreatedByID" json:"-"`
	Config      string `gorm:"type:jsonb" json:"config"`
	LastSyncAt  *string `json:"last_sync_at"`
	ErrorCount  int    `gorm:"default:0" json:"error_count"`
}

// IntegrationEvent logs integration events
type IntegrationEvent struct {
	gorm.Model
	IntegrationID uint   `gorm:"not null;index" json:"integration_id"`
	Integration   *Integration `gorm:"foreignKey:IntegrationID" json:"-"`
	EventType     string `json:"event_type"` // sync, error, webhook_received
	Status        string `json:"status"` // success, failed, pending
	Payload       string `gorm:"type:jsonb" json:"payload"`
	Error         string `json:"error"`
	Retries       int    `json:"retries"`
}

// IntegrationMapping represents a mapping between chat and external service
type IntegrationMapping struct {
	gorm.Model
	IntegrationID uint   `gorm:"not null;index" json:"integration_id"`
	Integration   *Integration `gorm:"foreignKey:IntegrationID" json:"-"`
	ChatID        uint   `gorm:"not null;index" json:"chat_id"`
	Chat          *Chat  `gorm:"foreignKey:ChatID" json:"-"`
	ExternalID    string `json:"external_id"` // Slack channel ID, Discord server ID, etc
	ExternalName  string `json:"external_name"`
	SyncMessages  bool   `gorm:"default:false" json:"sync_messages"`
	SyncUsers     bool   `gorm:"default:false" json:"sync_users"`
	Direction     string `json:"direction"` // one_way, two_way
}

// APIKeyConfig represents stored API keys for integrations
type APIKeyConfig struct {
	gorm.Model
	IntegrationID uint   `gorm:"not null;index" json:"integration_id"`
	Integration   *Integration `gorm:"foreignKey:IntegrationID" json:"-"`
	Key           string `gorm:"not null" json:"key"`
	Secret        string `gorm:"not null" json:"secret"`
	IsEncrypted   bool   `gorm:"default:true" json:"is_encrypted"`
	ExpiresAt     *string `json:"expires_at"`
	Scope         string `json:"scope"` // Permissions granted
	LastUsedAt    *string `json:"last_used_at"`
}
