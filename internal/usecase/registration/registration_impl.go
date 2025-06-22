package registration

import (
	"context"
	"employee-app/internal/model/entity"
	"employee-app/internal/repository"
	"employee-app/pkg/security"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type registrationUsecase struct {
	userRepo       repository.UserRepository
	employeeRepo   repository.EmployeeRepository
}

func New(userRepo repository.UserRepository, empRepo repository.EmployeeRepository) RegistrationUsecase {
	return &registrationUsecase{
		userRepo:     userRepo,
		employeeRepo: empRepo,
	}
}

func (u *registrationUsecase) RegisterEmployee(ctx context.Context, user entity.User, employee entity.Employee, password string) error {
	if user.Username == "" || password == "" {
		return errors.New("username and password are required")
	}
	
	// Generate salt and hash
	salt, err := security.GenerateSalt()
	if err != nil {
		return err
	}

	hashedPassword, err := security.HashPassword(password, salt)
	if err != nil {
		return err
	}

	err = u.userRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
		user.PasswordHash = hashedPassword
		user.Salt = salt
		userResp, err := u.userRepo.CreateTx(ctx, tx, &user)
		if err != nil {
			return err
		}
		
		employee.UserID = userResp.ID
		if err := u.employeeRepo.CreateTx(ctx, tx, &employee); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}