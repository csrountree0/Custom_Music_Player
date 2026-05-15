package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Artist struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"        json:"id"`
	OwnerID   uuid.UUID      `gorm:"type:uuid;not null;index"    json:"owner_id"`
	Name      string         `gorm:"not null"                    json:"name"`
	Bio       string         `                                   json:"bio,omitempty"`
	ImageURL  string         `                                   json:"image_url,omitempty"`
	CreatedAt time.Time      `                                   json:"created_at"`
	UpdatedAt time.Time      `                                   json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                       json:"-"`
}

func (a *Artist) BeforeCreate(_ *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

type Album struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey"        json:"id"`
	OwnerID      uuid.UUID      `gorm:"type:uuid;not null;index"    json:"owner_id"`
	ArtistID     *uuid.UUID     `gorm:"type:uuid"                   json:"artist_id,omitempty"`
	Title        string         `gorm:"not null"                    json:"title"`
	Year         int            `                                   json:"year,omitempty"`
	CoverArtPath string         `                                   json:"cover_art_path,omitempty"`
	CreatedAt    time.Time      `                                   json:"created_at"`
	UpdatedAt    time.Time      `                                   json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"                       json:"-"`
}

func (a *Album) BeforeCreate(_ *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

type Track struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"        json:"id"`
	OwnerID     uuid.UUID      `gorm:"type:uuid;not null;index"    json:"owner_id"`
	ArtistID    *uuid.UUID     `gorm:"type:uuid"                   json:"artist_id,omitempty"`
	AlbumID     *uuid.UUID     `gorm:"type:uuid"                   json:"album_id,omitempty"`
	Title       string         `gorm:"not null"                    json:"title"`
	DurationMs  int            `                                   json:"duration_ms,omitempty"`
	TrackNumber int            `                                   json:"track_number,omitempty"`
	DiscNumber  int            `                                   json:"disc_number,omitempty"`
	Year        int            `                                   json:"year,omitempty"`
	Visibility  string         `gorm:"not null;default:private"    json:"visibility"`
	Source      string         `gorm:"not null"                    json:"source"`
	SourceURL   string         `                                   json:"source_url,omitempty"`
	FilePath    string         `                                   json:"-"`
	FileSize    int64          `                                   json:"file_size,omitempty"`
	MimeType    string         `                                   json:"mime_type,omitempty"`
	PlayCount   int            `gorm:"default:0"                   json:"play_count"`
	CreatedAt   time.Time      `                                   json:"created_at"`
	UpdatedAt   time.Time      `                                   json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"                       json:"-"`
}

func (t *Track) BeforeCreate(_ *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type UploadJob struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"        json:"id"`
	OwnerID   uuid.UUID  `gorm:"type:uuid;not null;index"    json:"owner_id"`
	Status    string     `gorm:"not null"                    json:"status"`
	Source    string     `gorm:"not null"                    json:"source"`
	SourceURL string     `                                   json:"source_url,omitempty"`
	TrackID   *uuid.UUID `gorm:"type:uuid"                   json:"track_id,omitempty"`
	ErrorMsg  string     `                                   json:"error_msg,omitempty"`
	CreatedAt time.Time  `                                   json:"created_at"`
	UpdatedAt time.Time  `                                   json:"updated_at"`
}

func (j *UploadJob) BeforeCreate(_ *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}
