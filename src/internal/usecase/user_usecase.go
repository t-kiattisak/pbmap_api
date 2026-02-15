package usecase

import (
	"context"
	"pbmap_api/src/internal/domain/entities"
	"pbmap_api/src/internal/domain/repositories"
	"pbmap_api/src/internal/dto"

	"github.com/google/uuid"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUser(ctx context.Context, id uuid.UUID) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context) ([]entities.User, error)
	SyncUserFromSocial(ctx context.Context, input dto.CreateUserFromSocialInput) (*entities.User, error)
	UpsertDevice(ctx context.Context, device *entities.UserDevice) error
}

type userUsecase struct {
	userRepo   repositories.UserRepository
	deviceRepo repositories.DeviceRepository
}

func NewUserUsecase(userRepo repositories.UserRepository, deviceRepo repositories.DeviceRepository) UserUsecase {
	return &userUsecase{
		userRepo:   userRepo,
		deviceRepo: deviceRepo,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, user *entities.User) error {
	return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) GetUser(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return u.userRepo.FindByID(ctx, id)
}

func (u *userUsecase) UpdateUser(ctx context.Context, user *entities.User) error {
	return u.userRepo.Update(ctx, user)
}

func (u *userUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.userRepo.Delete(ctx, id)
}

func (u *userUsecase) ListUsers(ctx context.Context) ([]entities.User, error) {
	return u.userRepo.FindAll(ctx)
}

func (u *userUsecase) SyncUserFromSocial(ctx context.Context, input dto.CreateUserFromSocialInput) (*entities.User, error) {
	user, err := u.userRepo.FindBySocialID(ctx, input.Provider, input.ProviderID)
	if err == nil {
		return user, nil
	}

	newUser := &entities.User{
		DisplayName: input.DisplayName,
		Role:        "citizen",
	}

	if input.Email != "" {
		newUser.Email = &input.Email
	}
	newUser.SocialAccounts = []entities.UserSocialAccount{
		{
			Provider:   input.Provider,
			ProviderID: input.ProviderID,
		},
	}

	if err := u.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (u *userUsecase) UpsertDevice(ctx context.Context, device *entities.UserDevice) error {
	return u.deviceRepo.UpsertDevice(ctx, device)
}
