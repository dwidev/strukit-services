package services

import (
	"strukit-services/internal/models"
	"strukit-services/internal/repository"
	"strukit-services/pkg/logger"
	"strukit-services/pkg/token"
	"time"

	"github.com/google/uuid"
)

func NewAuth(token *token.Manager, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{tokenManager: token, UserRepository: userRepo}
}

type AuthService struct {
	tokenManager *token.Manager
	*repository.UserRepository
}

func (a *AuthService) LoginWithEmail(email string) (token *token.TokenResponse, err error) {
	now := time.Now()
	user := &models.User{BaseModel: models.BaseModel{ID: uuid.New()}}
	user, err = a.UserRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		newUser := &models.User{
			Email:       &email,
			LastLoginAt: &now,
		}
		user, err = a.UserRepository.CreateUser(newUser)
		if err != nil {
			return nil, err
		}
	}

	user.LastLoginAt = &now
	token, err = a.tokenManager.Generate(user)
	if err != nil {
		logger.Log.Errorf("error when generate user token, error :  %s", err)
		return nil, err
	}

	// TODO create login activiy in backgrou
	// go recordLogin

	return token, nil
}
