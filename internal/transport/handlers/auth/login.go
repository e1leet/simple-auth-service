package auth

import (
	"errors"
	"net/http"

	"github.com/e1leet/simple-auth-service/internal/domain/auth/service"
	userDAO "github.com/e1leet/simple-auth-service/internal/domain/user/dao"
	"github.com/e1leet/simple-auth-service/internal/utils"
	"github.com/e1leet/simple-auth-service/internal/utils/api"
	"github.com/go-chi/render"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	request := &LoginRequest{}

	if err := render.Bind(r, request); err != nil {
		h.logger.Error().Err(err).Send()
		_ = render.Render(w, r, api.NewErrorResponse(err.Error(), http.StatusInternalServerError))

		return
	}

	access, refresh, err := h.authService.Login(r.Context(), request.ToDomain())
	if err != nil {
		var response render.Renderer

		switch {
		case errors.Is(err, userDAO.ErrUserNotFound), errors.Is(err, service.ErrIncorrectPassword):
			h.logger.Warn().Err(err).Msg("incorrect username or password")

			response = api.NewErrorResponse("incorrect username or password", http.StatusForbidden)
		default:
			h.logger.Error().Err(err).Send()

			response = api.NewErrorResponse(err.Error(), http.StatusInternalServerError)
		}

		_ = render.Render(w, r, response)

		return
	}

	utils.SetTokens(w, access, refresh)
	_ = render.Render(w, r, NewTokenResponse(access, refresh))
}
