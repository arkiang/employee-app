package persistent

import (
	"context"
	"employee-app/internal/model"
	"employee-app/internal/model/entity"
	"employee-app/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

type payslipRepository struct {
	db *gorm.DB
}

func (r *payslipRepository) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

func (r *payslipRepository) CreatePayslipTx(ctx context.Context, tx *gorm.DB, payslip *entity.Payslip) (*entity.Payslip, error) {
	if err := tx.WithContext(ctx).Create(payslip).Error; err != nil {
		return nil, err
	}
	return payslip, nil
}

func (r *payslipRepository) GetByEmployeesAndPeriod(ctx context.Context, filter model.EmployeePeriodFilter) ([]*entity.Payslip, error) {
	var payslips []*entity.Payslip

	tx := r.db.WithContext(ctx).Model(&entity.Payslip{}).
		Preload("Employee.User").
		Preload("Attendances").
		Preload("Overtimes").
		Preload("Reimbursements")

	if filter.EmpIds != nil && len(*filter.EmpIds) > 0 {
		tx = tx.Where("employee_id IN ?", *filter.EmpIds)
	}
	if filter.Start != nil && filter.End != nil {
		tx = tx.Where("period_start >= ? AND period_end <= ?", *filter.Start, *filter.End)
	}

	tx = tx.Order(fmt.Sprintf("%s %s", filter.Base.GetSortByOrDefault("created_at"), filter.Base.GetSortBySQL()))
	
	if filter.Base.Limit != nil {
		tx = tx.Limit(*filter.Base.Limit).Offset(filter.Base.GetOffset())
	}

	if err := tx.Find(&payslips).Error; err != nil {
		return nil, err
	}

	return payslips, nil
}

func (r *payslipRepository) ListForPeriod(ctx context.Context, filter model.EmployeePeriodFilter) ([]*entity.Payslip, error) {
	var payslips []*entity.Payslip

	tx := r.db.WithContext(ctx).Model(&entity.Payslip{}).
		Preload("Employee.User").
		Preload("Attendances").
		Preload("Overtimes").
		Preload("Reimbursements")

	if filter.EmpIds != nil && len(*filter.EmpIds) > 0 {
		tx = tx.Where("employee_id IN ?", *filter.EmpIds)
	}
	if filter.Start != nil && filter.End != nil {
		tx = tx.Where("period_start >= ? AND period_end <= ?", *filter.Start, *filter.End)
	}

	tx = tx.Order(fmt.Sprintf("%s %s", filter.Base.GetSortByOrDefault("created_at"), filter.Base.GetSortBySQL()))
	
	if filter.Base.Limit != nil {
		tx = tx.Limit(*filter.Base.Limit).Offset(filter.Base.GetOffset())
	}

	if err := tx.Find(&payslips).Error; err != nil {
		return nil, err
	}

	return payslips, nil
}


// Payslip Attendance
func (r *payslipRepository) BulkInsertAttendancesTx(ctx context.Context, tx *gorm.DB, data []entity.PayslipAttendance) error {
	if len(data) == 0 {
		return nil
	}

	return tx.WithContext(ctx).Create(&data).Error
}

func (r *payslipRepository) GetAttendancesByPayslipID(ctx context.Context, payslipID uint) ([]entity.PayslipAttendance, error) {
	var records []entity.PayslipAttendance
	err := r.db.WithContext(ctx).
		Where("payslip_id = ?", payslipID).
		Find(&records).Error
	return records, err
}

// Payslip Overtime
func (r *payslipRepository) BulkInsertOvertimesTx(ctx context.Context, tx *gorm.DB, data []entity.PayslipOvertime) error {
	if len(data) == 0 {
		return nil
	}

	return tx.WithContext(ctx).Create(&data).Error
}

func (r *payslipRepository) GetOvertimesByPayslipID(ctx context.Context, payslipID uint) ([]entity.PayslipOvertime, error) {
	var records []entity.PayslipOvertime
	err := r.db.WithContext(ctx).
		Where("payslip_id = ?", payslipID).
		Find(&records).Error
	return records, err
}

// Payslip Reimbursement
func (r *payslipRepository) BulkInsertReimbursementsTx(ctx context.Context, tx *gorm.DB, data []entity.PayslipReimbursement) error {
	if len(data) == 0 {
		return nil
	}
	return tx.WithContext(ctx).Create(&data).Error
}

func (r *payslipRepository) GetReimbursementsByPayslipID(ctx context.Context, payslipID uint) ([]entity.PayslipReimbursement, error) {
	var records []entity.PayslipReimbursement
	err := r.db.WithContext(ctx).
		Where("payslip_id = ?", payslipID).
		Find(&records).Error
	return records, err
}

// New -.
func NewPayslip(db *gorm.DB) repository.PayslipRepository {
	if db == nil {
		return nil
	}
	return &payslipRepository{db: db}
}
