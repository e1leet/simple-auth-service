package auth

import (
	"github.com/e1leet/simple-auth-service/internal/config"
	"github.com/e1leet/simple-auth-service/internal/domain/auth/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	authService service.Service
	cfg         *config.Config
	logger      zerolog.Logger
}

func New(authService service.Service, cfg *config.Config) *Handler {
	return &Handler{
		authService: authService,
		cfg:         cfg,
		logger:      log.With().Str("component", "authHandler").Logger(),
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/register", h.Register)
		r.Post("/refresh-tokens", h.RefreshToken)
		r.Delete("/logout", h.Logout)
	})
}
