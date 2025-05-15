package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/shardy678/pet-freelance/backend/internal/models"
	"github.com/shardy678/pet-freelance/backend/internal/repository"
	"gorm.io/gorm"
)

var ErrSlotAlreadyBooked = errors.New("slot already booked")

type BookingService struct {
	bookingRepo *repository.BookingRepository
	slotRepo    *repository.AvailabilitySlotRepository
	db          *gorm.DB
}

func NewBookingService(
	bookingRepo *repository.BookingRepository,
	slotRepo *repository.AvailabilitySlotRepository,
	db *gorm.DB,
) *BookingService {
	return &BookingService{bookingRepo, slotRepo, db}
}

// BookSlot reserves a slot and creates a booking within a single transaction.
func (s *BookingService) BookSlot(
	ctx context.Context,
	offerID, slotID, ownerID uuid.UUID,
) (*models.Booking, error) {
	var booking *models.Booking

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Load and lock the slot
		slot, err := s.slotRepo.FindByID(ctx, slotID)
		if err != nil {
			return err
		}
		if slot.IsBooked {
			return ErrSlotAlreadyBooked
		}

		// 2. Mark the slot as booked
		slot.IsBooked = true
		if err := s.slotRepo.Update(ctx, slot); err != nil {
			return err
		}

		// 3. Create the booking record
		booking = &models.Booking{
			OfferID: offerID,
			SlotID:  slotID,
			OwnerID: ownerID,
			Status:  "pending",
		}
		if err := s.bookingRepo.Create(ctx, booking); err != nil {
			return err
		}

		return nil
	})

	return booking, err
}

func (s *BookingService) GetBooking(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	return s.bookingRepo.FindByID(ctx, id)
}

func (s *BookingService) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]models.Booking, error) {
	return s.bookingRepo.ListByOwner(ctx, ownerID)
}

func (s *BookingService) ListByOffer(ctx context.Context, offerID uuid.UUID) ([]models.Booking, error) {
	return s.bookingRepo.ListByOffer(ctx, offerID)
}
