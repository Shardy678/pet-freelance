package repository

import (
	"context"

	"github.com/shardy678/pet-freelance/backend/internal/models"
	"gorm.io/gorm"
)

type ActivityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{db}
}

func (r *ActivityRepository) Create(ctx context.Context, a *models.Activity) error {
	return r.db.WithContext(ctx).Create(a).Error
}

// List recent N activities for a user, ordered newest-first
func (r *ActivityRepository) ListByUser(ctx context.Context, userID any, limit int) ([]models.Activity, error) {
	var list []models.Activity
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(limit).
		Find(&list).Error
	return list, err
}
