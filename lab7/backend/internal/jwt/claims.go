package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType = int

const (
	AccessToken TokenType = iota
	RefreshToken
)

type Claims[T any] struct {
	jwt.RegisteredClaims
	Data T `json:"data"`
}

func CreateClaims[T any](config *JwtConfig, tokenType TokenType, data T, id string) (*Claims[T], error) {
	var expirationString string

	switch tokenType {
	case AccessToken:
		expirationString = config.ExpireTimeoutAccess
	case RefreshToken:
		expirationString = config.ExpireTimeoutRefresh
	}

	expiration, parseError := time.ParseDuration(expirationString)
	if parseError != nil {
		return nil, parseError
	}

	claims := Claims[T]{
		Data: data,
	}

	now := time.Now()

	claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    config.Issuer,
		ExpiresAt: &jwt.NumericDate{Time: now.Add(expiration)},
		IssuedAt:  &jwt.NumericDate{Time: now},
		ID:        id,
	}

	return &claims, nil
}
