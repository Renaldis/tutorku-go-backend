package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/renaldis/tutorku-backend/internal/handler"
	"github.com/renaldis/tutorku-backend/middleware"
)

func Setup(
	r *gin.Engine,
	authH *handler.AuthHandler,
	materialH *handler.MaterialHandler,
	chatH *handler.ChatHandler,
	featureH *handler.FeatureHandler,
	userH *handler.UserHandler,
) {
	// Tambahkan CORS di sini
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	api := r.Group("/api/v1")

	// Public
	auth := api.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
	}

	// n8n callback
	api.POST("/callback/ingestion", materialH.UpdateStatus)

	// Protected
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		materials := protected.Group("/materials")
		{
			materials.POST("/upload", materialH.Upload)
			materials.GET("", materialH.GetAll)
			materials.GET("/:id", materialH.GetById)
			materials.GET("/:id/download", materialH.Download)
			materials.GET("/:id/status", materialH.GetStatus)
			materials.DELETE("/:id", materialH.Delete)
		}

		chat := protected.Group("/chat")
		{
			chat.POST("", chatH.Chat)
			chat.GET("/sessions", chatH.GetSessions)
			chat.GET("/history/:session_id", chatH.GetHistory)
		}

		features := protected.Group("/features")
		{
			features.POST("/summarize", featureH.Summarize)
			features.POST("/quiz", featureH.GenerateQuiz)
			features.POST("/essay", featureH.EvaluateEssay)
		}

		users := protected.Group("/users")
		{
			users.PUT("/profile", userH.UpdateProfile)
			users.PUT("/password", userH.ChangePassword)
			users.GET("get-me", userH.GetMe)
		}
	}
}
