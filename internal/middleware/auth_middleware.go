package middleware

import (
	"authenticationProject/internal/services"
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func AuthMiddleware(tokenService *services.TokenService, logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("Missing Authorization header")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := tokenService.ValidateAccessToken(tokenStr)
			if err != nil {
				logger.Warn("Invalid access token: ", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", claims.UserID)
			ctx = context.WithValue(ctx, "userIP", claims.IP)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
