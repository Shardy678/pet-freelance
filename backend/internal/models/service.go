package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Name               string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"name"`
	Description        string         `gorm:"type:text" json:"description,omitempty"`
	BasePrice          float64        `gorm:"not null;default:0" json:"basePrice"`
	DefaultDurationMin int            `gorm:"not null;default:60" json:"defaultDurationMin"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (s *Service) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
