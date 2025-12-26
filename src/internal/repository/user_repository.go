package repository

import (
	"pbmap_api/src/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id uuid.UUID) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uuid.UUID) error
	FindAll() ([]domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("SocialAccounts").
		Preload("SpecialCredential").
		Preload("Devices").
		Preload("Sessions").
		First(&user, "id = ?", id).Error
	return &user, err
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Model(user).Updates(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.User{}, "id = ?", id).Error
}

func (r *userRepository) FindAll() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Find(&users).Error
	return users, err
}
