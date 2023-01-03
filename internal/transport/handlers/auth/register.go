package auth

import (
	"errors"
	"net/http"

	"github.com/e1leet/simple-auth-service/internal/domain/user/dao"
	"github.com/e1leet/simple-auth-service/internal/utils/api"
	"github.com/go-chi/render"
)

// Register godoc
//
//	@Summary		Register user
//	@Description	register user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			data	body	RegisterRequest	true	"data for register"
//	@Success		201
//	@Failure		422	{object}	api.ErrorResponse
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/auth/register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	request := &RegisterRequest{}

	if err := render.Bind(r, request); err != nil {
		h.logger.Error().Err(err).Send()
		_ = render.Render(w, r, api.NewErrorResponse(err.Error(), http.StatusInternalServerError))

		return
	}

	if err := h.authService.Register(r.Context(), request.ToDomain()); err != nil {
		h.logger.Error().Err(err).Send()

		var response render.Renderer

		switch {
		case errors.Is(err, dao.ErrUsernameAlreadyUsed):
			h.logger.Warn().
				Str("username", request.Username).
				Err(err).Send()

			response = api.NewErrorResponse("user already exists", http.StatusConflict)
		default:
			h.logger.Error().Err(err).Send()

			response = api.NewErrorResponse(err.Error(), http.StatusInternalServerError)
		}

		_ = render.Render(w, r, response)

		return
	}

	w.WriteHeader(http.StatusCreated)
}
