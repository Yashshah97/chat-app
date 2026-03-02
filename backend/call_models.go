package main

import "gorm.io/gorm"

// VoiceCall represents a voice/audio call
type VoiceCall struct {
	gorm.Model
	ChatID      uint    `gorm:"not null;index" json:"chat_id"`
	Chat        *Chat   `gorm:"foreignKey:ChatID" json:"-"`
	InitiatorID uint    `json:"initiator_id"`
	Initiator   *User   `gorm:"foreignKey:InitiatorID" json:"-"`
	Status      string  `gorm:"default:'ringing'" json:"status"` // ringing, active, ended, missed, rejected
	Duration    int     `json:"duration"` // seconds
	StartedAt   *string `json:"started_at"`
	EndedAt     *string `json:"ended_at"`
	ParticipantCount int `json:"participant_count"`
	Quality     string  `json:"quality"` // poor, fair, good, excellent
	RecordingURL *string `json:"recording_url"`
	IsRecorded  bool    `gorm:"default:false" json:"is_recorded"`
	CallType    string  `json:"call_type"` // one_to_one, group, conference
}

// CallParticipant represents a participant in a call
type CallParticipant struct {
	gorm.Model
	CallID    uint    `gorm:"not null;index" json:"call_id"`
	Call      *VoiceCall `gorm:"foreignKey:CallID" json:"-"`
	UserID    uint    `json:"user_id"`
	User      *User   `gorm:"foreignKey:UserID" json:"-"`
	JoinedAt  string  `json:"joined_at"`
	LeftAt    *string `json:"left_at"`
	Duration  int     `json:"duration"`
	IsMuted   bool    `gorm:"default:false" json:"is_muted"`
	Video     bool    `gorm:"default:false" json:"video"`
	Quality   string  `json:"quality"`
}

// CallLog tracks all call information
type CallLog struct {
	gorm.Model
	CallID    uint    `gorm:"not null;index" json:"call_id"`
	Call      *VoiceCall `gorm:"foreignKey:CallID" json:"-"`
	EventType string  `json:"event_type"` // started, muted, video_on, video_off, participant_joined, participant_left, ended
	Timestamp string  `json:"timestamp"`
	UserID    *uint   `json:"user_id"`
	Metadata  string  `json:"metadata"`
}
