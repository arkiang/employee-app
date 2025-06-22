package persistent

import (
	"context"
	common "employee-app/internal/common/model"
	"employee-app/internal/model/entity"
	"employee-app/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

type employeeRepository struct {
	db *gorm.DB
}

// New -.
func NewEmployee(db *gorm.DB) repository.EmployeeRepository {
	if db == nil {
		return nil
	}
	return &employeeRepository{db: db}
}

// Create implements repository.EmployeeRepository.
func (r *employeeRepository) CreateTx(ctx context.Context, tx *gorm.DB, spec *entity.Employee) error {
	return tx.WithContext(ctx).Create(spec).Error
}

// GetByID implements repository.EmployeeRepository.
func (r *employeeRepository) GetByID(ctx context.Context, id uint) (*entity.Employee, error) {
	var emp entity.Employee
	if err := r.db.WithContext(ctx).Preload("User").First(&emp, id).Error; err != nil {
		return nil, err
	}
	return &emp, nil
}

// GetByUserID implements repository.EmployeeRepository.
func (r *employeeRepository) GetByUserID(ctx context.Context, userId uint) (*entity.Employee, error) {
	var employee entity.Employee
	if err := r.db.WithContext(ctx).
		Preload("User").
		Where("user_id = ?", userId).
		First(&employee).Error; err != nil {
		return nil, err
	}
	return &employee, nil
}

// List implements repository.EmployeeRepository.
func (r *employeeRepository) List(ctx context.Context, filter common.CommonFilter) ([]*entity.Employee, error) {
	var employees []*entity.Employee

	tx := r.db.WithContext(ctx).Model(&entity.Employee{}).
		Preload("User").
		Order(fmt.Sprintf("%s %s", filter.GetSortByOrDefault("created_at"), filter.GetSortBySQL()))

	if filter.Limit != nil {
		tx.Limit(*filter.Limit).Offset(filter.GetOffset())
	}

	if err := tx.Find(&employees).Error; err != nil {
		return nil, err
	}

	return employees, nil
}

// Update implements repository.EmployeeRepository.
func (r *employeeRepository) Update(ctx context.Context, spec *entity.Employee) (*entity.Employee, error) {
	err := r.db.WithContext(ctx).
		Model(&entity.Employee{}).
		Where("id = ?", spec.ID).
		Updates(map[string]interface{}{
			"name":   spec.Name,
			"salary": spec.Salary,
		}).Error

	if err != nil {
		return nil, err
	}

	var updated entity.Employee
	if err := r.db.WithContext(ctx).
		Preload("User").
		First(&updated, spec.ID).Error; err != nil {
		return nil, err
	}

	return &updated, nil
}
