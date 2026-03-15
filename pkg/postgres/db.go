package postgres

import (
	"fmt"
	"log"

	"github.com/renaldis/tutorku-backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		config.Cfg.DBHost,
		config.Cfg.DBPort,
		config.Cfg.DBUser,
		config.Cfg.DBPass,
		config.Cfg.DBName,
	)

	logLevel := logger.Info
	if config.Cfg.AppEnv == "production" {
		logLevel = logger.Error
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	log.Println("✅ Database connected!")
	return db, nil
}
