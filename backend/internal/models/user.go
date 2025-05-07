package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents everyone: pet owners, freelancers, and admins.
type User struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primaryKey"` // UUID generated in Go
	Email              string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash       string         `gorm:"type:char(60);not null"`          // bcrypt(60) fits here
	Phone              *string        `gorm:"type:varchar(20);uniqueIndex"`    // nullable
	Role               string         `gorm:"type:varchar(20);index;not null"` // "freelancer" | "owner" | "admin"
	ProfilePhotoURL    *string        `gorm:"type:text"`
	IsEmailVerified    bool           `gorm:"default:false"`
	IsTwoFactorEnabled bool           `gorm:"default:false"`
	CreatedAt          time.Time      `gorm:"autoCreateTime"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `gorm:"index"` // optional soft‑delete
}

// BeforeCreate will set a UUID rather than relying on a DB extension.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
