package domain

import (
	"time"
)

type QuizAnswer struct {
	ID           string       `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	AttemptID    string       `json:"attempt_id" gorm:"type:uuid;not null;index"`
	QuestionID   string       `json:"question_id" gorm:"type:uuid;not null;index"`
	Question     QuizQuestion `json:"-" gorm:"foreignKey:QuestionID"`
	UserAnswer   string       `json:"user_answer"`
	IsCorrect    bool         `json:"is_correct"`
	EarnedPoints float64      `json:"earned_points"`
	CreatedAt    time.Time    `json:"created_at"`
}

type SubmitQuizRequest struct {
	Answers []SubmitAnswer `json:"answers" binding:"required"`
}

type SubmitAnswer struct {
	QuestionID string `json:"question_id" binding:"required"`
	UserAnswer string `json:"user_answer" binding:"required"`
}
