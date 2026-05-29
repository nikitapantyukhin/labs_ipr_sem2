package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type IJwtHandler[T any] interface {
	GenerateJwtPair(data T, id string) (accessToken string, refreshToken string, err error)
	RefreshAccessToken(data T, id string) (accessToken string, err error)
	Validate(token string, tokenType TokenType) (*Claims[T], error)
}

type Handler[T any] struct {
	cfg *JwtConfig
}

func (h Handler[T]) Validate(token string, tokenType TokenType) (*Claims[T], error) {
	parsedToken, parseError := jwt.ParseWithClaims(
		token,
		&Claims[T]{},
		func(token *jwt.Token) (any, error) {
			switch tokenType {
			case AccessToken:
				return []byte(h.cfg.AccessSecret), nil
			case RefreshToken:
				return []byte(h.cfg.RefreshSecret), nil
			}

			return nil, errors.New("invalid claims type")
		},
	)

	if parseError != nil {
		return nil, parseError
	}

	claims, ok := parsedToken.Claims.(*Claims[T])

	if !ok {
		return nil, errors.New("unknown claims type")
	}

	validator := jwt.NewValidator(
		jwt.WithIssuer(h.cfg.Issuer),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)

	if err := validator.Validate(claims); err != nil {
		return nil, err
	}

	return claims, nil
}

func (h Handler[T]) generateSingleToken(claims *Claims[T], tokenType TokenType) (string, error) {
	var signingKey []byte

	switch tokenType {
	case AccessToken:
		signingKey = []byte(h.cfg.AccessSecret)
	case RefreshToken:
		signingKey = []byte(h.cfg.RefreshSecret)
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(signingKey)

}

func (h Handler[T]) GenerateJwtPair(data T, id string) (accessToken string, refreshToken string, err error) {
	claimsAccess, accessTokenClaimsGenerationError := CreateClaims(h.cfg, AccessToken, data, id)
	if accessTokenClaimsGenerationError != nil {
		return "", "", accessTokenClaimsGenerationError
	}

	claimsRefresh, refreshTokenClaimsGenerationError := CreateClaims(h.cfg, RefreshToken, data, id)
	if refreshTokenClaimsGenerationError != nil {
		return "", "", refreshTokenClaimsGenerationError
	}

	accessToken, accessTokenGenerationError := h.generateSingleToken(claimsAccess, AccessToken)
	if accessTokenGenerationError != nil {
		return "", "", accessTokenGenerationError
	}

	refreshToken, refreshTokenGenerationError := h.generateSingleToken(claimsRefresh, RefreshToken)
	if refreshTokenGenerationError != nil {
		return "", "", refreshTokenGenerationError
	}

	return
}

func (h Handler[T]) RefreshAccessToken(data T, id string) (accessToken string, err error) {
	claims, refreshTokenClaimsGenerationError := CreateClaims(h.cfg, AccessToken, data, id)
	if refreshTokenClaimsGenerationError != nil {
		return "", refreshTokenClaimsGenerationError
	}

	accessToken, refreshTokenGenerationError := h.generateSingleToken(claims, AccessToken)
	if refreshTokenGenerationError != nil {
		return "", refreshTokenGenerationError
	}

	return
}

func CreateHandler[T any](cfg *JwtConfig) IJwtHandler[T] {
	return &Handler[T]{
		cfg: cfg,
	}
}
