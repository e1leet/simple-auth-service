package service

import "github.com/golang-jwt/jwt/v4"

type AccessToken string

type RefreshToken struct {
	Token     string
	ExpiresIn int64
}

func NewRefreshToken(token string, expiresIn int64) RefreshToken {
	return RefreshToken{
		Token:     token,
		ExpiresIn: expiresIn,
	}
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}
