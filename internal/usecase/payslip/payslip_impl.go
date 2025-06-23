package payslip

import (
	"context"
	"employee-app/internal/api/middleware"
	common "employee-app/internal/common/model"
	"employee-app/internal/model"
	"employee-app/internal/model/entity"
	"employee-app/internal/repository"
	"employee-app/internal/usecase/utils"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type payrollUsecase struct {
	employeeRepo       repository.EmployeeRepository
	attendanceRepo     repository.AttendanceRepository
	overtimeRepo       repository.OvertimeRepository
	reimbursementRepo  repository.ReimbursementRepository
	payslipRepo        repository.PayslipRepository
	payrollPeriodRepo  repository.PayrollPeriodRepository
}

func New(
	employeeRepo repository.EmployeeRepository,
	attendanceRepo repository.AttendanceRepository,
	overtimeRepo repository.OvertimeRepository,
	reimbursementRepo repository.ReimbursementRepository,
	payrollPeriodRepo repository.PayrollPeriodRepository,
	payslipRepo repository.PayslipRepository,
) PayslipUsecase {
	return &payrollUsecase{
		employeeRepo:      employeeRepo,
		attendanceRepo:    attendanceRepo,
		overtimeRepo:      overtimeRepo,
		reimbursementRepo: reimbursementRepo,
		payrollPeriodRepo: payrollPeriodRepo,
		payslipRepo:       payslipRepo,
	}
}

func (u *payrollUsecase) RunPayroll(ctx context.Context, periodID uint) error {
	val := ctx.Value(middleware.ContextUserIDKey)
	adminID, ok := val.(uint)
	if !ok {
		return errors.New("unauthorized or missing user ID")
	}

	// 1. Fetch payroll period
	period, err := u.payrollPeriodRepo.GetByID(ctx, periodID)
	if err != nil {
		return fmt.Errorf("failed to get period: %w", err)
	}

	if period.Processed {
		return fmt.Errorf("payroll has already been processed for this period")
	}

	// 3. Fetch employees
	employees, err := u.employeeRepo.List(ctx, common.CommonFilter{})
	if err != nil {
		return fmt.Errorf("failed to list employees: %w", err)
	}

	for _, emp := range employees {
		empIDStr := fmt.Sprint(emp.ID)

		filter := model.EmployeePeriodFilter{
			EmpIds: &[]string{empIDStr},
			Start:  &period.StartDate,
			End:    &period.EndDate,
		}

		// Get related records
		attendances, _ := u.attendanceRepo.List(ctx, filter)
		overtimes, _ := u.overtimeRepo.List(ctx, filter)
		reimbursements, _ := u.reimbursementRepo.List(ctx, filter)

		// Salary breakdown
		workingDays := 20
		dailySalary := emp.Salary.Div(decimal.NewFromInt32(int32(workingDays)))
		hourlySalary := dailySalary.Div(decimal.NewFromInt(8))
		presentDays := utils.CountUniqueDates(attendances)
		attendanceTotal := dailySalary.Mul(decimal.NewFromInt32(int32(presentDays)))
		overtimeHours := utils.SumOvertimeHours(overtimes)
		overtimeTotal := hourlySalary.Mul(decimal.NewFromInt32(int32(overtimeHours))).Mul(decimal.NewFromInt(2))
		reimburseTotal := utils.SumReimbursementAmounts(reimbursements)
		total := attendanceTotal.Add(overtimeTotal).Add(reimburseTotal)

		// Prepare records
		payslipEnt := entity.Payslip{
			EmployeeID:      emp.ID,
			PeriodStart:     period.StartDate,
			PeriodEnd:       period.EndDate,
			BaseSalary: 	 emp.Salary,
			AttendanceDays:  presentDays,
			AttendancePay:   attendanceTotal,
			OvertimeHours:   overtimeHours,
			OvertimePay:     overtimeTotal,
			ReimbursementTotal:  reimburseTotal,
			TakeHomePay:     total,
			CreatedBy:       adminID,
			UpdatedBy:       adminID,
		}

		// ⚙️ Wrap all ops per employee in a DB transaction
		err := u.payslipRepo.WithTransaction(ctx, func(tx *gorm.DB) error {
			payslip, err := u.payslipRepo.CreatePayslipTx(ctx, tx, &payslipEnt);
			if err != nil {
				return err
			}

			_, err = u.payrollPeriodRepo.MarkProcessedTx(ctx, tx, period.ID, adminID)
			if err != nil {
				return err
			}	

			var payslipAttendances []entity.PayslipAttendance
			for _, a := range attendances {
				payslipAttendances = append(payslipAttendances, entity.PayslipAttendance{
					PayslipID:  payslip.ID,
					Date: 		a.AttendanceDate,
					CheckIn: 	a.CheckInTime,
					CheckOut: 	a.CheckOutTime,
				})
			}

			var payslipOvertimes []entity.PayslipOvertime
			for _, o := range overtimes {
				payslipOvertimes = append(payslipOvertimes, entity.PayslipOvertime{
					PayslipID:  payslip.ID,
					Date: 		o.OvertimeDate,
					Hours: 		o.Hours,
				})
			}

			var payslipReimbursements []entity.PayslipReimbursement
			for _, r := range reimbursements {
				payslipReimbursements = append(payslipReimbursements, entity.PayslipReimbursement{
					PayslipID:  	payslip.ID,
					Date: 			r.ReimbursementDate,
					Amount:			r.Amount,
					Description:    r.Description,
				})
			}
			
			if len(payslipAttendances) > 0 {
				if err := u.payslipRepo.BulkInsertAttendancesTx(ctx, tx, payslipAttendances); err != nil {
					return err
				}
			}
			if len(payslipOvertimes) > 0 {
				if err := u.payslipRepo.BulkInsertOvertimesTx(ctx, tx, payslipOvertimes); err != nil {
					return err
				}
			}
			if len(payslipReimbursements) > 0 {
				if err := u.payslipRepo.BulkInsertReimbursementsTx(ctx, tx, payslipReimbursements); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to generate payslip for employee %d: %w", emp.ID, err)
		}
	}

	return nil
}

func (u *payrollUsecase) GetPayslipForEmployee(ctx context.Context, periodID uint) (*entity.Payslip, error) {
	val := ctx.Value(middleware.ContextUserIDKey)
	userID, ok := val.(uint)
	if !ok {
		return nil, errors.New("unauthorized or missing user ID")
	}

	employee, err := u.employeeRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee with userID %d: %w", userID, err)
	}

	filter := model.EmployeePeriodFilter{
		EmpIds: &[]string{fmt.Sprint(employee.ID)},
	}
	period, err := u.payrollPeriodRepo.GetByID(ctx, periodID)
	if err != nil {
		return nil, err
	}
	filter.Start = &period.StartDate
	filter.End = &period.EndDate

	payslips, err := u.payslipRepo.GetByEmployeesAndPeriod(ctx, filter)
	if err != nil || len(payslips) == 0 {
		return nil, err
	}
	return payslips[0], nil
}

func (u *payrollUsecase) GetPayslipSummary(ctx context.Context, periodID uint, filter common.CommonFilter) ([]*entity.Payslip, error) {
	period, err := u.payrollPeriodRepo.GetByID(ctx, periodID)
	if err != nil {
		return nil, err
	}
	return u.payslipRepo.ListForPeriod(ctx, model.EmployeePeriodFilter{
		Base:  filter,
		Start: &period.StartDate,
		End:   &period.EndDate,
	})
}