package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/renaldis/tutorku-backend/internal/domain"
	"github.com/renaldis/tutorku-backend/internal/service"
	"github.com/renaldis/tutorku-backend/pkg/response"
)

type QuizHandler struct {
	quizService *service.QuizService
}

func NewQuizHandler(s *service.QuizService) *QuizHandler {
	return &QuizHandler{quizService: s}
}

func (h *QuizHandler) GenerateQuiz(c *gin.Context) {
	userID := c.GetString("user_id")
	materialID := c.Param("id")

	var req domain.GenerateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	quiz, err := h.quizService.GenerateQuiz(userID, materialID, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, "Quiz berhasil digenerate", quiz)
}

func (h *QuizHandler) GetQuiz(c *gin.Context) {
	userID := c.GetString("user_id")
	quizID := c.Param("id")

	quiz, err := h.quizService.GetQuiz(quizID, userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.OK(c, "success", quiz)
}

func (h *QuizHandler) GetQuizzesByMaterial(c *gin.Context) {
	userID := c.GetString("user_id")
	materialID := c.Param("id")

	quizzes, err := h.quizService.GetQuizzesByMaterial(materialID, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "success", quizzes)
}

func (h *QuizHandler) StartAttempt(c *gin.Context) {
	userID := c.GetString("user_id")
	quizID := c.Param("id")

	attempt, err := h.quizService.StartAttempt(quizID, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, "Quiz dimulai", attempt)
}

func (h *QuizHandler) SubmitAttempt(c *gin.Context) {
	userID := c.GetString("user_id")
	attemptID := c.Param("attempt_id")

	var req domain.SubmitQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	attempt, err := h.quizService.SubmitAttempt(attemptID, userID, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "Quiz disubmit", attempt)
}

func (h *QuizHandler) GetAttemptsByQuiz(c *gin.Context) {
	userID := c.GetString("user_id")
	quizID := c.Param("id")

	attempts, err := h.quizService.GetAttemptsByQuiz(quizID, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "success", attempts)
}
