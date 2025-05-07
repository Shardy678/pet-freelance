package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceOffer struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	FreelancerID uuid.UUID `gorm:"type:uuid;not null;index"`
	ServiceID    uuid.UUID `gorm:"type:uuid;not null;index"`

	Title       string `gorm:"type:varchar(150);not null"`
	Description string `gorm:"type:text;not null"`

	Price     float32 `gorm:"not null"`
	Currency  string  `gorm:"type:char(3);not null"`
	PriceType string  `gorm:"type:varchar(20);not null"`

	DurationEstimateMin int `gorm:"not null; default:60"`

	IsActive  bool           `gorm:"not null;default:true"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
