package reimbursement

import (
	"context"
	"employee-app/internal/api/dto/common"
	"employee-app/internal/model"
	"employee-app/internal/model/entity"

	"github.com/shopspring/decimal"
)

type ReimbursementUsecase interface {
	SubmitReimbursement(ctx context.Context, empID uint, date common.DateOnly, amount decimal.Decimal, description *string) error
	GetReimbursementsForPeriod(ctx context.Context, spec model.EmployeePeriodFilter) ([]*entity.Reimbursement, error)
}