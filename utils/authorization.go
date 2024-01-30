package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type UserClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func NewToken(claims UserClaims) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return accessToken.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
}

func ParseToken(accessToken string) *jwt.Token {
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	return parsedAccessToken
}

func NewAccessToken(id uint) (string, error) {
	accessToken, err := NewToken(UserClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	})
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func NewRefreshToken(id uint) (string, error) {
	accessToken, err := NewToken(UserClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
	})
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func ValidateToken(token string) (*UserClaims, error) {
	if token == "" {
		return &UserClaims{}, errors.New("unauthorized - Missing access token")
	}

	parsedToken := ParseToken(token)

	if parsedToken == nil || !parsedToken.Valid {
		return &UserClaims{}, errors.New("unauthorized - Invalid access token")
	}

	claims, ok := parsedToken.Claims.(*UserClaims)
	if !ok {
		return &UserClaims{}, errors.New("unauthorized - Invalid token claims")
	}

	return claims, nil
}
