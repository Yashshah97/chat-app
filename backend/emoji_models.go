package main

import "gorm.io/gorm"

// EmojiPack represents a collection of emojis/stickers
type EmojiPack struct {
	gorm.Model
	Name        string `gorm:"not null;uniqueIndex" json:"name"`
	Description string `json:"description"`
	CreatedByID uint   `json:"created_by_id"`
	CreatedBy   *User  `gorm:"foreignKey:CreatedByID" json:"-"`
	IsPublic    bool   `gorm:"default:true" json:"is_public"`
	Version     string `json:"version"`
	IconURL     string `json:"icon_url"`
	Downloads   int64  `gorm:"default:0" json:"downloads"`
}

// Emoji represents a single emoji/sticker in a pack
type Emoji struct {
	gorm.Model
	PackID      uint       `gorm:"not null;index" json:"pack_id"`
	Pack        *EmojiPack `gorm:"foreignKey:PackID" json:"-"`
	Code        string     `gorm:"not null" json:"code"`
	Alias       string     `json:"alias"`
	ImageURL    string     `gorm:"not null" json:"image_url"`
	Category    string     `json:"category"`
	Tags        string     `json:"tags"` // comma-separated
	SearchScore int        `json:"search_score"`
}

// UserEmojiPack tracks which emoji packs a user has subscribed to
type UserEmojiPack struct {
	gorm.Model
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	User      *User      `gorm:"foreignKey:UserID" json:"-"`
	EmojiPack uint       `gorm:"not null;index" json:"emoji_pack_id"`
	Pack      *EmojiPack `gorm:"foreignKey:EmojiPack" json:"-"`
	IsFavorite bool      `gorm:"default:false" json:"is_favorite"`
}

// MessageEmoji tracks emoji reactions on messages
type MessageEmoji struct {
	gorm.Model
	MessageID uint   `gorm:"not null;index:idx_message_emoji" json:"message_id"`
	Message   *Message `gorm:"foreignKey:MessageID" json:"-"`
	UserID    uint   `gorm:"not null;index:idx_message_emoji" json:"user_id"`
	User      *User  `gorm:"foreignKey:UserID" json:"-"`
	EmojiCode string `gorm:"not null;index:idx_message_emoji" json:"emoji_code"`
	Count     int    `gorm:"default:1" json:"count"`
}

// EmojiPackReview represents user reviews of emoji packs
type EmojiPackReview struct {
	gorm.Model
	PackID     uint       `gorm:"not null;index" json:"pack_id"`
	Pack       *EmojiPack `gorm:"foreignKey:PackID" json:"-"`
	UserID     uint       `gorm:"not null;index" json:"user_id"`
	User       *User      `gorm:"foreignKey:UserID" json:"-"`
	Rating     int        `gorm:"not null" json:"rating"` // 1-5
	Comment    string     `json:"comment"`
	Helpful    int        `gorm:"default:0" json:"helpful"`
	Unhelpful  int        `gorm:"default:0" json:"unhelpful"`
}
