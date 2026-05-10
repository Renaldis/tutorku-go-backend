package domain

import "time"

type ChatSession struct {
	ID         string        `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     string        `json:"user_id" gorm:"type:uuid;not null;index"`
	MaterialID string        `json:"material_id" gorm:"type:uuid;not null;index"`
	Material   Material      `json:"material,omitempty" gorm:"foreignKey:MaterialID"`
	Messages   []ChatMessage `json:"messages,omitempty" gorm:"foreignKey:SessionID"`
	CreatedAt  time.Time     `json:"created_at"`
}

type ChatMessage struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SessionID string    `json:"session_id" gorm:"type:uuid;not null;index"`
	Role      string    `json:"role" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatRequest struct {
	MaterialID string `json:"material_id" binding:"required"`
	SessionID  string `json:"session_id"`
	Query      string `json:"query" binding:"required"`
}
