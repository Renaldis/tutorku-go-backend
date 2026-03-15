package repository

import (
	"github.com/renaldis/tutorku-backend/internal/domain"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CreateSession(session *domain.ChatSession) error {
	return r.db.Create(session).Error
}

func (r *ChatRepository) FindSession(sessionID, userID string) (*domain.ChatSession, error) {
	var session domain.ChatSession
	err := r.db.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error
	return &session, err
}

func (r *ChatRepository) GetMessages(sessionID string) ([]domain.ChatMessage, error) {
	var messages []domain.ChatMessage
	err := r.db.Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

func (r *ChatRepository) SaveMessage(message *domain.ChatMessage) error {
	return r.db.Create(message).Error
}

func (r *ChatRepository) GetRecentMessages(sessionID string, limit int) ([]domain.ChatMessage, error) {
	var messages []domain.ChatMessage
	err := r.db.Where("session_id = ?", sessionID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages, err
}

func (r *ChatRepository) GetSessionsByUser(userID string) ([]domain.ChatSession, error) {
	var sessions []domain.ChatSession
	err := r.db.Where("user_id = ?", userID).
		Preload("Material").
		Order("created_at DESC").
		Find(&sessions).Error
	return sessions, err
}
