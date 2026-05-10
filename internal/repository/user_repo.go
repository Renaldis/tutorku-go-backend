package repository

import (
	"github.com/renaldis/tutorku-backend/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User

	err := r.db.
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User

	err := r.db.
		Where("id = ?", id).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) GetStats(userID string) (*domain.UserStats, error) {
	stats := &domain.UserStats{}

	if err := r.db.
		Model(&domain.Material{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalMaterials).Error; err != nil {
		return nil, err
	}

	if err := r.db.
		Model(&domain.Quiz{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalQuizzes).Error; err != nil {
		return nil, err
	}

	if err := r.db.
		Model(&domain.QuizAttempt{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalAttempts).Error; err != nil {
		return nil, err
	}

	type AggregateResult struct {
		AverageScore   float64 `gorm:"column:average_score"`
		BestScore      float64 `gorm:"column:best_score"`
		TotalCorrect   int64   `gorm:"column:total_correct"`
		TotalQuestions int64   `gorm:"column:total_questions"`
	}

	var aggregate AggregateResult

	if err := r.db.Raw(`
		SELECT
			COALESCE(AVG(score), 0)           AS average_score,
			COALESCE(MAX(score), 0)           AS best_score,
			COALESCE(SUM(total_correct), 0)   AS total_correct,
			COALESCE(SUM(total_questions), 0) AS total_questions
		FROM quiz_attempts
		WHERE user_id = ?
		AND finished_at IS NOT NULL
	`, userID).Scan(&aggregate).Error; err != nil {
		return nil, err
	}

	stats.AverageScore = aggregate.AverageScore
	stats.BestScore = aggregate.BestScore
	stats.TotalCorrect = aggregate.TotalCorrect
	stats.TotalQuestions = aggregate.TotalQuestions

	return stats, nil
}
