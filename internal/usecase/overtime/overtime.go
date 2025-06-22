package overtime

import (
	"context"
	"time"
	"employee-app/internal/model/entity"
	"employee-app/internal/model"
)

type OvertimeUsecase interface {
	SubmitOvertime(ctx context.Context, empID uint, date time.Time, hours uint8) error
	GetOvertimeForPeriod(ctx context.Context, spec model.EmployeePeriodFilter) ([]*entity.Overtime, error)
}