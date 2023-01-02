package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/e1leet/simple-auth-service/internal/config"
	"github.com/e1leet/simple-auth-service/internal/domain/jwt/service"
	"github.com/golang-jwt/jwt/v4"
)

func SetTokens(w http.ResponseWriter, access service.AccessToken, refresh service.RefreshToken) {
	http.SetCookie(w, &http.Cookie{
		Name:     config.RefreshCookieName,
		Value:    refresh.Token,
		Path:     "/api/auth",
		MaxAge:   int(refresh.ExpiresIn),
		HttpOnly: true,
	})
	w.Header().Add("Authorization", fmt.Sprintf("Bearer %s", string(access)))
}

func ParseAccessTokenFromHeader(header string, jwtSecret string) (*service.Claims, error) {
	data := strings.Split(header, " ")
	if len(data) != 2 || data[0] != "Bearer" {
		return nil, fmt.Errorf("incorrect token")
	}

	claims := &service.Claims{}

	_, err := jwt.ParseWithClaims(data[1], claims, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !(ok && method == jwt.SigningMethodHS256) {
			return nil, fmt.Errorf("incorrect token")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return nil, err
	}

	return claims, nil
}
