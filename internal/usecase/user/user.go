package user

import (
	"context"
	"employee-app/internal/model/entity"
)

type UserUsecase interface {
	Login(ctx context.Context, username, password string) (*entity.User, error)
	GetByID(ctx context.Context, id uint) (*entity.User, error)
}