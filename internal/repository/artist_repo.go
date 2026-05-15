package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musicapp/backend/internal/models"
	"musicapp/backend/pkg/pagination"
)

type ArtistRepository interface {
	List(ownerID uuid.UUID, p pagination.Params) ([]models.Artist, int64, error)
	Create(artist *models.Artist) error
	FindByID(id, ownerID uuid.UUID) (*models.Artist, error)
	FindByNameAndOwner(name string, ownerID uuid.UUID) (*models.Artist, error)
	Update(artist *models.Artist) error
	Delete(id, ownerID uuid.UUID) error
}

type gormArtistRepository struct {
	db *gorm.DB
}

func NewArtistRepository(db *gorm.DB) ArtistRepository {
	return &gormArtistRepository{db: db}
}

func (r *gormArtistRepository) List(ownerID uuid.UUID, p pagination.Params) ([]models.Artist, int64, error) {
	var artists []models.Artist
	var total int64

	q := r.db.Model(&models.Artist{}).Where("owner_id = ?", ownerID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("counting artists: %w", err)
	}
	if err := q.Offset(p.Offset()).Limit(p.Limit).Find(&artists).Error; err != nil {
		return nil, 0, fmt.Errorf("listing artists: %w", err)
	}
	return artists, total, nil
}

func (r *gormArtistRepository) Create(artist *models.Artist) error {
	return r.db.Create(artist).Error
}

func (r *gormArtistRepository) FindByID(id, ownerID uuid.UUID) (*models.Artist, error) {
	var artist models.Artist
	err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).First(&artist).Error
	if err != nil {
		return nil, fmt.Errorf("finding artist: %w", err)
	}
	return &artist, nil
}

func (r *gormArtistRepository) FindByNameAndOwner(name string, ownerID uuid.UUID) (*models.Artist, error) {
	var artist models.Artist
	err := r.db.Where("name = ? AND owner_id = ?", name, ownerID).First(&artist).Error
	if err != nil {
		return nil, fmt.Errorf("finding artist by name: %w", err)
	}
	return &artist, nil
}

func (r *gormArtistRepository) Update(artist *models.Artist) error {
	return r.db.Save(artist).Error
}

func (r *gormArtistRepository) Delete(id, ownerID uuid.UUID) error {
	return r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&models.Artist{}).Error
}
