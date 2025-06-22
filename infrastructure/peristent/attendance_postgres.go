package persistent

import (
	"context"
	"employee-app/internal/model"
	"employee-app/internal/model/entity"
	"employee-app/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

type attendanceRepository struct {
	db *gorm.DB
}

// New -.
func NewAttendance(db *gorm.DB) repository.AttendanceRepository {
	if db == nil {
		return nil
	}
	return &attendanceRepository{db: db}
}

// Create implements repository.AttendanceRepository.
func (r *attendanceRepository) Create(ctx context.Context, spec entity.Attendance) (*entity.Attendance, error) {
	if err := r.db.WithContext(ctx).Create(&spec).Error; err != nil {
		return nil, err
	}
	return &spec, nil
}

// GetByID implements repository.AttendanceRepository.
func (r *attendanceRepository) GetByID(ctx context.Context, id uint) (*entity.Attendance, error) {
	var attendance entity.Attendance
	if err := r.db.WithContext(ctx).First(&attendance, id).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}

// List implements repository.AttendanceRepository.
func (r *attendanceRepository) List(ctx context.Context, filter model.EmployeePeriodFilter) ([]*entity.Attendance, error) {
	var attendances []*entity.Attendance

	tx := r.db.WithContext(ctx).Model(&entity.Attendance{})

	if filter.EmpIds != nil && len(*filter.EmpIds) > 0 {
		tx = tx.Where("employee_id IN ?", *filter.EmpIds)
	}

	// Filter by period (Start and/or End)
	if filter.Start != nil {
		tx = tx.Where("attendance_date >= ?", filter.Start)
	}
	if filter.End != nil {
		tx = tx.Where("attendance_date <= ?", filter.End)
	}
	
	tx = tx.Order(fmt.Sprintf("%s %s", filter.Base.GetSortByOrDefault("created_at"), filter.Base.GetSortBySQL()))
	
	if filter.Base.Limit != nil {
		tx = tx.Limit(*filter.Base.Limit).Offset(filter.Base.GetOffset())
	}

	if err := tx.Find(&attendances).Error; err != nil {
		return nil, err
	}

	return attendances, nil
}

// Update implements repository.AttendanceRepository.
func (r *attendanceRepository) Update(ctx context.Context, spec entity.Attendance) (*entity.Attendance, error) {
	if err := r.db.WithContext(ctx).Save(&spec).Error; err != nil {
		return nil, err
	}
	return &spec, nil
}
