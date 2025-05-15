package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Booking struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	OfferID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"offerId"`
	SlotID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"slotId"`
	OwnerID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"ownerId"`
	Status    string         `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (b *Booking) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
