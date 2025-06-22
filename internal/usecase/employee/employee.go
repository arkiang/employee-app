package employee

import (
	"context"
	common "employee-app/internal/common/model"
	"employee-app/internal/model/entity"
)

type EmployeeUsecase interface {
	UpdateEmployee(ctx context.Context, employee entity.Employee) (*entity.Employee, error)
	ListEmployees(ctx context.Context, filter common.CommonFilter) ([]*entity.Employee, error)
	GetEmployeeByID(ctx context.Context, id uint) (*entity.Employee, error)
}