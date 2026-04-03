package middleware

import (
	"context"
	"md_api/internal/config"
	"net/http"
	"strings"

	"md_api/internal/handlers"

	"github.com/golang-jwt/jwt"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
)

func AuthMiddleware(cfg *config.Config, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			handlers.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" || tokenString == authHeader {
			handlers.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			return
		}

		claims, err := validateJWT(tokenString, cfg)
		if err != nil {
			handlers.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims["user_id"])
		ctx = context.WithValue(ctx, EmailKey, claims["email"])

		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func validateJWT(tokenString string, cfg *config.Config) (map[string]any, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
