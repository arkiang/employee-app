package employee

import (
	"context"
	common "employee-app/internal/common/model"
	"employee-app/internal/model/entity"
	"employee-app/internal/repository"
)

type employeeUsecase struct {
	employeeRepo repository.EmployeeRepository
}

func New(employeeRepo repository.EmployeeRepository) EmployeeUsecase {
	return &employeeUsecase{
		employeeRepo: employeeRepo,
	}
}

func (uc *employeeUsecase) UpdateEmployee(ctx context.Context, employee entity.Employee) (*entity.Employee, error) {
	updatedEmployee, err := uc.employeeRepo.Update(ctx, &employee)
	if err != nil {
		return nil, err
	}
	
	return updatedEmployee, nil
}

func (uc *employeeUsecase) ListEmployees(ctx context.Context, filter common.CommonFilter) ([]*entity.Employee, error) {
	return uc.employeeRepo.List(ctx, filter)
}

func (uc *employeeUsecase) GetEmployeeByID(ctx context.Context, id uint) (*entity.Employee, error) {
	return uc.employeeRepo.GetByID(ctx, id)
}