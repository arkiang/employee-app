package repository

import (
	"context"
	"employee-app/internal/model/entity"
	"employee-app/internal/common/model"
)

type PayrollPeriodRepository interface {
	Create(ctx context.Context, spec *entity.PayrollPeriod) (*entity.PayrollPeriod, error)
	GetByID(ctx context.Context, id uint) (*entity.PayrollPeriod, error)
	List(ctx context.Context, filter common.CommonFilter) ([]*entity.PayrollPeriod, error)
}