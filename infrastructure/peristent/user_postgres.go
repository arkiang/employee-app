package persistent

import (
	"context"
	"fmt"

	common "employee-app/internal/common/model"
	"employee-app/internal/model/entity"
	"employee-app/internal/repository"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// New -.
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	if db == nil {
		return nil
	}
	return &userRepository{db: db}
}

// WithTransaction implements repository.UserRepository.
func (r *userRepository) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

func (r *userRepository) CreateTx(ctx context.Context, tx *gorm.DB, user *entity.User) (*entity.User, error) {
	if err := tx.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) List(ctx context.Context, filter common.CommonFilter) ([]*entity.User, error) {
	var users []*entity.User

	tx := r.db.WithContext(ctx).Model(&entity.User{}).
		Order(fmt.Sprintf("%s %s", filter.GetSortByOrDefault("created_at"), filter.GetSortBySQL()))

	if filter.Limit != nil {
		tx = tx.Limit(*filter.Limit).Offset(filter.GetOffset())
	}

	if err := tx.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}
