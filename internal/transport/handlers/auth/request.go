package auth

import (
	"net/http"

	"github.com/e1leet/simple-auth-service/internal/domain/user"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (req LoginRequest) Bind(r *http.Request) error {
	return nil
}

func (req LoginRequest) ToDomain() *user.User {
	return &user.User{
		Username: req.Username,
		Password: req.Password,
	}
}

type RegisterRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"passwordRepeat"`
}

func (req RegisterRequest) Bind(r *http.Request) error {
	return nil
}

func (req RegisterRequest) ToDomain() *user.User {
	return &user.User{
		Username: req.Username,
		Password: req.Password,
	}
}
