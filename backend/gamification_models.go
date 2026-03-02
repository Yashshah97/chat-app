package main

import "gorm.io/gorm"

// Badge represents user badges/achievements
type Badge struct {
	gorm.Model
	Name        string `gorm:"not null;uniqueIndex" json:"name"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	Category    string `json:"category"` // achievement, activity, role, special
	Condition   string `json:"condition"` // What to earn it
	Rarity      string `json:"rarity"` // common, rare, epic, legendary
}

// UserBadge represents badges earned by users
type UserBadge struct {
	gorm.Model
	UserID    uint   `gorm:"not null;index" json:"user_id"`
	User      *User  `gorm:"foreignKey:UserID" json:"-"`
	BadgeID   uint   `gorm:"not null;index" json:"badge_id"`
	Badge     *Badge `gorm:"foreignKey:BadgeID" json:"-"`
	EarnedAt  string `json:"earned_at"`
	Progress  int    `json:"progress"`
}

// Invitation represents chat/group invitations
type Invitation struct {
	gorm.Model
	ChatID      uint   `gorm:"not null;index" json:"chat_id"`
	Chat        *Chat  `gorm:"foreignKey:ChatID" json:"-"`
	InvitedByID uint   `json:"invited_by_id"`
	InvitedBy   *User  `gorm:"foreignKey:InvitedByID" json:"-"`
	InvitedUserID *uint `json:"invited_user_id"`
	InvitedUser *User  `gorm:"foreignKey:InvitedUserID" json:"-"`
	Email       string `gorm:"index" json:"email"`
	InviteToken string `gorm:"not null;uniqueIndex" json:"invite_token"`
	Status      string `gorm:"default:'pending'" json:"status"` // pending, accepted, declined, expired
	ExpiresAt   string `json:"expires_at"`
	AcceptedAt  *string `json:"accepted_at"`
}

// TrendingTopic represents trending topics
type TrendingTopic struct {
	gorm.Model
	Topic       string `gorm:"not null;index" json:"topic"`
	Mentions    int64  `json:"mentions"`
	Growth      float64 `json:"growth"` // percentage change
	Trend       string `json:"trend"` // up, down, stable
	LastUpdated string `json:"last_updated"`
	Rank        int    `json:"rank"`
	Category    string `json:"category"`
}

// SummaryStatistics represents aggregated statistics
type SummaryStatistics struct {
	gorm.Model
	Period          string `json:"period"` // day, week, month, year
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	TotalUsers      int64  `json:"total_users"`
	ActiveUsers     int64  `json:"active_users"`
	TotalChats      int64  `json:"total_chats"`
	TotalMessages   int64  `json:"total_messages"`
	NewUsers        int64  `json:"new_users"`
	ChatEngagement  float64 `json:"chat_engagement"`
	MessageGrowth   float64 `json:"message_growth"`
	RetentionRate   float64 `json:"retention_rate"`
	AverageSessionLength int `json:"average_session_length"`
}
