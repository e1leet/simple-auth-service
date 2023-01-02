package auth

import (
	"net/http"

	"github.com/e1leet/simple-auth-service/internal/domain/jwt/service"
	"github.com/go-chi/render"
)

type TokenResponse struct {
	Access  string `json:"accessToken"`
	Refresh string `json:"refreshToken"`
}

func NewTokenResponse(access service.AccessToken, refresh service.RefreshToken) TokenResponse {
	return TokenResponse{
		Access:  string(access),
		Refresh: refresh.Token,
	}
}

func (res TokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}
