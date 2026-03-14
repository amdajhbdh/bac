package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/bac-unified/api/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidToken       = errors.New("invalid token")
)

type AuthService struct {
	jwtSecret []byte
	db        *pgxpool.Pool
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func NewAuthService(db *pgxpool.Pool) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "bac-unified-secret-key-change-in-production"
	}

	return &AuthService{
		jwtSecret: []byte(secret),
		db:        db,
	}
}

func (s *AuthService) Register(username, email, password string) (*models.User, error) {
	// Check if user exists
	if s.db != nil {
		ctx := context.Background()
		var exists bool
		err := s.db.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM users WHERE username = $1 OR email = $2)",
			username, email).Scan(&exists)

		if err == nil && exists {
			return nil, ErrUserExists
		}
	}

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

	// Save to database if available
	if s.db != nil {
		ctx := context.Background()
		_, err = s.db.Exec(ctx, `
			INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			user.ID, user.Username, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt)

		if err != nil {
			slog.Warn("failed to save user to database", "error", err)
		}
	}

	return user, nil
}

func (s *AuthService) Login(username, password string) (*models.User, *TokenPair, error) {
	var user *models.User
	var storedHash string

	// Try to get user from database
	if s.db != nil {
		ctx := context.Background()
		var userID uuid.UUID
		var email, usernameDB string

		err := s.db.QueryRow(ctx, `
			SELECT id, username, email, password_hash 
			FROM users 
			WHERE username = $1 OR email = $1`, username).Scan(
			&userID, &usernameDB, &email, &storedHash)

		if err != nil {
			return nil, nil, ErrInvalidCredentials
		}

		user = &models.User{
			ID:       userID,
			Username: usernameDB,
			Email:    email,
		}
	} else {
		// Demo mode - accept any credentials
		user = &models.User{
			ID:       uuid.New(),
			Username: username,
			Email:    username + "@demo.local",
		}
		storedHash = "$2a$10$demo" // Won't match any password
	}

	// Verify password
	if storedHash != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
			return nil, nil, ErrInvalidCredentials
		}
	}

	// Generate tokens
	tokens, err := s.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

func (s *AuthService) GenerateTokenPair(userID uuid.UUID) (*TokenPair, error) {
	// Access token - short lived
	accessClaims := jwt.MapClaims{
		"user_id": userID.String(),
		"type":    "access",
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Refresh token - longer lived
	refreshClaims := jwt.MapClaims{
		"user_id": userID.String(),
		"type":    "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    time.Now().Add(time.Hour * 24),
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*TokenPair, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Verify it's a refresh token
	if claims["type"] != "refresh" {
		return nil, ErrInvalidToken
	}

	userIDStr := claims["user_id"].(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return s.GenerateTokenPair(userID)
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

func (s *AuthService) GetUser(userID uuid.UUID) (*models.User, error) {
	if s.db == nil {
		return &models.User{ID: userID}, nil
	}

	ctx := context.Background()
	var user models.User
	err := s.db.QueryRow(ctx, `
		SELECT id, username, email, full_name, school, region, role, 
		       points, level, streak_days, created_at
		FROM users WHERE id = $1`, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.FullName,
		&user.School, &user.Region, &user.Role, &user.Points,
		&user.Level, &user.StreakDays, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
