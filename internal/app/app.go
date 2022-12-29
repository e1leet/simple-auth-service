package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/e1leet/simple-auth-service/internal/config"
	"github.com/e1leet/simple-auth-service/pkg/shutdown"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type App struct {
	srv    http.Server
	cfg    *config.Config
	closer *shutdown.Closer
	logger zerolog.Logger
}

func New(cfg *config.Config) *App {
	logger := log.With().Str("component", "app").Logger()
	logger.Info().Msg("create app")

	r := chi.NewRouter()

	// TODO Create zerolog logger middleware
	r.Use(middleware.AllowContentType("application/json"))

	return &App{
		srv: http.Server{
			Addr:    cfg.Server.Addr,
			Handler: r,
		},
		cfg:    cfg,
		closer: &shutdown.Closer{},
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	a.closer.Add(a.srv.Shutdown)

	go func() {
		a.logger.Info().Msgf("running server on %s", a.srv.Addr)

		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal().Err(err).Send()
		}
	}()

	<-ctx.Done()

	a.logger.Info().Msg("shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		a.cfg.Server.ShutdownTimeout,
	)
	defer cancel()

	if err := a.closer.Close(shutdownCtx); err != nil {
		return fmt.Errorf("closer: %v", err)
	}

	return nil
}
