package persistent

import (
	"context"
	"employee-app/internal/model"
	"employee-app/internal/model/entity"
	"employee-app/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

type reimbursementRepository struct {
	db *gorm.DB
}

// Create implements repository.ReimbursementRepository.
func (r *reimbursementRepository) Create(ctx context.Context, spec *entity.Reimbursement) (*entity.Reimbursement, error) {
	if err := r.db.WithContext(ctx).Create(&spec).Error; err != nil {
		return nil, err
	}
	return spec, nil
}

// GetByID implements repository.ReimbursementRepository.
func (r *reimbursementRepository) GetByID(ctx context.Context, id uint) (*entity.Reimbursement, error) {
	var reimbursement entity.Reimbursement
	if err := r.db.WithContext(ctx).First(&reimbursement, id).Error; err != nil {
		return nil, err
	}
	return &reimbursement, nil
}

// List implements repository.ReimbursementRepository.
func (r *reimbursementRepository) List(ctx context.Context, filter model.EmployeePeriodFilter) ([]*entity.Reimbursement, error) {
	var reimbursements []*entity.Reimbursement

	tx := r.db.WithContext(ctx).Model(&entity.Reimbursement{})

	if filter.EmpIds != nil && len(*filter.EmpIds) > 0 {
		tx = tx.Where("employee_id IN ?", *filter.EmpIds)
	}

	// Filter by period (Start and/or End)
	if filter.Start != nil {
		tx = tx.Where("reimbursement_date >= ?", filter.Start)
	}
	if filter.End != nil {
		tx = tx.Where("reimbursement_date <= ?", filter.End)
	}
	
	tx = tx.Order(fmt.Sprintf("%s %s", filter.Base.GetSortByOrDefault("created_at"), filter.Base.GetSortBySQL()))
	
	if filter.Base.Limit != nil {
		tx = tx.Limit(*filter.Base.Limit).Offset(filter.Base.GetOffset())
	}

	if err := tx.Find(&reimbursements).Error; err != nil {
		return nil, err
	}

	return reimbursements, nil
}

// New -.
func NewReimbursement(db *gorm.DB) repository.ReimbursementRepository {
	if db == nil {
		return nil
	}
	return &reimbursementRepository{db: db}
}
