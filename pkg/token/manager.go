package token

import (
	"strukit-services/internal/models"
	"strukit-services/pkg/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewManager(SecretKeys ...string) *Manager {
	if len(SecretKeys) < 2 {
		logger.Log.Fatalf("SecretKeys cannot parse")
		return nil
	}

	return &Manager{
		SecretKeys: SecretKeys,
	}
}

type Manager struct {
	SecretKeys []string
}

func (t *Manager) claims(user *models.User, expired time.Duration) *TokenClaims {
	return &TokenClaims{
		UserID: user.ID,
		RegisteredClaims: &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expired)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func (t *Manager) Generate(user *models.User) (*TokenResponse, error) {
	accessClaims := t.claims(user, accessExpTime)
	Manager := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := Manager.SignedString([]byte(t.SecretKeys[0]))
	if err != nil {
		return nil, err
	}

	refreshClaims := t.claims(user, refreshExpTime)
	tokenRefresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := tokenRefresh.SignedString([]byte(t.SecretKeys[1]))
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (t *Manager) Parse(tokenRaw string, secretKey string) (*TokenClaims, error) {
	claims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenRaw, claims, func(t *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, err
}
