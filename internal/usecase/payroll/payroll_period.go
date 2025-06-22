package payroll

import (
	"context"
	common "employee-app/internal/common/model"
	commonDate "employee-app/internal/api/dto/common"
	"employee-app/internal/model/entity"
)

type PayrollPeriodUsecase interface {
	CreatePeriod(ctx context.Context, startDate, endDate commonDate.DateOnly, adminID uint) error
	ListPeriod(ctx context.Context, filter common.CommonFilter) ([]*entity.PayrollPeriod, error)
	GetPeriodByID(ctx context.Context, id uint) (*entity.PayrollPeriod, error)
}