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

type SummarizeRequest struct {
	MaterialID string `json:"material_id" binding:"required"`
	Mode       string `json:"mode" binding:"required,oneof=short detailed mindmap"`
}

type QuizRequest struct {
	MaterialID string `json:"material_id" binding:"required"`
	Type       string `json:"type" binding:"required,oneof=multiple_choice essay true_false"`
	Count      int    `json:"count" binding:"required,min=1,max=20"`
	Difficulty string `json:"difficulty" binding:"required,oneof=easy medium hard"`
}

type EssayRequest struct {
	MaterialID string `json:"material_id" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required,min=50"`
}
