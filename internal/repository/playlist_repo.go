package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musicapp/backend/internal/models"
	"musicapp/backend/pkg/pagination"
)

type TrackPosition struct {
	TrackID  uuid.UUID
	Position int
}

type PlaylistRepository interface {
	List(ownerID uuid.UUID, p pagination.Params) ([]models.Playlist, int64, error)
	Create(playlist *models.Playlist) error
	FindByID(id, ownerID uuid.UUID) (*models.Playlist, error)
	Update(playlist *models.Playlist) error
	Delete(id, ownerID uuid.UUID) error
	AddTrack(playlistID, trackID, addedBy uuid.UUID, position int) error
	RemoveTrack(playlistID, trackID uuid.UUID) error
	ListTracks(playlistID uuid.UUID) ([]models.PlaylistTrack, error)
	ReorderTracks(playlistID uuid.UUID, positions []TrackPosition) error
	MaxPosition(playlistID uuid.UUID) (int, error)
}

type gormPlaylistRepository struct {
	db *gorm.DB
}

func NewPlaylistRepository(db *gorm.DB) PlaylistRepository {
	return &gormPlaylistRepository{db: db}
}

func (r *gormPlaylistRepository) List(ownerID uuid.UUID, p pagination.Params) ([]models.Playlist, int64, error) {
	var playlists []models.Playlist
	var total int64

	q := r.db.Model(&models.Playlist{}).Where("owner_id = ?", ownerID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("counting playlists: %w", err)
	}
	if err := q.Offset(p.Offset()).Limit(p.Limit).Find(&playlists).Error; err != nil {
		return nil, 0, fmt.Errorf("listing playlists: %w", err)
	}
	return playlists, total, nil
}

func (r *gormPlaylistRepository) Create(playlist *models.Playlist) error {
	return r.db.Create(playlist).Error
}

func (r *gormPlaylistRepository) FindByID(id, ownerID uuid.UUID) (*models.Playlist, error) {
	var playlist models.Playlist
	err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).First(&playlist).Error
	if err != nil {
		return nil, fmt.Errorf("finding playlist: %w", err)
	}
	return &playlist, nil
}

func (r *gormPlaylistRepository) Update(playlist *models.Playlist) error {
	return r.db.Save(playlist).Error
}

func (r *gormPlaylistRepository) Delete(id, ownerID uuid.UUID) error {
	return r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&models.Playlist{}).Error
}

func (r *gormPlaylistRepository) AddTrack(playlistID, trackID, addedBy uuid.UUID, position int) error {
	pt := models.PlaylistTrack{
		PlaylistID: playlistID,
		TrackID:    trackID,
		AddedBy:    addedBy,
		Position:   position,
	}
	return r.db.Create(&pt).Error
}

func (r *gormPlaylistRepository) RemoveTrack(playlistID, trackID uuid.UUID) error {
	return r.db.Where("playlist_id = ? AND track_id = ?", playlistID, trackID).
		Delete(&models.PlaylistTrack{}).Error
}

func (r *gormPlaylistRepository) ListTracks(playlistID uuid.UUID) ([]models.PlaylistTrack, error) {
	var rows []models.PlaylistTrack
	err := r.db.Where("playlist_id = ?", playlistID).Order("position ASC").Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("listing playlist tracks: %w", err)
	}
	return rows, nil
}

func (r *gormPlaylistRepository) ReorderTracks(playlistID uuid.UUID, positions []TrackPosition) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, tp := range positions {
			err := tx.Model(&models.PlaylistTrack{}).
				Where("playlist_id = ? AND track_id = ?", playlistID, tp.TrackID).
				Update("position", tp.Position).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *gormPlaylistRepository) MaxPosition(playlistID uuid.UUID) (int, error) {
	var max int
	err := r.db.Model(&models.PlaylistTrack{}).
		Where("playlist_id = ?", playlistID).
		Select("COALESCE(MAX(position), 0)").
		Scan(&max).Error
	if err != nil {
		return 0, fmt.Errorf("getting max position: %w", err)
	}
	return max, nil
}
