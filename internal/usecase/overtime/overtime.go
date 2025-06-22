package overtime

import (
	"context"
	"employee-app/internal/api/dto/common"
	"employee-app/internal/model"
	"employee-app/internal/model/entity"
)

type OvertimeUsecase interface {
	SubmitOvertime(ctx context.Context, empID uint, date common.DateOnly, hours uint8) error
	GetOvertimeForPeriod(ctx context.Context, spec model.EmployeePeriodFilter) ([]*entity.Overtime, error)
}