package repository

import (
	"context"
	"employee-app/internal/model"
	"employee-app/internal/model/entity"

	"gorm.io/gorm"
)

type PayslipRepository interface {
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error

	CreatePayslipTx(ctx context.Context, tx *gorm.DB, payslip *entity.Payslip) (*entity.Payslip, error)
	GetByEmployeesAndPeriod(ctx context.Context, filter model.EmployeePeriodFilter) ([]*entity.Payslip, error)
	ListForPeriod(ctx context.Context, filter model.EmployeePeriodFilter) ([]*entity.Payslip, error)

	// PayslipAttendance operations
	BulkInsertAttendancesTx(ctx context.Context, tx *gorm.DB, data []entity.PayslipAttendance) error
	GetAttendancesByPayslipID(ctx context.Context, payslipID uint) ([]entity.PayslipAttendance, error)

	// PayslipOvertime operations
	BulkInsertOvertimesTx(ctx context.Context, tx *gorm.DB, data []entity.PayslipOvertime) error
	GetOvertimesByPayslipID(ctx context.Context, payslipID uint) ([]entity.PayslipOvertime, error)

	// PayslipReimbursement operations
	BulkInsertReimbursementsTx(ctx context.Context, tx *gorm.DB, data []entity.PayslipReimbursement) error
	GetReimbursementsByPayslipID(ctx context.Context, payslipID uint) ([]entity.PayslipReimbursement, error)
}