package main

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `gorm:"uniqueIndex" json:"username"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Password  string    `json:"-"`
	Avatar    string    `json:"avatar"`
	Status    string    `json:"status"` // online, offline, away
	Messages  []Message `gorm:"foreignKey:UserID" json:"-"`
	Chats     []Chat    `gorm:"many2many:chat_members;" json:"-"`
}

type Chat struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // private, group
	Messages  []Message `gorm:"foreignKey:ChatID" json:"messages,omitempty"`
	Members   []User    `gorm:"many2many:chat_members;" json:"members,omitempty"`
	CreatedBy uint      `json:"created_by"`
}

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ChatID    uint      `json:"chat_id"`
	Chat      Chat      `gorm:"foreignKey:ChatID" json:"-"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"-"`
	Content   string    `json:"content"`
	Type      string    `json:"type"` // text, image, file
	IsEdited  bool      `json:"is_edited"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName defines custom table names
func (User) TableName() string {
	return "users"
}

func (Chat) TableName() string {
	return "chats"
}

func (Message) TableName() string {
	return "messages"
}
