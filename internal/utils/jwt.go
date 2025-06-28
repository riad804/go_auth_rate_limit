package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/riad804/go_auth/internal/config"
)

type JWTWrapper struct {
	SecretKey       string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
}

type Claims struct {
	UserID string `json:"user_id"`
	OrgID  string `json:"org_id"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func NewJWTWrapper(cfg *config.Config) *JWTWrapper {
	return &JWTWrapper{
		SecretKey:       cfg.JWTSecret,
		AccessTokenExp:  cfg.AccessTokenExp,
		RefreshTokenExp: cfg.RefreshTokenExp,
	}
}

func (w *JWTWrapper) GenerateToken(userID, orgID string) (string, error) {
	claims := Claims{
		UserID: userID,
		OrgID:  orgID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(w.AccessTokenExp).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(w.SecretKey))
}

func (w *JWTWrapper) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(w.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (w *JWTWrapper) GenerateRefreshToken(userID string) (string, error) {
	claims := RefreshClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(w.RefreshTokenExp).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(w.SecretKey))
}

func (w *JWTWrapper) ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(w.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}
	return claims, nil
}
