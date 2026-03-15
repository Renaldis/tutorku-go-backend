package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/renaldis/tutorku-backend/internal/domain"
	"github.com/renaldis/tutorku-backend/internal/service"
	"github.com/renaldis/tutorku-backend/pkg/response"
)

type FeatureHandler struct {
	featureService *service.FeatureService
}

func NewFeatureHandler(s *service.FeatureService) *FeatureHandler {
	return &FeatureHandler{featureService: s}
}

func (h *FeatureHandler) Summarize(c *gin.Context) {
	userID := c.GetString("user_id")
	var req domain.SummarizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.featureService.Summarize(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, "success", result)
}

func (h *FeatureHandler) GenerateQuiz(c *gin.Context) {
	userID := c.GetString("user_id")
	var req domain.QuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.featureService.GenerateQuiz(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, "success", result)
}

func (h *FeatureHandler) EvaluateEssay(c *gin.Context) {
	userID := c.GetString("user_id")
	var req domain.EssayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.featureService.EvaluateEssay(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, "success", result)
}
