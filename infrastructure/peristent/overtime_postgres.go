package persistent

import (
	"context"
	"employee-app/internal/model"
	"employee-app/internal/model/entity"
	"employee-app/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

type overtimeRepository struct {
	db *gorm.DB
}

// New -.
func NewOvertime(db *gorm.DB) repository.OvertimeRepository {
	if db == nil {
		return nil
	}
	return &overtimeRepository{db: db}
}

// Create implements repository.OvertimeRepository.
func (r *overtimeRepository) Create(ctx context.Context, spec entity.Overtime) (*entity.Overtime, error) {
	if err := r.db.WithContext(ctx).Create(&spec).Error; err != nil {
		return nil, err
	}
	return &spec, nil
}

// GetByID implements repository.OvertimeRepository.
func (r *overtimeRepository) GetByID(ctx context.Context, id uint) (*entity.Overtime, error) {
	var overtime entity.Overtime
	if err := r.db.WithContext(ctx).First(&overtime, id).Error; err != nil {
		return nil, err
	}
	return &overtime, nil
}

// List implements repository.OvertimeRepository.
func (r *overtimeRepository) List(ctx context.Context, filter model.EmployeePeriodFilter) ([]*entity.Overtime, error) {
	var overtimes []*entity.Overtime

	tx := r.db.WithContext(ctx).Model(&entity.Overtime{})

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

	if err := tx.Find(&overtimes).Error; err != nil {
		return nil, err
	}

	return overtimes, nil
}
