// internal/infrastructure/sqlite/user_repo.go
package sqlite

import (
	"auth/internal/domain/entity"
	"auth/internal/domain/interfaces"
	"context"
	"errors"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserSQLiteRepository(db *gorm.DB) interfaces.UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) Save(ctx context.Context, u *entity.User) error {
	user, err := entity.NewUser(u.Email, u.PasswordHash, u.Role)
	if err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return interfaces.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (r *UserRepo) GetUser(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, interfaces.ErrUserNotFound
		}
		return nil, err
	}
	return &entity.User{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
	}, nil
}
