package repository

import (
	"context"
	common "employee-app/internal/common/model"
	"employee-app/internal/model/entity"

	"gorm.io/gorm"
)

type EmployeeRepository interface {
	CreateTx(ctx context.Context, tx *gorm.DB, spec *entity.Employee) error
	Update(ctx context.Context, spec *entity.Employee) (*entity.Employee, error)
	GetByID(ctx context.Context, id uint) (*entity.Employee, error)
	List(ctx context.Context, filter common.CommonFilter) ([]*entity.Employee, error)
}