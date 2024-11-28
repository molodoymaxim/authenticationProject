package main

import (
	"authenticationProject/configs"
	"authenticationProject/docs"
	"authenticationProject/internal/handlers"
	"authenticationProject/internal/repository"
	"authenticationProject/internal/repository/migration"
	"authenticationProject/internal/router"
	"authenticationProject/internal/services"
	"authenticationProject/internal/utils"
	"log"
	"net/http"
)

// @title           Authentication Service API
// @version         1.0
// @description     API для аутентификации пользователей с использованием JWT токенов.

// @license.name  MIT License
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	config := configs.LoadConfig()

	docs.SwaggerInfo.Host = "localhost:" + config.AppPort

	logger := utils.NewLogger(config.LogLevel)

	db, err := repository.NewDatabase(config.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	err = migration.ApplyMigrations(db, logger)
	if err != nil {
		logger.Fatal("Failed to apply migrations: ", err)
	}

	tokenService := services.NewTokenService(config.SecretKey)
	emailService := utils.NewEmailService(config.EmailAPIKey)
	tokenRepository := repository.NewTokenRepository(db)
	authHandler := handlers.NewAuthHandler(tokenService, tokenRepository, emailService, logger)

	r := router.NewRouter(authHandler, logger, tokenService)

	logger.Infof("Starting server on port %s", config.AppPort)
	log.Fatal(http.ListenAndServe(":"+config.AppPort, r))
}
