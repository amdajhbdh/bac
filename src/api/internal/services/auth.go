package services

import (
	"github.com/bac-unified/api/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtSecret []byte
}

func NewAuthService() *AuthService {
	return &AuthService{
		jwtSecret: []byte("bac-unified-secret-key-change-in-production"),
	}
}

func (s *AuthService) Register(username, email, password string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return user, nil
}

func (s *AuthService) GenerateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userIDStr := claims["user_id"].(string)
		userID, _ := uuid.Parse(userIDStr)
		return userID, nil
	}

	return uuid.Nil, jwt.ErrSignatureInvalid
}
