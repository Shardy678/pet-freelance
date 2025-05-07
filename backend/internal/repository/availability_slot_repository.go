package repository

import (
	"context"
	"time"

	"github.com/shardy678/pet-freelance/backend/internal/models"
	"gorm.io/gorm"
)

type AvailabilitySlotRepository struct {
	db *gorm.DB
}

func NewAvailabilitySlotRepository(db *gorm.DB) *AvailabilitySlotRepository {
	return &AvailabilitySlotRepository{db}
}

func (r *AvailabilitySlotRepository) Create(ctx context.Context, slot *models.AvailabilitySlot) error {
	return r.db.WithContext(ctx).Create(slot).Error
}

func (r *AvailabilitySlotRepository) FindByID(ctx context.Context, id any) (*models.AvailabilitySlot, error) {
	var slot models.AvailabilitySlot
	if err := r.db.WithContext(ctx).First(&slot, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *AvailabilitySlotRepository) ListByOffer(ctx context.Context, offerID any, onlyAvailable bool, from, to time.Time) ([]models.AvailabilitySlot, error) {
	q := r.db.WithContext(ctx).
		Where("offer_id = ?", offerID).
		Where("start_time >= ? AND start_time < ?", from, to).
		Order("start_time ASC")

	if onlyAvailable {
		q = q.Where("is_booked = ?", false)
	}

	var slots []models.AvailabilitySlot
	if err := q.Find(&slots).Error; err != nil {
		return nil, err
	}
	return slots, nil
}

func (r *AvailabilitySlotRepository) Update(ctx context.Context, slot *models.AvailabilitySlot) error {
	return r.db.WithContext(ctx).Save(slot).Error
}

func (r *AvailabilitySlotRepository) Delete(ctx context.Context, id any) error {
	return r.db.WithContext(ctx).Delete(&models.AvailabilitySlot{}, "id = ?", id).Error
}
