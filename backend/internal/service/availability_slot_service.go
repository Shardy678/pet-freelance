package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shardy678/pet-freelance/backend/internal/models"
	"github.com/shardy678/pet-freelance/backend/internal/repository"
)

type AvailabilitySlotService struct {
	repo *repository.AvailabilitySlotRepository
}

func NewAvailabilitySlotService(r *repository.AvailabilitySlotRepository) *AvailabilitySlotService {
	return &AvailabilitySlotService{repo: r}
}

func (s *AvailabilitySlotService) CreateSlot(ctx context.Context, offerID uuid.UUID, start, end time.Time) (*models.AvailabilitySlot, error) {
	slot := &models.AvailabilitySlot{
		ID:        uuid.New(),
		OfferID:   offerID,
		StartTime: start,
		EndTime:   end,
		IsBooked:  false,
	}
	if err := s.repo.Create(ctx, slot); err != nil {
		return nil, err
	}
	return slot, nil
}

func (s *AvailabilitySlotService) ListSlots(ctx context.Context, offerID uuid.UUID, onlyAvailable bool, from, to time.Time) ([]models.AvailabilitySlot, error) {
	return s.repo.ListByOffer(ctx, offerID, onlyAvailable, from, to)
}

func (s *AvailabilitySlotService) UpdateSlot(ctx context.Context, slotID uuid.UUID, start, end time.Time, isBooked bool) (*models.AvailabilitySlot, error) {
	slot, err := s.repo.FindByID(ctx, slotID)
	if err != nil {
		return nil, err
	}
	slot.StartTime = start
	slot.EndTime = end
	slot.IsBooked = isBooked
	if err := s.repo.Update(ctx, slot); err != nil {
		return nil, err
	}
	return slot, nil
}

func (s *AvailabilitySlotService) DeleteSlot(ctx context.Context, slotID uuid.UUID) error {
	return s.repo.Delete(ctx, slotID)
}
