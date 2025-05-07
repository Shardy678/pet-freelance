package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AvailabilitySlot struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	OfferID   uuid.UUID      `gorm:"type:uuid;not null;index"` // FK → service_offers.id
	StartTime time.Time      `gorm:"not null;index:idx_slot_time,priority:1"`
	EndTime   time.Time      `gorm:"not null"`
	IsBooked  bool           `gorm:"not null;default:false;index"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (s *AvailabilitySlot) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
