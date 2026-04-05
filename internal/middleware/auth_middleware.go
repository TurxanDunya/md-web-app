package middleware

import (
	"context"
	"log/slog"
	"md_api/internal/config"
	"net/http"
	"strings"
	"time"

	"md_api/internal/handlers"

	"github.com/golang-jwt/jwt"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
)

func RequestLogger(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		logger.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	}
}

func AuthMiddleware(cfg *config.Config, logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Warn("auth failed: missing authorization header",
				"method", r.Method, "path", r.URL.Path)
			handlers.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Authorization header is required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" || tokenString == authHeader {
			logger.Warn("auth failed: malformed bearer token",
				"method", r.Method, "path", r.URL.Path)
			handlers.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			return
		}

		claims, err := validateJWT(tokenString, cfg)
		if err != nil {
			logger.Warn("auth failed: invalid token",
				"method", r.Method, "path", r.URL.Path, "error", err)
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
