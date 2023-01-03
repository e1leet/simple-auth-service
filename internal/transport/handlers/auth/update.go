package auth

import (
	"errors"
	"net/http"

	"github.com/e1leet/simple-auth-service/internal/domain/jwt/service"
	"github.com/e1leet/simple-auth-service/internal/domain/session/dao"
	"github.com/e1leet/simple-auth-service/internal/utils"
	"github.com/e1leet/simple-auth-service/internal/utils/api"
	"github.com/go-chi/render"
)

// RefreshToken godoc
//
//	@Summary		Refresh token
//	@Description	Recreate access and refresh token
//	@Tags			auth
//	@Success		200	{object}	TokenResponse
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/auth/refresh-tokens [post]
//	@Security		ApiKeyAuth
func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	claims, err := utils.ParseAccessTokenFromHeader(authHeader, h.cfg.Security.JWTSecret)
	if err != nil {
		_ = render.Render(w, r, api.NewErrorResponse(err.Error(), http.StatusUnauthorized))
		return
	}

	refreshCookie, err := r.Cookie("refreshToken")
	if err != nil {
		h.logger.Error().Err(err).Send()
		_ = render.Render(w, r, api.NewErrorResponse(err.Error(), http.StatusUnauthorized))

		return
	}

	access, refresh, err := h.authService.UpdateAccessToken(r.Context(), claims.UserID, refreshCookie.Value)
	if err != nil {
		h.logger.Error().Err(err).Send()

		var response render.Renderer

		switch {
		case errors.Is(err, dao.ErrSessionNotFound), errors.Is(err, service.ErrRefreshTokenExpired):
			response = api.NewErrorResponse(err.Error(), http.StatusUnauthorized)
		default:
			response = api.NewErrorResponse(err.Error(), http.StatusUnauthorized)
		}

		_ = render.Render(w, r, response)

		return
	}

	utils.SetTokens(w, access, refresh)
	_ = render.Render(w, r, NewTokenResponse(access, refresh))
}
