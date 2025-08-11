package services

import (
	"strukit-services/internal/models"
	"strukit-services/pkg/logger"
	"strukit-services/pkg/token"

	"github.com/google/uuid"
)

func NewAuth(token *token.Token) *AuthService {
	return &AuthService{Token: token}
}

type AuthService struct {
	*token.Token
}

func (a *AuthService) LoginWithEmail(email string) (*token.TokenResponse, error) {
	// TODO: create check user already regist or not
	// 1. if already, generate token
	// 2. if not already, register user & generate token

	user := &models.UserModel{BaseModel: models.BaseModel{ID: uuid.New()}}
	token, err := a.Token.Generate(user)
	if err != nil {
		logger.Log.Errorf("error when generate user token, %s", err)
		return nil, err
	}

	return token, nil
}
