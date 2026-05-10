package repository

import (
	"github.com/renaldis/tutorku-backend/internal/domain"
	"gorm.io/gorm"
)

type QuizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

func (r *QuizRepository) CreateQuiz(quiz *domain.Quiz) error {
	return r.db.Create(quiz).Error
}

func (r *QuizRepository) GetQuizByID(quizID, userID string) (*domain.Quiz, error) {
	var quiz domain.Quiz
	err := r.db.Preload("Questions").
		Preload("Questions.Options").
		Where("id = ? AND user_id = ?", quizID, userID).
		First(&quiz).Error
	if err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (r *QuizRepository) GetQuizzesByMaterialID(materialID, userID string) ([]domain.Quiz, error) {
	var quizzes []domain.Quiz
	err := r.db.Preload("Questions").Where("material_id = ? AND user_id = ?", materialID, userID).
		Order("created_at desc").
		Find(&quizzes).Error
	return quizzes, err
}

func (r *QuizRepository) CreateAttempt(attempt *domain.QuizAttempt) error {
	return r.db.Create(attempt).Error
}

func (r *QuizRepository) GetAttemptByID(attemptID, userID string) (*domain.QuizAttempt, error) {
	var attempt domain.QuizAttempt
	err := r.db.Preload("Answers").
		Preload("Answers.Question").
		Preload("Answers.Question.Options").
		Where("id = ? AND user_id = ?", attemptID, userID).
		First(&attempt).Error
	if err != nil {
		return nil, err
	}
	return &attempt, nil
}

func (r *QuizRepository) UpdateAttempt(attempt *domain.QuizAttempt) error {
	return r.db.Save(attempt).Error
}

func (r *QuizRepository) GetAttemptsByQuizID(quizID, userID string) ([]domain.QuizAttempt, error) {
	var attempts []domain.QuizAttempt
	err := r.db.Where("quiz_id = ? AND user_id = ?", quizID, userID).
		Order("created_at desc").
		Find(&attempts).Error
	return attempts, err
}

func (r *QuizRepository) DeleteQuiz(quizID, userID string) error {
	result := r.db.Where("id = ? AND user_id = ?", quizID, userID).Delete(&domain.Quiz{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

