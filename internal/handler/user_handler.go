package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/renaldis/tutorku-backend/internal/domain"
	"github.com/renaldis/tutorku-backend/internal/service"
	"github.com/renaldis/tutorku-backend/pkg/response"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Tidak memiliki akses")
		return
	}
	user, err := h.userService.GetMe(userId.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, "Berhasil mengambil data", user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Tidak memiliki akses")
		return
	}

	var req domain.EditProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updatedUser, err := h.userService.UpdateProfile(userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "Profil berhasil diperbarui", updatedUser)
}

// === HANDLER CHANGE PASSWORD ===
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Sesi tidak valid atau Anda tidak memiliki akses")
		return
	}

	var req domain.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// === SULAP ERROR JELEK JADI BAGUS ===
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			// Kita ambil error pertama yang muncul
			for _, fe := range ve {
				switch fe.Tag() {
				case "required":
					response.BadRequest(c, fe.Field()+" tidak boleh kosong!")
					return
				case "min":
					response.BadRequest(c, fe.Field()+" minimal harus "+fe.Param()+" karakter!")
					return
				case "eqfield":
					if fe.Field() == "ConfirmPassword" {
						response.BadRequest(c, "Konfirmasi password tidak cocok dengan password baru!")
						return
					}
					response.BadRequest(c, fe.Field()+" tidak cocok!")
					return
				}
			}
		}

		// Kalau errornya bukan dari validator (misal JSON-nya berantakan salah ketik kurawal)
		response.BadRequest(c, "Format data tidak valid")
		return
	}

	err := h.userService.ChangePassword(userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Karena tidak ada data yang perlu dikembalikan, parameter data diisi nil
	response.OK(c, "Password berhasil diubah", nil)
}
