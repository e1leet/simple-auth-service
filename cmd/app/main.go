package main

import (
	"context"
	defaultLog "log"
	"os/signal"
	"syscall"

	"github.com/e1leet/simple-auth-service/internal/app"
	"github.com/e1leet/simple-auth-service/internal/config"
	"github.com/e1leet/simple-auth-service/internal/logging"
	"github.com/e1leet/simple-auth-service/internal/utils"
	"github.com/rs/zerolog/log"
)

func main() {
	cfgPath := utils.ConfigPath()

	cfg, err := config.LoadConfig(cfgPath)
	if err != nil {
		defaultLog.Fatal(err)
	}

	logging.ConfigureLogging(cfg.Log.Level)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a := app.New(cfg)

	if err := a.Run(ctx); err != nil {
		stop()
		log.Fatal().Err(err).Send() //nolint:gocritic
	}
}
