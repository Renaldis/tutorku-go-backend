package domain

import (
	"time"
)

type QuizAttempt struct {
	ID             string       `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID         string       `json:"user_id" gorm:"type:uuid;not null;index"`
	User           User         `json:"-" gorm:"foreignKey:UserID"`
	QuizID         string       `json:"quiz_id" gorm:"type:uuid;not null;index"`
	Quiz           Quiz         `json:"-" gorm:"foreignKey:QuizID"`
	Score          float64      `json:"score"`
	TotalCorrect   int          `json:"total_correct"`
	TotalQuestions int          `json:"total_questions"`
	StartedAt      time.Time    `json:"started_at"`
	FinishedAt     *time.Time   `json:"finished_at"`
	CreatedAt      time.Time    `json:"created_at"`
	Answers        []QuizAnswer `json:"answers" gorm:"foreignKey:AttemptID;constraint:OnDelete:CASCADE;"`
}

type StartAttemptRequest struct {
	QuizID string `json:"quiz_id" binding:"required"`
}
