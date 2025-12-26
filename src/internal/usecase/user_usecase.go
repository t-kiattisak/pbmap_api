package usecase

import (
	"pbmap_api/src/domain"
	"pbmap_api/src/internal/repository"

	"github.com/google/uuid"
)

type UserUsecase interface {
	CreateUser(user *domain.User) error
	GetUser(id uuid.UUID) (*domain.User, error)
	UpdateUser(user *domain.User) error
	DeleteUser(id uuid.UUID) error
	ListUsers() ([]domain.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo}
}

func (u *userUsecase) CreateUser(user *domain.User) error {
	return u.userRepo.Create(user)
}

func (u *userUsecase) GetUser(id uuid.UUID) (*domain.User, error) {
	return u.userRepo.FindByID(id)
}

func (u *userUsecase) UpdateUser(user *domain.User) error {
	return u.userRepo.Update(user)
}

func (u *userUsecase) DeleteUser(id uuid.UUID) error {
	return u.userRepo.Delete(id)
}

func (u *userUsecase) ListUsers() ([]domain.User, error) {
	return u.userRepo.FindAll()
}
