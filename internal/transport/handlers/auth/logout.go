package auth

import (
	"net/http"

	"github.com/e1leet/simple-auth-service/internal/config"
	"github.com/e1leet/simple-auth-service/internal/domain/jwt/service"
	"github.com/e1leet/simple-auth-service/internal/utils"
	"github.com/e1leet/simple-auth-service/internal/utils/api"
	"github.com/go-chi/render"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := r.Cookie(config.RefreshCookieName)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := h.authService.Logout(r.Context(), refreshToken.Value); err != nil {
		h.logger.Error().Err(err).Send()
		response := api.NewErrorResponse(err.Error(), http.StatusInternalServerError)
		_ = render.Render(w, r, response)
	}

	utils.SetTokens(w, "", service.RefreshToken{})
	w.WriteHeader(http.StatusOK)
}
