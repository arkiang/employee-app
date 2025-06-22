package repository

import (
	"context"
	"employee-app/internal/model"
	"employee-app/internal/model/entity"
)

type OvertimeRepository interface {
	Create(ctx context.Context, spec entity.Overtime) (*entity.Overtime, error)
	GetByID(ctx context.Context, id uint) (*entity.Overtime, error)
	List(ctx context.Context, filter model.EmployeePeriodFilter) ([]*entity.Overtime, error)
}