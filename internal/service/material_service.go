package service

import (
	"encoding/base64"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/renaldis/tutorku-backend/internal/domain"
	"github.com/renaldis/tutorku-backend/internal/repository"
	"github.com/renaldis/tutorku-backend/pkg/n8n"
)

type MaterialService struct {
	materialRepo *repository.MaterialRepository
	n8nClient    *n8n.Client
}

func NewMaterialService(materialRepo *repository.MaterialRepository, n8nClient *n8n.Client) *MaterialService {
	return &MaterialService{materialRepo: materialRepo, n8nClient: n8nClient}
}

func (s *MaterialService) Upload(userID, title, category, filename string, fileBytes []byte, fileSize int64) (*domain.Material, error) {
	material := &domain.Material{
		ID:       uuid.New().String(),
		UserID:   userID,
		Title:    title,
		Category: category,
		Filename: filename,
		FileSize: fileSize,
		Status:   domain.StatusProcessing,
	}

	if err := s.materialRepo.Create(material); err != nil {
		return nil, errors.New("gagal menyimpan materi")
	}

	// Kirim ke n8n secara async
	go func() {
		fileBase64 := base64.StdEncoding.EncodeToString(fileBytes)
		log.Printf("🚀 Sending to n8n: material_id=%s", material.ID)
		result, err := s.n8nClient.TriggerIngestion(n8n.IngestPayload{
			MaterialID: material.ID,
			UserID:     userID,
			FileBase64: fileBase64,
			Filename:   filename,
		})

		if err != nil {
			log.Printf("❌ n8n error: %v", err)
			return
		}
		log.Printf("✅ n8n response: %v", result)
	}()

	return material, nil
}

func (s *MaterialService) GetByUser(userID string) ([]domain.Material, error) {
	return s.materialRepo.FindByUser(userID)
}

func (s *MaterialService) GetByID(id, userID string) (*domain.Material, error) {
	return s.materialRepo.FindByID(id, userID)
}

func (s *MaterialService) UpdateStatus(id string, status domain.MaterialStatus) error {
	return s.materialRepo.UpdateStatus(id, status)
}

func (s *MaterialService) Delete(id, userID string) error {
	return s.materialRepo.Delete(id, userID)
}
