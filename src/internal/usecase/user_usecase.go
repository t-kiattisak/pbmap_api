package usecase

import (
	"context"
	"pbmap_api/src/domain"
	"pbmap_api/src/internal/dto"
	"pbmap_api/src/internal/repository"

	"github.com/google/uuid"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context) ([]domain.User, error)
	SyncUserFromSocial(ctx context.Context, input dto.CreateUserFromSocialInput) (*domain.User, error)
	UpsertDevice(ctx context.Context, device *domain.UserDevice) error
}

type userUsecase struct {
	userRepo   repository.UserRepository
	deviceRepo repository.DeviceRepository
}

func NewUserUsecase(userRepo repository.UserRepository, deviceRepo repository.DeviceRepository) UserUsecase {
	return &userUsecase{
		userRepo:   userRepo,
		deviceRepo: deviceRepo,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, user *domain.User) error {
	return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return u.userRepo.FindByID(ctx, id)
}

func (u *userUsecase) UpdateUser(ctx context.Context, user *domain.User) error {
	return u.userRepo.Update(ctx, user)
}

func (u *userUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.userRepo.Delete(ctx, id)
}

func (u *userUsecase) ListUsers(ctx context.Context) ([]domain.User, error) {
	return u.userRepo.FindAll(ctx)
}

func (u *userUsecase) SyncUserFromSocial(ctx context.Context, input dto.CreateUserFromSocialInput) (*domain.User, error) {
	user, err := u.userRepo.FindBySocialID(ctx, input.Provider, input.ProviderID)
	if err == nil {
		return user, nil
	}

	newUser := &domain.User{
		DisplayName: input.DisplayName,
		Role:        "citizen",
	}

	if input.Email != "" {
		newUser.Email = &input.Email
	}
	newUser.SocialAccounts = []domain.UserSocialAccount{
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

func (u *userUsecase) UpsertDevice(ctx context.Context, device *domain.UserDevice) error {
	return u.deviceRepo.UpsertDevice(ctx, device)
}
