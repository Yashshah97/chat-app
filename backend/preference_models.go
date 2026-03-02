package main

import "gorm.io/gorm"

// UserPreference stores user preferences
type UserPreference struct {
	gorm.Model
	UserID    uint   `gorm:"not null;uniqueIndex" json:"user_id"`
	User      *User  `gorm:"foreignKey:UserID" json:"-"`
	Theme     string `json:"theme"` // light, dark, auto
	Language  string `json:"language"`
	Privacy   string `json:"privacy"` // public, friends, private
	Notifications bool `gorm:"default:true" json:"notifications"`
	EmailDigest bool `gorm:"default:true" json:"email_digest"`
	ShowActivity bool `gorm:"default:true" json:"show_activity"`
	TwoFAEnabled bool `gorm:"default:false" json:"two_fa_enabled"`
}

// Survey represents user surveys
type Survey struct {
	gorm.Model
	Title       string `gorm:"not null" json:"title"`
	Description string `json:"description"`
	CreatedByID uint   `json:"created_by_id"`
	CreatedBy   *User  `gorm:"foreignKey:CreatedByID" json:"-"`
	Status      string `gorm:"default:'active'" json:"status"` // active, closed, archived
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Responses   int64  `gorm:"default:0" json:"responses"`
}

// SurveyQuestion represents survey questions
type SurveyQuestion struct {
	gorm.Model
	SurveyID    uint   `gorm:"not null;index" json:"survey_id"`
	Survey      *Survey `gorm:"foreignKey:SurveyID" json:"-"`
	Question    string `json:"question"`
	Type        string `json:"type"` // text, multiple_choice, rating, likert
	Required    bool   `gorm:"default:true" json:"required"`
	Order       int    `json:"order"`
	Options     string `json:"options"` // JSON array
}

// SurveyResponse represents survey responses
type SurveyResponse struct {
	gorm.Model
	SurveyID   uint   `gorm:"not null;index" json:"survey_id"`
	Survey     *Survey `gorm:"foreignKey:SurveyID" json:"-"`
	UserID     *uint  `json:"user_id"`
	User       *User  `gorm:"foreignKey:UserID" json:"-"`
	Answers    string `gorm:"type:jsonb" json:"answers"`
	CompletedAt string `json:"completed_at"`
}

// Recommendation represents content recommendations
type Recommendation struct {
	gorm.Model
	UserID      uint   `gorm:"not null;index" json:"user_id"`
	User        *User  `gorm:"foreignKey:UserID" json:"-"`
	ChatID      *uint  `json:"chat_id"`
	Chat        *Chat  `gorm:"foreignKey:ChatID" json:"-"`
	UserIDTarget *uint `json:"user_id_target"`
	Type        string `json:"type"` // chat, user, content
	Score       float64 `json:"score"`
	Reason      string `json:"reason"`
	Clicked     bool   `gorm:"default:false" json:"clicked"`
}
