package attendance

import (
	"context"
	"employee-app/internal/model"
	"employee-app/internal/model/entity"
)

type AttendanceUsecase interface {
	SubmitAttendance(ctx context.Context, empID uint, time int64, attendanceType model.AttendanceType) error
	GetAttendanceForPeriod(ctx context.Context, spec model.EmployeePeriodFilter) ([]*entity.Attendance, error)
}