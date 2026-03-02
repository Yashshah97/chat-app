package main

import "gorm.io/gorm"

// Webhook represents an outgoing webhook configuration
type Webhook struct {
	gorm.Model
	Name        string `gorm:"not null;index" json:"name"`
	URL         string `gorm:"not null" json:"url"`
	ChatID      *uint  `gorm:"index" json:"chat_id"`
	Chat        *Chat  `gorm:"foreignKey:ChatID" json:"-"`
	CreatedByID uint   `json:"created_by_id"`
	CreatedBy   *User  `gorm:"foreignKey:CreatedByID" json:"-"`
	Events      string `json:"events"` // comma-separated: message.created, message.edited, user.joined, etc
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	Secret      string `json:"secret"` // HMAC secret for verification
	Headers     string `json:"headers"` // JSON object of custom headers
	RateLimit   int    `gorm:"default:100" json:"rate_limit"` // per minute
	Retries     int    `gorm:"default:3" json:"retries"`
	RetryDelay  int    `gorm:"default:5" json:"retry_delay"` // seconds
	Timeout     int    `gorm:"default:30" json:"timeout"` // seconds
}

// WebhookEvent represents an event sent to a webhook
type WebhookEvent struct {
	gorm.Model
	WebhookID   uint   `gorm:"not null;index" json:"webhook_id"`
	Webhook     *Webhook `gorm:"foreignKey:WebhookID" json:"-"`
	EventType   string `json:"event_type"`
	Payload     string `gorm:"type:jsonb" json:"payload"`
	Status      string `gorm:"default:'pending'" json:"status"` // pending, sent, failed, delivered
	StatusCode  int    `json:"status_code"`
	Response    string `json:"response"`
	Error       string `json:"error"`
	Attempts    int    `gorm:"default:0" json:"attempts"`
	NextRetryAt *string `json:"next_retry_at"`
	SentAt      *string `json:"sent_at"`
}

// WebhookLog represents detailed logs of webhook calls
type WebhookLog struct {
	gorm.Model
	WebhookID   uint   `gorm:"not null;index" json:"webhook_id"`
	Webhook     *Webhook `gorm:"foreignKey:WebhookID" json:"-"`
	EventID     *uint  `json:"event_id"`
	Event       *WebhookEvent `gorm:"foreignKey:EventID" json:"-"`
	Method      string `json:"method"`
	URL         string `json:"url"`
	RequestBody string `json:"request_body"`
	ResponseBody string `json:"response_body"`
	StatusCode  int    `json:"status_code"`
	Duration    int    `json:"duration"` // milliseconds
	Success     bool   `json:"success"`
	Error       string `json:"error"`
	IPAddress   string `json:"ip_address"`
}

// IncomingWebhook represents a webhook integration that receives data
type IncomingWebhook struct {
	gorm.Model
	ChatID      uint   `gorm:"not null;index" json:"chat_id"`
	Chat        *Chat  `gorm:"foreignKey:ChatID" json:"-"`
	CreatedByID uint   `json:"created_by_id"`
	CreatedBy   *User  `gorm:"foreignKey:CreatedByID" json:"-"`
	Name        string `json:"name"`
	Token       string `gorm:"not null;uniqueIndex" json:"token"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	AllowedIP   string `json:"allowed_ip"` // comma-separated IPs
	WebhookURL  string `json:"webhook_url"`
	Format      string `json:"format"` // slack, discord, generic
}
