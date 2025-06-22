package attendance

import (
	"context"
	"time"
	"employee-app/internal/model/entity"
	"employee-app/internal/model"
)

type AttendanceUsecase interface {
	SubmitAttendance(ctx context.Context, empID uint, time time.Time, attendanceType model.AttendanceType) error
	GetAttendanceForPeriod(ctx context.Context, spec model.EmployeePeriodFilter) ([]*entity.Attendance, error)
}