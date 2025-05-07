package repository

import (
	"context"

	"github.com/shardy678/pet-freelance/backend/internal/models"
	"gorm.io/gorm"
)

type ServiceOfferRepository struct{ db *gorm.DB }

func NewServiceOfferRepository(db *gorm.DB) *ServiceOfferRepository {
	return &ServiceOfferRepository{db}
}

func (r *ServiceOfferRepository) Create(ctx context.Context, o *models.ServiceOffer) error {
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *ServiceOfferRepository) FindByID(ctx context.Context, id any) (*models.ServiceOffer, error) {
	var o models.ServiceOffer
	if err := r.db.WithContext(ctx).First(&o, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *ServiceOfferRepository) ListAll(ctx context.Context) ([]models.ServiceOffer, error) {
	var list []models.ServiceOffer
	err := r.db.WithContext(ctx).
		Where("is_active= ?", true).
		Order("created_at desc").
		Find(&list).Error
	return list, err
}
