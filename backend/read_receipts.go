package main

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// ReadReceipt represents a message read receipt
type ReadReceipt struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MessageID uint      `json:"message_id"`
	Message   Message   `gorm:"foreignKey:MessageID" json:"message,omitempty"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ReadAt    time.Time `json:"read_at"`
	CreatedAt time.Time `json:"created_at"`
}

// ReadReceiptService handles read receipts
type ReadReceiptService struct {
	db *gorm.DB
}

// NewReadReceiptService creates a new read receipt service
func NewReadReceiptService(db *gorm.DB) *ReadReceiptService {
	return &ReadReceiptService{db: db}
}

// MarkAsRead marks a message as read by a user
func (s *ReadReceiptService) MarkAsRead(messageID, userID uint) (*ReadReceipt, error) {
	log.Printf("Marking message %d as read by user %d", messageID, userID)

	readReceipt := &ReadReceipt{}

	// Check if receipt exists
	result := s.db.Where("message_id = ? AND user_id = ?", messageID, userID).First(readReceipt)

	if result.Error == gorm.ErrRecordNotFound {
		// Create new receipt
		readReceipt = &ReadReceipt{
			MessageID: messageID,
			UserID:    userID,
			ReadAt:    time.Now(),
			CreatedAt: time.Now(),
		}
		if err := s.db.Create(readReceipt).Error; err != nil {
			log.Printf("Error creating read receipt: %v", err)
			return nil, err
		}
	} else if result.Error != nil {
		log.Printf("Error fetching read receipt: %v", result.Error)
		return nil, result.Error
	}

	return readReceipt, nil
}

// GetReadReceipts gets all read receipts for a message
func (s *ReadReceiptService) GetReadReceipts(messageID uint) ([]ReadReceipt, error) {
	log.Printf("Getting read receipts for message %d", messageID)

	var receipts []ReadReceipt
	result := s.db.Where("message_id = ?", messageID).
		Preload("User").
		Find(&receipts)

	if result.Error != nil {
		log.Printf("Error fetching read receipts: %v", result.Error)
		return nil, result.Error
	}

	return receipts, nil
}

// GetUnreadMessages gets unread messages for a user
func (s *ReadReceiptService) GetUnreadMessages(userID uint) (int64, error) {
	log.Printf("Getting unread message count for user %d", userID)

	var count int64
	result := s.db.Model(&Message{}).
		Where("user_id != ?", userID).
		Where("id NOT IN (?)",
			s.db.Select("message_id").
				From(&ReadReceipt{}).
				Where("user_id = ?", userID)).
		Count(&count)

	if result.Error != nil {
		log.Printf("Error counting unread messages: %v", result.Error)
		return 0, result.Error
	}

	return count, nil
}

// MarkChatAsRead marks all messages in a chat as read
func (s *ReadReceiptService) MarkChatAsRead(chatID, userID uint) error {
	log.Printf("Marking chat %d as read for user %d", chatID, userID)

	// Get all messages in chat
	var messages []Message
	if err := s.db.Where("chat_id = ?", chatID).Find(&messages).Error; err != nil {
		log.Printf("Error fetching messages: %v", err)
		return err
	}

	// Mark each as read
	now := time.Now()
	for _, msg := range messages {
		receipt := &ReadReceipt{
			MessageID: msg.ID,
			UserID:    userID,
			ReadAt:    now,
			CreatedAt: now,
		}

		// Use FirstOrCreate to avoid duplicates
		s.db.Where("message_id = ? AND user_id = ?", msg.ID, userID).
			FirstOrCreate(receipt)
	}

	log.Printf("Marked %d messages as read", len(messages))
	return nil
}
