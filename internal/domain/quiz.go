package domain

import (
	"time"

	"gorm.io/gorm"
)

type Quiz struct {
	ID          string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      string         `json:"user_id" gorm:"type:uuid;not null;index"`
	User        User           `json:"-" gorm:"foreignKey:UserID"`
	MaterialID  string         `json:"material_id" gorm:"type:uuid;not null;index"`
	Material    Material       `json:"-" gorm:"foreignKey:MaterialID"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	GeneratedBy string         `json:"generated_by" gorm:"default:'ai'"`
	Questions   []QuizQuestion `json:"questions" gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE;"`
	Attempts    []QuizAttempt  `json:"attempts" gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE;"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type QuizQuestion struct {
	ID            string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	QuizID        string         `json:"quiz_id" gorm:"type:uuid;not null;index"`
	Question      string         `json:"question" gorm:"not null"`
	Type          string         `json:"type" gorm:"not null"` // multiple_choice, true_false, essay
	CorrectAnswer string         `json:"correct_answer"`
	Explanation   string         `json:"explanation"`
	OrderNo       int            `json:"order_no"`
	Options       []QuizOption   `json:"options" gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE;"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type QuizOption struct {
	ID         string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	QuestionID string    `json:"question_id" gorm:"type:uuid;not null;index"`
	OptionKey  string    `json:"key" gorm:"not null"`
	OptionText string    `json:"text" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
}

type GenerateQuizRequest struct {
	Type       string `json:"type" binding:"required"` // multiple_choice, true_false, essay
	Difficulty string `json:"difficulty" binding:"required"` // easy, medium, hard
	Count      int    `json:"count" binding:"required,min=1,max=20"`
}

// Struct for handling AI response
type AIQuizResponse struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Questions   []AIQuizQuestion `json:"questions"`
}

type AIQuizQuestion struct {
	Question      string            `json:"question"`
	Type          string            `json:"type"`
	Options       map[string]string `json:"options"` // for multiple_choice and true_false
	KeyPoints     []string          `json:"key_points"` // for essay
	SampleAnswer  string            `json:"sample_answer"` // for essay
	CorrectAnswer string            `json:"correct_answer"`
	Explanation   string            `json:"explanation"`
	Difficulty    string            `json:"difficulty"`
}
