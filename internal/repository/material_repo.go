package repository

import (
	"errors"

	"github.com/renaldis/tutorku-backend/internal/domain"
	"gorm.io/gorm"
)

type MaterialRepository struct {
	db *gorm.DB
}

func NewMaterialRepository(db *gorm.DB) *MaterialRepository {
	return &MaterialRepository{db: db}
}

func (r *MaterialRepository) Create(material *domain.Material) error {
	return r.db.Create(material).Error
}

func (r *MaterialRepository) FindByUser(userID string) ([]domain.Material, error) {
	var materials []domain.Material
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&materials).Error
	return materials, err
}

func (r *MaterialRepository) FindByID(id, userID string) (*domain.Material, error) {
	var material domain.Material
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&material).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("materi tidak ditemukan")
	}
	return &material, err
}

func (r *MaterialRepository) UpdateStatus(id string, status domain.MaterialStatus) error {
	return r.db.Model(&domain.Material{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *MaterialRepository) Delete(id, userID string) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).
		Delete(&domain.Material{})
	if result.RowsAffected == 0 {
		return errors.New("materi tidak ditemukan")
	}
	return result.Error
}
