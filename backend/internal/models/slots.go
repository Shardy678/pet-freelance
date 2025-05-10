package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AvailabilitySlot struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OfferID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"offerId"`
	StartTime time.Time      `gorm:"not null;index:idx_slot_time,priority:1" json:"startTime"`
	EndTime   time.Time      `gorm:"not null" json:"endTime"`
	IsBooked  bool           `gorm:"not null;default:false;index" json:"isBooked"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (s *AvailabilitySlot) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
