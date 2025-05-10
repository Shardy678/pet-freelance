package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceOffer struct {
	ID                  uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	FreelancerID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"freelancerId"`
	ServiceID           uuid.UUID      `gorm:"type:uuid;not null;index" json:"serviceId"`
	Title               string         `gorm:"type:varchar(150);not null" json:"title"`
	Description         string         `gorm:"type:text;not null" json:"description"`
	Price               float32        `gorm:"not null" json:"price"`
	Currency            string         `gorm:"type:char(3);not null" json:"currency"`
	PriceType           string         `gorm:"type:varchar(20);not null" json:"priceType"`
	DurationEstimateMin int            `gorm:"not null;default:60" json:"durationEstimateMin"`
	IsActive            bool           `gorm:"not null;default:true" json:"isActive"`
	CreatedAt           time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt           time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (o *ServiceOffer) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}
