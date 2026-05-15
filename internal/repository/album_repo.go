package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musicapp/backend/internal/models"
	"musicapp/backend/pkg/pagination"
)

type AlbumRepository interface {
	List(ownerID uuid.UUID, p pagination.Params) ([]models.Album, int64, error)
	Create(album *models.Album) error
	FindByID(id, ownerID uuid.UUID) (*models.Album, error)
	FindByTitleArtistOwner(title string, artistID *uuid.UUID, ownerID uuid.UUID) (*models.Album, error)
	Update(album *models.Album) error
	Delete(id, ownerID uuid.UUID) error
}

type gormAlbumRepository struct {
	db *gorm.DB
}

func NewAlbumRepository(db *gorm.DB) AlbumRepository {
	return &gormAlbumRepository{db: db}
}

func (r *gormAlbumRepository) List(ownerID uuid.UUID, p pagination.Params) ([]models.Album, int64, error) {
	var albums []models.Album
	var total int64

	q := r.db.Model(&models.Album{}).Where("owner_id = ?", ownerID)
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("counting albums: %w", err)
	}
	if err := q.Offset(p.Offset()).Limit(p.Limit).Find(&albums).Error; err != nil {
		return nil, 0, fmt.Errorf("listing albums: %w", err)
	}
	return albums, total, nil
}

func (r *gormAlbumRepository) Create(album *models.Album) error {
	return r.db.Create(album).Error
}

func (r *gormAlbumRepository) FindByID(id, ownerID uuid.UUID) (*models.Album, error) {
	var album models.Album
	err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).First(&album).Error
	if err != nil {
		return nil, fmt.Errorf("finding album: %w", err)
	}
	return &album, nil
}

func (r *gormAlbumRepository) FindByTitleArtistOwner(title string, artistID *uuid.UUID, ownerID uuid.UUID) (*models.Album, error) {
	var album models.Album
	q := r.db.Where("title = ? AND owner_id = ?", title, ownerID)
	if artistID == nil {
		q = q.Where("artist_id IS NULL")
	} else {
		q = q.Where("artist_id = ?", *artistID)
	}
	if err := q.First(&album).Error; err != nil {
		return nil, fmt.Errorf("finding album by title: %w", err)
	}
	return &album, nil
}

func (r *gormAlbumRepository) Update(album *models.Album) error {
	return r.db.Save(album).Error
}

func (r *gormAlbumRepository) Delete(id, ownerID uuid.UUID) error {
	return r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&models.Album{}).Error
}
