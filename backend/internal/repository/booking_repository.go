package repository

import (
	"context"

	"github.com/shardy678/pet-freelance/backend/internal/models"
	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db}
}

func (r *BookingRepository) Create(ctx context.Context, b *models.Booking) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *BookingRepository) FindByID(ctx context.Context, id any) (*models.Booking, error) {
	var b models.Booking
	if err := r.db.WithContext(ctx).First(&b, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookingRepository) ListByOwner(ctx context.Context, ownerID any) ([]models.Booking, error) {
	var list []models.Booking
	err := r.db.WithContext(ctx).
		Where("owner_id = ?", ownerID).
		Order("created_at desc").
		Find(&list).Error
	return list, err
}

func (r *BookingRepository) ListByOffer(ctx context.Context, offerID any) ([]models.Booking, error) {
	var list []models.Booking
	err := r.db.WithContext(ctx).
		Where("offer_id = ?", offerID).
		Order("created_at desc").
		Find(&list).Error
	return list, err
}
