package repository

import (
	"context"
	"employee-app/internal/model"
	"employee-app/internal/model/entity"
)

type ReimbursementRepository interface {
	Create(ctx context.Context, spec entity.Reimbursement) (*entity.Reimbursement, error)
	GetByID(ctx context.Context, id uint) (*entity.Reimbursement, error)
	List(ctx context.Context, filter model.EmployeePeriodFilter) ([]*entity.Reimbursement, error)
}