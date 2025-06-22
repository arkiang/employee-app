package repository

import (
	"context"
	"employee-app/internal/model"
	"employee-app/internal/model/entity"
)

type AttendanceRepository interface {
	Create(ctx context.Context, spec entity.Attendance) (*entity.Attendance, error)
	Update(ctx context.Context, spec entity.Attendance) (*entity.Attendance, error)
	GetByID(ctx context.Context, id uint) (*entity.Attendance, error)
	List(ctx context.Context, filter model.EmployeePeriodFilter) ([]*entity.Attendance, error)
}