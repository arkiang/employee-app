package reimbursement

import (
	"context"
	"time"
	"employee-app/internal/model/entity"
	"employee-app/internal/model"
	"github.com/shopspring/decimal"
)

type ReimbursementUsecase interface {
	SubmitReimbursement(ctx context.Context, empID uint, date time.Time, amount decimal.Decimal, description *string) error
	GetReimbursementsForPeriod(ctx context.Context, spec model.EmployeePeriodFilter) ([]*entity.Reimbursement, error)
}