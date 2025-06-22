package payroll

import (
	"context"
	common "employee-app/internal/common/model"
	"employee-app/internal/model/entity"
	"time"
)

type PayrollPeriodUsecase interface {
	CreatePeriod(ctx context.Context, startDate, endDate time.Time, adminID uint) error
	ListPeriod(ctx context.Context, filter common.CommonFilter) ([]*entity.PayrollPeriod, error)
	GetPeriodByID(ctx context.Context, id uint) (*entity.PayrollPeriod, error)
}