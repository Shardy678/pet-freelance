package repository

import (
	"context"

	"github.com/shardy678/pet-freelance/backend/internal/models"
	"gorm.io/gorm"
)

type ServiceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

// Create inserts a new Service record
func (r *ServiceRepository) Create(ctx context.Context, svc *models.Service) error {
	return r.db.WithContext(ctx).Create(svc).Error
}

// FindByID looks up a Service by its UUID
func (r *ServiceRepository) FindByID(ctx context.Context, id any) (*models.Service, error) {
	var svc models.Service
	if err := r.db.WithContext(ctx).First(&svc, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &svc, nil
}

// ListAll returns all non-deleted Services, ordered by name
func (r *ServiceRepository) ListAll(ctx context.Context) ([]models.Service, error) {
	var list []models.Service
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL").
		Order("name").
		Find(&list).Error
	return list, err
}
