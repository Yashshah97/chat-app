package main

import (
	"time"
)

// Notification represents a notification sent to a user
type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type      string    `json:"type"` // message, mention, reaction, admin, system
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	RefID     *uint     `json:"ref_id,omitempty"` // references message_id, chat_id, etc
	RefType   string    `json:"ref_type"` // message, chat, user, admin
	IsRead    bool      `json:"is_read"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	Priority  string    `json:"priority"` // low, normal, high, critical
}

// NotificationDelivery tracks notification delivery across channels
type NotificationDelivery struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	NotificationID   uint      `json:"notification_id"`
	Notification     Notification `gorm:"foreignKey:NotificationID" json:"notification,omitempty"`
	Channel          string    `json:"channel"` // push, email, sms, in_app
	Status           string    `json:"status"` // pending, sent, failed, delivered
	SentAt           *time.Time `json:"sent_at,omitempty"`
	DeliveredAt      *time.Time `json:"delivered_at,omitempty"`
	FailureReason    string    `json:"failure_reason"`
	RetryCount       int       `json:"retry_count"`
	ExternalID       string    `json:"external_id"` // provider's tracking ID
}

// NotificationTemplate represents reusable notification templates
type NotificationTemplate struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `gorm:"uniqueIndex" json:"name"`
	Type      string    `json:"type"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"` // supports {{variables}}
	Variables []string  `gorm:"-" json:"variables"` // extracted variable names
	Channels  string    `json:"channels"` // comma-separated: push,email,sms
	IsActive  bool      `json:"is_active"`
}

// TableName defines custom table names
func (Notification) TableName() string {
	return "notifications"
}

func (NotificationDelivery) TableName() string {
	return "notification_deliveries"
}

func (NotificationTemplate) TableName() string {
	return "notification_templates"
}
