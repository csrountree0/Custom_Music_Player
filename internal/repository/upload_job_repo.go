package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musicapp/backend/internal/models"
)

type UploadJobRepository interface {
	Create(job *models.UploadJob) error
	FindByID(id, ownerID uuid.UUID) (*models.UploadJob, error)
	UpdateStatus(id uuid.UUID, status, errorMsg string) error
	SetTrackID(id, trackID uuid.UUID) error
}

type gormUploadJobRepository struct {
	db *gorm.DB
}

func NewUploadJobRepository(db *gorm.DB) UploadJobRepository {
	return &gormUploadJobRepository{db: db}
}

func (r *gormUploadJobRepository) Create(job *models.UploadJob) error {
	return r.db.Create(job).Error
}

func (r *gormUploadJobRepository) FindByID(id, ownerID uuid.UUID) (*models.UploadJob, error) {
	var job models.UploadJob
	err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).First(&job).Error
	if err != nil {
		return nil, fmt.Errorf("finding upload job: %w", err)
	}
	return &job, nil
}

func (r *gormUploadJobRepository) UpdateStatus(id uuid.UUID, status, errorMsg string) error {
	return r.db.Model(&models.UploadJob{}).Where("id = ?", id).
		Updates(map[string]any{"status": status, "error_msg": errorMsg}).Error
}

func (r *gormUploadJobRepository) SetTrackID(id, trackID uuid.UUID) error {
	return r.db.Model(&models.UploadJob{}).Where("id = ?", id).
		Updates(map[string]any{"track_id": trackID}).Error
}
