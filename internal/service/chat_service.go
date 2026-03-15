package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/renaldis/tutorku-backend/internal/domain"
	"github.com/renaldis/tutorku-backend/internal/repository"
	"github.com/renaldis/tutorku-backend/pkg/n8n"
)

type ChatService struct {
	chatRepo     *repository.ChatRepository
	materialRepo *repository.MaterialRepository
	n8nClient    *n8n.Client
}

func NewChatService(chatRepo *repository.ChatRepository, materialRepo *repository.MaterialRepository, n8nClient *n8n.Client) *ChatService {
	return &ChatService{chatRepo: chatRepo, materialRepo: materialRepo, n8nClient: n8nClient}
}

func (s *ChatService) Chat(userID string, req domain.ChatRequest) (map[string]interface{}, string, error) {
	// Validasi materi milik user
	material, err := s.materialRepo.FindByID(req.MaterialID, userID)
	if err != nil {
		return nil, "", err
	}
	if material.Status != domain.StatusReady {
		return nil, "", errors.New("materi belum selesai diproses")
	}

	// Buat atau ambil session
	sessionID := req.SessionID
	if sessionID == "" {
		session := &domain.ChatSession{
			ID:         uuid.New().String(),
			UserID:     userID,
			MaterialID: req.MaterialID,
		}
		s.chatRepo.CreateSession(session)
		sessionID = session.ID
	}

	// Ambil history chat terakhir
	recentMessages, _ := s.chatRepo.GetRecentMessages(sessionID, 10)
	var history []n8n.ChatHistory
	for _, msg := range recentMessages {
		history = append(history, n8n.ChatHistory{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Simpan pesan user
	s.chatRepo.SaveMessage(&domain.ChatMessage{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		Role:      "user",
		Content:   req.Query,
	})

	// Call n8n RAG
	result, err := s.n8nClient.QueryRAG(n8n.ChatPayload{
		MaterialID:  req.MaterialID,
		UserID:      userID,
		Query:       req.Query,
		ChatHistory: history,
	})
	if err != nil {
		return nil, sessionID, errors.New("gagal mendapatkan jawaban")
	}

	// Simpan jawaban AI
	aiAnswer := ""
	if answer, ok := result["answer"].(string); ok {
		aiAnswer = answer
		s.chatRepo.SaveMessage(&domain.ChatMessage{
			ID:        uuid.New().String(),
			SessionID: sessionID,
			Role:      "assistant",
			Content:   aiAnswer,
		})
	}

	return result, sessionID, nil
}

func (s *ChatService) GetHistory(sessionID, userID string) ([]domain.ChatMessage, error) {
	// Validasi session milik user
	_, err := s.chatRepo.FindSession(sessionID, userID)
	if err != nil {
		return nil, errors.New("sesi tidak ditemukan")
	}
	return s.chatRepo.GetMessages(sessionID)
}

func (s *ChatService) GetSessions(userID string) ([]domain.ChatSession, error) {
	return s.chatRepo.GetSessionsByUser(userID)
}
