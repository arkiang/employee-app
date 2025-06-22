package payroll

import (
	"context"
	common "employee-app/internal/common/model"
	"employee-app/internal/model/entity"
	"employee-app/internal/repository"
	"errors"
	"fmt"
	"time"
)

type payrollPeriodUsecase struct {
	payrollRepo repository.PayrollPeriodRepository
}

// Constructor
func New(repo repository.PayrollPeriodRepository) PayrollPeriodUsecase {
	return &payrollPeriodUsecase{
		payrollRepo: repo,
	}
}

// CreatePeriod ensures no overlapping payroll period exists and creates a new one
func (u *payrollPeriodUsecase) CreatePeriod(ctx context.Context, startDate, endDate time.Time, adminID uint) error {
	if endDate.Before(startDate) {
		return errors.New("end date cannot be before start date")
	}

	// Fetch all periods to check for overlaps
	existingPeriods, err := u.payrollRepo.List(ctx, common.CommonFilter{})
	if err != nil {
		return fmt.Errorf("failed to fetch payroll periods: %w", err)
	}

	for _, p := range existingPeriods {
		if (startDate.Before(p.EndDate) && endDate.After(p.StartDate)) || startDate.Equal(p.StartDate) || endDate.Equal(p.EndDate) {
			return errors.New("new payroll period overlaps with existing one")
		}
	}

	period := entity.PayrollPeriod{
		StartDate: startDate,
		EndDate:   endDate,
		CreatedBy: adminID,
		UpdatedBy: adminID,
	}

	_, err = u.payrollRepo.Create(ctx, period)
	return err
}

func (u *payrollPeriodUsecase) ListPeriod(ctx context.Context, filter common.CommonFilter) ([]*entity.PayrollPeriod, error) {
	periods, err := u.payrollRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	
	return periods, nil
}

// GetPeriodByID returns a specific payroll period
func (u *payrollPeriodUsecase) GetPeriodByID(ctx context.Context, id uint) (*entity.PayrollPeriod, error) {
	return u.payrollRepo.GetByID(ctx, id)
}