package overtime

import (
	"context"
	"employee-app/internal/api/dto/common"
	"employee-app/internal/common/constant"
	"employee-app/internal/model"
	"employee-app/internal/model/entity"
	"employee-app/internal/repository"
	"errors"
	"fmt"
	"time"
)

type overtimeUsecase struct {
	overtimeRepo   repository.OvertimeRepository
	attendanceRepo repository.AttendanceRepository
}

// Constructor
func New(overtimeRepo repository.OvertimeRepository, attendanceRepo repository.AttendanceRepository) OvertimeUsecase {
	return &overtimeUsecase{
		overtimeRepo:   overtimeRepo,
		attendanceRepo: attendanceRepo,
	}
}

func (u *overtimeUsecase) SubmitOvertime(ctx context.Context, empID uint, date common.DateOnly, hours uint8) error {
	val := ctx.Value(constant.UserId)
	userID, ok := val.(uint)
	if !ok {
		return errors.New("unauthorized or missing user ID")
	}

	if hours <= 0 || hours > 3 {
		return errors.New("overtime hours must be greater than 0 and cannot exceed 3 hours")
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return fmt.Errorf("failed to load timezone: %w", err)
	}
	
	tLocal := date.In(loc)
	overtimeDate := time.Date(
		tLocal.Year(), tLocal.Month(), tLocal.Day(),
		0, 0, 0, 0,
		loc,
	)

	// Reject weekends
	weekday := overtimeDate.Weekday()
	if weekday != time.Saturday && weekday != time.Sunday {
			//Check if the employee has already checked in on that day, assuming checking out is not required for overtime submission
		attendances, err := u.attendanceRepo.List(ctx, model.EmployeePeriodFilter{
			EmpIds: &[]string{fmt.Sprint(empID)},
			Start:  &overtimeDate,
			End:    &overtimeDate,
		})
		if err != nil {
			return fmt.Errorf("failed to check attendance: %w", err)
		}
		if len(attendances) == 0 || attendances[0].CheckInTime.IsZero() {
			return errors.New("cannot submit overtime when not checked in")
		}
	}

	// Check if overtime already exists for this day
	existing, err := u.overtimeRepo.List(ctx, model.EmployeePeriodFilter{
		EmpIds: &[]string{fmt.Sprint(empID)},
		Start:  &overtimeDate,
		End:    &overtimeDate,
	})
	if err != nil {
		return fmt.Errorf("failed to check existing overtime: %w", err)
	}
	if len(existing) > 0 {
		return errors.New("overtime already submitted for this day")
	}

	// ✅ Create and save overtime record
	overtime := entity.Overtime{
		EmployeeID:   empID,
		OvertimeDate: overtimeDate,
		Hours:        hours,
		CreatedBy:    userID,
		UpdatedBy:    userID,
	}

	_, err = u.overtimeRepo.Create(ctx, &overtime)
	return err
}

func (u *overtimeUsecase) GetOvertimeForPeriod(ctx context.Context, spec model.EmployeePeriodFilter) ([]*entity.Overtime, error) {
	overtimes, err := u.overtimeRepo.List(ctx, spec)
	if err != nil {
		return nil, err
	}
	return overtimes, nil
}