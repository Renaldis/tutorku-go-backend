package handler

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/renaldis/tutorku-backend/internal/domain"
	"github.com/renaldis/tutorku-backend/internal/service"
	"github.com/renaldis/tutorku-backend/pkg/response"
)

type MaterialHandler struct {
	materialService *service.MaterialService
}

func NewMaterialHandler(s *service.MaterialService) *MaterialHandler {
	return &MaterialHandler{materialService: s}
}

func (h *MaterialHandler) Upload(c *gin.Context) {
	userID := c.GetString("user_id")

	var req domain.UploadMaterialRequest
	if err := c.ShouldBind(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "File PDF tidak ditemukan")
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		response.InternalError(c, "Gagal membaca file")
		return
	}

	material, err := h.materialService.Upload(
		userID, req.Title, req.Category,
		header.Filename, fileBytes, header.Size,
	)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, "Materi sedang diproses", material)
}

func (h *MaterialHandler) Download(c *gin.Context) {
	userID := c.GetString("user_id")
	materialID := c.Param("id")

	material, err := h.materialService.GetByID(materialID, userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	filePath := "uploads/materials/" + material.ID + ".pdf"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		response.NotFound(c, "File tidak ditemukan di server")
		return
	}

	c.FileAttachment(filePath, material.Filename)
}

func (h *MaterialHandler) GetAll(c *gin.Context) {
	userID := c.GetString("user_id")
	materials, err := h.materialService.GetByUser(userID)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data materi")
		return
	}
	response.OK(c, "success", materials)
}
func (h *MaterialHandler) GetById(c *gin.Context) {
	userID := c.GetString("user_id")
	materialID := c.Param("id")

	material, err := h.materialService.GetByID(materialID, userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, "success", material)
}

func (h *MaterialHandler) GetStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	materialID := c.Param("id")

	material, err := h.materialService.GetByID(materialID, userID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, "success", gin.H{"status": material.Status})
}

func (h *MaterialHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	materialID := c.Param("id")

	if err := h.materialService.Delete(materialID, userID); err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, "Materi berhasil dihapus", nil)
}

// Callback dari n8n setelah ingestion selesai
func (h *MaterialHandler) UpdateStatus(c *gin.Context) {
	var req domain.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.materialService.UpdateStatus(req.MaterialID, req.Status); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.OK(c, "Status diperbarui", nil)
}
