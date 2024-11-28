package router

import (
	"authenticationProject/internal/handlers"
	"authenticationProject/internal/middleware"
	"authenticationProject/internal/services"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"

	_ "authenticationProject/docs"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(authHandler *handlers.AuthHandler, logger *logrus.Logger, tokenService *services.TokenService) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	r.Mount("/swagger", httpSwagger.WrapHandler)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/token", authHandler.GenerateTokens)
		r.Post("/refresh", authHandler.RefreshTokens)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(tokenService, logger))
		r.Get("/protected", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value("userID").(string)
			w.Write([]byte("Hello, user " + userID))
		})
	})

	return r
}
