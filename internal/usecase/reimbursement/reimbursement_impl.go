package reimbursement

import (
	"context"
	"employee-app/internal/api/dto/common"
	"employee-app/internal/api/middleware"
	"employee-app/internal/model"
	"employee-app/internal/model/entity"
	"employee-app/internal/repository"
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type reimbursementUsecase struct {
	reimbursementRepo repository.ReimbursementRepository
}

func New(reimbursementRepo repository.ReimbursementRepository) ReimbursementUsecase {
	return &reimbursementUsecase{
		reimbursementRepo: reimbursementRepo,
	}
}

func (u *reimbursementUsecase) SubmitReimbursement(ctx context.Context, empID uint, date common.DateOnly, amount decimal.Decimal, description *string) error {
	val := ctx.Value(middleware.ContextUserIDKey)
	userID, ok := val.(uint)
	if !ok {
		return errors.New("unauthorized or missing user ID")
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("reimbursement amount must be greater than 0")
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return fmt.Errorf("failed to load timezone: %w", err)
	}
	tLocal := date.In(loc)
	reimbursementDate := time.Date(
		tLocal.Year(), tLocal.Month(), tLocal.Day(),
		0, 0, 0, 0,
		loc,
	)

	r := entity.Reimbursement{
		EmployeeID:   empID,
		ReimbursementDate: reimbursementDate,
		Amount:       amount,
		Description:  description,
		CreatedBy:    userID,
		UpdatedBy:    userID,
	}

	_, err = u.reimbursementRepo.Create(ctx, &r)
	return err
}

func (u *reimbursementUsecase) GetReimbursementsForPeriod(ctx context.Context, spec model.EmployeePeriodFilter) ([]*entity.Reimbursement, error) {
	return u.reimbursementRepo.List(ctx, spec)
}