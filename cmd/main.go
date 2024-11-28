package main

import (
	"authenticationProject/configs"
	"authenticationProject/internal/handlers"
	"authenticationProject/internal/repository"
	"authenticationProject/internal/repository/migration"
	"authenticationProject/internal/router"
	"authenticationProject/internal/services"
	"authenticationProject/internal/utils"
	"log"
	"net/http"
)

func main() {
	config := configs.LoadConfig()

	// Инициализация логгера
	logger := utils.NewLogger(config.LogLevel)

	// Инициализация базы данных
	db, err := repository.NewDatabase(config.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	err = migration.ApplyMigrations(db, logger)
	if err != nil {
		logger.Fatal("Failed to apply migrations: ", err)
	}

	// Инициализация сервисов и хендлеров
	tokenService := services.NewTokenService(config.SecretKey)
	emailService := utils.NewEmailService(config.EmailAPIKey)
	tokenRepository := repository.NewTokenRepository(db)
	authHandler := handlers.NewAuthHandler(tokenService, tokenRepository, emailService, logger)

	// Создание роутера
	r := router.NewRouter(authHandler, logger, tokenService)

	logger.Infof("Starting server on port %s", config.AppPort)
	log.Fatal(http.ListenAndServe(":"+config.AppPort, r))
}
