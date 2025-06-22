package user

import (
	"context"
	"employee-app/internal/model/entity"
	"employee-app/internal/repository"
	"employee-app/pkg/security"
	"errors"
)

type userUsecase struct {
	userRepo     repository.UserRepository
	employeeRepo repository.EmployeeRepository
}

func New(userRepo repository.UserRepository, employeeRepo repository.EmployeeRepository) UserUsecase {
	return &userUsecase{
		userRepo:     userRepo,
		employeeRepo: employeeRepo,
	}
}

func (u *userUsecase) Login(ctx context.Context, username, password string) (*entity.User, error) {
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found/invalid credentials")
	}

	// Check password
	if !security.VerifyPassword(password, user.Salt, user.PasswordHash) {
		return nil, errors.New("user not found/invalid credentials")
	}

	return user, nil
}

func (u *userUsecase) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	return u.userRepo.GetByID(ctx, id)
}