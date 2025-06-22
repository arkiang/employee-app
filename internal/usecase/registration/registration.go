package registration

import "context"
import "employee-app/internal/model/entity"

type RegistrationUsecase interface {
	RegisterEmployee(ctx context.Context, user entity.User, employee entity.Employee, password string) error
}