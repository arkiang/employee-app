package attendance

import (
	"context"
	"employee-app/internal/common/constant"
	common "employee-app/internal/common/model"
	"employee-app/internal/model"
	"employee-app/internal/model/entity"
	"employee-app/internal/repository"
	"errors"
	"fmt"
	"time"
)

type attendanceUsecase struct {
	attendanceRepo repository.AttendanceRepository
}

func New(attendanceRepo repository.AttendanceRepository) AttendanceUsecase {
	return &attendanceUsecase{
		attendanceRepo: attendanceRepo,
	}
}

func (u *attendanceUsecase) GetAttendanceForPeriod(ctx context.Context, spec model.EmployeePeriodFilter) ([]*entity.Attendance, error) {
	attendances, err := u.attendanceRepo.List(ctx, spec)
	if err != nil {
		return nil, err
	}
	return attendances, nil
}

func (u *attendanceUsecase) SubmitAttendance(ctx context.Context, empID uint, t int64, attendanceType model.AttendanceType) error {
	val := ctx.Value(constant.UserId)
	userID, ok := val.(uint)
	if !ok {
		return errors.New("unauthorized or missing user ID")
	}

	// Load Asia/Jakarta timezone
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return errors.New("failed to load timezone Asia/Jakarta")
	}

	// Convert time to Asia/Jakarta
	tLocal := time.Unix(t, 0).In(loc)
	// Truncate to start of day (local time)
	attendanceDate := time.Date(
		tLocal.Year(), tLocal.Month(), tLocal.Day(),
		0, 0, 0, 0,
		loc,
	)

	// Reject weekends
	weekday := attendanceDate.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return errors.New("cannot submit attendance on weekends ")
	}

	// Check for duplicates (same day & type)
	filter := model.EmployeePeriodFilter{
		Base:   common.CommonFilter{},
		EmpIds: &[]string{fmt.Sprint(empID)},
		Start:  &attendanceDate,
		End:    &attendanceDate,
	}

	attendances, err := u.attendanceRepo.List(ctx, filter)
	if err != nil {
		return err
	}

	if len(attendances) > 0 {
		// Attendance record already exists
		existing := attendances[0]

		if attendanceType == model.CHECK_IN && !existing.CheckInTime.IsZero() {
			return errors.New("check-in already submitted for today")
		}
		if attendanceType == model.CHECK_OUT && existing.CheckOutTime != nil {
			return errors.New("check-out already submitted for today")
		}

		// Update record with new check-in or check-out
		if attendanceType == model.CHECK_OUT {
			existing.CheckOutTime = &tLocal
		} 
		existing.UpdatedBy = userID

		_, err = u.attendanceRepo.Update(ctx, existing)
		return err
	}

	// Create new record if not exist and it's a check-in
	if attendanceType == model.CHECK_OUT {
		return errors.New("cannot check out without checking in first")
	}

	// Save attendance (store time in UTC for consistency)
	newAttendance := entity.Attendance{
		EmployeeID:     empID,
		AttendanceDate: attendanceDate,
		CheckInTime:    tLocal,
		CreatedBy:      userID,
		UpdatedBy:      userID,
	}
	_, err = u.attendanceRepo.Create(ctx, &newAttendance)
	return err
}
