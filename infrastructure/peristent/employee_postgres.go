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
func (r *employeeRepository) CreateTx(ctx context.Context, tx *gorm.DB, spec entity.Employee) error {
	return tx.WithContext(ctx).Create(spec).Error
}

// GetByID implements repository.EmployeeRepository.
func (r *employeeRepository) GetByID(ctx context.Context, id uint) (*entity.Employee, error) {
	var emp entity.Employee
	if err := r.db.WithContext(ctx).First(&emp, id).Error; err != nil {
		return nil, err
	}
	return &emp, nil
}

// List implements repository.EmployeeRepository.
func (r *employeeRepository) List(ctx context.Context, filter common.CommonFilter) ([]*entity.Employee, error) {
	var employees []*entity.Employee

	tx := r.db.WithContext(ctx).Model(&entity.Employee{}).
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
func (r *employeeRepository) Update(ctx context.Context, spec entity.Employee) (*entity.Employee, error) {
	if err := r.db.WithContext(ctx).Save(spec).Error; err != nil {
		return nil, err
	}
	return &spec, nil
}
