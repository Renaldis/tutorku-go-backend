package domain

import (
	"time"

	"gorm.io/gorm"
)

type MaterialStatus string

const (
	StatusProcessing MaterialStatus = "processing"
	StatusReady      MaterialStatus = "ready"
	StatusFailed     MaterialStatus = "failed"
)

type Material struct {
	ID        string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    string         `json:"user_id" gorm:"type:uuid;not null;index"`
	User      User           `json:"-" gorm:"foreignKey:UserID"`
	Title     string         `json:"title" gorm:"not null"`
	Category  string         `json:"category"`
	Filename  string         `json:"filename" gorm:"not null"`
	FileSize  int64          `json:"file_size"`
	Status    MaterialStatus `json:"status" gorm:"default:'processing'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UploadMaterialRequest struct {
	Title    string `form:"title" binding:"required"`
	Category string `form:"category"`
}

type UpdateStatusRequest struct {
	MaterialID string         `json:"material_id" binding:"required"`
	Status     MaterialStatus `json:"status" binding:"required"`
}
