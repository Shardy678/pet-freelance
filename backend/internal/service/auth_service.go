package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shardy678/pet-freelance/backend/internal/config"
	"github.com/shardy678/pet-freelance/backend/internal/models"
	"github.com/shardy678/pet-freelance/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	users *repository.UserRepository
	cfg   *config.AppConfig
}

func NewAuthService(u *repository.UserRepository, cfg *config.AppConfig) *AuthService {
	return &AuthService{u, cfg}
}

func hashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

func verifyPassword(hash, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}

func (s *AuthService) Register(ctx context.Context, email, pw, role string) (*models.User, error) {
	h, err := hashPassword(pw)
	if err != nil {
		return nil, err
	}
	user := &models.User{Email: email, PasswordHash: h, Role: role}
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

var ErrInvalidCredentials = errors.New("invalid email or password")

func (s *AuthService) Login(ctx context.Context, email, pw string) (string, error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}
	if err := verifyPassword(u.PasswordHash, pw); err != nil {
		return "", ErrInvalidCredentials
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  u.ID.String(),
		"role": u.Role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(s.cfg.JWTSecret))
}
