package persistent

import (
	"context"
	common "employee-app/internal/common/model"
	"employee-app/internal/model/entity"
	"employee-app/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

type payrollPeriodRepository struct {
	db *gorm.DB
}

// Create implements repository.PayrollPeriodRepository.
func (r *payrollPeriodRepository) Create(ctx context.Context, spec *entity.PayrollPeriod) (*entity.PayrollPeriod, error) {
	if err := r.db.WithContext(ctx).Create(&spec).Error; err != nil {
		return nil, err
	}
	return spec, nil
}

func (r *payrollPeriodRepository) MarkProcessedTx(ctx context.Context, tx *gorm.DB, id, userId uint) (*entity.PayrollPeriod, error) {
	err := tx.WithContext(ctx).
		Model(&entity.PayrollPeriod{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"processed" : true,
			"updated_by": userId,
		}).Error

	if err != nil {
		return nil, err
	}

	var updated entity.PayrollPeriod
	if err := r.db.WithContext(ctx).
		First(&updated, id).Error; err != nil {
		return nil, err
	}

	return &updated, nil
}

// GetByID implements repository.PayrollPeriodRepository.
func (r *payrollPeriodRepository) GetByID(ctx context.Context, id uint) (*entity.PayrollPeriod, error) {
	var payrollPeriod entity.PayrollPeriod
	if err := r.db.WithContext(ctx).First(&payrollPeriod, id).Error; err != nil {
		return nil, err
	}
	return &payrollPeriod, nil
}

// List implements repository.PayrollPeriodRepository.
func (r *payrollPeriodRepository) List(ctx context.Context, filter common.CommonFilter) ([]*entity.PayrollPeriod, error) {
	var payrollPeriods []*entity.PayrollPeriod

	tx := r.db.WithContext(ctx).Model(&entity.PayrollPeriod{}).
		Order(fmt.Sprintf("%s %s", filter.GetSortByOrDefault("created_at"), filter.GetSortBySQL()))

	if filter.Limit != nil {
		tx = tx.Limit(*filter.Limit).Offset(filter.GetOffset())
	}
		
	if err := tx.Find(&payrollPeriods).Error; err != nil {
		return nil, err
	}

	return payrollPeriods, nil
}

// New -.
func NewPayrollPeriod(db *gorm.DB) repository.PayrollPeriodRepository {
	if db == nil {
		return nil
	}
	return &payrollPeriodRepository{db: db}
}
