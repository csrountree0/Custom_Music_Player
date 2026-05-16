package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Playlist struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey"        json:"id"`
	OwnerID         uuid.UUID      `gorm:"type:uuid;not null;index"    json:"owner_id"`
	Name            string         `gorm:"not null"                    json:"name"`
	Description     string         `                                   json:"description,omitempty"`
	IsPublic        bool           `gorm:"default:false"               json:"is_public"`
	IsCollaborative bool           `gorm:"default:false"               json:"is_collaborative"`
	CoverURL        string         `                                   json:"cover_url,omitempty"`
	FollowerCount   int            `gorm:"default:0"                   json:"follower_count"`
	CreatedAt       time.Time      `                                   json:"created_at"`
	UpdatedAt       time.Time      `                                   json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index"                       json:"-"`
}

func (p *Playlist) BeforeCreate(_ *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type PlaylistTrack struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"     json:"id"`
	PlaylistID uuid.UUID `gorm:"type:uuid;not null;index" json:"playlist_id"`
	TrackID    uuid.UUID `gorm:"type:uuid;not null"       json:"track_id"`
	AddedBy    uuid.UUID `gorm:"type:uuid;not null"       json:"added_by"`
	AddedAt    time.Time `                                json:"added_at"`
	Position   int       `gorm:"not null"                 json:"position"`
}

func (pt *PlaylistTrack) BeforeCreate(_ *gorm.DB) error {
	if pt.ID == uuid.Nil {
		pt.ID = uuid.New()
	}
	if pt.AddedAt.IsZero() {
		pt.AddedAt = time.Now()
	}
	return nil
}
