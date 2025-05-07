package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name               string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	Description        string    `gorm:"type:text"`
	BasePrice          float64   `gorm:"not null;default:0"`
	DefaultDurationMin int       `gorm:"not null;default:60"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (s *Service) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
