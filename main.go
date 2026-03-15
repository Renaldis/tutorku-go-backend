package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/renaldis/tutorku-backend/config"
	"github.com/renaldis/tutorku-backend/internal/domain"
	"github.com/renaldis/tutorku-backend/internal/handler"
	"github.com/renaldis/tutorku-backend/internal/repository"
	"github.com/renaldis/tutorku-backend/internal/service"
	"github.com/renaldis/tutorku-backend/pkg/n8n"
	"github.com/renaldis/tutorku-backend/pkg/postgres"
	"github.com/renaldis/tutorku-backend/routes"
	"gorm.io/gorm"
)

func main() {
	config.Load()

	db, err := postgres.Connect()
	if err != nil {
		log.Fatal("❌ Gagal koneksi database:", err)
	}
	runMigrations(db)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	materialRepo := repository.NewMaterialRepository(db)
	chatRepo := repository.NewChatRepository(db)

	// n8n
	n8nClient := n8n.NewClient()

	// Services
	authSvc := service.NewAuthService(userRepo)
	materialSvc := service.NewMaterialService(materialRepo, n8nClient)
	chatSvc := service.NewChatService(chatRepo, materialRepo, n8nClient)
	featureSvc := service.NewFeatureService(materialRepo, n8nClient)

	// Handlers
	authH := handler.NewAuthHandler(authSvc)
	materialH := handler.NewMaterialHandler(materialSvc)
	chatH := handler.NewChatHandler(chatSvc)
	featureH := handler.NewFeatureHandler(featureSvc)

	if config.Cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(gin.Recovery())
	routes.Setup(r, authH, materialH, chatH, featureH)

	addr := fmt.Sprintf(":%s", config.Cfg.AppPort)
	log.Printf("🚀 TutorKu Backend running on %s", addr)
	r.Run(addr)
}

func runMigrations(db *gorm.DB) {
	db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto")
	db.AutoMigrate(
		&domain.User{},
		&domain.Material{},
		&domain.ChatSession{},
		&domain.ChatMessage{},
	)
	log.Println("✅ Migration completed!")
}
