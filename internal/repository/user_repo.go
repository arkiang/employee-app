package repository

import (
	"context"
	"employee-app/internal/common/model"
	"employee-app/internal/model/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error

	CreateTx(ctx context.Context, tx *gorm.DB, spec entity.User) (*entity.User, error)
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	List(ctx context.Context, filter common.CommonFilter) ([]*entity.User, error)
}