package handlers

import (
	"md_api/internal/config"
	"md_api/internal/models"
	"md_api/internal/repository"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Token string `json:"token"`
}

func RegisterUserHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request RegisterUserRequest
		if err := readJSON(r, &request); err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if request.Email == "" || request.Password == "" {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "email and password are required"})
			return
		}

		bcryptedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
			return
		}

		user, err := repository.CreateUser(pool, &models.User{
			Email:    request.Email,
			Password: string(bcryptedPassword),
		})
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		WriteJSON(w, http.StatusCreated, user)
	}
}

func LoginUserHandler(pool *pgxpool.Pool, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request LoginUserRequest
		if err := readJSON(r, &request); err != nil {
			WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		user, err := repository.GetUserByEmail(pool, request.Email)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
			return
		}

		token, err := generateJWT(user.ID, user.Email, cfg)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
			return
		}

		WriteJSON(w, http.StatusOK, LoginUserResponse{Token: token})
	}
}

func generateJWT(userID string, email string, cfg *config.Config) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
