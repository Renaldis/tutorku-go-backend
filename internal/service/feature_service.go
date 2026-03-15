package service

import (
	"errors"

	"github.com/renaldis/tutorku-backend/internal/domain"
	"github.com/renaldis/tutorku-backend/internal/repository"
	"github.com/renaldis/tutorku-backend/pkg/n8n"
)

type FeatureService struct {
	materialRepo *repository.MaterialRepository
	n8nClient    *n8n.Client
}

func NewFeatureService(materialRepo *repository.MaterialRepository, n8nClient *n8n.Client) *FeatureService {
	return &FeatureService{materialRepo: materialRepo, n8nClient: n8nClient}
}

func (s *FeatureService) validateMaterial(materialID, userID string) error {
	material, err := s.materialRepo.FindByID(materialID, userID)
	if err != nil {
		return err
	}
	if material.Status != domain.StatusReady {
		return errors.New("materi belum selesai diproses")
	}
	return nil
}

func (s *FeatureService) Summarize(userID string, req domain.SummarizeRequest) (map[string]interface{}, error) {
	if err := s.validateMaterial(req.MaterialID, userID); err != nil {
		return nil, err
	}
	return s.n8nClient.Summarize(n8n.SummarizePayload{
		MaterialID: req.MaterialID,
		UserID:     userID,
		Mode:       req.Mode,
	})
}

func (s *FeatureService) GenerateQuiz(userID string, req domain.QuizRequest) (map[string]interface{}, error) {
	if err := s.validateMaterial(req.MaterialID, userID); err != nil {
		return nil, err
	}
	return s.n8nClient.GenerateQuiz(n8n.QuizPayload{
		MaterialID: req.MaterialID,
		UserID:     userID,
		Type:       req.Type,
		Count:      req.Count,
		Difficulty: req.Difficulty,
	})
}

func (s *FeatureService) EvaluateEssay(userID string, req domain.EssayRequest) (map[string]interface{}, error) {
	if err := s.validateMaterial(req.MaterialID, userID); err != nil {
		return nil, err
	}
	return s.n8nClient.EvaluateEssay(n8n.EssayPayload{
		MaterialID: req.MaterialID,
		UserID:     userID,
		Title:      req.Title,
		Content:    req.Content,
	})
}
