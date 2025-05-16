package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shardy678/pet-freelance/backend/internal/models"
	"github.com/shardy678/pet-freelance/backend/internal/repository"
	"gorm.io/gorm"
)

var ErrSlotAlreadyBooked = errors.New("slot already booked")

type BookingService struct {
	bookingRepo *repository.BookingRepository
	slotRepo    *repository.AvailabilitySlotRepository
	activitySvc *ActivityService
	db          *gorm.DB
}

func NewBookingService(
	bookingRepo *repository.BookingRepository,
	slotRepo *repository.AvailabilitySlotRepository,
	activitySvc *ActivityService,
	db *gorm.DB,
) *BookingService {
	return &BookingService{bookingRepo, slotRepo, activitySvc, db}
}

// BookSlot reserves a slot and creates a booking within a single transaction.
func (s *BookingService) BookSlot(
	ctx context.Context,
	offerID, slotID, ownerID uuid.UUID,
) (*models.Booking, error) {
	var booking *models.Booking

	// 1) Transactionally reserve the slot & create booking
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		slot, err := s.slotRepo.FindByID(ctx, slotID)
		if err != nil {
			return err
		}
		if slot.IsBooked {
			return ErrSlotAlreadyBooked
		}
		slot.IsBooked = true
		if err := s.slotRepo.Update(ctx, slot); err != nil {
			return err
		}
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
	if err != nil {
		return nil, err
	}

	// 2) Emit a Recent-Activity record
	timeStr := booking.CreatedAt.Format("Jan 2, 2006 at 15:04")
	title := "Grooming appointment confirmed"
	message := fmt.Sprintf(
		"Your booking is confirmed for %s.",
		timeStr,
	)
	// fire-and-forget
	if emitErr := s.activitySvc.Emit(ctx, ownerID, title, message, "appointment"); emitErr != nil {
		fmt.Printf("warning: could not emit activity: %v\n", emitErr)
	}

	return booking, nil
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
