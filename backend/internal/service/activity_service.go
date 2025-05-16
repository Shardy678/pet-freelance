package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shardy678/pet-freelance/backend/internal/models"
	"github.com/shardy678/pet-freelance/backend/internal/repository"
)

type ActivityService struct {
	repo *repository.ActivityRepository
}

func NewActivityService(repo *repository.ActivityRepository) *ActivityService {
	return &ActivityService{repo}
}

// Emit an activity
func (s *ActivityService) Emit(ctx context.Context, userID uuid.UUID, title, message, typ string) error {
	a := &models.Activity{
		UserID:  userID,
		Title:   title,
		Message: message,
		Type:    typ,
	}
	return s.repo.Create(ctx, a)
}

// List up to `limit` recent activities
func (s *ActivityService) List(ctx context.Context, userID uuid.UUID, limit int) ([]models.Activity, error) {
	return s.repo.ListByUser(ctx, userID, limit)
}
