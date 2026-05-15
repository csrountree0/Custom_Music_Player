package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musicapp/backend/internal/models"
	"musicapp/backend/pkg/pagination"
)

type TrackRepository interface {
	List(ownerID uuid.UUID, p pagination.Params) ([]models.Track, int64, error)
	Create(track *models.Track) error
	FindByID(id, ownerID uuid.UUID) (*models.Track, error)
	Update(track *models.Track) error
	Delete(id, ownerID uuid.UUID) error
}

type gormTrackRepository struct {
	db *gorm.DB
}

func NewTrackRepository(db *gorm.DB) TrackRepository {
	return &gormTrackRepository{db: db}
}

func (r *gormTrackRepository) List(ownerID uuid.UUID, p pagination.Params) ([]models.Track, int64, error) {
	var tracks []models.Track
	var total int64

	q := r.db.Model(&models.Track{}).Where("owner_id = ?", ownerID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("counting tracks: %w", err)
	}
	if err := q.Offset(p.Offset()).Limit(p.Limit).Find(&tracks).Error; err != nil {
		return nil, 0, fmt.Errorf("listing tracks: %w", err)
	}
	return tracks, total, nil
}

func (r *gormTrackRepository) Create(track *models.Track) error {
	return r.db.Create(track).Error
}

func (r *gormTrackRepository) FindByID(id, ownerID uuid.UUID) (*models.Track, error) {
	var track models.Track
	err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).First(&track).Error
	if err != nil {
		return nil, fmt.Errorf("finding track: %w", err)
	}
	return &track, nil
}

func (r *gormTrackRepository) Update(track *models.Track) error {
	return r.db.Save(track).Error
}

func (r *gormTrackRepository) Delete(id, ownerID uuid.UUID) error {
	return r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&models.Track{}).Error
}
