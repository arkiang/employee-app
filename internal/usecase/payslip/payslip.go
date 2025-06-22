package payslip

import (
	"context"
	common "employee-app/internal/common/model"
	"employee-app/internal/model/entity"
)

type PayslipUsecase interface {
	RunPayroll(ctx context.Context, periodID uint) error
	GetPayslipForEmployee(ctx context.Context, empID uint, periodID uint) (*entity.Payslip, error)
	GetPayslipSummary(ctx context.Context, periodID uint, filter common.CommonFilter) ([]*entity.Payslip, error)
}