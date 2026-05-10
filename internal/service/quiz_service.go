package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/renaldis/tutorku-backend/internal/domain"
	"github.com/renaldis/tutorku-backend/internal/repository"
	"github.com/renaldis/tutorku-backend/pkg/n8n"
)

type QuizService struct {
	quizRepo     *repository.QuizRepository
	materialRepo *repository.MaterialRepository
	n8nClient    *n8n.Client
}

func NewQuizService(quizRepo *repository.QuizRepository, materialRepo *repository.MaterialRepository, n8nClient *n8n.Client) *QuizService {
	return &QuizService{
		quizRepo:     quizRepo,
		materialRepo: materialRepo,
		n8nClient:    n8nClient,
	}
}

func (s *QuizService) GenerateQuiz(userID, materialID string, req domain.GenerateQuizRequest) (*domain.Quiz, error) {
	// Validate material
	material, err := s.materialRepo.FindByID(materialID, userID)
	if err != nil {
		return nil, err
	}
	if material.Status != domain.StatusReady {
		return nil, errors.New("materi belum selesai diproses")
	}

	// Call AI to generate quiz
	aiResult, err := s.n8nClient.GenerateQuiz(n8n.QuizPayload{
		MaterialID: materialID,
		UserID:     userID,
		Type:       req.Type,
		Count:      req.Count,
		Difficulty: req.Difficulty,
	})
	if err != nil {
		return nil, err
	}

	// Parse AI Result
	// The AI result usually returns a "questions" array, but let's assume it matches AIQuizResponse
	// or we extract the array directly depending on the AI format.
	// Looking at the previous Python code, it returns:
	// {"material_id": "...", "questions": [{...}]}

	questionsData, ok := aiResult["questions"].([]interface{})
	if !ok {
		return nil, errors.New("format respons AI tidak valid: gagal menemukan pertanyaan")
	}

	questionsBytes, err := json.Marshal(questionsData)
	if err != nil {
		return nil, errors.New("gagal memproses pertanyaan dari AI")
	}

	var aiQuestions []domain.AIQuizQuestion
	if err := json.Unmarshal(questionsBytes, &aiQuestions); err != nil {
		return nil, errors.New("gagal mem-parsing pertanyaan dari AI")
	}

	// Build Quiz
	quiz := &domain.Quiz{
		ID:          uuid.New().String(),
		UserID:      userID,
		MaterialID:  materialID,
		Title:       "Quiz untuk " + material.Title, // We can customize this
		Description: "Dihasilkan secara otomatis oleh AI",
		GeneratedBy: "ai",
	}

	// Build Questions and Options
	for i, aiQ := range aiQuestions {
		q := domain.QuizQuestion{
			ID:      uuid.New().String(),
			QuizID:  quiz.ID,
			Question: aiQ.Question,
			Type:    req.Type,
			OrderNo: i + 1,
		}

		switch req.Type {
		case "multiple_choice", "true_false":
			q.CorrectAnswer = aiQ.CorrectAnswer
			q.Explanation = aiQ.Explanation

			// Sort keys agar urutan option selalu A, B, C, D (map Go tidak berurutan)
			keys := make([]string, 0, len(aiQ.Options))
			for key := range aiQ.Options {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			for _, key := range keys {
				q.Options = append(q.Options, domain.QuizOption{
					ID:         uuid.New().String(),
					QuestionID: q.ID,
					OptionKey:  key,
					OptionText: aiQ.Options[key],
				})
			}

		case "essay":
			// Untuk essay: simpan sample_answer sebagai CorrectAnswer
			// dan key_points (digabung) sebagai Explanation sebagai referensi penilaian
			q.CorrectAnswer = aiQ.SampleAnswer
			if len(aiQ.KeyPoints) > 0 {
				keyPointsBytes, _ := json.Marshal(aiQ.KeyPoints)
				q.Explanation = string(keyPointsBytes)
			} else {
				q.Explanation = aiQ.Explanation
			}
			// Essay tidak punya options
		}

		quiz.Questions = append(quiz.Questions, q)
	}

	// Save to DB
	if err := s.quizRepo.CreateQuiz(quiz); err != nil {
		return nil, err
	}

	return quiz, nil
}

func (s *QuizService) GetQuiz(quizID, userID string) (*domain.Quiz, error) {
	return s.quizRepo.GetQuizByID(quizID, userID)
}

func (s *QuizService) GetQuizzesByMaterial(materialID, userID string) ([]domain.Quiz, error) {
	return s.quizRepo.GetQuizzesByMaterialID(materialID, userID)
}

func (s *QuizService) StartAttempt(quizID, userID string) (*domain.QuizAttempt, error) {
	quiz, err := s.quizRepo.GetQuizByID(quizID, userID)
	if err != nil {
		return nil, err
	}

	attempt := &domain.QuizAttempt{
		ID:             uuid.New().String(),
		UserID:         userID,
		QuizID:         quiz.ID,
		TotalQuestions: len(quiz.Questions),
	}

	if err := s.quizRepo.CreateAttempt(attempt); err != nil {
		return nil, err
	}

	return attempt, nil
}

func (s *QuizService) SubmitAttempt(attemptID, userID string, req domain.SubmitQuizRequest) (*domain.QuizAttempt, error) {
	attempt, err := s.quizRepo.GetAttemptByID(attemptID, userID)
	if err != nil {
		return nil, err
	}

	if attempt.FinishedAt != nil {
		return nil, errors.New("quiz sudah disubmit")
	}

	quiz, err := s.quizRepo.GetQuizByID(attempt.QuizID, userID)
	if err != nil {
		return nil, err
	}

	// Create a map for quick question lookup
	questionMap := make(map[string]domain.QuizQuestion)
	for _, q := range quiz.Questions {
		questionMap[q.ID] = q
	}

	totalCorrect := 0
	totalEarnedPoints := 0.0
	var answers []domain.QuizAnswer

	for _, ansReq := range req.Answers {
		q, exists := questionMap[ansReq.QuestionID]
		if !exists {
			continue
		}

		isCorrect := false
		earnedPoints := 0.0

		switch q.Type {
		case "multiple_choice", "true_false":
			if ansReq.UserAnswer == q.CorrectAnswer {
				isCorrect = true
				earnedPoints = 1.0
				totalCorrect++
			}

		case "essay":
			// Kirim jawaban essay ke AI (chatbot /essay) untuk dievaluasi
			essayResult, err := s.n8nClient.EvaluateEssay(n8n.EssayPayload{
				MaterialID: quiz.MaterialID,
				UserID:     userID,
				Title:      q.Question,
				Content:    ansReq.UserAnswer,
			})
			if err == nil {
				// Chatbot mengembalikan: {"material_id": ..., "user_id": ..., "evaluation": {...}}
				if evaluation, ok := essayResult["evaluation"].(map[string]interface{}); ok {
					if scoreRaw, ok := evaluation["score"]; ok {
						var score float64
						switch v := scoreRaw.(type) {
						case float64:
							score = v
						case int:
							score = float64(v)
						}
						// Normalisasi skor 0-100 ke 0-1 poin
						earnedPoints = score / 100.0
						if score >= 60 {
							isCorrect = true
							totalCorrect++
						}
						// Simpan feedback AI ke UserAnswer sebagai JSON tambahan
						ansReq.UserAnswer = fmt.Sprintf("%s", ansReq.UserAnswer)
					}
				}
			}
		}

		totalEarnedPoints += earnedPoints
		answers = append(answers, domain.QuizAnswer{
			ID:           uuid.New().String(),
			AttemptID:    attempt.ID,
			QuestionID:   q.ID,
			UserAnswer:   ansReq.UserAnswer,
			IsCorrect:    isCorrect,
			EarnedPoints: earnedPoints,
		})
	}
	_ = totalEarnedPoints // dapat digunakan untuk perhitungan skor lebih lanjut

	// Update Attempt
	now := time.Now()
	attempt.FinishedAt = &now
	attempt.TotalCorrect = totalCorrect

	if attempt.TotalQuestions > 0 {
		attempt.Score = float64(totalCorrect) / float64(attempt.TotalQuestions) * 100.0
	} else {
		attempt.Score = 0
	}

	attempt.Answers = answers

	if err := s.quizRepo.UpdateAttempt(attempt); err != nil {
		return nil, err
	}

	return attempt, nil
}

func (s *QuizService) GetAttemptsByQuiz(quizID, userID string) ([]domain.QuizAttempt, error) {
	return s.quizRepo.GetAttemptsByQuizID(quizID, userID)
}

func (s *QuizService) DeleteQuiz(quizID, userID string) error {
	return s.quizRepo.DeleteQuiz(quizID, userID)
}
