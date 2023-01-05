package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/e1leet/simple-auth-service/internal/config"
	"github.com/e1leet/simple-auth-service/internal/domain/auth/service"
	jwtService "github.com/e1leet/simple-auth-service/internal/domain/jwt/service"
	sessionDAO "github.com/e1leet/simple-auth-service/internal/domain/session/dao"
	userDAO "github.com/e1leet/simple-auth-service/internal/domain/user/dao"
	"github.com/e1leet/simple-auth-service/internal/transport/handlers/auth"
	"github.com/e1leet/simple-auth-service/internal/transport/middlewares"
	"github.com/e1leet/simple-auth-service/internal/utils/password/manager"
	"github.com/e1leet/simple-auth-service/pkg/client/postgresql"
	"github.com/e1leet/simple-auth-service/pkg/shutdown"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
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

	closer := &shutdown.Closer{}

	r := chi.NewRouter()

	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middlewares.LoggerMiddleware())

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})

	postgresClient, err := postgresql.NewClient(context.Background(), cfg.Postgres.URI, 5)
	if err != nil {
		logger.Fatal().Err(err).Send()
	}

	closer.Add(func(ctx context.Context) error {
		postgresClient.Close()
		return nil
	})

	usr := userDAO.NewPostgresql(postgresClient)
	session := sessionDAO.NewMemory()
	password := manager.New(cfg.Security.PasswordSalt)
	jwt := jwtService.New(
		session,
		cfg.Security.JWTSecret,
		cfg.Security.AccessExpiresIn,
		cfg.Security.RefreshExpiresIn,
	)
	authService := service.New(jwt, usr, session, password)

	authHandler := auth.New(authService, cfg)
	authHandler.RegisterRoutes(r)

	return &App{
		srv: http.Server{
			Addr:    cfg.Server.Addr,
			Handler: r,
		},
		cfg:    cfg,
		closer: closer,
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
