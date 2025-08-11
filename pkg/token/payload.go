package token

import "github.com/golang-jwt/jwt/v5"

type TokenClaims struct {
	UserID string
	*jwt.RegisteredClaims
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
