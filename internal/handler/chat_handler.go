package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/renaldis/tutorku-backend/internal/domain"
	"github.com/renaldis/tutorku-backend/internal/service"
	"github.com/renaldis/tutorku-backend/pkg/response"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(s *service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: s}
}

func (h *ChatHandler) Chat(c *gin.Context) {
	userID := c.GetString("user_id")

	var req domain.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, sessionID, err := h.chatService.Chat(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "success", gin.H{
		"session_id": sessionID,
		"answer":     result["answer"],
	})
}

func (h *ChatHandler) GetHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	sessionID := c.Param("session_id")

	messages, err := h.chatService.GetHistory(sessionID, userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, "success", messages)
}

func (h *ChatHandler) GetSessions(c *gin.Context) {
	userID := c.GetString("user_id")
	sessions, err := h.chatService.GetSessions(userID)
	if err != nil {
		response.InternalError(c, "Gagal mengambil sesi chat")
		return
	}
	response.OK(c, "success", sessions)
}
