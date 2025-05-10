package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Email              string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash       string         `gorm:"type:char(60);not null" json:"-"`
	Phone              *string        `gorm:"type:varchar(20);uniqueIndex" json:"phone,omitempty"`
	Role               string         `gorm:"type:varchar(20);index;not null" json:"role"`
	ProfilePhotoURL    *string        `gorm:"type:text" json:"profilePhotoUrl,omitempty"`
	IsEmailVerified    bool           `gorm:"default:false" json:"isEmailVerified"`
	IsTwoFactorEnabled bool           `gorm:"default:false" json:"isTwoFactorEnabled"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
