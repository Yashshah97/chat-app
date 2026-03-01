package main

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// MessageReaction represents a reaction to a message
type MessageReaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MessageID uint      `json:"message_id"`
	Message   Message   `gorm:"foreignKey:MessageID" json:"message,omitempty"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
}

// TypingIndicator represents a user's typing status
type TypingIndicator struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ChatID    uint      `json:"chat_id"`
	Chat      Chat      `gorm:"foreignKey:ChatID" json:"chat,omitempty"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	IsTyping  bool      `json:"is_typing"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MessageReactionService handles message reactions
type MessageReactionService struct {
	db *gorm.DB
}

// NewMessageReactionService creates a new message reaction service
func NewMessageReactionService(db *gorm.DB) *MessageReactionService {
	return &MessageReactionService{db: db}
}

// AddReaction adds a reaction to a message
func (s *MessageReactionService) AddReaction(messageID, userID uint, emoji string) (*MessageReaction, error) {
	log.Printf("Adding reaction %s to message %d by user %d", emoji, messageID, userID)

	reaction := &MessageReaction{
		MessageID: messageID,
		UserID:    userID,
		Emoji:     emoji,
		CreatedAt: time.Now(),
	}

	result := s.db.Create(reaction)
	if result.Error != nil {
		log.Printf("Error adding reaction: %v", result.Error)
		return nil, result.Error
	}

	return reaction, nil
}

// RemoveReaction removes a reaction from a message
func (s *MessageReactionService) RemoveReaction(reactionID uint) error {
	log.Printf("Removing reaction %d", reactionID)

	result := s.db.Delete(&MessageReaction{}, reactionID)
	if result.Error != nil {
		log.Printf("Error removing reaction: %v", result.Error)
		return result.Error
	}

	return nil
}

// GetMessageReactions gets all reactions for a message
func (s *MessageReactionService) GetMessageReactions(messageID uint) ([]MessageReaction, error) {
	log.Printf("Getting reactions for message %d", messageID)

	var reactions []MessageReaction
	result := s.db.Where("message_id = ?", messageID).
		Preload("User").
		Find(&reactions)

	if result.Error != nil {
		log.Printf("Error fetching reactions: %v", result.Error)
		return nil, result.Error
	}

	return reactions, nil
}

// TypingIndicatorService handles typing indicators
type TypingIndicatorService struct {
	db *gorm.DB
}

// NewTypingIndicatorService creates a new typing indicator service
func NewTypingIndicatorService(db *gorm.DB) *TypingIndicatorService {
	return &TypingIndicatorService{db: db}
}

// UpdateTypingStatus updates user's typing status
func (s *TypingIndicatorService) UpdateTypingStatus(chatID, userID uint, isTyping bool) (*TypingIndicator, error) {
	log.Printf("Updating typing status for chat %d, user %d: %v", chatID, userID, isTyping)

	indicator := &TypingIndicator{}

	// Check if exists
	result := s.db.Where("chat_id = ? AND user_id = ?", chatID, userID).First(indicator)

	if result.Error == gorm.ErrRecordNotFound {
		// Create new
		indicator = &TypingIndicator{
			ChatID:    chatID,
			UserID:    userID,
			IsTyping:  isTyping,
			UpdatedAt: time.Now(),
		}
		if err := s.db.Create(indicator).Error; err != nil {
			log.Printf("Error creating typing indicator: %v", err)
			return nil, err
		}
	} else if result.Error != nil {
		log.Printf("Error fetching typing indicator: %v", result.Error)
		return nil, result.Error
	} else {
		// Update existing
		indicator.IsTyping = isTyping
		indicator.UpdatedAt = time.Now()
		if err := s.db.Save(indicator).Error; err != nil {
			log.Printf("Error updating typing indicator: %v", err)
			return nil, err
		}
	}

	return indicator, nil
}

// GetTypingUsers gets users currently typing in a chat
func (s *TypingIndicatorService) GetTypingUsers(chatID uint) ([]TypingIndicator, error) {
	log.Printf("Getting typing users for chat %d", chatID)

	var indicators []TypingIndicator
	result := s.db.Where("chat_id = ? AND is_typing = ?", chatID, true).
		Preload("User").
		Find(&indicators)

	if result.Error != nil {
		log.Printf("Error fetching typing users: %v", result.Error)
		return nil, result.Error
	}

	return indicators, nil
}

// ClearTypingStatus clears typing status for a user
func (s *TypingIndicatorService) ClearTypingStatus(chatID, userID uint) error {
	log.Printf("Clearing typing status for chat %d, user %d", chatID, userID)

	result := s.db.Where("chat_id = ? AND user_id = ?", chatID, userID).
		Update("is_typing", false)

	if result.Error != nil {
		log.Printf("Error clearing typing status: %v", result.Error)
		return result.Error
	}

	return nil
}
