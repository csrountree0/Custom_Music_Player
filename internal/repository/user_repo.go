package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"musicapp/backend/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *gormUserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Where("deleted_at IS NULL").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, fmt.Errorf("finding user by id: %w", err)
	}
	return &user, nil
}

func (r *gormUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		return nil, fmt.Errorf("finding user by email: %w", err)
	}
	return &user, nil
}

func (r *gormUserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}
