package repositories

import (
	"context"
	"pbmap_api/src/internal/domain/entities"
	"pbmap_api/src/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	return GetDB(ctx, r.db).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := GetDB(ctx, r.db).Preload("SocialAccounts").
		Preload("SpecialCredential").
		Preload("Devices").
		Preload("Sessions").
		First(&user, "id = ?", id).Error
	return &user, err
}

func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	return GetDB(ctx, r.db).Model(user).Updates(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return GetDB(ctx, r.db).Delete(&entities.User{}, "id = ?", id).Error
}

func (r *userRepository) FindAll(ctx context.Context) ([]entities.User, error) {
	var users []entities.User
	err := GetDB(ctx, r.db).Find(&users).Error
	return users, err
}

func (r *userRepository) FindBySocialID(ctx context.Context, provider, providerID string) (*entities.User, error) {
	var user entities.User
	err := GetDB(ctx, r.db).
		Joins("JOIN user_social_accounts ON user_social_accounts.user_id = users.id").
		Where("user_social_accounts.provider = ? AND user_social_accounts.provider_id = ?", provider, providerID).
		Preload("SocialAccounts").
		Preload("SpecialCredential").
		Preload("Devices").
		Preload("Sessions").
		First(&user).Error
	return &user, err
}
