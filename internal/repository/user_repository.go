package repository

import (
	"context"
	"errors"
	"fmt"
	"strukit-services/internal/models"
	"strukit-services/pkg/logger"

	"gorm.io/gorm"
)

func NewUser(base *BaseRepository) *UserRepository {
	return &UserRepository{
		BaseRepository: base,
	}
}

type UserRepository struct {
	*BaseRepository
}

func (u *UserRepository) UpdatePasswordByUserID(ctx context.Context, password string) (err error) {
	uId := u.UserID(ctx)
	res := u.db.Model(&models.User{}).Where("id = ?", uId).Update("password_hash", password)
	if err = res.Error; err != nil {
		return fmt.Errorf("[UserRepository.CreatePassword] error create password, error : %w", err)
	}

	return nil
}

func (u *UserRepository) GetUserByEmail(email string) (user *models.User, err error) {
	user = new(models.User)
	res := u.db.First(user, "email = ?", email)
	if err = res.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		logger.Log.Errorf("error get userby email, error : %s", err)
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) CreateUser(user *models.User) (*models.User, error) {
	if err := u.db.Create(user).Error; err != nil {
		logger.Log.Errorf("error when create user, error : %s", err)
		return nil, err
	}

	return user, nil
}
