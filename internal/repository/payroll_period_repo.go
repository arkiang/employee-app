package repository

import (
	"context"
	common "employee-app/internal/common/model"
	"employee-app/internal/model/entity"

	"gorm.io/gorm"
)

type PayrollPeriodRepository interface {
	Create(ctx context.Context, spec *entity.PayrollPeriod) (*entity.PayrollPeriod, error)
	GetByID(ctx context.Context, id uint) (*entity.PayrollPeriod, error)
	List(ctx context.Context, filter common.CommonFilter) ([]*entity.PayrollPeriod, error)
	MarkProcessedTx(ctx context.Context, tx *gorm.DB, id, userId uint) (*entity.PayrollPeriod, error)
}