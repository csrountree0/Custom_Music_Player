package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey"  json:"id"`
	Email        string         `gorm:"uniqueIndex;not null"  json:"email"`
	Username     string         `gorm:"uniqueIndex;not null"  json:"username"`
	PasswordHash string         `gorm:"not null"              json:"-"`
	DisplayName  string         `                             json:"display_name,omitempty"`
	AvatarURL    string         `                             json:"avatar_url,omitempty"`
	Country      string         `gorm:"size:2"                json:"country,omitempty"`
	CreatedAt    time.Time      `                             json:"created_at"`
	UpdatedAt    time.Time      `                             json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"                 json:"-"`
}

func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

func (rt *RefreshToken) BeforeCreate(_ *gorm.DB) error {
	if rt.ID == uuid.Nil {
		rt.ID = uuid.New()
	}
	return nil
}
