package service

import (
	"encoding/base64"
	"errors"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

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
	// Decode URL-encoded string (misal: "materi%20dbms.pdf" -> "materi dbms.pdf")
	decodedFilename, err := url.QueryUnescape(filename)
	if err != nil {
		decodedFilename = filename // fallback jika gagal decode
	}

	// Sanitasi nama file: ubah spasi jadi dash (-) dan jadikan huruf kecil semua
	safeFilename := strings.ToLower(strings.ReplaceAll(decodedFilename, " ", "-"))

	material := &domain.Material{
		ID:       uuid.New().String(),
		UserID:   userID,
		Title:    title,
		Category: category,
		Filename: safeFilename,
		FileSize: fileSize,
		Status:   domain.StatusProcessing,
	}

	if err := s.materialRepo.Create(material); err != nil {
		return nil, errors.New("gagal menyimpan materi")
	}

	// Save file locally for download
	uploadDir := "uploads/materials"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("❌ Gagal membuat direktori upload: %v", err)
	} else {
		filePath := filepath.Join(uploadDir, material.ID+".pdf")
		if err := os.WriteFile(filePath, fileBytes, 0644); err != nil {
			log.Printf("❌ Gagal menyimpan file secara lokal: %v", err)
		}
	}

	// Kirim ke n8n secara async
	go func() {
		fileBase64 := base64.StdEncoding.EncodeToString(fileBytes)
		log.Printf("🚀 Sending to n8n: material_id=%s", material.ID)
		result, err := s.n8nClient.TriggerIngestion(n8n.IngestPayload{
			MaterialID: material.ID,
			UserID:     userID,
			FileBase64: fileBase64,
			Filename:   safeFilename,
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
